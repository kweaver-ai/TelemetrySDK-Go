package field

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestNewSpanFromPool(t *testing.T) {
	s0 := NewSpanFromPool(func(LogSpan) {}, context.Background())
	s1 := NewSpanFromPool(func(LogSpan) {}, context.Background())

	assert.NotEqual(t, s0, s1)

	s0.Free()
	s01 := NewSpanFromPool(func(LogSpan) {}, context.Background())
	s01.Free()
	assert.Equal(t, s0, s01)
	assert.NotEqual(t, s1, s01)
}

func TestGenID(t *testing.T) {
	// 同一个上下文，trace中的 traceID，spanID 和 log span中的traceID spanID一致
	tp1 := tracesdk.NewTracerProvider()
	tr1 := tp1.Tracer("123")
	ctx1, span := tr1.Start(context.Background(), "fdsaf")
	defer func(tp1 *tracesdk.TracerProvider, ctx context.Context) {
		_ = tp1.Shutdown(ctx)
	}(tp1, nil)
	defer span.End()
	s0 := NewSpanFromPool(func(LogSpan) {}, ctx1)
	assert.Equal(t, span.SpanContext().TraceID().String(), s0.TraceID())
	assert.Equal(t, span.SpanContext().SpanID().String(), s0.SpanID())
	// 不同上下文则不一致
	s1 := NewSpanFromPool(func(LogSpan) {}, context.Background())
	assert.NotEqual(t, span.SpanContext().TraceID().String(), s1.TraceID())
	assert.NotEqual(t, span.SpanContext().SpanID().String(), s1.SpanID())
}

func TestSignal(t *testing.T) {
	lock := &sync.Mutex{}
	count := 0
	s0 := NewSpanFromPool(func(s LogSpan) {
		lock.Lock()
		count += 1
		lock.Unlock()
		s.Free()
	}, context.Background())

	s0.Signal()

	s1 := NewSpanFromPool(func(s LogSpan) {
		lock.Lock()
		count += 1
		lock.Unlock()
		s.Free()
	}, context.Background())

	s1.Signal()

	start := time.Now()

	s2 := NewSpanFromPool(func(LogSpan) {
		cost := time.Since(start)
		assert.True(t, cost < 1*time.Microsecond, "parent span signal() complete before children span")
	}, context.Background())

	time.Sleep(1 * time.Millisecond)
	assert.Equal(t, 2, count)

	s2.Free()

}

func TestLogSpanRecord(t *testing.T) {

	s0 := NewSpanFromPool(func(LogSpan) {}, context.Background())

	now := time.Now()
	arrayField := MallocArrayField(4)
	arrayField.Append(IntField(1))
	arrayField.Append(Float64Field(1))
	arrayField.Append(StringField("test string in array"))
	arrayField.Append(TimeField(now))

	structField := MallocStructField(4)
	structField.Set("int", IntField(2))
	structField.Set("float", Float64Field(2))
	structField.Set("string", StringField("test string in struct"))
	structField.Set("time", TimeField(now))

	records := []Field{
		IntField(0),
		Float64Field(0),
		StringField("test string"),
		TimeField(now),
		arrayField,
		structField,
	}

	for _, f := range records {
		s0.SetRecord(f)
		assert.Equal(t, f, s0.GetRecord())
	}
}

func TestLogSpanV1_SetOption(t *testing.T) {
	s0 := NewSpanFromPool(func(LogSpan) {}, context.Background())
	s0.SetOption(nil)
	s0.SetOption()
	assert.Equal(t, s0.GetContext(), context.Background())
	s0.Signal()
}

func TestLogSpanV1Attribute(t *testing.T) {
	s0 := NewSpanFromPool(func(LogSpan) {}, context.Background())
	attr := NewAttribute("test", StringField("testattr"))

	assert.Equal(t, s0.GetAttributes(), MallocStructField(0))

	s0.SetOption(WithAttribute(attr))

	record := MallocMapField()
	record.Append(attr.Type, attr.Message)

	assert.Equal(t, s0.GetAttributes(), Field(record))

	s0.Signal()

}

func TestLogSpanV1Context(t *testing.T) {
	ctx := context.Background()
	s0 := NewSpanFromPool(func(LogSpan) {}, nil)
	s0.SetOption(WithContext(ctx))
	assert.Equal(t, ctx, s0.GetContext())

	tp1 := tracesdk.NewTracerProvider()
	tr1 := tp1.Tracer("123")
	ctx1, span := tr1.Start(context.Background(), "fdsaf")
	defer func(tp1 *tracesdk.TracerProvider, ctx context.Context) {
		_ = tp1.Shutdown(ctx)
	}(tp1, nil)
	defer span.End()

	s1 := NewSpanFromPool(func(LogSpan) {}, nil)

	s1.SetOption(WithContext(ctx1))
	assert.Equal(t, ctx1, s1.GetContext())
	assert.NotEqual(t, ctx, s1.GetContext())
}

func TestLogSpanV1IsNilContext(t *testing.T) {
	s0 := &logSpanV1{}
	s0.init()
	assert.True(t, s0.IsNilContext())
	s0.SetOption(WithContext(context.Background()))
	assert.True(t, s0.IsNilContext())

	tp1 := tracesdk.NewTracerProvider()
	tr1 := tp1.Tracer("123")
	ctx1, span := tr1.Start(context.Background(), "fdsaf")
	defer func(tp1 *tracesdk.TracerProvider, ctx context.Context) {
		_ = tp1.Shutdown(ctx)
	}(tp1, nil)
	defer span.End()
	s0.SetOption(WithContext(ctx1))
	assert.False(t, s0.IsNilContext())
}

func TestLogSpanV1LogLevel(t *testing.T) {
	s0 := NewSpanFromPool(func(LogSpan) {}, nil)
	s0.SetLogLevel(StringField("Trace"))
	assert.Equal(t, StringField("Trace"), s0.GetLogLevel())
}
