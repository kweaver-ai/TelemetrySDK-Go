// Package log
// Provide some interfaces of logger
// looger can use to record string log of any object Log in the form of field
// logger can use to log internal thread/task with InteranlSpan.
// A LogSpan use to describe a thread job info or a task info. You can get a thread context info ever recorded.
// And then aggregation Span by TraceID or describe a Trace tree by SpanID and TraceID
package log

import (
	"context"
	"math/rand"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/runtime"
)

type Logger interface {
	// SetSample set the log sample
	// sample number from 0 to 1.0
	SetSample(sample float32)

	// SetLevel set Logger level
	// level can be log.AllLevel,log.debugLevel,...,log.OffLevel
	SetLevel(level int)

	// Close logger
	Close()

	// SetRuntime for logger
	SetRuntime(*runtime.Runtime)
	// TraceField do a trace log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	TraceField(message field.Field, typ string, options ...field.LogOptionFunc)

	// DebugField do a debug log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	DebugField(message field.Field, typ string, options ...field.LogOptionFunc)

	// InfoField do a Info log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	InfoField(message field.Field, typ string, options ...field.LogOptionFunc)

	// WarnField do a Warn log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	WarnField(message field.Field, typ string, options ...field.LogOptionFunc)

	// ErrorField do a Error log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	ErrorField(message field.Field, typ string, options ...field.LogOptionFunc)

	// FatalField do a Fatal log a object into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	FatalField(message field.Field, typ string, options ...field.LogOptionFunc)

	// Trace do a trace string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Trace(message string, options ...field.LogOptionFunc)

	// Debug do a Debug string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// but not signal the LogSpan
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Debug(message string, options ...field.LogOptionFunc)

	// Info do a Info string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Info(message string, options ...field.LogOptionFunc)

	// Warn do a Warn string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Warn(message string, options ...field.LogOptionFunc)

	// Error do a Error string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Error(message string, options ...field.LogOptionFunc)

	// Fatal do a Fatal string log into LogSpan,
	// if LogSpan is not nil, this interface will log the info,
	// if LogSpan is nil, this interface will create a LogSpan
	// to log the info and signal the LogSpan.
	Fatal(message string, options ...field.LogOptionFunc)
}

type SamplerLogger struct {
	// logger sample
	Sample float32
	// logger info
	LogLevel int
	runtime  *runtime.Runtime
	ctx      context.Context
}

func NewSamplerLogger(opts ...LoggerStartOption) Logger {
	cfg := newLoggerStartConfig(opts...)
	return &SamplerLogger{
		Sample:   cfg.Sample,
		LogLevel: cfg.LogLevel,
	}
}

func NewDefaultSamplerLogger() *SamplerLogger {
	return &SamplerLogger{
		Sample:   1.0,
		LogLevel: InfoLevel,
	}
}

func newRecord(typ string, message field.Field) field.Field {
	record := field.MallocStructField(2)
	record.Set(typ, message)
	record.Set("Type", field.StringField(typ))
	return record
}

// Close logger
func (s *SamplerLogger) Close() {
	if s.runtime != nil {
		s.runtime.Signal()
	}
}

// SetRuntime for logger
// this function will signal older runtime, but will not start new runtime
func (s *SamplerLogger) SetRuntime(r *runtime.Runtime) {
	if s.runtime != nil {
		s.runtime.Signal()
	}
	s.runtime = r
}

func (s *SamplerLogger) getLogSpan() field.LogSpan {
	if s.runtime != nil {
		return s.runtime.Children(s.ctx)
	}
	return nil
}

func (s *SamplerLogger) sampleCheck() bool {
	if s.Sample >= 1.0 {
		return true
	}

	if s.Sample <= 0 {
		return false
	}

	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return newRand.Float32() <= s.Sample
}

// SetSample set the log sample
// sample number from 0 to 1.0
func (s *SamplerLogger) SetSample(sample float32) {
	s.Sample = sample
}

// SetLevel set Logger level
// level can be log.AllLevel,log.debugLevel,...,log.OffLevel
func (s *SamplerLogger) SetLevel(level int) {
	s.LogLevel = level
}

func (s *SamplerLogger) writeLogField(typ string, message, level field.Field, options ...field.LogOptionFunc) {
	l := s.getLogSpan()
	if l == nil {
		return
	}
	defer l.Signal()

	l.SetLogLevel(level)
	record := newRecord(typ, message)
	l.SetRecord(record)
	l.SetOption(options...)

}

func (s *SamplerLogger) writeLog(message string, level field.Field, options ...field.LogOptionFunc) {
	l := s.getLogSpan()
	if l == nil {
		return
	}
	defer l.Signal()

	l.SetLogLevel(level)
	record := field.MallocStructField(1)
	record.Set("Message", field.StringField(message))
	l.SetRecord(record)
	l.SetOption(options...)
}

// TraceField do a trace log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) TraceField(message field.Field, typ string, options ...field.LogOptionFunc) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLogField(typ, message, TraceLevelString, options...)

}

// DebugField do a debug log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) DebugField(message field.Field, typ string, options ...field.LogOptionFunc) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLogField(typ, message, DebugLevelString, options...)
}

// InfoField do a Info log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) InfoField(message field.Field, typ string, options ...field.LogOptionFunc) {
	if InfoLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLogField(typ, message, InfoLevelString, options...)
}

// WarnField do a Warn log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) WarnField(message field.Field, typ string, options ...field.LogOptionFunc) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLogField(typ, message, WarnLevelString, options...)
}

// ErrorField do a Error log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) ErrorField(message field.Field, typ string, options ...field.LogOptionFunc) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLogField(typ, message, ErrorLevelString, options...)
}

// FatalField do a Fatal log a object into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) FatalField(message field.Field, typ string, options ...field.LogOptionFunc) {
	if FatalLevel < s.LogLevel {
		return
	}

	s.writeLogField(typ, message, FatalLevelString, options...)
}

// Trace do a trace string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Trace(message string, options ...field.LogOptionFunc) {
	if TraceLevel < s.LogLevel || !s.sampleCheck() {
		return
	}

	s.writeLog(message, TraceLevelString, options...)
}

// Debug do a Debug string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Debug(message string, options ...field.LogOptionFunc) {
	if DebugLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, DebugLevelString, options...)
}

// Info do a Info string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Info(message string, options ...field.LogOptionFunc) {
	if InfoLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, InfoLevelString, options...)
}

// Warn do a Warn string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Warn(message string, options ...field.LogOptionFunc) {
	if WarnLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, WarnLevelString, options...)
}

// Error do a Error string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// but not signal the LogSpan
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Error(message string, options ...field.LogOptionFunc) {
	if ErrorLevel < s.LogLevel || !s.sampleCheck() {
		return
	}
	s.writeLog(message, ErrorLevelString, options...)
}

// Fatal do a Fatal string log into LogSpan,
// if LogSpan is not nil, this interface will log the info,
// if LogSpan is nil, this interface will create a LogSpan
// to log the info and signal the LogSpan.
func (s *SamplerLogger) Fatal(message string, options ...field.LogOptionFunc) {
	if FatalLevel < s.LogLevel {
		return
	}
	s.writeLog(message, FatalLevelString, options...)
}
