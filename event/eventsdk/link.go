package eventsdk

import "go.opentelemetry.io/otel/trace"

// Link 和 Trace 关联，记录 TraceID 和 SpanID 。
type Link interface {
	// GetTraceID 返回 TraceID 。
	GetTraceID() string
	// GetSpanID 返回 SpanID 。
	GetSpanID() string
	// Valid 校验是否合法。
	Valid() bool
	// private 禁止用户自己实现接口。
	private()
}

// link 用于关联 Trace 信息， event 和 Trace 一对一。
type link struct {
	TraceID string `json:"TraceID"`
	SpanID  string `json:"SpanID"`
}

// newLink 创建新的 link 。
func newLink(spanContext trace.SpanContext) *link {
	if spanContext.IsValid() {
		return &link{
			TraceID: spanContext.TraceID().String(),
			SpanID:  spanContext.SpanID().String(),
		}
	}
	return nil
}

func (l *link) GetTraceID() string {
	return l.TraceID
}

func (l *link) GetSpanID() string {
	return l.SpanID
}

func (l *link) Valid() bool {
	return len(l.GetTraceID()) == 32 && len(l.GetSpanID()) == 16
}

func (l *link) private() {
	// private 禁止用户自己实现接口。
}
