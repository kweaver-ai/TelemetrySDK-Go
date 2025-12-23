package eventsdk

import (
	"reflect"
	"testing"
	"time"
)

var config = &eventProviderConfig{
	Exporters:     make(map[string]EventExporter),
	FlushInternal: 5 * time.Second,
	MaxEvent:      9,
}

func TestWithExporters(t *testing.T) {
	type args struct {
		exporters []EventExporter
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{[]EventExporter{GetDefaultExporter()}},
			Exporters(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exporters(tt.args.exporters...); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("Exporters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithFlushInternal(t *testing.T) {
	type args struct {
		flushInternal time.Duration
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{time.Minute},
			FlushInternal(time.Minute),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FlushInternal(tt.args.flushInternal); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("FlushInternal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMaxEvent(t *testing.T) {
	type args struct {
		maxEvent int
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{99},
			MaxEvent(99),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxEvent(tt.args.maxEvent); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("MaxEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServiceInfo(t *testing.T) {
	type args struct {
		ServiceName     string
		ServiceVersion  string
		ServiceInstance string
	}
	tests := []struct {
		name string
		args args
		want EventProviderOption
	}{
		{
			"",
			args{
				ServiceName:     "",
				ServiceVersion:  "",
				ServiceInstance: "",
			},
			ServiceInfo("", "", ""),
		},
		{
			"",
			args{
				ServiceName:     "123",
				ServiceVersion:  "456",
				ServiceInstance: "789",
			},
			ServiceInfo("123", "456", "789"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ServiceInfo(tt.args.ServiceName, tt.args.ServiceVersion, tt.args.ServiceInstance); !reflect.DeepEqual(got.apply(config), tt.want.apply(config)) {
				t.Errorf("ServiceInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventProviderOptionFuncApply(t *testing.T) {
	type args struct {
		cfg *eventProviderConfig
	}
	tests := []struct {
		name string
		o    eventProviderOptionFunc
		args args
		want *eventProviderConfig
	}{
		{
			"",
			eventProviderOptionFunc(func(providerConfig *eventProviderConfig) *eventProviderConfig {
				return providerConfig
			}),
			args{config},
			config,
		},
		{
			"",
			eventProviderOptionFunc(func(providerConfig *eventProviderConfig) *eventProviderConfig {
				providerConfig.MaxEvent = 12
				return providerConfig
			}),
			args{MaxEvent(12).apply(config)},
			config,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.apply(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
