package config

import (
	"reflect"
	"testing"
	"time"
)

func TestWithAnyRobotURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"通过HTTP上报数据",
			args{"http://www.rockyrori.cn/1234"},
			Option(func(cfg *Config) *Config {
				cfg.HTTPConfig.Insecure = true
				cfg.HTTPConfig.Endpoint = "www.rockyrori.cn"
				cfg.HTTPConfig.Path = "/1234"
				return cfg
			}),
		},
		{
			"通过HTTPS上报数据",
			args{"https://www.rockyrori.cn:80"},
			Option(func(cfg *Config) *Config {
				cfg.HTTPConfig.Insecure = false
				cfg.HTTPConfig.Endpoint = "www.rockyrori.cn:80"
				cfg.HTTPConfig.Path = ""
				return cfg
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAnyRobotURL(tt.args.URL); !reflect.DeepEqual(got(DefaultConfig()), tt.want(DefaultConfig())) {
				t.Errorf("WithAnyRobotURL() = %v, want %v", got(DefaultConfig()), tt.want(DefaultConfig()))
			}
		})
	}
}

func TestWithCompression(t *testing.T) {
	type args struct {
		compression Compression
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"不压缩数据",
			args{0},
			Option(func(cfg *Config) *Config {
				cfg.HTTPConfig.Compression = 0
				return cfg
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompression(tt.args.compression); !reflect.DeepEqual(got(DefaultConfig()), tt.want(DefaultConfig())) {
				t.Errorf("WithCompression() = %v, want %v", got(DefaultConfig()), tt.want(DefaultConfig()))
			}
		})
	}
}

func TestWithHeader(t *testing.T) {
	type args struct {
		headers map[string]string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"自定义header",
			args{nil},
			Option(func(cfg *Config) *Config {
				cfg.HTTPConfig.Headers = nil
				return cfg
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHeader(tt.args.headers); !reflect.DeepEqual(got(DefaultConfig()), tt.want(DefaultConfig())) {
				t.Errorf("WithHeader() = %v, want %v", got(DefaultConfig()), tt.want(DefaultConfig()))
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"自定义连接超时时间",
			args{123},
			Option(func(cfg *Config) *Config {
				cfg.HTTPConfig.Timeout = 123
				return cfg
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimeout(tt.args.duration); !reflect.DeepEqual(got(DefaultConfig()), tt.want(DefaultConfig())) {
				t.Errorf("WithTimeout() = %v, want %v", got(DefaultConfig()), tt.want(DefaultConfig()))
			}
		})
	}
}
