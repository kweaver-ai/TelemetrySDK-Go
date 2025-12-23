package common

//
//import (
//	"go.opentelemetry.io/otel/sdk/instrumentation"
//	"go.opentelemetry.io/otel/sdk/metric/metricdata"
//	"reflect"
//	"testing"
//	"time"
//)
//
//func TestAnyRobotGaugeFromGaugeFloat(t *testing.T) {
//	type args struct {
//		gauge metricdata.Gauge[float64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Gauge
//	}{
//		{
//			"转换Gauge[float64]",
//			args{gauge: metricdata.Gauge[float64]{}},
//			AnyRobotGaugeFromGaugeFloat(metricdata.Gauge[float64]{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotGaugeFromGaugeFloat(tt.args.gauge); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotGaugeFromGaugeFloat() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotGaugeFromGaugeInt(t *testing.T) {
//	type args struct {
//		gauge metricdata.Gauge[int64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Gauge
//	}{
//		{
//			"转换Gauge[int64]",
//			args{gauge: metricdata.Gauge[int64]{}},
//			AnyRobotGaugeFromGaugeInt(metricdata.Gauge[int64]{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotGaugeFromGaugeInt(tt.args.gauge); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotGaugeFromGaugeInt() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotHistogramFromHistogram(t *testing.T) {
//	type args struct {
//		histogram metricdata.Histogram
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Histogram
//	}{
//		{
//			"转换Histogram",
//			args{histogram: metricdata.Histogram{}},
//			AnyRobotHistogramFromHistogram(metricdata.Histogram{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotHistogramFromHistogram(tt.args.histogram); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotHistogramFromHistogram() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotMetricFromMetric(t *testing.T) {
//	type args struct {
//		metric *metricdata.Metrics
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Metrics
//	}{
//		{
//			"转换Data为Gauge[int64]",
//			args{metric: &metricdata.Metrics{Data: metricdata.Gauge[int64]{}}},
//			AnyRobotMetricFromMetric(&metricdata.Metrics{Data: metricdata.Gauge[int64]{}}),
//		},
//		{
//			"转换Data为Gauge[float64]",
//			args{metric: &metricdata.Metrics{Data: metricdata.Gauge[float64]{}}},
//			AnyRobotMetricFromMetric(&metricdata.Metrics{Data: metricdata.Gauge[float64]{}}),
//		},
//		{
//			"转换Data为Sum[int64]",
//			args{metric: &metricdata.Metrics{Data: metricdata.Sum[int64]{}}},
//			AnyRobotMetricFromMetric(&metricdata.Metrics{Data: metricdata.Sum[int64]{}}),
//		},
//		{
//			"转换Data为Sum[float64]",
//			args{metric: &metricdata.Metrics{Data: metricdata.Sum[float64]{}}},
//			AnyRobotMetricFromMetric(&metricdata.Metrics{Data: metricdata.Sum[float64]{}}),
//		},
//		{
//			"转换Data为Histogram",
//			args{metric: &metricdata.Metrics{Data: metricdata.Histogram{}}},
//			AnyRobotMetricFromMetric(&metricdata.Metrics{Data: metricdata.Histogram{}}),
//		},
//		{
//			"测试空返回值",
//			args{metric: &metricdata.Metrics{}},
//			nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotMetricFromMetric(tt.args.metric); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotMetricFromMetric() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotMetricsFromMetrics(t *testing.T) {
//	type args struct {
//		metrics []metricdata.Metrics
//	}
//	tests := []struct {
//		name string
//		args args
//		want []*Metrics
//	}{
//		{
//			"转换Metrics",
//			args{metrics: []metricdata.Metrics{{Data: metricdata.Histogram{}}}},
//			AnyRobotMetricsFromMetrics([]metricdata.Metrics{{Data: metricdata.Histogram{}}}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotMetricsFromMetrics(tt.args.metrics); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotMetricsFromMetrics() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotScopeMetricFromScopeMetric(t *testing.T) {
//	type args struct {
//		scopeMetric *metricdata.ScopeMetrics
//	}
//	tests := []struct {
//		name string
//		args args
//		want *ScopeMetric
//	}{
//		{
//			"转换ScopeMetric",
//			args{scopeMetric: &metricdata.ScopeMetrics{}},
//			AnyRobotScopeMetricFromScopeMetric(&metricdata.ScopeMetrics{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotScopeMetricFromScopeMetric(tt.args.scopeMetric); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotScopeMetricFromScopeMetric() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotScopeMetricsFromScopeMetrics(t *testing.T) {
//	type args struct {
//		scopeMetrics []metricdata.ScopeMetrics
//	}
//	tests := []struct {
//		name string
//		args args
//		want []*ScopeMetric
//	}{
//		{
//			"转换ScopeMetrics",
//			args{scopeMetrics: []metricdata.ScopeMetrics{{Scope: instrumentation.Scope{Name: "ar"}}}},
//			AnyRobotScopeMetricsFromScopeMetrics([]metricdata.ScopeMetrics{{Scope: instrumentation.Scope{Name: "ar"}}}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotScopeMetricsFromScopeMetrics(tt.args.scopeMetrics); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotScopeMetricsFromScopeMetrics() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotSumFromSumFloat(t *testing.T) {
//	type args struct {
//		sum metricdata.Sum[float64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Sum
//	}{
//		{
//			"转换Sum[float64]",
//			args{sum: metricdata.Sum[float64]{}},
//			AnyRobotSumFromSumFloat(metricdata.Sum[float64]{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotSumFromSumFloat(tt.args.sum); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotSumFromSumFloat() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAnyRobotSumFromSumInt(t *testing.T) {
//	type args struct {
//		sum metricdata.Sum[int64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Sum
//	}{
//		{
//			"转换Sum[int64]",
//			args{sum: metricdata.Sum[int64]{}},
//			AnyRobotSumFromSumInt(metricdata.Sum[int64]{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := AnyRobotSumFromSumInt(tt.args.sum); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("AnyRobotSumFromSumInt() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestFloatDataPoint(t *testing.T) {
//	type args struct {
//		dp metricdata.DataPoint[float64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want *DataPoint
//	}{
//		{
//			"转换DataPoint",
//			args{dp: metricdata.DataPoint[float64]{}},
//			FloatDataPoint(metricdata.DataPoint[float64]{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := FloatDataPoint(tt.args.dp); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("FloatDataPoint() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestFloatDataPoints(t *testing.T) {
//	type args struct {
//		dps []metricdata.DataPoint[float64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want []*DataPoint
//	}{
//		{
//			"转换DataPoints",
//			args{dps: []metricdata.DataPoint[float64]{{Value: 1.0}}},
//			FloatDataPoints([]metricdata.DataPoint[float64]{{Value: 1.0}}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := FloatDataPoints(tt.args.dps); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("FloatDataPoints() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestHistogramDataPoints(t *testing.T) {
//	type args struct {
//		hdps []metricdata.HistogramDataPoint
//	}
//	tests := []struct {
//		name string
//		args args
//		want []*HistogramDataPoint
//	}{
//		{
//			"转换Histogram_DataPoints",
//			args{hdps: []metricdata.HistogramDataPoint{{Count: 1}}},
//			HistogramDataPoints([]metricdata.HistogramDataPoint{{Count: 1}}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := HistogramDataPoints(tt.args.hdps); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("HistogramDataPoints() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestIntDataPoint(t *testing.T) {
//	type args struct {
//		dp metricdata.DataPoint[int64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want *DataPoint
//	}{
//		{
//			"转换DataPoint",
//			args{dp: metricdata.DataPoint[int64]{}},
//			IntDataPoint(metricdata.DataPoint[int64]{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := IntDataPoint(tt.args.dp); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("IntDataPoint() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestIntDataPoints(t *testing.T) {
//	type args struct {
//		dps []metricdata.DataPoint[int64]
//	}
//	tests := []struct {
//		name string
//		args args
//		want []*DataPoint
//	}{
//		{
//			"转换DataPoints",
//			args{dps: []metricdata.DataPoint[int64]{{Value: 1}}},
//			IntDataPoints([]metricdata.DataPoint[int64]{{Value: 1}}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := IntDataPoints(tt.args.dps); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("IntDataPoints() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestOmitZeroTime(t *testing.T) {
//	type args struct {
//		startTime time.Time
//	}
//	var timing = time.Now()
//	tests := []struct {
//		name string
//		args args
//		want *time.Time
//	}{
//		{
//			"时间零值",
//			args{startTime: time.Time{}},
//			nil,
//		},
//		{
//			"时间非零值",
//			args{startTime: timing},
//			&timing,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := OmitZeroTime(tt.args.startTime); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("OmitZeroTime() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestSingleHistogramDataPoint(t *testing.T) {
//	type args struct {
//		hdp metricdata.HistogramDataPoint
//	}
//	tests := []struct {
//		name string
//		args args
//		want *HistogramDataPoint
//	}{
//		{
//			"转换HistogramDataPoint",
//			args{hdp: metricdata.HistogramDataPoint{}},
//			SingleHistogramDataPoint(metricdata.HistogramDataPoint{}),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := SingleHistogramDataPoint(tt.args.hdp); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("SingleHistogramDataPoint() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
