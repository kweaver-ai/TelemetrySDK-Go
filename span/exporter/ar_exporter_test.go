package exporter

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func channelWithStop() chan struct{} {
	stopCh := make(chan struct{})
	close(stopCh)
	return stopCh
}

func TestGetRealTimeExporter(t *testing.T) {
	tests := []struct {
		name string
		want LogExporter
	}{
		{
			"",
			&exporter{
				name:     "RealTimeExporter",
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRealTimeExporter(); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
				t.Errorf("GetDefaultExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExporterExportLogs(t *testing.T) {
	type fields struct {
		name     string
		stopCh   chan struct{}
		stopOnce *sync.Once
	}
	type args struct {
		ctx  context.Context
		logs []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"common",
			fields{
				name:     "common",
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				ctx:  contextWithDone(),
				logs: []byte("[{\"@timestamp\":\"2022-12-16T10:04:15.759143319+08:00\",\"Attributes\":{\"Type\":\"test\"},\"Body\":{\"Message\":\"this is a test\"},\"Link\":{\"SpanId\":\"0000000000000000\",\"TraceId\":\"00000000000000000000000000000000\"},\"Resource\":{\"service\":{\"name\":\"test\",\"version\":\"1.0.1\"},\"telemetry\":{\"sdk\":{\"language\":\"go\",\"name\":\"opentelemetry\",\"version\":\"1.9.0\"}}},\"SeverityText\":\"Info\",\"category\":\"log\"}]"),
			},
			true,
		}, {
			"err",
			fields{
				name:     "err",
				stopCh:   channelWithStop(),
				stopOnce: &sync.Once{},
			},
			args{
				ctx:  context.Background(),
				logs: []byte("[{\"@timestamp\":\"2022-12-16T10:04:15.759143319+08:00\",\"Attributes\":{\"Type\":\"test\"},\"Body\":{\"Message\":\"this is a test\"},\"Link\":{\"SpanId\":\"0000000000000000\",\"TraceId\":\"00000000000000000000000000000000\"},\"Resource\":{\"service\":{\"name\":\"test\",\"version\":\"1.0.1\"},\"telemetry\":{\"sdk\":{\"language\":\"go\",\"name\":\"opentelemetry\",\"version\":\"1.9.0\"}}},\"SeverityText\":\"Info\",\"category\":\"log\"}]"),
			},
			false,
		}, {
			"err1",
			fields{
				name:     "err1",
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				ctx:  context.Background(),
				logs: []byte("[{\"@timestamp\":\"2022-12-16T10:04:15.759143319+08:00\",\"Attributes\":{\"Type\":\"test\"},\"Body\":{\"Message\":\"this is a test\"},\"Link\":{\"SpanId\":\"0000000000000000\",\"TraceId\":\"00000000000000000000000000000000\"},\"Resource\":{\"service\":{\"name\":\"test\",\"version\":\"1.0.1\"},\"telemetry\":{\"sdk\":{\"language\":\"go\",\"name\":\"opentelemetry\",\"version\":\"1.9.0\"}}},\"SeverityText\":\"Info\",\"category\":\"log\"}]"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &exporter{
				name:     tt.fields.name,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if err := e.ExportLogs(tt.args.ctx, tt.args.logs); (err != nil) != tt.wantErr {
				t.Errorf("ExportEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExporterName(t *testing.T) {
	type fields struct {
		name     string
		stopCh   chan struct{}
		stopOnce *sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				name:     "NAME",
				stopCh:   nil,
				stopOnce: &sync.Once{},
			},
			"NAME",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &exporter{
				name:     tt.fields.name,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if got := e.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	type fields struct {
		name     string
		stopCh   chan struct{}
		stopOnce *sync.Once
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{
				name:     "",
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{contextWithDone()},
			true,
		}, {
			"",
			fields{
				name:     "",
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{context.Background()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &exporter{
				name:     tt.fields.name,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if err := e.Shutdown(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
