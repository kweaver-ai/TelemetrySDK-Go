package observable_code

import (
	"fmt"
	"time"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/ar_log"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
)

func main() {
	ar_log.InitBusinessLogger()
	fmt.Println("hello world")

	msg := make(map[string]interface{})
	msg["文件名称"] = "文件A"
	ar_log.BLogger.InfoField(field.MallocJsonField(msg), "数据产品元数据")
	// 日志打印有延迟，等待一会儿
	time.Sleep(5 * time.Second)
}
