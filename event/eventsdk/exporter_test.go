package eventsdk

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestGetDefaultExporter(t *testing.T) {
	tests := []struct {
		name string
		want EventExporter
	}{
		{
			"",
			&exporter{
				name:     "DefaultExporter",
				stopCh:   make(chan struct{}),
				stopOnce: sync.Once{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultExporter(); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
				t.Errorf("GetDefaultExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExport(t *testing.T) {
	type args struct {
		events []Event
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"",
			args{[]Event{}},
			false,
		},
		{
			"",
			args{make([]Event, 1)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := export(tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("export() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExporterExportEvents(t *testing.T) {
	type fields struct {
		name     string
		stopCh   chan struct{}
		stopOnce *sync.Once
	}
	type args struct {
		ctx    context.Context
		events []Event
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
			args{
				ctx:    contextWithDone(),
				events: nil,
			},
			true,
		},
		{
			"",
			fields{
				name:     "",
				stopCh:   channelWithStop(),
				stopOnce: &sync.Once{},
			},
			args{
				ctx:    context.Background(),
				events: nil,
			},
			false,
		},
		{
			"",
			fields{
				name:     "",
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				ctx:    context.Background(),
				events: nil,
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
			if err := e.ExportEvents(tt.args.ctx, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("ExportMetrics() error = %v, wantErr %v", err, tt.wantErr)
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
		},
		{
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
