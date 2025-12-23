package exporter

import "context"

// LogExporter 异步模式专用。
type LogExporter interface {
	Name() string
	ExportLogs(ctx context.Context, p []byte) error
	Shutdown(ctx context.Context) error
}

// SyncExporter 同步模式专用。
type SyncExporter interface {
	LogExporter
	Sync()
}
