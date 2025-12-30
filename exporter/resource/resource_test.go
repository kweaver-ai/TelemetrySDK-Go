package resource

import (
	"reflect"
	"testing"

	"github.com/shirou/gopsutil/v3/host"
	"go.opentelemetry.io/otel/attribute"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"

	"github.com/kweaver-ai/TelemetrySDK-Go/event/v2/eventsdk"
)

func TestEventResource(t *testing.T) {
	tests := []struct {
		name string
		want eventsdk.EventProviderOption
	}{
		{
			"填充Event资源信息",
			EventResource(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EventResource(); !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("EventResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceInstance(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"获取服务实例ID",
			"UnknownServiceInstance",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceInstanceID(); got != tt.want {
				t.Errorf("GetServiceInstanceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"获取服务名",
			defaultServiceName(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceName(); got != tt.want {
				t.Errorf("GetServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"获取服务版本",
			"UnknownServiceVersion",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceVersion(); got != tt.want {
				t.Errorf("GetServiceVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricResource(t *testing.T) {
	tests := []struct {
		name string
		want *sdkresource.Resource
	}{
		{
			"填充Metric资源信息",
			MetricResource(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MetricResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetServiceInstance(t *testing.T) {
	type args struct {
		serviceInstance string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"设置服务实例ID",
			args{
				"abcd1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetServiceInstance(tt.args.serviceInstance)
		})
	}
}

func TestSetServiceName(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"设置服务名",
			args{
				"TelemetrySDK",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetServiceName(tt.args.serviceName)
		})
	}
}

func TestSetServiceVersion(t *testing.T) {
	type args struct {
		serviceVersion string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"设置服务版本",
			args{
				"2.5.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetServiceVersion(tt.args.serviceVersion)
		})
	}
}

func TestTraceResource(t *testing.T) {
	tests := []struct {
		name string
		want *sdkresource.Resource
	}{
		{
			"填充Trace资源信息",
			TraceResource(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TraceResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TraceResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultServiceName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"获取默认的服务名",
			defaultServiceName(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultServiceName(); got != tt.want {
				t.Errorf("defaultServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostIP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"获取主机IP",
			getHostIP(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostIP(); got != tt.want {
				t.Errorf("getHostIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostInfo(t *testing.T) {
	tests := []struct {
		name string
		want *host.InfoStat
	}{
		{
			"获取主机信息",
			getHostInfo(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHostInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInnerAttributes(t *testing.T) {
	tests := []struct {
		name string
		want []attribute.KeyValue
	}{
		{
			"Resource中的公共部分",
			innerAttributes(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := innerAttributes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("innerAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
