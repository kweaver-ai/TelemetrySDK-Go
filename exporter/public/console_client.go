package public

import (
	"context"
	"os"
)

// ConsoleClient 客户端结构体。
type ConsoleClient struct {
	filepath string
	stopCh   chan struct{}
}

// Path 获取上报地址。
func (c *ConsoleClient) Path() string {
	return c.filepath
}

// Stop 关闭发送器。
func (c *ConsoleClient) Stop(ctx context.Context) error {
	close(c.stopCh)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// UploadData 批量发送可观测性数据。
func (c *ConsoleClient) UploadData(ctx context.Context, data []byte) error {
	// 退出逻辑关闭了发送。
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.stopCh:
		return nil
	default:

	}
	// 控制台输出。
	output := os.Stdout
	if _, err := output.Write(data); err != nil {
		return err
	}
	return nil
}

// NewConsoleClient 创建Exporter的控制台发送客户端。
func NewConsoleClient() Client {
	return &ConsoleClient{filepath: "CONSOLE", stopCh: make(chan struct{})}
}
