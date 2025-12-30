package eventsdk

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/event/v2/custom_errors"
)

// EventProvider 批量发送数据到 AnyRobot Feed Ingester 的Event数据接收器。
type EventProvider interface {
	// Shutdown 关闭 Event 生产和发送。
	Shutdown() error
	// ForceFlush 立即发送 []Event 。
	ForceFlush() error

	// private 禁止用户自己实现接口。
	private()
}

// eventProvider 全局唯一，控制 Event 批量发送。
type eventProvider struct {
	Ctx       context.Context
	Cancel    context.CancelFunc
	In        chan Event
	Events    []Event
	Size      int
	StopOnce  *sync.Once
	Ticker    *time.Ticker
	MaxEvent  int
	Exporters map[string]EventExporter
}

// NewEventProvider 根据配置项，新建 EventProvider 。
func NewEventProvider(opts ...EventProviderOption) EventProvider {
	// 获取默认配置，更新传入配置。
	cfg := newEventProviderConfig(opts...)
	// 根据配置创建 EventExporter 。
	ctx, cancel := context.WithCancel(context.Background())
	return &eventProvider{
		Ctx:       ctx,
		Cancel:    cancel,
		In:        make(chan Event, 100),
		Events:    make([]Event, 0, cfg.MaxEvent+1),
		Size:      0,
		StopOnce:  &sync.Once{},
		Ticker:    time.NewTicker(cfg.FlushInternal),
		MaxEvent:  cfg.MaxEvent,
		Exporters: cfg.Exporters,
	}
}

func (ep *eventProvider) Shutdown() error {
	// 只返回其中一个错误。
	var returnErr error = errors.New(custom_errors.AlreadyShutdown)
	ep.StopOnce.Do(func() {
		// 读取最后一个数据之后发送。
		returnErr = ep.ForceFlush()
		// 关闭信号。
		ep.Cancel()
		close(ep.In)
	})
	return returnErr
}

func (ep *eventProvider) ForceFlush() error {
	// 不存在发送目标，直接丢弃数据。
	if len(ep.Exporters) == 0 {
		ep.Events = make([]Event, 0, ep.MaxEvent+1)
		return nil
	}
	// 只返回其中一个错误。
	var returnErr error = nil
	// 往所有发送地址发送相同的数据。
	for _, e := range ep.Exporters {
		if err := e.ExportEvents(ep.Ctx, ep.Events); err != nil {
			// 如果多余一个错误则记日志。
			if returnErr == nil {
				returnErr = err
			} else {
				log.Println(err)
			}
		}
	}
	// 发送完之后清空队列。
	ep.Size = 0
	ep.Events = make([]Event, 0, ep.MaxEvent+1)
	return returnErr
}

// GetEventProvider 获取全局唯一 EventProvider 。
func GetEventProvider() EventProvider {
	return globalEventProvider
}

// SetEventProvider 设置全局唯一 EventProvider 。
func SetEventProvider(ep EventProvider) {
	// 先关闭上一个 EventProvider 。
	_ = globalEventProvider.Shutdown()
	globalEventProvider = ep
	// 设置了 EventProvider 就开启发送。
	go func() {
		ep.(*eventProvider).sendEvents()
	}()
}

// loadEvent 缓存 Event 等待定时发送或强制发送。
func (ep *eventProvider) loadEvent(event Event) {
	// 如果 EventProvider 正常运行才能添加 Event 。
	if ep.Ctx.Err() == nil {
		select {
		// 如果有空位，往里添加。
		case ep.In <- event:
		default:
			// 如果阻塞了丢弃数据。
		}

	}
}

// sendEvents 无限等待定时发送或超过上限发送。
func (ep *eventProvider) sendEvents() {
	for {
		select {
		// 每次进来一个Event就进入这里。
		case e, ok := <-ep.In:
			// 关闭之后退出循环。
			if !ok {
				// 发完最后的数据关闭Exporter。
				for _, eventExporter := range ep.Exporters {
					if err := eventExporter.Shutdown(ep.Ctx); err != nil {
						log.Println(err)
					}
				}
				return
			}
			ep.Events = append(ep.Events, e)
			ep.Size++
			// 超过上限发送。
			if ep.Size >= ep.MaxEvent {
				if err := ep.ForceFlush(); err != nil {
					log.Println(err)
				}
			}
		// 定时发送。
		case <-ep.Ticker.C:
			if err := ep.ForceFlush(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (ep *eventProvider) private() {
	// private 禁止用户自己实现接口。
}
