package ar_trace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/common"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/config"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/public"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/resource"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/version"
)

// 跨包实现接口占位用。
var _ sdktrace.SpanExporter = (*TraceExporter)(nil)

// Tracer 是一个全局变量，用于在业务代码中生产Span。
var Tracer = otel.GetTracerProvider().Tracer(
	version.TraceInstrumentationName,
	trace.WithInstrumentationVersion(version.TelemetrySDKVersion),
	trace.WithSchemaURL(version.TraceInstrumentationURL),
)

var tp = (*sdktrace.TracerProvider)(nil)
var te = &TraceExporter{}

// TraceExporter 导出数据到AnyRobot Feed Ingester的 Event 数据接收器。
type TraceExporter struct {
	*public.Exporter
}

var (
	// traceClientModifyLock 修改 traceClient 对象的锁
	traceClientModifyLock = sync.Mutex{}
)

// ExportSpans 批量发送AnyRobotSpans到AnyRobot Feed Ingester的Trace数据接收器。
func (e *TraceExporter) ExportSpans(ctx context.Context, traces []sdktrace.ReadOnlySpan) error {
	if len(traces) == 0 {
		return nil
	}
	arTrace := common.AnyRobotTraceFromReadOnlyTrace(traces)
	file := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(arTrace); err != nil {
		return err
	}
	return e.ExportData(ctx, file.Bytes())
}

// NewExporter 创建已启动的 TraceExporter 。
func NewExporter(c public.Client) *TraceExporter {
	return &TraceExporter{
		public.NewExporter(c),
	}
}

// TraceResource 传入 Trace 的默认Resource。
func TraceResource() *sdkresource.Resource {
	return resource.TraceResource()
}

// InitTracer 初始化函数，初始化全局链路数据记录器
// cfgType 配置类型，可选：cm、yaml
// cfgName 配置文件名称，如配置类型为cm，则该参数为configmap的名称；如配置类型为yaml，则该参数为yaml文件的路径，比如./ob-app.yaml
// serverName 微服务名称
func InitTracer(cfgType string, cfgName string, serverName string) {
	// 先初始化一个不记录链路数据的全局链路数据记录器
	InitSilentTracer(serverName)
	if cfgType == "cm" { // 如果配置为configmap形式
		// 监听configmap的内容，更新全局链路数据记录器的配置
		config.CmName = cfgName
		if kubeClient := config.InitKubeClient(); kubeClient != nil {
			watchConfigMap(kubeClient)
		}
	} else if cfgType == "yaml" { // 如果配置为yaml文件形式
		config.CfgFileNameTrace = cfgName
		// 初始化配置
		config.NewTraceConfig()
		UpdateTracerClient(config.YamlTraceCfg.Enabled, config.YamlTraceCfg.Endpoint)
		config.TraceVP.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("Trace config file changed:%s, update tracer client\n", e)
			config.LoadTraceConfig()
			UpdateTracerClient(config.YamlTraceCfg.Enabled, config.YamlTraceCfg.Endpoint)
		})
	}

}

// InitSilentTracer 初始化全局链路数据记录器，不记录数据
func InitSilentTracer(serverName string) {
	serverInstance := os.Getenv("HOSTNAME")

	if serverName != "" {
		resource.SetServiceName(serverName)
	}
	resource.SetServiceInstance(serverInstance)

	traceClient := public.NewSilentClient()
	te = NewExporter(traceClient)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(te,
			sdktrace.WithMaxExportBatchSize(1000)),
		sdktrace.WithResource(TraceResource()))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

// UpdateTracerClient 更新全局链路数据记录器的数据发送客户端
func UpdateTracerClient(traceEnabled string, traceEndpoint string) {
	traceClientModifyLock.Lock()
	defer traceClientModifyLock.Unlock()

	if traceEnabled == "true" && traceEndpoint != "" {
		traceClient := public.NewHTTPClient(public.WithAnyRobotURL(traceEndpoint),
			public.WithCompression(1), public.WithTimeout(10*time.Second),
			public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
		StopTracerClient()
		te.SetClient(traceClient)
	} else if traceEnabled == "true" && traceEndpoint == "" {
		traceClient := public.NewConsoleClient()
		StopTracerClient()
		te.SetClient(traceClient)
	} else {
		traceClient := public.NewSilentClient()
		StopTracerClient()
		te.SetClient(traceClient)
	}
}

// StopTracerClient 关闭全局链路数据记录器的数据发送客户端
func StopTracerClient() {
	// 关闭旧的client
	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := te.GetClient().Stop(ctx)
	if err != nil {
		log.Printf("[TelemetrySDK]Error shutting down tracer client: %v", err)
	}
}

// GetTracer 获取全局链路数据记录器
func GetTracer() *sdktrace.TracerProvider {
	return tp
}

// ShutdownTracer 关闭全局链路数据记录器
func ShutdownTracer() {
	if tp == nil {
		return
	}

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("[TelemetrySDK]Error shutting down tracer provider: %v", err)
	}

	tp = nil
}

// InitARTracer 初始化上报到AnyRobot的链路数据记录器
// ServerName 微服务名称
// ServerVersion 微服务版本
// ServerInstance 微服务实例标识
func InitARTracer() *sdktrace.TracerProvider {
	serverName := os.Getenv("TELEMETRY_SERVICE_NAME")
	serverVersion := os.Getenv("TELEMETRY_SERVICE_VERSION")
	serverInstance := os.Getenv("HOSTNAME")
	traceEnabled := os.Getenv("TELEMETRY_TRACE_ENABLED")
	traceUrl := os.Getenv("TELEMETRY_TRACE_ENDPOINT")

	if traceEnabled == "true" {
		resource.SetServiceName(serverName)
		resource.SetServiceVersion(serverVersion)
		resource.SetServiceInstance(serverInstance)

		traceClient := public.NewHTTPClient(public.WithAnyRobotURL(traceUrl),
			public.WithCompression(1), public.WithTimeout(10*time.Second),
			public.WithRetry(true, 5*time.Second, 30*time.Second, 1*time.Minute))
		traceExporter := NewExporter(traceClient)
		tracerProvider := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(traceExporter,
				sdktrace.WithMaxExportBatchSize(1000)),
			sdktrace.WithResource(TraceResource()))

		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

		return tracerProvider
	} else {
		return nil
	}
}

// StopARTracer 关闭 ARTracer
func StopARTracer(tp *sdktrace.TracerProvider) {
	if tp == nil {
		return
	}

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("[TelemetrySDK]Error shutting down tracer provider: %v", err)
	}
}

// StartInternalSpanSimple 简单一点的，性能损耗小一点的内部方法调用trace埋点
func StartInternalSpanSimple(ctx context.Context, spanName string) (context.Context, trace.Span) {
	if c, ok := ctx.(*gin.Context); ok { // 如果是gin.Context，要转换一下
		ctx = c.Request.Context()
	}

	ctx, span := otel.GetTracerProvider().Tracer(version.TraceInstrumentationName,
		trace.WithInstrumentationVersion(version.TelemetrySDKVersion),
		trace.WithSchemaURL(version.TraceInstrumentationURL)).Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindInternal))
	return ctx, span
}

// StartInternalSpan 内部方法调用trace埋点
func StartInternalSpan(ctx context.Context) (context.Context, trace.Span) {
	if c, ok := ctx.(*gin.Context); ok { // 如果是gin.Context，要转换一下
		ctx = c.Request.Context()
	}

	pc, file, linkNo, ok := runtime.Caller(1)
	if !ok {
		log.Printf("[TelemetrySDK]start span error")
		ctx, span := otel.GetTracerProvider().Tracer(version.TraceInstrumentationName,
			trace.WithInstrumentationVersion(version.TelemetrySDKVersion),
			trace.WithSchemaURL(version.TraceInstrumentationURL)).Start(ctx, "unKnow", trace.WithSpanKind(trace.SpanKindInternal))
		return ctx, span
	} else {
		funcPaths := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		spanName := funcPaths[len(funcPaths)-1]
		ctx, span := otel.GetTracerProvider().Tracer(version.TraceInstrumentationName,
			trace.WithInstrumentationVersion(version.TelemetrySDKVersion),
			trace.WithSchemaURL(version.TraceInstrumentationURL)).Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindInternal))
		span.SetAttributes(attribute.String("func.path", fmt.Sprintf("%s:%v", file, linkNo)))
		return ctx, span
	}
}

// EndSpan 关闭span
func EndSpan(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "OK")
	}
	span.End()
}

func watchConfigMap(clientset *kubernetes.Clientset) {
	configMapClient := clientset.CoreV1().ConfigMaps("")
	fmt.Printf("[TelemetrySDK]%s: Starting to watch ConfigMaps...\n", time.Now().Format("2006-01-02 15:04:05"))

	go func() {
		fmt.Printf("[TelemetrySDK]%s: ConfigMap Watcher Goroutine Start\n", time.Now().Format("2006-01-02 15:04:05"))
		var tc config.CmTraceConfig

		// 无限循环，防止监听器异常退出
		for {
			// 监听指定ConfigMap
			fmt.Printf("[TelemetrySDK]%s: Create ConfigMap Watcher\n", time.Now().Format("2006-01-02 15:04:05"))
			watcher, err := configMapClient.Watch(context.TODO(), metav1.ListOptions{FieldSelector: fmt.Sprintf("metadata.name=%s", config.CmName)})
			if err != nil {
				fmt.Printf("[TelemetrySDK]%s: Failed to watch ConfigMaps: %+v\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
				return
			}

			fmt.Printf("[TelemetrySDK]%s: Start Watch ConfigMap\n", time.Now().Format("2006-01-02 15:04:05"))
			for event := range watcher.ResultChan() {
				fmt.Printf("[TelemetrySDK]%s: ConfigMap Event: %s\n", time.Now().Format("2006-01-02 15:04:05"), event.Type)
				switch event.Type {
				case watch.Added:
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Added: %s\n", time.Now().Format("2006-01-02 15:04:05"), event.Object.(*corev1.ConfigMap).Name)

					err := yaml.Unmarshal([]byte(event.Object.(*corev1.ConfigMap).Data[config.CmMapKeyTrace]), &tc)
					if err != nil {
						fmt.Printf("[TelemetrySDK]%s: error: %+v\n", time.Now().Format("2006-01-02 15:04:05"), err)
					}

					fmt.Printf("[TelemetrySDK]%s: Trace Config Content: %+v\n", time.Now().Format("2006-01-02 15:04:05"), &tc)

					UpdateTracerClient(config.GetTraceEnabled(&tc), tc.Endpoint)
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Add Event Complate.\n", time.Now().Format("2006-01-02 15:04:05"))
				case watch.Modified:
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Modified: %s\n", time.Now().Format("2006-01-02 15:04:05"), event.Object.(*corev1.ConfigMap).Name)

					err := yaml.Unmarshal([]byte(event.Object.(*corev1.ConfigMap).Data[config.CmMapKeyTrace]), &tc)
					if err != nil {
						fmt.Printf("[TelemetrySDK]%s: error: %+v\n", time.Now().Format("2006-01-02 15:04:05"), err)
					}

					fmt.Printf("[TelemetrySDK]%s: Trace Config Content: %+v\n", time.Now().Format("2006-01-02 15:04:05"), &tc)

					UpdateTracerClient(config.GetTraceEnabled(&tc), tc.Endpoint)
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Modify Event Complate.\n", time.Now().Format("2006-01-02 15:04:05"))
				case watch.Deleted:
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Deleted: %s\n", time.Now().Format("2006-01-02 15:04:05"), event.Object.(*corev1.ConfigMap).Name)
					UpdateTracerClient("false", "")
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Delete Event Complate.\n", time.Now().Format("2006-01-02 15:04:05"))
				}
			}
		}

	}()

}
