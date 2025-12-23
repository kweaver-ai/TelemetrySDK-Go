package config

import (
	"crypto/tls"
	"net/url"
	"time"
)

// Compression 定义了Trace压缩方式。
type Compression int

const (
	// NoCompression 0：不压缩。
	NoCompression Compression = iota
	// GzipCompression 1：Gzip压缩。
	GzipCompression
)

type HTTPConfig struct {
	Insecure    bool
	Endpoint    string
	Path        string
	Compression Compression
	Timeout     time.Duration
	Headers     map[string]string
	TLSCfg      *tls.Config
	IsSync      bool
}

// 以下各项为发送器配置项目。

// WithAnyRobotURL 设置上报地址。
func WithAnyRobotURL(URL string) Option {
	AnyRobotURL, _ := url.Parse(URL)
	if AnyRobotURL.Scheme == "http" {
		return func(cfg *Config) *Config {
			cfg.HTTPConfig.Insecure = true
			cfg.HTTPConfig.Endpoint = AnyRobotURL.Host
			cfg.HTTPConfig.Path = AnyRobotURL.Path
			return cfg
		}
	}
	return func(cfg *Config) *Config {
		cfg.HTTPConfig.Insecure = false
		cfg.HTTPConfig.Endpoint = AnyRobotURL.Host
		cfg.HTTPConfig.Path = AnyRobotURL.Path
		return cfg
	}
}

// WithCompression 设置压缩方式。
func WithCompression(compression Compression) Option {
	return func(cfg *Config) *Config {
		cfg.HTTPConfig.Compression = compression
		return cfg
	}
}

// WithTimeout 设置连接超时时间。
func WithTimeout(duration time.Duration) Option {
	return func(cfg *Config) *Config {
		cfg.HTTPConfig.Timeout = duration
		return cfg
	}
}

// WithHeader 设置请求头。
func WithHeader(headers map[string]string) Option {
	return func(cfg *Config) *Config {
		cfg.HTTPConfig.Headers = headers
		return cfg
	}
}
