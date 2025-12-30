package log

import (
	"math/rand"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/open_standard"
)

// SyncLogger 同步发送模式的日志器。
type SyncLogger interface {
	SetSample(sample float32)
	SetLevel(level int)
	SetWriter(writer open_standard.SyncWriter)
	Close()
	TraceField(message field.Field, typ string, opts ...field.LogOptionFunc) error
	DebugField(message field.Field, typ string, opts ...field.LogOptionFunc) error
	InfoField(message field.Field, typ string, opts ...field.LogOptionFunc) error
	WarnField(message field.Field, typ string, opts ...field.LogOptionFunc) error
	ErrorField(message field.Field, typ string, opts ...field.LogOptionFunc) error
	FatalField(message field.Field, typ string, opts ...field.LogOptionFunc) error
	Trace(message string, opts ...field.LogOptionFunc) error
	Debug(message string, opts ...field.LogOptionFunc) error
	Info(message string, opts ...field.LogOptionFunc) error
	Warn(message string, opts ...field.LogOptionFunc) error
	Error(message string, opts ...field.LogOptionFunc) error
	Fatal(message string, opts ...field.LogOptionFunc) error
}

// syncLogger 同步发送模式的日志器。
type syncLogger struct {
	logLevel int
	sample   float32
	writer   open_standard.Writer
}

// NewSyncLogger 创建同步发送模式的日志器，创建时可传入参数。记日志的方法有返回值，返回error=nil代表发送成功，返回error!=nil代表发送失败。
func NewSyncLogger(opts ...LoggerStartOption) SyncLogger {
	cfg := newLoggerStartConfig(opts...)
	return &syncLogger{
		sample:   cfg.Sample,
		logLevel: cfg.LogLevel,
	}
}

// SetLevel 设置日志级别，从0~7，0代表全部输出，7代表关闭输出。
func (logger *syncLogger) SetLevel(level int) {
	if level < AllLevel || level > OffLevel {
		return
	}
	logger.logLevel = level
}

// SetSample 设置采样等级，从0.0~1.0，0.0代表不采样，1.0代表全采样。
func (logger *syncLogger) SetSample(sample float32) {
	if sample < 0.0 || sample > 1.0 {
		return
	}
}

// SetWriter 设置日志器写入器。
func (logger *syncLogger) SetWriter(writer open_standard.SyncWriter) {
	logger.writer = writer
}

// Close 释放Writer。
func (logger *syncLogger) Close() {
	if logger.writer != nil {
		_ = logger.writer.Close()
	}
}

// TraceField Trace 级别的日志，记录结构体。
func (logger *syncLogger) TraceField(message field.Field, typ string, opts ...field.LogOptionFunc) error {
	if TraceLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(typ, message, TraceLevelString, opts...)
}

// DebugField Debug 级别的日志，记录结构体。
func (logger *syncLogger) DebugField(message field.Field, typ string, opts ...field.LogOptionFunc) error {
	if DebugLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(typ, message, DebugLevelString, opts...)
}

// InfoField Info 级别的日志，记录结构体。
func (logger *syncLogger) InfoField(message field.Field, typ string, opts ...field.LogOptionFunc) error {
	if InfoLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(typ, message, InfoLevelString, opts...)
}

// WarnField Warn 级别的日志，记录结构体。
func (logger *syncLogger) WarnField(message field.Field, typ string, opts ...field.LogOptionFunc) error {
	if WarnLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(typ, message, WarnLevelString, opts...)
}

// ErrorField Error 级别的日志，记录结构体。
func (logger *syncLogger) ErrorField(message field.Field, typ string, opts ...field.LogOptionFunc) error {
	if ErrorLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLogField(typ, message, ErrorLevelString, opts...)
}

// FatalField Fatal 级别的日志，记录结构体。
func (logger *syncLogger) FatalField(message field.Field, typ string, opts ...field.LogOptionFunc) error {
	if FatalLevel < logger.logLevel {
		return nil
	}
	return logger.writeLogField(typ, message, FatalLevelString, opts...)
}

// Trace Trace 级别的日志，记录字符串。
func (logger *syncLogger) Trace(message string, opts ...field.LogOptionFunc) error {
	if TraceLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, TraceLevelString, opts...)
}

// Debug Debug 级别的日志，记录字符串。
func (logger *syncLogger) Debug(message string, opts ...field.LogOptionFunc) error {
	if DebugLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, DebugLevelString, opts...)
}

// Info Info 级别的日志，记录字符串。
func (logger *syncLogger) Info(message string, opts ...field.LogOptionFunc) error {
	if InfoLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, InfoLevelString, opts...)
}

// Warn Warn 级别的日志，记录字符串。
func (logger *syncLogger) Warn(message string, opts ...field.LogOptionFunc) error {
	if WarnLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, WarnLevelString, opts...)
}

// Error Error 级别的日志，记录字符串。
func (logger *syncLogger) Error(message string, opts ...field.LogOptionFunc) error {
	if ErrorLevel < logger.logLevel || !logger.sampleCheck() {
		return nil
	}
	return logger.writeLog(message, ErrorLevelString, opts...)
}

// Fatal Fatal 级别的日志，记录字符串。
func (logger *syncLogger) Fatal(message string, opts ...field.LogOptionFunc) error {
	if FatalLevel < logger.logLevel {
		return nil
	}
	return logger.writeLog(message, FatalLevelString, opts...)
}

// sampleCheck 检查采样率决定是否记录当前日志。
func (logger *syncLogger) sampleCheck() bool {
	// 全采样
	if logger.sample >= 1.0 {
		return true
	}
	// 全丢弃
	if logger.sample <= 0 {
		return false
	}
	// 生成0.0~1.0之间的随机数
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Float32() <= logger.sample
}

// writeLogField 写非结构化日志。
func (logger *syncLogger) writeLogField(typ string, message, level field.Field, opts ...field.LogOptionFunc) error {
	logSpan := field.SyncLog()
	logSpan.SetLogLevel(level)
	record := newRecord(typ, message)
	logSpan.SetRecord(record)
	logSpan.SetOption(opts...)
	return logger.writer.Write([]field.LogSpan{logSpan})
}

// writeLog 写结构化日志。
func (logger *syncLogger) writeLog(message string, level field.Field, opts ...field.LogOptionFunc) error {
	logSpan := field.SyncLog()
	logSpan.SetLogLevel(level)
	record := field.MallocStructField(1)
	record.Set("Message", field.StringField(message))
	logSpan.SetRecord(record)
	logSpan.SetOption(opts...)
	return logger.writer.Write([]field.LogSpan{logSpan})
}
