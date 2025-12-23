package eventsdk

import (
	"reflect"
	"testing"
	"time"
)

func TestDefaultEventProviderConfig(t *testing.T) {
	tests := []struct {
		name string
		want *eventProviderConfig
	}{
		{
			"",
			&eventProviderConfig{
				Exporters:     make(map[string]EventExporter),
				FlushInternal: DefaultInternal,
				MaxEvent:      DefaultMaxEvent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultEventProviderConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultEventProviderConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEventProviderConfig(t *testing.T) {
	type args struct {
		opts []EventProviderOption
	}
	tests := []struct {
		name string
		args args
		want *eventProviderConfig
	}{
		{
			"",
			args{nil},
			&eventProviderConfig{
				Exporters:     make(map[string]EventExporter),
				FlushInternal: 10 * time.Second,
				MaxEvent:      49,
			},
		},
		{
			"",
			args{[]EventProviderOption{MaxEvent(19)}},
			MaxEvent(19).apply(&eventProviderConfig{
				Exporters:     make(map[string]EventExporter),
				FlushInternal: 10 * time.Second,
				MaxEvent:      49,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEventProviderConfig(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEventProviderConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
