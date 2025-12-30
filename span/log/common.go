package log

import (
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
)

// log level text format
const (
	TraceLevelString = field.StringField("Trace")
	DebugLevelString = field.StringField("Debug")
	InfoLevelString  = field.StringField("Info")
	WarnLevelString  = field.StringField("Warn")
	ErrorLevelString = field.StringField("Error")
	FatalLevelString = field.StringField("Fatal")
)

// log level int format
const (
	AllLevel = iota
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	OffLevel
)
