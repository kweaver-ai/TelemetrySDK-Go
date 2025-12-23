package portal_service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/ar_log"
	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/ar_trace"
	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/resource"
)

func main() {
	// 通过configmap设置配置，适用于k8s集群内部应用
	// 第二个参数为configmap的名称，爱数产品各个微服务用{产品名称}-telemetry-sdk，产品名称可选：anyshare、anydata、anyfabric、anyrobot、anybackup
	// 第三个参数为各个微服务的名称
	ar_trace.InitTracer("cm", "anyshare-telemetry-sdk", "my-service-1")
	// 通过yaml文件设置配置，适用于k8s集群外部应用
	// 第二个参数为go程序执行目录下yaml文件名称，yaml文件格式可参考api_service目录下的ob-app-config-trace.yaml
	// 第三个参数为各个微服务的名称
	//ar_trace.InitTracer("yaml", "ob-app-config-trace", "my-service-1")

	// 设置微服务版本
	resource.SetServiceVersion("1.0.0")

	// 服务停止时先把内存中的链路数据立马发送出去
	defer ar_trace.ShutdownTracer()

	// 通过configmap设置配置，适用于k8s集群内部应用
	// 第二个参数为configmap的名称，爱数产品各个微服务用{产品名称}-telemetry-sdk，产品名称可选：anyshare、anydata、anyfabric、anyrobot、anybackup
	// 第三个参数为各个微服务的名称
	ar_log.InitLogger("cm", "anyshare-telemetry-sdk", "my-service-1")
	// 通过yaml文件设置配置，适用于k8s集群外部应用
	// 第二个参数为go程序执行目录下yaml文件名称，yaml文件格式可参考api_service目录下的ob-app-config-log.yaml
	// 第三个参数为各个微服务的名称
	//ar_log.InitLogger("yaml", "ob-app-config-log", "my-service-1")

	r := gin.Default()
	r.Use(otelgin.Middleware("my-server-portal"))
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := getUser(id, c)
		maskedName := desensitizeUserName(name, c)
		c.String(http.StatusOK, maskedName)
	})
	_ = r.Run(":50080")
}

// desensitizeUserName 用户名称脱敏
func desensitizeUserName(name string, ctx context.Context) string {
	var err error
	ctx, _ = ar_trace.StartInternalSpanSimple(ctx, "用户名称脱敏")
	// 结束span时，err如果不为空的话，则span状态设置为error
	defer func() { ar_trace.EndSpan(ctx, err) }()

	// 将程序错误赋值给err，便于结束span时根据err是否为空设置span状态
	if len(name) == 0 {
		err = errors.New("用户名称为空字符串")
	}

	runes := []rune(name)

	for i := 1; i < len(runes); i++ {
		runes[i] = '*'
	}

	// 输出支持与trace关联的log
	ar_log.Info(ctx, "用户名称脱敏成功")

	return string(runes)
}

// getUser 调用其他服务，根据用户ID获取用户名称
func getUser(id string, ctx context.Context) string {
	var err error
	ctx, _ = ar_trace.StartInternalSpanSimple(ctx, "根据用户ID获取用户名称")
	// 结束span时，err如果不为空的话，则span状态设置为error
	defer func() { ar_trace.EndSpan(ctx, err) }()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	url := fmt.Sprintf("http://127.0.0.1:50081/users/%s", id)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return ""
	}

	return string(body)
}
