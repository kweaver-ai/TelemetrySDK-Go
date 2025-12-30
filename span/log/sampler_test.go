package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/encoder"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/open_standard"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/runtime"

	"github.com/stretchr/testify/assert"
)

func testLogString(l *SamplerLogger) {
	attr := field.NewAttribute("test", nil)

	l.Trace(string(TraceLevelString))
	l.Debug(string(DebugLevelString))
	l.Info(string(InfoLevelString))
	l.Warn(string(WarnLevelString))
	l.Error(string(ErrorLevelString))
	l.Fatal(string(FatalLevelString))

	l.Trace(string(TraceLevelString), nil)
	l.Debug(string(DebugLevelString), nil)
	l.Info(string(InfoLevelString), nil)
	l.Warn(string(WarnLevelString), nil)
	l.Error(string(ErrorLevelString), nil)
	l.Fatal(string(FatalLevelString), nil)

	l.Trace(string(TraceLevelString), field.WithAttribute(attr))
	l.Debug(string(DebugLevelString), field.WithAttribute(attr))
	l.Info(string(InfoLevelString), field.WithAttribute(attr))
	l.Warn(string(WarnLevelString), field.WithAttribute(attr))
	l.Error(string(ErrorLevelString), field.WithAttribute(attr))
	l.Fatal(string(FatalLevelString), field.WithAttribute(attr))
}

func testLogField(l *SamplerLogger) {
	ctx := context.Background()
	l.TraceField(TraceLevelString, "test", field.WithContext(ctx))
	l.DebugField(DebugLevelString, "test", field.WithContext(ctx))
	l.InfoField(InfoLevelString, "test", field.WithContext(ctx))
	l.WarnField(WarnLevelString, "test", field.WithContext(ctx))
	l.ErrorField(ErrorLevelString, "test", field.WithContext(ctx))
	l.FatalField(FatalLevelString, "test", field.WithContext(ctx))

	l.TraceField(TraceLevelString, "test")
	l.DebugField(DebugLevelString, "test")
	l.InfoField(InfoLevelString, "test")
	l.WarnField(WarnLevelString, "test")
	l.ErrorField(ErrorLevelString, "test")
	l.FatalField(FatalLevelString, "test")
}

func testLogLevel(l *SamplerLogger, level int) {
	l.SetLevel(level)
	testLogString(l)
	testLogField(l)
}

func TestSamplerLoggerSpan(t *testing.T) {
	buf := io.Discard
	l := NewDefaultSamplerLogger()
	run := runtime.NewRuntime(
		open_standard.OpenTelemetryWriter(
			encoder.NewJsonEncoder(buf),
			field.IntField(0)),
		field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	testLogLevel(l, TraceLevel)
	testLogLevel(l, DebugLevel)
	testLogLevel(l, InfoLevel)
	testLogLevel(l, WarnLevel)
	testLogLevel(l, ErrorLevel)
	testLogLevel(l, FatalLevel)

	l.Close()
}

func TestSampleCheck(t *testing.T) {
	count := 1000000
	total := 0
	l := NewDefaultSamplerLogger()
	l.Sample = 1
	for i := 0; i < count; i += 1 {
		if l.sampleCheck() {
			total += 1
		}
	}
	assert.Equal(t, count, total)

	l.Sample = 0
	total = 0
	for i := 0; i < count; i += 1 {
		if l.sampleCheck() {
			total += 1
		}
	}
	assert.Equal(t, 0, total)

	l.Sample = 0.8
	for i := 0; i < count; i += 1 {
		if l.sampleCheck() {
			total += 1
		}
	}

	// allow +- 0.5
	if total > int(float32(count)*(l.Sample+0.5)) || total < int(float32(count)*(l.Sample-0.5)) {
		t.Errorf("TestSampleCheck error, total: %d, sample: %d", total, int(float32(count)*l.Sample))
	}
}

func TestSamplerLoggerNil(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewDefaultSamplerLogger()
	run := runtime.NewRuntime(
		open_standard.OpenTelemetryWriter(
			encoder.NewJsonEncoder(buf),
			field.IntField(0)),
		field.NewSpanFromPool)
	l.SetRuntime(run)
	l.LogLevel = AllLevel
	go run.Run()

	// test nil span to record
	testLogField(l)

	l.Close()
	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() > 0, "record shouldn't drop")

	// test runtime is nil
	run = runtime.NewRuntime(&open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(buf),
	}, field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	l.runtime = nil
	buf.Reset()
	testLogField(l)
	l.Close()
	run.Signal()
	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() == 0, "record should drop")

	testLogString(l)
	assert.True(t, buf.Len() == 0, "record should drop")

}

func TestSamplerLoggerClose(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	l := NewDefaultSamplerLogger()
	run := runtime.NewRuntime(
		open_standard.OpenTelemetryWriter(
			encoder.NewJsonEncoder(buf),
			field.IntField(0)),
		field.NewSpanFromPool)
	l.SetRuntime(run)
	go run.Run()

	// runtime will wait sub LogSpan completed

	go l.Close()
	testLogField(l)

	// run.Signal()
	// assert.Equal(t, nil, enc.Close())
	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() > 0, "record drop")

	buf.Reset()

	// will not log after runtime is closed
	testLogField(l)

	time.Sleep(1 * time.Second)
	assert.True(t, buf.Len() == 0, "record error")
}

func TestSamplerLogger(t *testing.T) {
	// 0. create logger and start runtime
	buf := bytes.NewBuffer(nil)
	l := NewDefaultSamplerLogger()
	run := runtime.NewRuntime(
		open_standard.OpenTelemetryWriter(
			encoder.NewJsonEncoder(buf),
			field.IntField(0)),
		field.NewSpanFromPool)
	l.SetRuntime(run)
	l.LogLevel = AllLevel
	go run.Run()

	// 1.1 log message into roor LogSpan
	l.Debug("debug string message", nil)
	l.DebugField(field.StringField("debug field message"), "test", nil)

	attrs := field.MallocStructField(3)
	attrs.Set("work", field.StringField("test"))
	attrs.Set("testFunc", field.StringField("TestSamplerLogger"))
	attrs.Set("testSpan", field.StringField("root"))
	// set attr

	attr := field.NewAttribute("attr", attrs)
	l.Info("infomessage", field.WithAttribute(attr))

	// final close runtime and clean work space
	l.Close()
	// run.Signal()

	time.Sleep(1 * time.Second)

	capacity := map[string]interface{}{}
	betty := buf.Bytes()
	left := 0
	i := 0
	n := 0
	for ; i < len(betty); i += 1 {
		if betty[i] == '\n' {
			if err := json.Unmarshal(betty[left:i], &capacity); err != nil {
				t.Error(err)
				t.FailNow()
			} else {
				n += 1
				fmt.Println(string(betty[left:i]))
				fmt.Println()
			}
			left = i + 1
		}
	}
	if left < len(betty) {
		if err := json.Unmarshal(betty[left:i], &capacity); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}

func TestSamplerLoggerSetSample(t *testing.T) {
	type fields struct {
		Sample   float32
		LogLevel int
		runtime  *runtime.Runtime
		ctx      context.Context
	}
	type args struct {
		sample float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				Sample:   0.0,
				LogLevel: 2,
				runtime:  nil,
				ctx:      nil,
			},
			args{0.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SamplerLogger{
				Sample:   tt.fields.Sample,
				LogLevel: tt.fields.LogLevel,
				runtime:  tt.fields.runtime,
				ctx:      tt.fields.ctx,
			}
			s.SetSample(tt.args.sample)
		})
	}
}

func TestSamplerLoggerAllLogLevel(t *testing.T) {
	type fields struct {
		Sample   float32
		LogLevel int
		runtime  *runtime.Runtime
		ctx      context.Context
	}
	type args struct {
		message      string
		messageField field.Field
		typ          string
		options      []field.LogOptionFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				Sample:   0,
				LogLevel: 7,
				runtime:  nil,
				ctx:      nil,
			},
			args{
				message:      "123",
				messageField: field.StringField("456"),
				typ:          "type_test",
				options:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SamplerLogger{
				Sample:   tt.fields.Sample,
				LogLevel: tt.fields.LogLevel,
				runtime:  tt.fields.runtime,
				ctx:      tt.fields.ctx,
			}
			s.Trace(tt.args.message, tt.args.options...)
			s.Debug(tt.args.message, tt.args.options...)
			s.Info(tt.args.message, tt.args.options...)
			s.Warn(tt.args.message, tt.args.options...)
			s.Error(tt.args.message, tt.args.options...)
			s.Fatal(tt.args.message, tt.args.options...)
			s.TraceField(tt.args.messageField, tt.args.typ, tt.args.options...)
			s.DebugField(tt.args.messageField, tt.args.typ, tt.args.options...)
			s.InfoField(tt.args.messageField, tt.args.typ, tt.args.options...)
			s.WarnField(tt.args.messageField, tt.args.typ, tt.args.options...)
			s.ErrorField(tt.args.messageField, tt.args.typ, tt.args.options...)
			s.FatalField(tt.args.messageField, tt.args.typ, tt.args.options...)
		})
	}
}

func TestNewSamplerLogger(t *testing.T) {
	l := NewSamplerLogger()
	l.Close()
}
