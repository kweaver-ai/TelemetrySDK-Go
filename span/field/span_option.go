package field

import (
	"context"
)

type LogOptionFunc func(*logSpanV1)

func WithAttribute(attr *attribute) LogOptionFunc {
	return func(l *logSpanV1) {
		if attr == nil || attr.Message == nil {
			return
		}
		if l.attributes == nil {
			l.attributes = MallocMapField()
		}
		l.attributes.Append(attr.Type, attr.Message)
	}
}

func WithContext(ctx context.Context) LogOptionFunc {
	return func(l *logSpanV1) {
		if ctx != nil {
			l.ctx = ctx
		}
	}
}
