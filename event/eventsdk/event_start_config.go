package eventsdk

import (
	"time"

	"go.opentelemetry.io/otel/trace"
)

type eventStartConfig struct {
	EventID     string
	EventType   string
	Time        time.Time
	Level       Level
	Attributes  []Attribute
	Subject     string
	SpanContext trace.SpanContext
	Data        interface{}
}

// DefaultEventType 默认的非空事件类型
const DefaultEventType = "Default.EventType"

func defaultEventStartConfig() *eventStartConfig {
	return &eventStartConfig{
		EventID:     generateID(),
		EventType:   DefaultEventType,
		Time:        time.Now(),
		Level:       INFO,
		Attributes:  []Attribute{},
		Subject:     "",
		SpanContext: trace.SpanContext{},
		Data:        nil,
	}
}

func newEventStartConfig(opts ...EventStartOption) *eventStartConfig {
	cfg := defaultEventStartConfig()
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}
