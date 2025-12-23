package eventsdk

import (
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestLinkGetSpanID(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			if got := l.GetSpanID(); got != tt.want {
				t.Errorf("GetSpanID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkGetTraceID(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			if got := l.GetTraceID(); got != tt.want {
				t.Errorf("GetTraceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkPrivate(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			l.private()
		})
	}
}

func TestLinkValid(t *testing.T) {
	type fields struct {
		TraceID string
		SpanID  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				TraceID: "5bd0ecc145cc8639007721df27ecda50",
				SpanID:  "4cbf3e2c1e8517e7",
			},
			true,
		},
		{
			"",
			fields{
				TraceID: "",
				SpanID:  "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{
				TraceID: tt.fields.TraceID,
				SpanID:  tt.fields.SpanID,
			}
			if got := l.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLink(t *testing.T) {
	type args struct {
		spanContext trace.SpanContext
	}
	tests := []struct {
		name string
		args args
		want *link
	}{
		{
			"",
			args{spanContext: trace.SpanContext{}},
			nil,
		},
		{
			"",
			args{spanContext: trace.NewSpanContext(
				trace.SpanContextConfig{
					TraceID: [16]byte{0x1},
					SpanID:  [8]byte{0x1},
				})},
			&link{
				TraceID: "01000000000000000000000000000000",
				SpanID:  "0100000000000000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLink(tt.args.spanContext); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
