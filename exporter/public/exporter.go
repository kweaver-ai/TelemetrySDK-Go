package public

import (
	"context"
	"sync"
)

// Exporter 导出可观测性数据到 AnyRobot Feed Ingester 的数据接收器。
type Exporter struct {
	name     string
	client   Client
	stopCh   chan struct{}
	stopOnce sync.Once
}

// Name Exporter身份证，同名视为同一个发送器，本质为上报地址。
func (e *Exporter) Name() string {
	return e.name
}

// Shutdown 关闭Exporter，关闭HTTP连接，丢弃导出缓存。
func (e *Exporter) Shutdown(ctx context.Context) error {
	var err error = nil
	e.stopOnce.Do(func() {
		close(e.stopCh)
		err = e.client.Stop(ctx)
	})
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}

// ExportData 批量发送可观测性数据到 AnyRobot Feed Ingester 的数据接收器。
func (e *Exporter) ExportData(ctx context.Context, data []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-e.stopCh:
		return nil
	default:
	}
	if len(data) == 0 {
		return nil
	}
	return e.client.UploadData(ctx, data)
}

// SetClient 设置客户端
func (e *Exporter) SetClient(client Client) {
	e.client = client
}

// GetClient 获取客户端
func (e *Exporter) GetClient() Client {
	return e.client
}

// NewExporter 创建已启动的Exporter。
func NewExporter(client Client) *Exporter {
	return &Exporter{
		name:     client.Path(),
		client:   client,
		stopCh:   make(chan struct{}),
		stopOnce: sync.Once{},
	}
}

type SyncExporter struct {
	*Exporter
}

func (s *SyncExporter) Sync() {
	// 仅实现接口用，无功能。
}

// NewSyncExporter 创建已启动的Exporter。
func NewSyncExporter(client Client) *SyncExporter {
	return &SyncExporter{
		NewExporter(client),
	}
}
