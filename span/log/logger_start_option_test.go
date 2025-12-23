package log

import (
	"reflect"
	"testing"
)

func TestWithLevel(t *testing.T) {
	type args struct {
		logLevel int
		cfg      *loggerStartConfig
	}
	tests := []struct {
		name string
		args args
		want LoggerStartOption
	}{
		{"TestWithLevel",
			args{AllLevel,
				defaultLoggerStartConfig(),
			},
			loggerStartOptionFunc(func(cfg *loggerStartConfig) *loggerStartConfig {
				cfg.LogLevel = AllLevel
				return cfg
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithLevel(tt.args.logLevel); !reflect.DeepEqual(got.apply(tt.args.cfg).LogLevel, tt.want.apply(tt.args.cfg).LogLevel) {
				t.Errorf("WithLevel(%v), want %v", got, tt.want)
			}
		})
	}
}

func TestWithSample(t *testing.T) {
	type args struct {
		sample float32
		cfg    *loggerStartConfig
	}
	tests := []struct {
		name string
		args args
		want LoggerStartOption
	}{
		{"TestWithLevel",
			args{1.0,
				defaultLoggerStartConfig(),
			},
			loggerStartOptionFunc(func(cfg *loggerStartConfig) *loggerStartConfig {
				cfg.Sample = 1.0
				return cfg
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSample(tt.args.sample); !reflect.DeepEqual(got.apply(tt.args.cfg).Sample, tt.want.apply(tt.args.cfg).Sample) {
				t.Errorf("WithSample(%v), want %v", got, tt.want)
			}
		})
	}
}
