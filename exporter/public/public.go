package public

import (
	"log"
	"net/url"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/config"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/custom_errors"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/resource"
)

// SetServiceInfo 设置服务信息，包括服务名、版本号、实例ID。
func SetServiceInfo(name string, version string, instance string) {
	resource.SetServiceName(name)
	resource.SetServiceVersion(version)
	resource.SetServiceInstance(instance)
}

// WithAnyRobotURL 设置 httpClient 数据上报地址。
func WithAnyRobotURL(URL string) config.Option {
	if _, err := url.Parse(URL); err != nil {
		log.Fatalln(custom_errors.InvalidURL)
		return config.EmptyOption()
	}
	return config.WithAnyRobotURL(URL)
}

// WithCompression 设置可观测性数据压缩方式：0代表无压缩，1代表GZIP压缩。
func WithCompression(compression int) config.Option {
	if compression >= 2 || compression < 0 {
		log.Fatalln(custom_errors.InvalidCompression)
		return config.EmptyOption()
	}
	return config.WithCompression(config.Compression(compression))
}

// WithTimeout 设置 httpClient 连接超时时间。
func WithTimeout(duration time.Duration) config.Option {
	if duration > 60*time.Second || duration < 0 {
		log.Fatalln(custom_errors.DurationTooLong)
		return config.EmptyOption()
	}
	return config.WithTimeout(duration)
}

// WithHeader 设置 httpClient 用户自定义请求头。
func WithHeader(headers map[string]string) config.Option {
	return config.WithHeader(headers)
}

// WithRetry 设置 httpClient 重发机制，如果显著干扰到业务运行了，请增加重发间隔maxInterval，减少最大重发时间maxElapsedTime，甚至关闭重发enabled=false。
func WithRetry(enabled bool, internal time.Duration, maxInterval time.Duration, maxElapsedTime time.Duration) config.Option {
	if enabled && (internal > 10*time.Minute || maxInterval > 20*time.Minute || maxElapsedTime > 60*time.Minute) {
		log.Fatalln(custom_errors.RetryTooLong)
		return config.EmptyOption()
	}
	retry := &config.RetryConfig{
		Enabled:         enabled,
		InitialInterval: internal,
		MaxInterval:     maxInterval,
		MaxElapsedTime:  maxElapsedTime,
	}
	return config.WithRetry(retry)
}
