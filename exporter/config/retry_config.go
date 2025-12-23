package config

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/custom_errors"
)

// RetryConfig 当可观测性数据发送失败时，根据重发机制来重新发送，保障数据不漏。
type RetryConfig struct {
	// Enabled 是否启用重发机制。
	Enabled bool
	// InitialInterval 第一次重发与上一次发送的时间间隔。
	InitialInterval time.Duration
	// MaxInterval 两次重发的最长时间间隔。
	MaxInterval time.Duration
	// MaxElapsedTime 重发最长持续的时间。
	MaxElapsedTime time.Duration
}

// RetryFunc 包含重发机制的请求方法。
type RetryFunc func(context.Context, func(context.Context) error) error

// RetryableError 可以触发重发机制的错误。
type RetryableError struct {
	Throttle int64
}

// Error 实现error接口。
func (e RetryableError) Error() string {
	return custom_errors.RetryFailure
}

// evaluate 通过 RetryableError 类型来判断是否可重发。
func evaluate(err error) (bool, time.Duration) {
	if err == nil {
		return false, 0
	}

	rErr, ok := err.(RetryableError)
	if !ok {
		return false, 0
	}

	return true, time.Duration(rErr.Throttle)
}

// RetryFunc 带有重发机制的HTTP请求，如果错误为可重发错误则重发AnyRobotSpans，否则丢弃数据。
func (r RetryConfig) RetryFunc() RetryFunc {
	// 不开启重发，则直接执行内部发送逻辑。内部的fn function(context),error是真正的http请求逻辑，在client处实现。
	if !r.Enabled {
		return func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}
	}
	// 计算重发时间间隔。
	offset := getBackoff(r)
	offset.Reset()

	// 启用了重发，则返回一个嵌套的函数是为了重发时复用HTTP连接，外部的 retryFunc(context,function)error 是包装了retry的逻辑，内部的fn function(context),error是真正的http请求逻辑，在client处实现。
	retryFunc := func(ctx context.Context, fn func(context.Context) error) error {
		for {
			err := fn(ctx)
			if err == nil {
				return nil
			}

			retryable, throttle := evaluate(err)
			if !retryable {
				return err
			}

			backOff := offset.NextBackOff()
			if backOff == backoff.Stop {
				log.Println(custom_errors.ExceedRetryElapsedTime)
				return errors.New(custom_errors.ExceedRetryElapsedTime)
			}

			var delay time.Duration
			if backOff > throttle {
				delay = backOff
			} else {
				elapsed := offset.GetElapsedTime()
				if offset.MaxElapsedTime != 0 && elapsed+throttle > offset.MaxElapsedTime {
					log.Println(custom_errors.ExceedRetryElapsedTime)
					return errors.New(custom_errors.ExceedRetryElapsedTime)
				}
				delay = throttle
			}

			if err := wait(ctx, delay); err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return retryFunc
}

// getBackoff 返回计算重发时间的结构体，每次发送间隔指数增加。
func getBackoff(r RetryConfig) *backoff.ExponentialBackOff {
	return &backoff.ExponentialBackOff{
		InitialInterval:     r.InitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         r.MaxInterval,
		MaxElapsedTime:      r.MaxElapsedTime,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
}

// wait 等待间隔时间。
func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		select {
		case <-timer.C:
		default:
			return ctx.Err()
		}
	case <-timer.C:
	}

	return nil
}

// WithRetry 设置重发。
func WithRetry(rc *RetryConfig) Option {
	return func(cfg *Config) *Config {
		cfg.RetryConfig = rc
		return cfg
	}
}
