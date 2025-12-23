package public

import (
	"context"
)

// SilentClient 不发送数据的客户端结构体
type SilentClient struct {
}

// NewSilentClient 创建Exporter的不发送数据的客户端
func NewSilentClient() Client {
	return &SilentClient{}
}

// Path 获取上报地址，没啥用，为了实现Client接口。
func (c *SilentClient) Path() string {
	return ""
}

// Stop 关闭发送器，没啥用，为了实现Client接口。
func (c *SilentClient) Stop(ctx context.Context) error {
	return nil
}

// UploadData 不发送数据。
func (c *SilentClient) UploadData(ctx context.Context, data []byte) error {
	return nil
}
