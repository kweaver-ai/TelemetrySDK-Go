package version

// 每次拉release都要修改这里的版本号。

const (
	TelemetrySDKVersion = "2.10.0"

	TraceInstrumentationName = "TelemetrySDK-Go/exporter/ar_trace"
	TraceInstrumentationURL  = "https://github.com/kweaver-ai/TelemetrySDK-Go?path=/exporter/ar_trace"

	LogInstrumentationName = "TelemetrySDK-Go/exporter/ar_log"
	//LogInstrumentationURL  = "https://github.com/kweaver-ai/TelemetrySDK-Go?path=/exporter/ar_log"

	MetricInstrumentationName = "TelemetrySDK-Go/exporter/ar_metric"
	MetricInstrumentationURL  = "https://github.com/kweaver-ai/TelemetrySDK-Go?path=/exporter/ar_metric"
)
