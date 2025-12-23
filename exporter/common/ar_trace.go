package common

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// AnyRobotSpan 实现ReadOnlySpan接口的结构体，属性无增删。
type AnyRobotSpan struct {
	Name                 string                `json:"Name"`
	SpanContext          trace.SpanContext     `json:"SpanContext"`
	Parent               trace.SpanContext     `json:"Parent"`
	SpanKind             trace.SpanKind        `json:"SpanKind"`
	StartTime            int64                 `json:"StartTime"`
	EndTime              int64                 `json:"EndTime"`
	Attributes           []*Attribute          `json:"Attributes"`
	Links                []*Link               `json:"Links"`
	Events               []*Event              `json:"Events"`
	Status               ArStatus              `json:"Status"`
	InstrumentationScope instrumentation.Scope `json:"InstrumentationScope"`
	Resource             *Resource             `json:"Resource"`
	DroppedAttributes    int                   `json:"DroppedAttributes"`
	DroppedEvents        int                   `json:"DroppedEvents"`
	DroppedLinks         int                   `json:"DroppedLinks"`
	ChildSpanCount       int                   `json:"ChildSpanCount"`
}

// AnyRobotSpanFromReadOnlySpan 单条span转换为*AnyRobotSpan。
func AnyRobotSpanFromReadOnlySpan(span sdktrace.ReadOnlySpan) *AnyRobotSpan {
	if span == nil {
		return &AnyRobotSpan{}
	}
	return &AnyRobotSpan{Name: span.Name(),
		SpanContext:          span.SpanContext(),
		Parent:               span.Parent(),
		SpanKind:             span.SpanKind(),
		StartTime:            span.StartTime().UnixNano(),
		EndTime:              span.EndTime().UnixNano(),
		Attributes:           AnyRobotAttributesFromKeyValues(span.Attributes()),
		Links:                AnyRobotLinksFromLinks(span.Links()),
		Events:               AnyRobotEventsFromEvents(span.Events()),
		Status:               transformStatus(span.Status()),
		InstrumentationScope: span.InstrumentationScope(),
		Resource:             AnyRobotResourceFromResource(span.Resource()),
		DroppedAttributes:    span.DroppedAttributes(),
		DroppedEvents:        span.DroppedEvents(),
		DroppedLinks:         span.DroppedLinks(),
		ChildSpanCount:       span.ChildSpanCount()}
}

// AnyRobotTraceFromReadOnlyTrace 批量trace转换为[]*AnyRobotTrace。
func AnyRobotTraceFromReadOnlyTrace(trace []sdktrace.ReadOnlySpan) []*AnyRobotSpan {
	if trace == nil {
		return make([]*AnyRobotSpan, 0)
	}
	arTrace := make([]*AnyRobotSpan, 0, len(trace))
	for _, span := range trace {
		arTrace = append(arTrace, AnyRobotSpanFromReadOnlySpan(span))
	}
	return arTrace
}

// 用于转换StatusCode的输出格式
type ArStatus struct {
	// Code is an identifier of a Spans state classification.
	Code uint32
	// Description is a user hint about why that status was set. It is only
	// applicable when Code is Error.
	Description string
}

func transformStatus(s sdktrace.Status) ArStatus {
	return ArStatus{uint32(s.Code), s.Description}
}
