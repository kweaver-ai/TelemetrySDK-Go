package public

import (
	"context"
)

// Client 负责连接数据接收器，并且负责转换数据格式并发送可观测性数据，内部为net/http/public。
type Client interface {
	// Path 用来获取上报地址。
	Path() string
	// Stop 用来关闭连接，它只会被调用一次因此不用担心幂等性问题，但是可能存在并发调用，需要上层 Exporter 通过sync.Once来控制。
	Stop(ctx context.Context) error
	// UploadData 用来发送任意数据，可能会并发调用。
	UploadData(ctx context.Context, data []byte) error
}

type SyncClient interface {
	Client
	// Sync 修改发送方式为同步发送。
	Sync()
}
