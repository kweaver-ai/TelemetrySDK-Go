package eventsdk

import (
	"reflect"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func TestWithAttributes(t *testing.T) {
	type args struct {
		attrs []Attribute
	}
	tests := []struct {
		name string
		args args
		want EventStartOption
	}{
		{
			"",
			args{},
			WithAttributes(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAttributes(tt.args.attrs...); !reflect.DeepEqual(got.apply(defaultEventStartConfig()).Attributes, tt.want.apply(defaultEventStartConfig()).Attributes) {
				t.Errorf("WithAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithData(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want EventStartOption
	}{
		{
			"",
			args{nil},
			withData(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withData(tt.args.data); !reflect.DeepEqual(got.apply(defaultEventStartConfig()).Data, tt.want.apply(defaultEventStartConfig()).Data) {
				t.Errorf("withData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEventType(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		args args
		want EventStartOption
	}{
		{
			"",
			args{"type"},
			WithEventType("type"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEventType(tt.args.eventType); !reflect.DeepEqual(got.apply(defaultEventStartConfig()).EventType, tt.want.apply(defaultEventStartConfig()).EventType) {
				t.Errorf("WithEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithLevel(t *testing.T) {
	type args struct {
		level Level
	}
	tests := []struct {
		name string
		args args
		want EventStartOption
	}{
		{
			"",
			args{WARN},
			withLevel(WARN),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withLevel(tt.args.level); !reflect.DeepEqual(got.apply(defaultEventStartConfig()).Level, tt.want.apply(defaultEventStartConfig()).Level) {
				t.Errorf("withLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSpanContext(t *testing.T) {
	type args struct {
		spanContext trace.SpanContext
	}
	tests := []struct {
		name string
		args args
		want EventStartOption
	}{
		{
			"",
			args{trace.SpanContext{}},
			WithSpanContext(trace.SpanContext{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSpanContext(tt.args.spanContext); !reflect.DeepEqual(got.apply(defaultEventStartConfig()).SpanContext, tt.want.apply(defaultEventStartConfig()).SpanContext) {
				t.Errorf("WithSpanContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithSubject(t *testing.T) {
	type args struct {
		subject string
	}
	tests := []struct {
		name string
		args args
		want EventStartOption
	}{
		{
			"",
			args{"object"},
			WithSubject("object"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSubject(tt.args.subject); !reflect.DeepEqual(got.apply(defaultEventStartConfig()).Subject, tt.want.apply(defaultEventStartConfig()).Subject) {
				t.Errorf("WithSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want EventStartOption
	}{
		{
			"",
			args{time.Now()},
			WithTime(time.Now()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTime(tt.args.t); !reflect.DeepEqual(got.apply(defaultEventStartConfig()).Level, tt.want.apply(defaultEventStartConfig()).Level) {
				t.Errorf("WithTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventStartOptionFuncApply(t *testing.T) {
	type args struct {
		cfg *eventStartConfig
	}
	tests := []struct {
		name string
		o    eventStartOptionFunc
		args args
		want *eventStartConfig
	}{
		{
			"",
			func(cfg *eventStartConfig) *eventStartConfig {
				return cfg
			},
			args{defaultEventStartConfig()},
			defaultEventStartConfig(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.apply(tt.args.cfg); !reflect.DeepEqual(got.Level, tt.want.Level) {
				t.Errorf("apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
