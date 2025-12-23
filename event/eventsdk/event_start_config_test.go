package eventsdk

import (
	"reflect"
	"testing"
)

func TestDefaultEventStartConfig(t *testing.T) {
	tests := []struct {
		name string
		want *eventStartConfig
	}{
		{
			"",
			defaultEventStartConfig(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultEventStartConfig(); !reflect.DeepEqual(got.EventType, tt.want.EventType) {
				t.Errorf("defaultEventStartConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEventStartConfig(t *testing.T) {
	type args struct {
		opts []EventStartOption
	}
	tests := []struct {
		name string
		args args
		want *eventStartConfig
	}{
		{
			"",
			args{},
			defaultEventStartConfig(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEventStartConfig(tt.args.opts...); !reflect.DeepEqual(got.EventType, tt.want.EventType) {
				t.Errorf("newEventStartConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
