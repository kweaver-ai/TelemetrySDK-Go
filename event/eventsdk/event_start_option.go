package eventsdk

import (
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// EventStartOption Event初始化选项。
type EventStartOption interface {
	// apply 更改Event默认配置。
	apply(*eventStartConfig) *eventStartConfig
}

// eventStartOptionFunc 执行 EventStartOption 的方法。
type eventStartOptionFunc func(*eventStartConfig) *eventStartConfig

func (o eventStartOptionFunc) apply(cfg *eventStartConfig) *eventStartConfig {
	return o(cfg)
}

//当前版本不允许修改eventID。
//func WithEventID(eventID string) EventStartOption {
//	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
//		if len(eventID) == 26 {
//			cfg.EventID = eventID
//		}
//		return cfg
//	})
//}

// WithEventType 设置事件类型。
func WithEventType(eventType string) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		if strings.TrimSpace(eventType) != "" {
			cfg.EventType = eventType
		}
		return cfg
	})
}

// WithTime 设置事件发生时间。
func WithTime(t time.Time) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		if t.After(time.Time{}) {
			cfg.Time = t
		}
		return cfg
	})
}

// withLevel 设置事件级别。
func withLevel(level Level) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Level = level
		return cfg
	})
}

// WithAttributes 设置资源信息。
func WithAttributes(attrs ...Attribute) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Attributes = attrs
		return cfg
	})
}

// WithSubject 设置事件操作对象。
func WithSubject(subject string) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Subject = subject
		return cfg
	})
}

// WithSpanContext 设置关联的Trace信息。
func WithSpanContext(spanContext trace.SpanContext) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.SpanContext = spanContext
		return cfg
	})
}

func withData(data interface{}) EventStartOption {
	return eventStartOptionFunc(func(cfg *eventStartConfig) *eventStartConfig {
		cfg.Data = data
		return cfg
	})
}
