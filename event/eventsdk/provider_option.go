package eventsdk

import (
	"strings"
	"time"
)

// EventProviderOption EventProvider初始化选项。
type EventProviderOption interface {
	// apply 更改EventProvider默认配置。
	apply(*eventProviderConfig) *eventProviderConfig
}

// eventProviderOptionFunc 执行 EventProviderOption 的方法。
type eventProviderOptionFunc func(*eventProviderConfig) *eventProviderConfig

func (o eventProviderOptionFunc) apply(cfg *eventProviderConfig) *eventProviderConfig {
	return o(cfg)
}

// Exporters 批量设置 EventExporter 。
func Exporters(exporters ...EventExporter) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		for _, e := range exporters {
			cfg.Exporters[e.Name()] = e
		}
		return cfg
	})
}

// ServiceInfo 记录服务信息。
func ServiceInfo(ServiceName string, ServiceVersion string, ServiceInstance string) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		if strings.TrimSpace(ServiceName) != "" {
			globalServiceName = ServiceName
		}
		if strings.TrimSpace(globalServiceVersion) != "" {
			globalServiceVersion = ServiceVersion
		}
		if strings.TrimSpace(globalServiceInstance) != "" {
			globalServiceInstance = ServiceInstance
		}
		return cfg
	})
}

// FlushInternal 设置发送间隔。
func FlushInternal(flushInternal time.Duration) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		cfg.FlushInternal = flushInternal
		return cfg
	})
}

// MaxEvent 设置Event发送上限。
func MaxEvent(maxEvent int) EventProviderOption {
	return eventProviderOptionFunc(func(cfg *eventProviderConfig) *eventProviderConfig {
		cfg.MaxEvent = maxEvent
		return cfg
	})
}
