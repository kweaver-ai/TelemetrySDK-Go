package config

import (
	"reflect"
	"testing"
)

func TestEmptyOption(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		{
			"不改变配置",
			func(cfg *Config) *Config {
				return cfg
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EmptyOption(); !reflect.DeepEqual(got(DefaultConfig()), tt.want(DefaultConfig())) {
				t.Errorf("EmptyOption() = %v, want %v", got(DefaultConfig()), tt.want(DefaultConfig()))
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			"没有配置项",
			args{opts: nil},
			Config{
				HTTPConfig:  DefaultHTTPConfig(),
				RetryConfig: DefaultRetryConfig(),
			},
		}, {
			"有配置项",
			args{opts: []Option{EmptyOption()}},
			Config{
				HTTPConfig:  DefaultHTTPConfig(),
				RetryConfig: DefaultRetryConfig(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(tt.args.opts...); !reflect.DeepEqual(got.HTTPConfig.Path, tt.want.HTTPConfig.Path) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptionApplyOption(t *testing.T) {
	type fields struct {
		Fn func(*Config) *Config
	}
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Config
	}{
		{
			"执行配置更改",
			fields{EmptyOption()},
			args{DefaultConfig()},
			DefaultConfig(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Option(tt.fields.Fn)
			if got := h(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
