package common

import (
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var event = sdktrace.Event{
	Name:                  "testEvent",
	Attributes:            []attribute.KeyValue{attribute.String("123", "123")},
	DroppedAttributeCount: 12,
	Time:                  time.Now(),
}

func TestAnyRobotEventFromEvent(t *testing.T) {
	type args struct {
		event sdktrace.Event
	}
	tests := []struct {
		name string
		args args
		want *Event
	}{
		{
			"转换空Event",
			args{},
			AnyRobotEventFromEvent(sdktrace.Event{}),
		},
		{
			"转换非空Event",
			args{event},
			AnyRobotEventFromEvent(event),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotEventFromEvent(tt.args.event); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotEventFromEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyRobotEventsFromEvents(t *testing.T) {
	type args struct {
		events []sdktrace.Event
	}
	tests := []struct {
		name string
		args args
		want []*Event
	}{
		{
			"转换空Events",
			args{},
			[]*Event{},
		},
		{
			"转换非空Events",
			args{[]sdktrace.Event{event}},
			AnyRobotEventsFromEvents([]sdktrace.Event{event}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotEventsFromEvents(tt.args.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotEventsFromEvents() = %v, want %v", got, tt.want)
			}
		})
	}
}
