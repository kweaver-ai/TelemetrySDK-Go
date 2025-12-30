package encoder

import (
	"bytes"
	"context"
	"log"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/exporter"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
)

// SyncEncoder 同步模式专用，只能上报到一个地址。
type SyncEncoder interface {
	Encoder
	sync()
}

type SyncJsonEncoder struct {
	syncExporter exporter.SyncExporter
}

// NewSyncEncoder 创建2.6.0后的同步发送模式JSON编码器。
func NewSyncEncoder(e exporter.SyncExporter) SyncEncoder {
	return &SyncJsonEncoder{syncExporter: e}

}

// Write 区分本地输出和上报AnyRobot，本地输出仅提取关键信息输出，上报AnyRobot为数组形式，有返回值error判断是否发送成功。
func (js *SyncJsonEncoder) Write(f field.Field) error {
	//复用JsonEncoder的发送逻辑
	ctx, cancel := context.WithCancel(context.Background())
	eps := make(map[string]exporter.LogExporter)
	eps[js.syncExporter.Name()] = js.syncExporter

	//每次发送初始化一个JsonEncoder
	encoder := &JsonEncoder{
		w:            nil,
		logExporters: eps,
		End:          _lineFeed,

		bufReal:    bytes.NewBuffer(make([]byte, 0, 4096)),
		ctx:        ctx,
		cancelFunc: cancel,
	}
	encoder.buf = encoder.bufReal
	//发送结束等待Golang垃圾回收机制释放JsonEncoder

	//以下部分为JsonEncoder.Write相似代码
	stdoutExporter, ok := encoder.logExporters["RealTimeExporter"]
	if ok {
		//将数组遍历使用RealTimeExporter单个输出，并提取关键信息输出。
		encoder.dealRealTimeExporter(stdoutExporter, f)
	}
	//不是RealTimeExporter的剩余exporter比如arExporter将输出整个数组，下面将整个数组转为byte
	err := encoder.write(f)
	if err != nil {
		log.Println(field.GenerateSpecificError(err))
	}
	// _, res := encoder.WriteBytes(encoder.End)
	_, writeBytesErr := encoder.WriteBytes(encoder.End)
	if writeBytesErr != nil {
		log.Println(field.GenerateSpecificError(writeBytesErr))
	}
	//调用不是RealTimeExporter的剩余exporter比如arExporter的输出方法将整个数组进行输出
	flushWithExportersErr := encoder.flushWithExporters()
	if flushWithExportersErr != nil {
		log.Println(field.GenerateSpecificError(flushWithExportersErr))
		return flushWithExportersErr
	}
	return nil
}

// Close 关闭JSON编码器。
func (js *SyncJsonEncoder) Close() error {
	return nil
}

func (js *SyncJsonEncoder) sync() {}
