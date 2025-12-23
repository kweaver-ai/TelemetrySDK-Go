package config

import (
	"time"
)

// Config Exporter的配置项。
type Config struct {
	HTTPConfig  *HTTPConfig
	RetryConfig *RetryConfig
}

// DefaultConfig 默认的配置项。
func DefaultConfig() *Config {
	return &Config{
		HTTPConfig:  DefaultHTTPConfig(),
		RetryConfig: DefaultRetryConfig(),
	}
}

// DefaultHTTPConfig 默认的HTTP配置。
func DefaultHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Insecure:    true,
		Endpoint:    "localhost:5678",
		Path:        "/api/feed_ingester/v1/jobs/{jobid}/data",
		Compression: 1,
		Timeout:     10 * time.Second,
		Headers:     nil,
		TLSCfg:      nil,
	}
}

// DefaultRetryConfig 默认的重发机制。
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		Enabled:         true,
		InitialInterval: 5 * time.Second,
		MaxInterval:     30 * time.Second,
		MaxElapsedTime:  time.Minute,
	}
}

// NewConfig 返回执行 Option 后的配置项。
func NewConfig(opts ...Option) *Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}

// Option client的配置项结构体。
type Option func(*Config) *Config

func (o Option) apply(cfg *Config) *Config {
	return o(cfg)
}

// EmptyOption 空的配置项，不改变配置，用于配置错误发生时候。
func EmptyOption() Option {
	return func(cfg *Config) *Config {
		return cfg
	}
}
