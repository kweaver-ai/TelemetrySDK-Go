package common

//
//import (
//	"go.opentelemetry.io/otel/sdk/metric/metricdata"
//	"reflect"
//	"testing"
//)
//
//func TestAnyRobotMetricFromResourceMetric(t *testing.T) {
//	type args struct {
//		resourceMetric *metricdata.ResourceMetrics
//	}
//	tests := []struct {
//		name string
//		args args
//		want *AnyRobotMetric
//	}{
//		{
//			"测试单条*metricdata.ResourceMetrics转换",
//			args{&metricdata.ResourceMetrics{}},
//			AnyRobotMetricFromResourceMetric(&metricdata.ResourceMetrics{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotMetricFromResourceMetric(tt.args.resourceMetric); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotMetricFromResourceMetric() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotMetricsFromResourceMetrics(t *testing.T) {
//	type args struct {
//		resourceMetrics []*metricdata.ResourceMetrics
//	}
//	tests := []struct {
//		name string
//		args args
//		want []*AnyRobotMetric
//	}{
//		{
//			"测试批量[]*metricdata.ResourceMetrics转换",
//			args{[]*metricdata.ResourceMetrics{&metricdata.ResourceMetrics{}}},
//			AnyRobotMetricsFromResourceMetrics([]*metricdata.ResourceMetrics{&metricdata.ResourceMetrics{}}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotMetricsFromResourceMetrics(tt.args.resourceMetrics); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotMetricsFromResourceMetrics() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
