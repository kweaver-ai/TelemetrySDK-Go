package common

import (
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var link = sdktrace.Link{
	SpanContext:           trace.SpanContext{},
	Attributes:            []attribute.KeyValue{attribute.String("123", "123")},
	DroppedAttributeCount: 123,
}

func TestAnyRobotLinkFromLink(t *testing.T) {
	type args struct {
		link sdktrace.Link
	}
	tests := []struct {
		name string
		args args
		want Link
	}{
		{
			"转换空关联信息",
			args{sdktrace.Link{}},
			*AnyRobotLinkFromLink(sdktrace.Link{}),
		},
		{
			"转换非空关联信息",
			args{link},
			*AnyRobotLinkFromLink(link),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := *AnyRobotLinkFromLink(tt.args.link); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotLinkFromLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyRobotLinksFromLinks(t *testing.T) {
	type args struct {
		links []sdktrace.Link
	}
	tests := []struct {
		name string
		args args
		want []*Link
	}{
		{
			"转换空Links",
			args{},
			[]*Link{},
		},
		{
			"转换非空Links",
			args{[]sdktrace.Link{link}},
			AnyRobotLinksFromLinks([]sdktrace.Link{link}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotLinksFromLinks(tt.args.links); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotLinksFromLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
