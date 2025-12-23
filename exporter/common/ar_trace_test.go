package common

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func MockReadOnlySpan() sdktrace.ReadOnlySpan {
	_, span := sdktrace.NewTracerProvider().Tracer("").Start(context.Background(), "")
	if ros, ok := span.(sdktrace.ReadWriteSpan); ok {
		return ros
	}
	return nil
}

var ros = MockReadOnlySpan()

func MockReadOnlySpans() []sdktrace.ReadOnlySpan {
	return []sdktrace.ReadOnlySpan{ros}
}

func TestAnyRobotSpanFromReadOnlySpan(t *testing.T) {
	type args struct {
		ros sdktrace.ReadOnlySpan
	}
	tests := []struct {
		name string
		args args
		want *AnyRobotSpan
	}{
		{
			"转换空ReadOnlySpan",
			args{nil},
			AnyRobotSpanFromReadOnlySpan(nil),
		},
		{
			"转换非空ReadOnlySpan",
			args{ros},
			AnyRobotSpanFromReadOnlySpan(ros),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.ros)
			if got := AnyRobotSpanFromReadOnlySpan(tt.args.ros); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotSpanFromReadOnlySpan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyRobotSpansFromReadOnlySpans(t *testing.T) {
	type args struct {
		ross []sdktrace.ReadOnlySpan
	}
	tests := []struct {
		name string
		args args
		want []*AnyRobotSpan
	}{
		{
			"转换空ReadOnlySpan",
			args{nil},
			AnyRobotTraceFromReadOnlyTrace(nil),
		},
		{
			"转换非空ReadOnlySpans",
			args{MockReadOnlySpans()},
			AnyRobotTraceFromReadOnlyTrace(MockReadOnlySpans()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotTraceFromReadOnlyTrace(tt.args.ross); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotTraceFromReadOnlyTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}
