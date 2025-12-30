package ar_trace

import (
	"context"
	"reflect"
	"testing"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/public"
)

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestExporterExportSpans(t *testing.T) {
	type fields struct {
		Exporter *public.Exporter
	}
	type args struct {
		ctx    context.Context
		traces []sdktrace.ReadOnlySpan
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"发送空数据",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx:    context.Background(),
				traces: nil,
			},
			false,
		},
		{
			"发送Trace",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx:    context.Background(),
				traces: make([]sdktrace.ReadOnlySpan, 1),
			},
			false,
		},
		{
			"已关闭的Exporter不能发Trace",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx:    contextWithDone(),
				traces: make([]sdktrace.ReadOnlySpan, 1),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &TraceExporter{
				Exporter: tt.fields.Exporter,
			}
			if err := e.ExportSpans(tt.args.ctx, tt.args.traces); (err != nil) != tt.wantErr {
				t.Errorf("ExportSpans() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewExporter(t *testing.T) {
	type args struct {
		c public.Client
	}
	tests := []struct {
		name string
		args args
		want *TraceExporter
	}{
		{
			"StdoutClient的TraceExporter",
			args{c: public.NewStdoutClient("./AnyRobotTrace.txt")},
			NewExporter(public.NewStdoutClient("./AnyRobotTrace.txt")),
		},
		{
			"HTTPClient的TraceExporter",
			args{c: public.NewHTTPClient()},
			NewExporter(public.NewHTTPClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.c); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTraceResource(t *testing.T) {
	tests := []struct {
		name string
		want *sdkresource.Resource
	}{
		{
			"Trace的默认Resource",
			TraceResource(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TraceResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TraceResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
