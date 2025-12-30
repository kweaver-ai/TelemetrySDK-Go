package eventsdk

import (
	"net"

	"github.com/shirou/gopsutil/v3/host"

	"github.com/kweaver-ai/TelemetrySDK-Go/event/v2/version"
)

// Attribute 对外暴露的 attribute 接口。
type Attribute interface {
	// Valid 校验 attribute 是否合法。
	Valid() bool
	// GetKey 返回 attribute 的键。
	GetKey() string
	// GetValue 返回 attribute 的值。
	GetValue() interface{}

	// private 禁止用户自己实现接口。
	private()
}

// attribute 自定义 event attribute 和 Trace 输出格式一致。
type attribute struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Data"`
}

// NewAttribute 创建新的 attribute 。
func NewAttribute(key string, v interface{}) Attribute {
	return &attribute{
		Key:   key,
		Value: v,
	}
}

func (a *attribute) Valid() bool {
	return a.keyNotNil()
}

// keyNotNil 校验 attribute.Key 不为空，即有含义。
func (a *attribute) keyNotNil() bool {
	return len(a.Key) > 0
}

func (a *attribute) GetKey() string {
	return a.Key
}

func (a *attribute) GetValue() interface{} {
	return a.Value
}

func (a *attribute) private() {
	// private 禁止用户自己实现接口。
}

// getHostIP 获取主机IP。
func getHostIP() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && !ipNet.IP.IsLinkLocalUnicast() {
			if ipNet.IP.To4() != nil || ipNet.IP.To16() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}

// getHostInfo 获取主机信息。
func getHostInfo() *host.InfoStat {
	info, _ := host.Info()
	return info
}

// getDefaultAttributes 获取默认资源信息。
func getDefaultAttributes() map[string]interface{} {
	// 获取本机IP
	ip := getHostIP()
	info := getHostInfo()
	result := make(map[string]interface{})
	// 主机信息
	hostMap := make(map[string]string, 3)
	result["host"] = hostMap
	hostIP := ip
	hostMap["ip"] = hostIP
	hostArch := info.KernelArch
	hostMap["arch"] = hostArch
	hostName := info.Hostname
	hostMap["name"] = hostName
	// 操作系统信息
	osMap := make(map[string]string, 3)
	result["os"] = osMap
	osType := info.OS
	osMap["type"] = osType
	osVersion := info.PlatformVersion
	osMap["version"] = osVersion
	osDescription := info.Platform
	osMap["description"] = osDescription
	// 版本信息
	sdkMap := make(map[string]string, 3)
	telemetryMap := make(map[string]interface{}, 1)
	telemetryMap["sdk"] = sdkMap
	result["telemetry"] = telemetryMap
	sdkLanguage := "go"
	sdkMap["language"] = sdkLanguage
	sdkName := version.EventInstrumentationName
	sdkMap["name"] = sdkName
	sdkVersion := version.EventInstrumentationVersion
	sdkMap["version"] = sdkVersion
	return result
}

var defaultAttributes = getDefaultAttributes()

func copyDefaultAttributes() map[string]interface{} {
	copyMap := make(map[string]interface{})
	for k, v := range defaultAttributes {
		copyMap[k] = v
	}
	// 服务信息
	serviceMap := make(map[string]string, 3)
	copyMap["service"] = serviceMap
	serviceMap["name"] = globalServiceName
	serviceMap["version"] = globalServiceVersion
	serviceMap["instance"] = globalServiceInstance
	return copyMap
}
