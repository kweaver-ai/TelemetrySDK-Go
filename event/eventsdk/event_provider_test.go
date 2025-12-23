package eventsdk

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

var _, cancel = context.WithCancel(context.Background())
var cfg = newEventProviderConfig()

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestGetEventProvider(t *testing.T) {
	tests := []struct {
		name string
		want EventProvider
	}{
		{
			"",
			globalEventProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEventProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEventProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEventProvider(t *testing.T) {
	type args struct {
		opts []EventProviderOption
	}
	tests := []struct {
		name string
		args args
		want EventProvider
	}{
		{
			"",
			args{nil},
			&eventProvider{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventProvider(tt.args.opts...); !reflect.DeepEqual(got.Shutdown(), tt.want.Shutdown()) {
				t.Errorf("NewEventProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetEventProvider(t *testing.T) {
	type args struct {
		ep EventProvider
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{globalEventProvider},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetEventProvider(tt.args.ep)
		})
	}
}

func TestEventProviderForceFlush(t *testing.T) {
	exporterMap := make(map[string]EventExporter)
	exporterMap["DefaultExporter"] = GetDefaultExporter()
	type fields struct {
		Ctx       context.Context
		Cancel    context.CancelFunc
		In        chan Event
		Events    []Event
		Size      int
		StopOnce  *sync.Once
		Ticker    *time.Ticker
		MaxEvent  int
		Exporters map[string]EventExporter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"",
			fields{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
			false,
		},
		{
			"",
			fields{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: map[string]EventExporter{"DefaultExporter": GetDefaultExporter()},
			},
			false,
		},
		{
			"",
			fields{
				Ctx:       contextWithDone(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: map[string]EventExporter{"DefaultExporter": GetDefaultExporter()},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				Ctx:       tt.fields.Ctx,
				Cancel:    tt.fields.Cancel,
				In:        tt.fields.In,
				Events:    tt.fields.Events,
				Size:      tt.fields.Size,
				StopOnce:  &sync.Once{},
				Ticker:    tt.fields.Ticker,
				MaxEvent:  tt.fields.MaxEvent,
				Exporters: tt.fields.Exporters,
			}
			if err := ep.ForceFlush(); (err != nil) != tt.wantErr {
				t.Errorf("ForceFlush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventProviderShutdown(t *testing.T) {
	exporterMap := make(map[string]EventExporter)
	exporterMap["DefaultExporter"] = GetDefaultExporter()
	type fields struct {
		Ctx       context.Context
		Cancel    context.CancelFunc
		In        chan Event
		Events    []Event
		Size      int
		StopOnce  *sync.Once
		Ticker    *time.Ticker
		MaxEvent  int
		Exporters map[string]EventExporter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"",
			fields{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
			false,
		},
		{
			"",
			fields{
				Ctx:       contextWithDone(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				Ctx:       tt.fields.Ctx,
				Cancel:    tt.fields.Cancel,
				In:        tt.fields.In,
				Events:    tt.fields.Events,
				Size:      tt.fields.Size,
				StopOnce:  &sync.Once{},
				Ticker:    tt.fields.Ticker,
				MaxEvent:  tt.fields.MaxEvent,
				Exporters: tt.fields.Exporters,
			}
			if err := ep.Shutdown(); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventProvideLoadEvent(t *testing.T) {
	type fields struct {
		Ctx       context.Context
		Cancel    context.CancelFunc
		In        chan Event
		Events    []Event
		Size      int
		StopOnce  *sync.Once
		Ticker    *time.Ticker
		MaxEvent  int
		Exporters map[string]EventExporter
	}
	type args struct {
		event Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{

		{
			"",
			fields{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
			args{nil},
		},
		{
			"",
			fields{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        make(chan Event),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
			args{NewEvent()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				Ctx:       tt.fields.Ctx,
				Cancel:    tt.fields.Cancel,
				In:        tt.fields.In,
				Events:    tt.fields.Events,
				Size:      tt.fields.Size,
				StopOnce:  &sync.Once{},
				Ticker:    tt.fields.Ticker,
				MaxEvent:  tt.fields.MaxEvent,
				Exporters: tt.fields.Exporters,
			}
			ep.loadEvent(tt.args.event)
		})
	}
}

func TestEventProviderPrivate(t *testing.T) {
	type fields struct {
		Ctx       context.Context
		Cancel    context.CancelFunc
		In        chan Event
		Events    []Event
		Size      int
		StopOnce  *sync.Once
		Ticker    *time.Ticker
		MaxEvent  int
		Exporters map[string]EventExporter
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        make(chan Event, 10),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				Ctx:       tt.fields.Ctx,
				Cancel:    tt.fields.Cancel,
				In:        tt.fields.In,
				Events:    tt.fields.Events,
				Size:      tt.fields.Size,
				StopOnce:  &sync.Once{},
				Ticker:    tt.fields.Ticker,
				MaxEvent:  tt.fields.MaxEvent,
				Exporters: tt.fields.Exporters,
			}
			ep.private()
		})
	}
}

func channelWithStop() chan struct{} {
	stopCh := make(chan struct{})
	close(stopCh)
	return stopCh
}

func closedChannel() chan Event {
	channel := make(chan Event, 10)
	close(channel)
	return channel
}

func TestEventProviderSendEvents(t *testing.T) {
	type fields struct {
		Ctx       context.Context
		Cancel    context.CancelFunc
		In        chan Event
		Events    []Event
		Size      int
		StopOnce  *sync.Once
		Ticker    *time.Ticker
		MaxEvent  int
		Exporters map[string]EventExporter
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				Ctx:       context.Background(),
				Cancel:    cancel,
				In:        closedChannel(),
				Events:    make([]Event, 0, cfg.MaxEvent+1),
				Size:      0,
				StopOnce:  &sync.Once{},
				Ticker:    time.NewTicker(cfg.FlushInternal),
				MaxEvent:  cfg.MaxEvent,
				Exporters: cfg.Exporters,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				Ctx:       tt.fields.Ctx,
				Cancel:    tt.fields.Cancel,
				In:        tt.fields.In,
				Events:    tt.fields.Events,
				Size:      tt.fields.Size,
				StopOnce:  &sync.Once{},
				Ticker:    tt.fields.Ticker,
				MaxEvent:  tt.fields.MaxEvent,
				Exporters: tt.fields.Exporters,
			}
			ep.sendEvents()
		})
	}
}
