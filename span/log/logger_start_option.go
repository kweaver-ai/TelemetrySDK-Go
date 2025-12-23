package log

// LoggerStartOption Logger 初始化配置选项。
type LoggerStartOption interface {
	apply(*loggerStartConfig) *loggerStartConfig
}

// loggerStartOptionFunc 执行 LoggerStartOption 的方法。
type loggerStartOptionFunc func(*loggerStartConfig) *loggerStartConfig

func (fn loggerStartOptionFunc) apply(cfg *loggerStartConfig) *loggerStartConfig {
	return fn(cfg)
}

// emptyOption 当传入的配置项出错时，返回emptyOption。
func emptyOption() LoggerStartOption {
	return loggerStartOptionFunc(func(cfg *loggerStartConfig) *loggerStartConfig {
		return cfg
	})
}

// WithLevel 设置日志级别，从0~7，0代表全部输出，7代表关闭输出。
func WithLevel(logLevel int) LoggerStartOption {
	if logLevel < AllLevel || logLevel > OffLevel {
		return emptyOption()
	}
	return loggerStartOptionFunc(func(cfg *loggerStartConfig) *loggerStartConfig {
		cfg.LogLevel = logLevel
		return cfg
	})
}

// WithSample 设置采样等级，从0.0~1.0，0.0代表不采样，1.0代表全采样。
func WithSample(sample float32) LoggerStartOption {
	if sample < 0.0 || sample > 1.0 {
		return emptyOption()
	}
	return loggerStartOptionFunc(func(cfg *loggerStartConfig) *loggerStartConfig {
		cfg.Sample = sample
		return cfg
	})
}
