package ar_log

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/config"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/public"
	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/resource"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/encoder"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/exporter"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
	spanLog "github.com/kweaver-ai/TelemetrySDK-Go/span/v2/log"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/open_standard"
	sdkRuntime "github.com/kweaver-ai/TelemetrySDK-Go/span/v2/runtime"
)

// 跨包实现接口占位用。
var _ exporter.LogExporter = (*SpanExporter)(nil)
var _ exporter.SyncExporter = (*syncExporter)(nil)
var (
	// Logger 全局程序日志记录器
	Logger spanLog.Logger
	// loggerLock 修改 Logger 对象的锁
	loggerLock = sync.Mutex{}
	// BLogger 全局业务日志记录器
	BLogger spanLog.Logger
)

// SpanExporter 导出数据到AnyRobot Feed Ingester的 Log 数据接收器。
type SpanExporter struct {
	*public.Exporter
}

// ExportLogs 批量发送 log 到AnyRobot Feed Ingester的 Log 数据接收器。
func (e *SpanExporter) ExportLogs(ctx context.Context, logs []byte) error {
	return e.ExportData(ctx, logs)
}

// NewExporter 创建已启动的 LogExporter。
func NewExporter(c public.Client) *SpanExporter {
	return &SpanExporter{
		public.NewExporter(c),
	}
}

// syncExporter 同步导出数据到AnyRobot Feed Ingester的 Log 数据接收器。
type syncExporter struct {
	*public.SyncExporter
}

// ExportLogs 同步发送 log 到AnyRobot Feed Ingester的 Log 数据接收器。
func (s *syncExporter) ExportLogs(ctx context.Context, logs []byte) error {
	return s.ExportData(ctx, logs)
}

// NewSyncExporter 创建已启动的 LogExporter。
func NewSyncExporter(c public.SyncClient) *syncExporter {
	return &syncExporter{
		public.NewSyncExporter(c),
	}
}

// InitLogger 初始化上报到AnyRobot的程序日志记录器
// cfgType 配置类型，可选：cm、yaml
// cfgName 配置文件名称，如配置类型为cm，则该参数为configmap的名称；如配置类型为yaml，则该参数为yaml文件的路径，比如./ob-app.yaml
// serverName 微服务名称
func InitLogger(cfgType string, cfgName string, serverName string) {

	if cfgType == "cm" { // 如果配置为configmap形式
		// 监听configmap的内容，更新全局链路数据记录器的配置
		config.CmName = cfgName
		kubeClient := config.InitKubeClient()
		if kubeClient != nil {
			watchConfigMap(kubeClient)
		}
		//获取 nameSpace
		currentNameSpace, err := getNameSpace()
		if err != nil {
			fmt.Println("初始化SDK失败,必须设置环境变量PRODUCT_NAME或者在容器中运行！")
			return
		}
		//初始化Logger
		initLoggerFromConfigMap(context.Background(), kubeClient, currentNameSpace, cfgName, config.CmMapKeyLog, serverName)
	} else if cfgType == "yaml" { // 如果配置为yaml文件形式
		config.CfgFileNameLog = cfgName
		// 初始化配置
		config.NewLogConfig()
		initARLogger(config.YamlLogCfg, "")
		config.LogVP.OnConfigChange(func(e fsnotify.Event) {
			fmt.Printf("Log config file changed:%s, update logger\n", e)
			config.LoadLogConfig()
			initARLogger(config.YamlLogCfg, "")
		})
	}

}

// initLoggerFromConfigMap 从configMap中加载配置信息
func initLoggerFromConfigMap(ctx context.Context, client *kubernetes.Clientset, nameSpace string, cfgName, configMapKey, serverName string) {
	var (
		logConfig = &config.YamlLogConfig{Enabled: "false", Exporters: &config.ExportersTypConfig{}}
		lc        = config.CmLogConfig{Exporters: &config.ExportersTypConfig{}}
	)
	//加载配置
	data, err := loadConfigMapData(ctx, client, nameSpace, cfgName, configMapKey)
	if err != nil {
		fmt.Printf("[TelemetrySDK] initLoggerFromConfigMap loadConfigMapData error: %v", err)
	}

	if err = yaml.Unmarshal([]byte(data), &lc); err != nil {
		fmt.Printf("[TelemetrySDK] initLoggerFromConfigMap Unmarshal error: %v", err)
	}

	logConfig.Enabled = config.GetLogEnabled(&lc)
	logConfig.Endpoint = lc.Endpoint
	logConfig.Level = lc.Level
	logConfig.Exporters = lc.Exporters

	initARLogger(logConfig, serverName)
}

// getNameSpace 获取当前POD的nameSpace
func getNameSpace() (string, error) {
	if productName := os.Getenv("PRODUCT_NAME"); len(productName) > 0 {
		return productName, nil
	}

	namespaceFile := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespace, err := os.ReadFile(namespaceFile)
	if err != nil {
		return "", err
	}
	return string(namespace), nil
}

// Debug 拼接上文件、行号、函数名。用于日志记录时把位置信息带上
func Debug(ctx context.Context, msg string) {
	pc, filename, line, _ := runtime.Caller(1)
	Logger.Debug(fmt.Sprintf("%s:%d:%s: %v", filename, line, strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), "."), msg),
		field.WithContext(ctx))
}

func Info(ctx context.Context, msg string) {
	pc, filename, line, _ := runtime.Caller(1)
	Logger.Info(fmt.Sprintf("%s:%d:%s: %v", filename, line, strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), "."), msg),
		field.WithContext(ctx))
}

func Warn(ctx context.Context, msg string) {
	pc, filename, line, _ := runtime.Caller(1)
	Logger.Warn(fmt.Sprintf("%s:%d:%s: %v", filename, line, strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), "."), msg),
		field.WithContext(ctx))
}

func Error(ctx context.Context, msg string) {
	pc, filename, line, _ := runtime.Caller(1)
	Logger.Error(fmt.Sprintf("%s:%d:%s: %v", filename, line, strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), "."), msg),
		field.WithContext(ctx))
}

func Fatal(ctx context.Context, msg string) {
	pc, filename, line, _ := runtime.Caller(1)
	Logger.Fatal(fmt.Sprintf("%s:%d:%s: %v", filename, line, strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), "."), msg),
		field.WithContext(ctx))
}

// initARLogger 初始化上报到AnyRobot的日志记录器
// logEnabled 日志开关
// logEndpoint 日志上报地址，为空则打印标准输出
// logLevel 日志等级
// ServerName 微服务名称
func initARLogger(logConfig *config.YamlLogConfig, serverName string) {
	loggerLock.Lock()
	defer loggerLock.Unlock()

	if logConfig == nil {
		return
	}

	var (
		logEnabled = logConfig.Enabled
		logLevel   = logConfig.Level
	)

	serverInstance := os.Getenv("HOSTNAME")
	if logEnabled != "true" {
		logLevel = "off"
	}

	// 初始化ar_log
	var ARLogger = spanLog.NewSamplerLogger(spanLog.WithSample(1.0), spanLog.WithLevel(getLogLevel(logLevel)))

	// 设置微服务相关信息
	if serverName != "" {
		resource.SetServiceName(serverName)
	}
	resource.SetServiceInstance(serverInstance)

	var systemLogWriter open_standard.Writer

	//初始化
	systemLogExporters := initExporters(logConfig)
	if len(systemLogExporters) <= 0 {
		ARLogger.Error(fmt.Sprintf("initARLogger 初始化initExporters数据有误，长度:%d", len(systemLogExporters)))
		return
	}

	systemLogWriter = open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporters...),
		resource.LogResource())
	systemLogRunner := sdkRuntime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)

	go systemLogRunner.Run()
	ARLogger.SetLevel(getLogLevel(logLevel))
	ARLogger.SetRuntime(systemLogRunner)

	ARLogger.Info("AnyRobot Logger init success")

	oldLogger := Logger
	Logger = ARLogger
	if oldLogger != nil {
		oldLogger.Close()
	}
}

// getLogLevel Log配置转换为spanlog配置，默认不填的日志级别为warn
func getLogLevel(level string) int {
	switch level {
	case "all":
		return spanLog.AllLevel
	case "trace":
		return spanLog.TraceLevel
	case "debug":
		return spanLog.DebugLevel
	case "info":
		return spanLog.InfoLevel
	case "warn":
		return spanLog.WarnLevel
	case "error":
		return spanLog.ErrorLevel
	case "fatal":
		return spanLog.FatalLevel
	case "off":
		return spanLog.OffLevel
	default:
		return spanLog.WarnLevel
	}
}

// InitBusinessLogger 初始化业务日志记录器
func InitBusinessLogger() spanLog.Logger {
	// 设置微服务相关信息
	resource.SetServiceName(os.Getenv("TELEMETRY_SERVICE_NAME"))
	resource.SetServiceVersion(os.Getenv("TELEMETRY_SERVICE_VERSION"))
	resource.SetServiceInstance(os.Getenv("HOSTNAME"))

	var businessLogger = spanLog.NewSamplerLogger(spanLog.WithSample(1.0), spanLog.WithLevel(spanLog.AllLevel))
	systemLogExporter := exporter.GetRealTimeExporter()
	systemLogWriter := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoderWithExporters(systemLogExporter),
		resource.LogResource())
	systemLogRunner := sdkRuntime.NewRuntime(systemLogWriter, field.NewSpanFromPool)
	systemLogRunner.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// 运行SystemLogger日志器。
	go systemLogRunner.Run()
	businessLogger.SetRuntime(systemLogRunner)

	businessLogger.Info("AnyRobot BLogger init success")
	return businessLogger
}

// loadConfigMapData 从configMap中获取配置数据
func loadConfigMapData(ctx context.Context, cs kubernetes.Interface, nameSpace, configMapName, configMapKey string) (string, error) {
	if len(nameSpace) <= 0 || len(configMapName) <= 0 || len(configMapKey) <= 0 {
		return "", errors.New("nameSpace or configMapName or configMapKey is empty")
	}

	configMap, err := cs.CoreV1().ConfigMaps(nameSpace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	configMapData, has := configMap.Data[configMapKey]
	if !has {
		return "", fmt.Errorf("从命名空间:%s 获取configMapName:%s,其中 %s key不存在", nameSpace, configMapName, configMapKey)
	}
	return configMapData, nil
}

var K8SClientNamespace = ""

func SetK8SClientNamespace(namespace string) {
	K8SClientNamespace = namespace
}

func watchConfigMap(client *kubernetes.Clientset) {
	configMapClient := client.CoreV1().ConfigMaps(K8SClientNamespace)
	fmt.Printf("[TelemetrySDK]%s: Starting to watch ConfigMaps...\n", time.Now().Format("2006-01-02 15:04:05"))

	go func() {
		fmt.Printf("[TelemetrySDK]%s: ConfigMap Watcher Goroutine Start\n", time.Now().Format("2006-01-02 15:04:05"))
		var lc config.CmLogConfig

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

					err := yaml.Unmarshal([]byte(event.Object.(*corev1.ConfigMap).Data[config.CmMapKeyLog]), &lc)
					if err != nil {
						fmt.Printf("[TelemetrySDK]%s: error: %+v\n", time.Now().Format("2006-01-02 15:04:05"), err)
					}

					by, _ := json.Marshal(&lc)
					fmt.Printf("[TelemetrySDK]%s: Log Config Content: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(by))

					logConfig := &config.YamlLogConfig{
						Enabled:   config.GetLogEnabled(&lc),
						Endpoint:  lc.Endpoint,
						Level:     lc.Level,
						Exporters: lc.Exporters,
					}
					initARLogger(logConfig, "")
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Add Event Complate.\n", time.Now().Format("2006-01-02 15:04:05"))
				case watch.Modified:
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Modified: %s\n", time.Now().Format("2006-01-02 15:04:05"), event.Object.(*corev1.ConfigMap).Name)

					err := yaml.Unmarshal([]byte(event.Object.(*corev1.ConfigMap).Data[config.CmMapKeyLog]), &lc)
					if err != nil {
						fmt.Printf("[TelemetrySDK]%s: error: %+v\n", time.Now().Format("2006-01-02 15:04:05"), err)
					}

					by, _ := json.Marshal(&lc)
					fmt.Printf("[TelemetrySDK]%s: Log Config Content: %s\n", time.Now().Format("2006-01-02 15:04:05"), string(by))

					logConfig := &config.YamlLogConfig{
						Enabled:   config.GetLogEnabled(&lc),
						Endpoint:  lc.Endpoint,
						Level:     lc.Level,
						Exporters: lc.Exporters,
					}
					initARLogger(logConfig, "")
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Modify Event Complate.\n", time.Now().Format("2006-01-02 15:04:05"))
				case watch.Deleted:
					logConfig := &config.YamlLogConfig{
						Enabled:   "false",
						Endpoint:  "",
						Level:     "",
						Exporters: &config.ExportersTypConfig{},
					}
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Deleted: %s\n", time.Now().Format("2006-01-02 15:04:05"), event.Object.(*corev1.ConfigMap).Name)
					initARLogger(logConfig, "")
					fmt.Printf("[TelemetrySDK]%s: ConfigMap Delete Event Complate.\n", time.Now().Format("2006-01-02 15:04:05"))
				}
				fmt.Printf("[TelemetrySDK]%s: Deal Event: %s Complate\n", time.Now().Format("2006-01-02 15:04:05"), event.Type)
			}
		}
	}()
}

// initExportersFromEndpoint 兼容老版本endpont的配置
func initExportersFromEndpoint(logConfig *config.YamlLogConfig, systemLogExporters []exporter.LogExporter) []exporter.LogExporter {
	var (
		//systemLogExporters []exporter.LogExporter
		logEndpoint = logConfig.Endpoint
	)
	//老版本中 未配置endpoint 的表示输出到控制台，配置了endpoint表示输出到http
	if logEndpoint == "" {
		// 设置日志打印标准输出
		systemLogExporters = append(systemLogExporters, exporter.GetRealTimeExporter())
	} else {
		systemLogClient := public.NewHTTPClient(public.WithAnyRobotURL(logEndpoint),
			public.WithCompression(1),
			public.WithTimeout(10*time.Second),
			public.WithRetry(true, 5*time.Second, 20*time.Second, 1*time.Minute))

		systemLogExporters = append(systemLogExporters, NewExporter(systemLogClient))
	}
	return systemLogExporters
}

// initExporters 初始化Exporters
func initExporters(logConfig *config.YamlLogConfig) []exporter.LogExporter {
	//初始化默认silent
	var systemLogExporters []exporter.LogExporter
	systemLogExporters = append(systemLogExporters, NewExporter(public.NewSilentClient()))

	//兼容老版本的配置，没有配置exporter表示是老的配置 继续走老的业务逻辑
	if logConfig.Exporters == nil {
		return initExportersFromEndpoint(logConfig, systemLogExporters)
	}
	//激活所有已启用的exporter
	if logConfig.Exporters.FileExporters != nil && logConfig.Exporters.FileExporters.Enable {
		systemLogExporters = append(systemLogExporters, initFileExporter(logConfig.Exporters.FileExporters))
	}
	if logConfig.Exporters.ConsoleExporter != nil && logConfig.Exporters.ConsoleExporter.Enable {
		systemLogExporters = append(systemLogExporters, initConsoleExporter(logConfig.Exporters.ConsoleExporter))
	}
	if logConfig.Exporters.HttpExporters != nil && logConfig.Exporters.HttpExporters.Enable {
		systemLogExporters = append(systemLogExporters, initHttpExporter(logConfig.Exporters.HttpExporters))
	}
	if logConfig.Exporters.ProtonMqExporters != nil && logConfig.Exporters.ProtonMqExporters.Enable {
		systemLogExporters = append(systemLogExporters, initProtonMqExporter(logConfig.Exporters.ProtonMqExporters))
	}

	return systemLogExporters
}

// initHttpExportersClient 2024-03-25 最新版本的配置
func initHttpExporter(config *config.HttpExporterTyp) exporter.LogExporter {
	// 设置日志通过HTTP上报
	var logEndpoint = config.Config.Endpoint
	systemLogClient := public.NewHTTPClient(public.WithAnyRobotURL(logEndpoint),
		public.WithCompression(1),
		public.WithTimeout(10*time.Second),
		public.WithRetry(true, 5*time.Second, 20*time.Second, 1*time.Minute))
	return NewExporter(systemLogClient)
}

// initProtonMqExporter 初始化protonmq输出
func initProtonMqExporter(config *config.ProtonMqExporterTyp) exporter.LogExporter {
	systemLogClient, err := public.NewProtonMqClient(config)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return NewExporter(systemLogClient)
}

// initFileExporter 初始化文件输出
func initFileExporter(config *config.FileExporterTyp) exporter.LogExporter {
	if config == nil {
		return nil
	}
	stdoutPath := config.Config.Path
	systemLogClient := public.NewFileClient(stdoutPath)
	return NewExporter(systemLogClient)
}

// initConsoleExporter 初始化console输出
func initConsoleExporter(config *config.ConsoleExporterTyp) exporter.LogExporter {
	return exporter.GetRealTimeExporter()
}
