package field

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestWithAttribute(t *testing.T) {
	attr := NewAttribute("test", StringField("test_demo"))
	opt := WithAttribute(attr)
	info := &logSpanV1{}
	opt(info)
	assert.NotEqual(t, attr, info.attributes)
	record := MallocMapField()
	record.Append(attr.Type, attr.Message)
	assert.Equal(t, record, info.attributes)

	opt1 := WithAttribute(nil)
	info1 := &logSpanV1{}
	opt1(info1)
	assert.Equal(t, MallocStructField(0), info1.GetAttributes())
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	opt := WithContext(ctx)
	info := &logSpanV1{}
	opt(info)
	assert.Equal(t, ctx, info.ctx)

	tp1 := tracesdk.NewTracerProvider()
	tr1 := tp1.Tracer("123")
	ctx1, span := tr1.Start(context.Background(), "fdsaf")
	defer func(tp1 *tracesdk.TracerProvider, ctx context.Context) {
		_ = tp1.Shutdown(ctx)
	}(tp1, nil)
	defer span.End()
	assert.NotEqual(t, info.ctx, ctx1)
}
