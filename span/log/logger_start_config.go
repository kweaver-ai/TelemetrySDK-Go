package log

type loggerStartConfig struct {
	Sample   float32
	LogLevel int
}

// defaultLoggerStartConfig 默认的 Logger 初始化配置
func defaultLoggerStartConfig() *loggerStartConfig {
	return &loggerStartConfig{
		Sample:   1.0,
		LogLevel: InfoLevel,
	}
}

// newLoggerStartConfig 根据配置项新建 Logger 配置。
func newLoggerStartConfig(opts ...LoggerStartOption) *loggerStartConfig {
	cfg := defaultLoggerStartConfig()
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}
