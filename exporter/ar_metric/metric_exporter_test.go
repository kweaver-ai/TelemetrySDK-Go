package ar_metric

//
//import (
//	"context"
//	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/public"
//	"go.opentelemetry.io/otel/sdk/metric"
//	"go.opentelemetry.io/otel/sdk/metric/aggregation"
//	"go.opentelemetry.io/otel/sdk/metric/metricdata"
//	sdkresource "go.opentelemetry.io/otel/sdk/resource"
//	"reflect"
//	"testing"
//)
//
//func contextWithDone() context.Context {
//	ctx, cancel := context.WithCancel(context.Background())
//	cancel()
//	return ctx
//}
//
//func TestMetricExporterAggregation(t *testing.T) {
//	type fields struct {
//		Exporter *public.Exporter
//	}
//	type args struct {
//		k metric.InstrumentKind
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   aggregation.Aggregation
//	}{
//		{
//			"选择聚合类型",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{k: metric.InstrumentKind(1)},
//			NewExporter(public.NewStdoutClient("")).Aggregation(1),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &MetricExporter{
//				Exporter: tt.fields.Exporter,
//			}
//			if got := e.Aggregation(tt.args.k); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Aggregation() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMetricExporterExport(t *testing.T) {
//	type fields struct {
//		Exporter *public.Exporter
//	}
//	type args struct {
//		ctx  context.Context
//		data metricdata.ResourceMetrics
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			"正常导出数据",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{
//				ctx:  context.Background(),
//				data: metricdata.ResourceMetrics{},
//			},
//			false,
//		},
//		{
//			"已关闭的Exporter不能导出数据",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{
//				ctx:  contextWithDone(),
//				data: metricdata.ResourceMetrics{},
//			},
//			true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &MetricExporter{
//				Exporter: tt.fields.Exporter,
//			}
//			if err := e.Export(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
//				t.Errorf("Export() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestMetricExporterExportMetrics(t *testing.T) {
//	type fields struct {
//		Exporter *public.Exporter
//	}
//	type args struct {
//		ctx     context.Context
//		metrics []*metricdata.ResourceMetrics
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			"正常批量发送数据",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{
//				ctx:     context.Background(),
//				metrics: []*metricdata.ResourceMetrics{},
//			},
//			false,
//		},
//		{
//			"已关闭的Exporter不能批量发送数据",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{
//				ctx:     contextWithDone(),
//				metrics: []*metricdata.ResourceMetrics{{}},
//			},
//			true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &MetricExporter{
//				Exporter: tt.fields.Exporter,
//			}
//			if err := e.ExportMetrics(tt.args.ctx, tt.args.metrics); (err != nil) != tt.wantErr {
//				t.Errorf("ExportMetrics() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestMetricExporterForceFlush(t *testing.T) {
//	type fields struct {
//		Exporter *public.Exporter
//	}
//	type args struct {
//		ctx context.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			"正常强制发送",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{
//				ctx: context.Background(),
//			},
//			false,
//		},
//		{
//			"已关闭的Exporter不能强制发送",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{
//				ctx: contextWithDone(),
//			},
//			true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &MetricExporter{
//				Exporter: tt.fields.Exporter,
//			}
//			if err := e.ForceFlush(tt.args.ctx); (err != nil) != tt.wantErr {
//				t.Errorf("ForceFlush() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestMetricExporterTemporality(t *testing.T) {
//	type fields struct {
//		Exporter *public.Exporter
//	}
//	type args struct {
//		k metric.InstrumentKind
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   metricdata.Temporality
//	}{
//		{
//			"选择计量方式",
//			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
//			args{
//				k: metric.InstrumentKind(1),
//			},
//			metricdata.Temporality(1),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &MetricExporter{
//				Exporter: tt.fields.Exporter,
//			}
//			if got := e.Temporality(tt.args.k); got != tt.want {
//				t.Errorf("Temporality() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMetricResource(t *testing.T) {
//	tests := []struct {
//		name string
//		want *sdkresource.Resource
//	}{
//		{
//			"Metric的默认Resource",
//			MetricResource(),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := MetricResource(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("MetricResource() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNewExporter(t *testing.T) {
//	type args struct {
//		c public.Client
//	}
//	tests := []struct {
//		name string
//		args args
//		want *MetricExporter
//	}{
//		{
//			"StdoutClient的MetricExporter",
//			args{c: public.NewStdoutClient("./AnyRobotMetric.txt")},
//			NewExporter(public.NewStdoutClient("./AnyRobotMetric.txt")),
//		},
//		{
//			"HTTPClient的MetricExporter",
//			args{c: public.NewHTTPClient()},
//			NewExporter(public.NewHTTPClient()),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewExporter(tt.args.c); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
//				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
