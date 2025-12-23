package examplelog

import (
	"context"
	"time"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/ar_log"
)

func main() {
	ar_log.InitLogger("yaml", "ob-app-config-log", "my-service-2")

	ar_log.Info(context.Background(), "this is log")

	time.Sleep(time.Second * 5)
}
