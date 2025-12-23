package public

import (
	"context"
	"os"
	"strings"
)

// FileClient 客户端结构体。
type FileClient struct {
	filepath string
	stopCh   chan struct{}
}

// Path 获取上报地址。
func (c *FileClient) Path() string {
	return c.filepath
}

// Stop 关闭发送器。
func (c *FileClient) Stop(ctx context.Context) error {
	close(c.stopCh)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// UploadData 批量发送可观测性数据。
func (c *FileClient) UploadData(ctx context.Context, data []byte) error {
	// 退出逻辑关闭了发送。
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.stopCh:
		return nil
	default:

	}
	// 写入本地文件，每次追加。
	output, Err := os.OpenFile(c.filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if Err != nil {
		return Err
	}
	if _, err := output.Write(data); err != nil {
		return err
	}
	return nil
}

// NewFileClient 创建Exporter的本地文件发送客户端。
func NewFileClient(stdoutPath string) Client {
	if strings.TrimSpace(stdoutPath) == "" {
		return &FileClient{filepath: "./ObservableData.json", stopCh: make(chan struct{})}
	}
	return &FileClient{filepath: stdoutPath, stopCh: make(chan struct{})}
}
