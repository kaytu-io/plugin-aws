package rds_cluster

import (
	"fmt"
	"github.com/kaytu-io/kaytu/pkg/plugin/proto/src/golang"
	"github.com/kaytu-io/kaytu/pkg/plugin/sdk"
	"github.com/kaytu-io/kaytu/pkg/style"
	"github.com/kaytu-io/kaytu/pkg/utils"
	"github.com/kaytu-io/plugin-aws/plugin/aws"
	"github.com/kaytu-io/plugin-aws/plugin/kaytu"
	"github.com/kaytu-io/plugin-aws/plugin/processor/ec2_instance"
	util "github.com/kaytu-io/plugin-aws/utils"
	"strings"
	"sync/atomic"
)

type Processor struct {
	provider                *aws.AWS
	metricProvider          *aws.CloudWatch
	identification          map[string]string
	items                   util.ConcurrentMap[string, RDSClusterItem]
	publishOptimizationItem func(item *golang.ChartOptimizationItem)
	publishResultSummary    func(summary *golang.ResultSummary)
	kaytuAcccessToken       string
	jobQueue                *sdk.JobQueue
	configuration           *kaytu.Configuration
	lazyloadCounter         *atomic.Uint32
	observabilityDays       int

	summary            *util.ConcurrentMap[string, ec2_instance.EC2InstanceSummary]
	defaultPreferences []*golang.PreferenceItem
}

func NewProcessor(provider *aws.AWS, metricProvider *aws.CloudWatch, identification map[string]string, publishOptimizationItem func(item *golang.ChartOptimizationItem), publishResultSummary func(summary *golang.ResultSummary), kaytuAcccessToken string, jobQueue *sdk.JobQueue, configurations *kaytu.Configuration, lazyloadCounter *atomic.Uint32, observabilityDays int, summary *util.ConcurrentMap[string, ec2_instance.EC2InstanceSummary], preferences []*golang.PreferenceItem) *Processor {
	r := &Processor{
		provider:                provider,
		metricProvider:          metricProvider,
		identification:          identification,
		items:                   util.NewMap[string, RDSClusterItem](),
		publishOptimizationItem: publishOptimizationItem,
		publishResultSummary:    publishResultSummary,
		kaytuAcccessToken:       kaytuAcccessToken,
		jobQueue:                jobQueue,
		configuration:           configurations,
		lazyloadCounter:         lazyloadCounter,
		observabilityDays:       observabilityDays,
		summary:                 summary,
		defaultPreferences:      preferences,
	}

	jobQueue.Push(NewListAllRegionsJob(r))
	return r
}

func (m *Processor) ReEvaluate(id string, items []*golang.PreferenceItem) {
	v, _ := m.items.Get(id)
	v.Preferences = items
	m.items.Set(id, v)
	m.jobQueue.Push(NewOptimizeRDSJob(m, v))
}

func (m *Processor) ExportNonInteractive() *golang.NonInteractiveExport {
	return nil
}

func (m *Processor) ExportCsv() []*golang.CSVRow {
	var rows []*golang.CSVRow

	m.summary.Range(func(id string, _ ec2_instance.EC2InstanceSummary) bool {
		if _, ok := m.items.Get(id); !ok {
			fmt.Println("Skipping item", id)
			return true
		}
		cluster, _ := m.items.Get(id)
		for _, i := range cluster.Instances {
			var platform string
			if i.Engine != nil {
				platform = *i.Engine
			}
			hashedId := utils.HashString(*i.DBInstanceIdentifier)
			rightSizing := cluster.Wastage.RightSizing[hashedId]

			var computeAdditionalDetails []string
			var computeRightSizingCost, computeSaving, computeRecSpec string
			if rightSizing.Recommended != nil {
				computeRightSizingCost = utils.FormatPriceFloat(rightSizing.Recommended.ComputeCost)
				computeSaving = utils.FormatPriceFloat(rightSizing.Current.ComputeCost - rightSizing.Recommended.ComputeCost)
				computeRecSpec = rightSizing.Recommended.InstanceType

				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("Instance Size:: Current: %s - Recommended: %s", rightSizing.Current.InstanceType,
						rightSizing.Recommended.InstanceType))
				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("Engine:: Current: %s - Recommended: %s", rightSizing.Current.Engine,
						rightSizing.Recommended.Engine))
				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("Engine Version:: Current: %s - Recommended: %s", rightSizing.Current.EngineVersion,
						rightSizing.Recommended.EngineVersion))
				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("Cluster Type:: Current: %s - Recommended: %s", rightSizing.Current.ClusterType,
						rightSizing.Recommended.ClusterType))
				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("vCPU:: Current: %d - Avg: %s - Recommended: %d", rightSizing.Current.VCPU,
						utils.Percentage(rightSizing.VCPU.Avg), rightSizing.Recommended.VCPU))
				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("Processor(s):: Current: %s - Recommended: %s", rightSizing.Current.Processor,
						rightSizing.Recommended.Processor))
				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("Architecture:: Current: %s - Recommended: %s", rightSizing.Current.Architecture,
						rightSizing.Recommended.Architecture))
				computeAdditionalDetails = append(computeAdditionalDetails,
					fmt.Sprintf("Memory:: Current: %d GB - Avg: %s - Recommended: %d GB", rightSizing.Current.MemoryGb,
						utils.MemoryUsagePercentageByFreeSpace(rightSizing.FreeMemoryBytes.Avg, float64(rightSizing.Current.MemoryGb)),
						rightSizing.Recommended.MemoryGb))
			}
			computeRow := []string{m.identification["account"], cluster.Region, "RDS Instance Compute", fmt.Sprintf("%s-compute", *i.DBInstanceIdentifier),
				*i.DBInstanceIdentifier, platform, "730 hours", utils.FormatPriceFloat(rightSizing.Current.ComputeCost),
				computeRightSizingCost, computeSaving, rightSizing.Current.InstanceType, computeRecSpec, *i.DBInstanceIdentifier,
				rightSizing.Description, strings.Join(computeAdditionalDetails, "---")}
			rows = append(rows, &golang.CSVRow{Row: computeRow})

			var storageAdditionalDetails []string
			var storageRightSizingCost, storageSaving, storageRecSpec string
			if rightSizing.Recommended != nil {
				storageRightSizingCost = utils.FormatPriceFloat(rightSizing.Recommended.StorageCost)
				storageSaving = utils.FormatPriceFloat(rightSizing.Current.StorageCost - rightSizing.Recommended.StorageCost)
				storageRecSpec = fmt.Sprintf("%s/%s/%s IOPS", *rightSizing.Recommended.StorageType,
					utils.SizeByteToGB(rightSizing.Recommended.StorageSize), utils.PInt32ToString(rightSizing.Recommended.StorageIops))

				storageAdditionalDetails = append(storageAdditionalDetails,
					fmt.Sprintf("Type:: Current: %s - Recommended: %s", utils.PString(rightSizing.Current.StorageType),
						utils.PString(rightSizing.Recommended.StorageType)))
				storageAdditionalDetails = append(storageAdditionalDetails,
					fmt.Sprintf("Size:: Current: %s - Avg : %s - Recommended: %s", utils.SizeByteToGB(rightSizing.Current.StorageSize),
						utils.StorageUsagePercentageByFreeSpace(rightSizing.FreeStorageBytes.Avg, rightSizing.Current.StorageSize),
						utils.SizeByteToGB(rightSizing.Current.StorageSize)))
				storageAdditionalDetails = append(storageAdditionalDetails,
					fmt.Sprintf("IOPS:: Current: %s - Avg: %s - Recommended: %s", utils.PInt32ToString(rightSizing.Current.StorageIops),
						fmt.Sprintf("%s io/s", utils.PFloat64ToString(rightSizing.StorageIops.Avg)),
						utils.PInt32ToString(rightSizing.Recommended.StorageIops)))
				storageAdditionalDetails = append(storageAdditionalDetails,
					fmt.Sprintf("Throughput:: Current: %s - Avg: %s - Recommended: %s", utils.PStorageThroughputMbps(rightSizing.Current.StorageThroughput),
						utils.PStorageThroughputMbps(rightSizing.StorageThroughput.Avg), utils.PStorageThroughputMbps(rightSizing.Recommended.StorageThroughput)))
				storageAdditionalDetails = append(storageAdditionalDetails,
					fmt.Sprintf("VolumeTypeChange:: %v", utils.PString(rightSizing.Current.StorageType) != utils.PString(rightSizing.Recommended.StorageType)))
				storageAdditionalDetails = append(storageAdditionalDetails,
					fmt.Sprintf("VolumeSizeChange:: %v", *rightSizing.Current.StorageSize != *rightSizing.Recommended.StorageSize))
			}
			storageRow := []string{m.identification["account"], cluster.Region, "RDS Instance Storage", fmt.Sprintf("%s-storage", *i.DBInstanceIdentifier),
				*i.DBInstanceIdentifier, "N/A", "730 hours", utils.FormatPriceFloat(rightSizing.Current.StorageCost),
				storageRightSizingCost, storageSaving, fmt.Sprintf("%s/%s/%s IOPS", *rightSizing.Current.StorageType,
					utils.SizeByteToGB(rightSizing.Current.StorageSize), utils.PInt32ToString(rightSizing.Current.StorageIops)), storageRecSpec, *i.DBInstanceIdentifier,
				rightSizing.Description, strings.Join(storageAdditionalDetails, "---")}
			rows = append(rows, &golang.CSVRow{Row: storageRow})
		}
		return true
	})
	return rows
}

func (m *Processor) HasItem(id string) bool {
	_, ok := m.items.Get(id)
	return ok
}

func (m *Processor) ResultsSummary() *golang.ResultSummary {
	summary := &golang.ResultSummary{}
	var totalCost, savings float64

	m.summary.Range(func(_ string, item ec2_instance.EC2InstanceSummary) bool {
		totalCost += item.CurrentRuntimeCost
		savings += item.Savings
		return true
	})
	summary.Message = fmt.Sprintf("Current runtime cost: %s, Savings: %s",
		style.CostStyle.Render(fmt.Sprintf("%s", utils.FormatPriceFloat(totalCost))), style.SavingStyle.Render(fmt.Sprintf("%s", utils.FormatPriceFloat(savings))))
	return summary
}

func (m *Processor) UpdateSummary(itemId string) {
	i, ok := m.items.Get(itemId)
	if ok && i.Wastage.RightSizing != nil {
		totalSaving := 0.0
		totalCurrentCost := 0.0

		for _, instance := range i.Wastage.RightSizing {
			totalSaving += instance.Current.ComputeCost - instance.Recommended.ComputeCost
			totalCurrentCost += instance.Current.ComputeCost
			totalSaving += instance.Current.StorageCost - instance.Recommended.StorageCost
			totalCurrentCost += instance.Current.StorageCost
		}

		m.summary.Set(itemId, ec2_instance.EC2InstanceSummary{
			CurrentRuntimeCost: totalCurrentCost,
			Savings:            totalSaving,
		})
	}
	m.publishResultSummary(m.ResultsSummary())
}
