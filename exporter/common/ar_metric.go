package common

import (
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// AnyRobotMetric 自定义Metric，改造了Resource、ScopeMetrics。
type AnyRobotMetric struct {
	Resource     *Resource      `json:"Resource"`
	ScopeMetrics []*ScopeMetric `json:"ScopeMetrics"`
}

// AnyRobotMetricFromResourceMetric 单条 *metricdata.ResourceMetrics 转换为 *AnyRobotMetric 。
func AnyRobotMetricFromResourceMetric(resourceMetric *metricdata.ResourceMetrics) *AnyRobotMetric {
	return &AnyRobotMetric{
		Resource:     AnyRobotResourceFromResource(resourceMetric.Resource),
		ScopeMetrics: AnyRobotScopeMetricsFromScopeMetrics(resourceMetric.ScopeMetrics),
	}
}

// AnyRobotMetricsFromResourceMetrics 批量 []*metricdata.ResourceMetrics 转换为 []*AnyRobotMetric 。
func AnyRobotMetricsFromResourceMetrics(resourceMetrics []*metricdata.ResourceMetrics) []*AnyRobotMetric {
	arMetric := make([]*AnyRobotMetric, 0, len(resourceMetrics))
	for _, value := range resourceMetrics {
		arMetric = append(arMetric, AnyRobotMetricFromResourceMetric(value))
	}
	return arMetric
}
