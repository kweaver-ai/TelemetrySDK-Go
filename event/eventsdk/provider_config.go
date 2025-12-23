package eventsdk

import "time"

// eventProviderConfig EventProvider 初始化配置。
type eventProviderConfig struct {
	Exporters     map[string]EventExporter
	FlushInternal time.Duration
	MaxEvent      int
}

const DefaultInternal = 10 * time.Second
const DefaultMaxEvent = 49

func defaultEventProviderConfig() *eventProviderConfig {
	return &eventProviderConfig{
		Exporters:     make(map[string]EventExporter),
		FlushInternal: DefaultInternal,
		MaxEvent:      DefaultMaxEvent,
	}
}

func newEventProviderConfig(opts ...EventProviderOption) *eventProviderConfig {
	cfg := defaultEventProviderConfig()
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}
