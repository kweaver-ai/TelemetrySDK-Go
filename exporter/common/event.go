package common

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Event 自定义 Event 统一Attribute。
type Event struct {
	Name                  string       `json:"Name"`
	Attributes            []*Attribute `json:"Attributes"`
	DroppedAttributeCount int          `json:"DroppedAttributeCount"`
	Time                  int64        `json:"Time"`
}

// AnyRobotEventFromEvent 单条sdktrace.Event转换为*Event。
func AnyRobotEventFromEvent(event sdktrace.Event) *Event {
	return &Event{
		Name:                  event.Name,
		Attributes:            AnyRobotAttributesFromKeyValues(event.Attributes),
		DroppedAttributeCount: event.DroppedAttributeCount,
		Time:                  event.Time.UnixNano(),
	}
}

// AnyRobotEventsFromEvents 批量sdktrace.Event转换为[]*Event。
func AnyRobotEventsFromEvents(events []sdktrace.Event) []*Event {
	if events == nil {
		return make([]*Event, 0)
	}
	arEvents := make([]*Event, 0, len(events))
	for i := 0; i < len(events); i++ {
		arEvents = append(arEvents, AnyRobotEventFromEvent(events[i]))
	}
	return arEvents
}
