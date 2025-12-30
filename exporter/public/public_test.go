package public

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/config"
)

func TestSetServiceInfo(t *testing.T) {
	type args struct {
		name     string
		version  string
		instance string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"设置service信息",
			args{
				name:     "name",
				version:  "2.5.0",
				instance: "abcd1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetServiceInfo(tt.args.name, tt.args.version, tt.args.instance)
		})
	}
}

func TestWithAnyRobotURL(t *testing.T) {
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		URL string
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"设置正确的上报地址",
			args{"https://www.baidu.com/"},
			WithAnyRobotURL("https://www.baidu.com/"),
		},
		{
			"设置错误的上报地址",
			args{"1:2:3:4:5"},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAnyRobotURL(tt.args.URL); !reflect.DeepEqual(got(config.DefaultConfig()), tt.want(config.DefaultConfig())) {
				t.Errorf("WithAnyRobotURL() = %v, want %v", got(config.DefaultConfig()), tt.want(config.DefaultConfig()))
			}
		})
	}
}

func TestWithCompression(t *testing.T) {
	fmt.Println()
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		compression int
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"设置压缩方式为JSON",
			args{
				0,
			},
			WithCompression(0),
		},
		{
			"设置压缩方式为GZIP",
			args{
				1,
			},
			WithCompression(1),
		},
		{
			"设置压缩方式为非法",
			args{
				2,
			},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCompression(tt.args.compression); !reflect.DeepEqual(got(config.DefaultConfig()), tt.want(config.DefaultConfig())) {
				t.Errorf("WithCompression() = %v, want %v", got(config.DefaultConfig()), tt.want(config.DefaultConfig()))
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
		want config.Option
	}{
		{
			"设置空请求头信息",
			args{
				nil,
			},
			WithHeader(nil),
		},
		{
			"设置自定义请求头信息",
			args{
				map[string]string{"key": "value"},
			},
			WithHeader(map[string]string{"key": "value"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithHeader(tt.args.headers); !reflect.DeepEqual(got(config.DefaultConfig()), tt.want(config.DefaultConfig())) {
				t.Errorf("WithHeader() = %v, want %v", got(config.DefaultConfig()), tt.want(config.DefaultConfig()))
			}
		})
	}
}

func TestWithRetry(t *testing.T) {
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		enabled        bool
		internal       time.Duration
		maxInterval    time.Duration
		maxElapsedTime time.Duration
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"关闭重发机制",
			args{
				enabled:        false,
				internal:       1,
				maxInterval:    2,
				maxElapsedTime: 3,
			},
			WithRetry(false, 1, 2, 3),
		},
		{
			"开启重发机制",
			args{
				enabled:        true,
				internal:       4,
				maxInterval:    5,
				maxElapsedTime: 6,
			},
			WithRetry(true, 4, 5, 6),
		},
		{
			"重发机制超时",
			args{
				enabled:        true,
				internal:       1 * time.Hour,
				maxInterval:    2 * time.Hour,
				maxElapsedTime: 3 * time.Hour,
			},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetry(tt.args.enabled, tt.args.internal, tt.args.maxInterval, tt.args.maxElapsedTime); !reflect.DeepEqual(got(config.DefaultConfig()), tt.want(config.DefaultConfig())) {
				t.Errorf("WithRetry() = %v, want %v", got(config.DefaultConfig()), tt.want(config.DefaultConfig()))
			}
		})
	}
}

func TestWithTimeout(t *testing.T) {
	sth := gomonkey.ApplyFunc(log.Fatalln, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer sth.Reset()

	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want config.Option
	}{
		{
			"设置超时时间",
			args{2 * time.Second},
			WithTimeout(2 * time.Second),
		},
		{
			"超时时间非法",
			args{61 * time.Second},
			config.EmptyOption(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTimeout(tt.args.duration); !reflect.DeepEqual(got(config.DefaultConfig()), tt.want(config.DefaultConfig())) {
				t.Errorf("WithTimeout() = %v, want %v", got(config.DefaultConfig()), tt.want(config.DefaultConfig()))
			}
		})
	}
}
