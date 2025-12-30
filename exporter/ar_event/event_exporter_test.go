package ar_event

import (
	"context"
	"reflect"
	"testing"

	"github.com/kweaver-ai/TelemetrySDK-Go/event/v2/eventsdk"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/public"
)

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestEventExporterExportEvents(t *testing.T) {
	type fields struct {
		Exporter *public.Exporter
	}
	type args struct {
		ctx    context.Context
		events []eventsdk.Event
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
				events: nil,
			},
			false,
		},
		{
			"发送Event",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx:    context.Background(),
				events: make([]eventsdk.Event, 1),
			},
			false,
		},
		{
			"已关闭的Exporter不能发Event",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx:    contextWithDone(),
				events: make([]eventsdk.Event, 1),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventExporter{
				Exporter: tt.fields.Exporter,
			}
			if err := e.ExportEvents(tt.args.ctx, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("ExportEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventResource(t *testing.T) {
	tests := []struct {
		name string
		want eventsdk.EventProviderOption
	}{
		{
			"Event的默认Resource",
			EventResource(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EventResource(); !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("EventResource() = %v, want %v", got, tt.want)
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
		want *EventExporter
	}{
		{
			"StdoutClient的EventExporter",
			args{c: public.NewStdoutClient("./AnyRobotEvent.txt")},
			NewExporter(public.NewStdoutClient("./AnyRobotEvent.txt")),
		},
		{
			"HTTPClient的EventExporter",
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
