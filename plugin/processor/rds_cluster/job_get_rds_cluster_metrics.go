package rds_cluster

import (
	"fmt"
	types2 "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/kaytu-io/kaytu/pkg/utils"
	"github.com/kaytu-io/plugin-aws/plugin/aws"
	preferences2 "github.com/kaytu-io/plugin-aws/plugin/preferences"
	"strings"
	"time"
)

type GetRDSClusterMetricsJob struct {
	cluster   types.DBCluster
	instances []types.DBInstance
	region    string

	processor *Processor
}

func NewGetRDSInstanceMetricsJob(processor *Processor, region string, cluster types.DBCluster, instances []types.DBInstance) *GetRDSClusterMetricsJob {
	return &GetRDSClusterMetricsJob{
		processor: processor,
		cluster:   cluster,
		instances: instances,
		region:    region,
	}
}

func (j *GetRDSClusterMetricsJob) Id() string {
	return fmt.Sprintf("get_rds_cluster_metrics_%s", *j.cluster.DBClusterIdentifier)
}
func (j *GetRDSClusterMetricsJob) Description() string {
	return fmt.Sprintf("Getting metrics of %s", *j.cluster.DBClusterIdentifier)
}
func (j *GetRDSClusterMetricsJob) Run() error {
	startTime := time.Now().Add(-24 * 7 * time.Hour)
	endTime := time.Now()

	allMetrics := map[string]map[string][]types2.Datapoint{}
	for _, instance := range j.instances {
		isAurora := j.cluster.DBClusterIdentifier != nil && strings.Contains(strings.ToLower(*j.cluster.Engine), "aurora")
		allMetrics[utils.HashString(*instance.DBInstanceIdentifier)] = map[string][]types2.Datapoint{}
		cwMetrics, err := j.processor.metricProvider.GetMetrics(
			j.region,
			"AWS/RDS",
			[]string{
				"CPUUtilization",
				"FreeableMemory",
				"FreeStorageSpace",
			},
			map[string][]string{
				"DBInstanceIdentifier": {*instance.DBInstanceIdentifier},
			},
			startTime, endTime,
			time.Hour,
			[]types2.Statistic{
				types2.StatisticAverage,
				types2.StatisticMaximum,
				types2.StatisticMinimum,
			},
		)
		if err != nil {
			return err
		}

		cwPerSecondMetrics, err := j.processor.metricProvider.GetMetrics(
			j.region,
			"AWS/RDS",
			[]string{
				"NetworkReceiveThroughput",
				"NetworkTransmitThroughput",
			},
			map[string][]string{
				"DBInstanceIdentifier": {*instance.DBInstanceIdentifier},
			},
			startTime, endTime,
			time.Hour,
			[]types2.Statistic{
				types2.StatisticSum,
			},
		)
		if err != nil {
			return err
		}
		for k, val := range cwPerSecondMetrics {
			cwPerSecondMetrics[k] = aws.GetDatapointsAvgFromSum(val, int32(time.Hour/time.Second))
		}

		var volumeThroughput map[string][]types2.Datapoint
		var iopsMetrics map[string][]types2.Datapoint
		var clusterMetrics map[string][]types2.Datapoint
		if !isAurora {
			volumeThroughput, err = j.processor.metricProvider.GetMetrics(
				j.region,
				"AWS/RDS",
				[]string{
					"ReadThroughput",
					"WriteThroughput",
				},
				map[string][]string{
					"DBInstanceIdentifier": {*instance.DBInstanceIdentifier},
				},
				startTime, endTime,
				time.Hour,
				[]types2.Statistic{
					types2.StatisticSum,
				},
			)
			if err != nil {
				return err
			}
			iopsMetrics, err = j.processor.metricProvider.GetDayByDayMetrics(
				j.region,
				"AWS/RDS",
				[]string{
					"ReadIOPS",
					"WriteIOPS",
				},
				map[string][]string{
					"DBInstanceIdentifier": {*instance.DBInstanceIdentifier},
				},
				7,
				time.Minute,
				[]types2.Statistic{
					types2.StatisticSum,
				},
			)
			if err != nil {
				return err
			}
		} else {
			volumeThroughput, err = j.processor.metricProvider.GetMetrics(
				j.region,
				"AWS/RDS",
				[]string{
					"ReadThroughput",
					"WriteThroughput",
				},
				map[string][]string{
					"DBClusterIdentifier": {*instance.DBClusterIdentifier},
				},
				startTime, endTime,
				time.Hour,
				[]types2.Statistic{
					types2.StatisticSum,
				},
			)
			if err != nil {
				return err
			}
			iopsMetrics, err = j.processor.metricProvider.GetDayByDayMetrics(
				j.region,
				"AWS/RDS",
				[]string{
					"ReadIOPS",
					"WriteIOPS",
				},
				map[string][]string{
					"DBClusterIdentifier": {*instance.DBInstanceIdentifier},
				},
				7,
				time.Minute,
				[]types2.Statistic{
					types2.StatisticSum,
				},
			)
			if err != nil {
				return err
			}
			clusterMetrics, err = j.processor.metricProvider.GetMetrics(
				j.region,
				"AWS/RDS",
				[]string{
					"VolumeBytesUsed",
				},
				map[string][]string{
					"DBClusterIdentifier": {*instance.DBClusterIdentifier},
				},
				startTime, endTime,
				time.Hour,
				[]types2.Statistic{
					types2.StatisticAverage,
					types2.StatisticMaximum,
				},
			)
			if err != nil {
				return err
			}
		}
		for k, val := range volumeThroughput {
			volumeThroughput[k] = aws.GetDatapointsAvgFromSum(val, int32(time.Hour/time.Second))
		}
		for k, val := range iopsMetrics {
			iopsMetrics[k] = aws.GetDatapointsAvgFromSum(val, int32(time.Minute/time.Second))
		}

		hashedIdentifier := utils.HashString(*instance.DBInstanceIdentifier)
		for k, v := range cwMetrics {
			allMetrics[hashedIdentifier][k] = v
		}
		for k, v := range cwPerSecondMetrics {
			allMetrics[hashedIdentifier][k] = v
		}
		for k, v := range iopsMetrics {
			allMetrics[hashedIdentifier][k] = v
		}
		for k, v := range volumeThroughput {
			allMetrics[hashedIdentifier][k] = v
		}
		if clusterMetrics != nil {
			for k, v := range clusterMetrics {
				allMetrics[hashedIdentifier][k] = v
			}
		}
	}

	oi := RDSClusterItem{
		Cluster:             j.cluster,
		Instances:           j.instances,
		Region:              j.region,
		OptimizationLoading: true,
		Preferences:         preferences2.DefaultRDSPreferences,
		Skipped:             false,
		LazyLoadingEnabled:  false,
		SkipReason:          "",
		Metrics:             allMetrics,
	}

	j.processor.items[*oi.Cluster.DBClusterIdentifier] = oi
	j.processor.publishOptimizationItem(oi.ToOptimizationItem())
	if !oi.Skipped && !oi.LazyLoadingEnabled {
		j.processor.jobQueue.Push(NewOptimizeRDSJob(j.processor, oi))
	}
	return nil
}
