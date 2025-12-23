package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/custom_errors"
)

type key string

func deliver(ctx context.Context) error {
	switch ctx.Value(key("StatusCode")) {
	case 200:
		return nil
	case 413:
		return errors.New("413")
	case 500:
		return RetryableError{Throttle: 10}
	case 503:
		return RetryableError{Throttle: 1}
	default:
		return nil
	}
}

func result(ctx context.Context) error {
	switch ctx.Value(key("StatusCode")) {
	case 200:
		return nil
	case 413:
		return errors.New("413")
	case 500:
		return errors.New(custom_errors.ExceedRetryElapsedTime)
	case 503:
		return RetryableError{}
	default:
		return nil
	}
}

func TestRetryConfigRetryFunc(t *testing.T) {
	sth := gomonkey.ApplyFunc(wait, func(ctx context.Context, delay time.Duration) error {
		return RetryableError{}
	})
	defer sth.Reset()

	other := gomonkey.ApplyFunc(log.Println, func(v ...interface{}) {
		fmt.Println(v)
	})
	defer other.Reset()

	type fields struct {
		Enabled         bool
		InitialInterval time.Duration
		MaxInterval     time.Duration
		MaxElapsedTime  time.Duration
		ctx             context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   RetryFunc
	}{
		{
			"关闭重发",
			fields{
				Enabled:         false,
				InitialInterval: 0,
				MaxInterval:     0,
				MaxElapsedTime:  0,
				ctx:             context.Background(),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		},
		{
			"开启重发，200不重发",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     0,
				MaxElapsedTime:  0,
				ctx:             context.WithValue(context.Background(), key("StatusCode"), 200),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		},
		{
			"开启重发，413不重发",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     0,
				MaxElapsedTime:  0,
				ctx:             context.WithValue(context.Background(), key("StatusCode"), 413),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		},
		{
			"开启重发，500重发",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     1 * time.Second,
				MaxElapsedTime:  2 * time.Nanosecond,
				ctx:             context.WithValue(context.Background(), key("StatusCode"), 500),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		},
		{
			"开启重发，503重发",
			fields{
				Enabled:         true,
				InitialInterval: 0,
				MaxInterval:     1 * time.Second,
				MaxElapsedTime:  2 * time.Second,
				ctx:             context.WithValue(context.Background(), key("StatusCode"), 503),
			},
			func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RetryConfig{
				Enabled:         tt.fields.Enabled,
				InitialInterval: tt.fields.InitialInterval,
				MaxInterval:     tt.fields.MaxInterval,
				MaxElapsedTime:  tt.fields.MaxElapsedTime,
			}
			// if got := r.RetryFunc(); !reflect.DeepEqual(got(tt.fields.ctx, deliver), tt.want(tt.fields.ctx, result)) {
			// 	t.Errorf("RetryFunc() = %v, want %v", got(tt.fields.ctx, deliver), tt.want(tt.fields.ctx, result))
			// }
			got := r.RetryFunc()
			if got(tt.fields.ctx, deliver) != nil && tt.want(tt.fields.ctx, result) != nil {
				assert.Equal(t, tt.want(tt.fields.ctx, result).Error(), got(tt.fields.ctx, deliver).Error())
			} else {
				assert.Equal(t, tt.want(tt.fields.ctx, result), got(tt.fields.ctx, deliver))
			}
		})
	}
}

func TestRetryableErrorError(t *testing.T) {
	type fields struct {
		Throttle int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"可重发错误的报错",
			fields{0},
			"TelemetrySDK-Exporter(Go).Error: 数据正在重发",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := RetryableError{
				Throttle: tt.fields.Throttle,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetry(t *testing.T) {
	type args struct {
		rc *RetryConfig
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			"配置重发逻辑",
			args{DefaultRetryConfig()},
			WithRetry(DefaultRetryConfig()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithRetry(tt.args.rc); !reflect.DeepEqual(got(DefaultConfig()), tt.want(DefaultConfig())) {
				t.Errorf("WithRetry() = %v, want %v", got(DefaultConfig()), tt.want(DefaultConfig()))
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 time.Duration
	}{
		{
			"不可重发",
			args{nil},
			false,
			0,
		},
		{
			"普通error不可重发",
			args{errors.New("something")},
			false,
			0,
		},
		{
			"RetryableError可以重发",
			args{RetryableError{15}},
			true,
			15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := evaluate(tt.args.err)
			if got != tt.want {
				t.Errorf("evaluate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("evaluate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetBackoff(t *testing.T) {
	type args struct {
		r RetryConfig
	}
	tests := []struct {
		name string
		args args
		want *backoff.ExponentialBackOff
	}{
		{
			"计算重发时间",
			args{},
			getBackoff(RetryConfig{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBackoff(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBackoff() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func contextWithDone() context.Context {
//	ctx, cancel := context.WithCancel(context.Background())
//	cancel()
//	return ctx
//}
//
//func TestWait(t *testing.T) {
//	type args struct {
//		ctx   context.Context
//		delay time.Duration
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{
//			"等待重发",
//			args{
//				ctx:   context.Background(),
//				delay: 1 * time.Nanosecond,
//			},
//			false,
//		},
//		{
//			"已关闭Client放弃重发",
//			args{
//				ctx:   contextWithDone(),
//				delay: 1 * time.Second,
//			},
//			true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := wait(tt.args.ctx, tt.args.delay); (err != nil) != tt.wantErr {
//				t.Errorf("wait() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
