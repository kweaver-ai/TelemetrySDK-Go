package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// configmap相关配置
var CmName = ""
var CmMapKeyTrace = "trace-sdk-config.yaml"
var CmMapKeyLog = "log-sdk-config.yaml"

// CmTraceConfig 链路数据记录器配置，结构体映射到YAML数据结构
type CmTraceConfig struct {
	Enabled       string `yaml:"enabled"`
	Endpoint      string `yaml:"endpoint"`
	EnabledAllPod string `yaml:"enabledAllPod"`
	EnabledPods   string `yaml:"enabledPods"`
}

// CmLogConfig 程序日志记录器配置，结构体映射到YAML数据结构
type CmLogConfig struct {
	Enabled       string              `yaml:"enabled"`
	Endpoint      string              `yaml:"endpoint"`
	Level         string              `yaml:"level"`
	EnabledAllPod string              `yaml:"enabledAllPod"`
	EnabledPods   string              `yaml:"enabledPods"`
	Exporters     *ExportersTypConfig `yaml:"exporters"` //新的输出配置
}

func InitKubeClient() *kubernetes.Clientset {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[TelemetrySDK]创建kubernetes api客户端失败：%v\n", err)
		}
	}()

	// 使用Pod内的Service Account来创建一个kubernetes api客户端
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("[TelemetrySDK]在kubernetes集群主机创建kubernetes api客户端\n")
		// 当在集群外部调试时，使用kubeconfig文件
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("[TelemetrySDK]在kubernetes集群主机创建kubernetes api客户端失败：%v\n", err.Error())
		}
	} else {
		fmt.Printf("[TelemetrySDK]在kubernetes集群内部创建kubernetes api客户端\n")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("[TelemetrySDK]创建kubernetes api客户端失败：%v\n", err.Error())
	}

	return client
}

// GetTraceEnabled 获取链路数据记录器开关配置。如果功能开关为false，则不开启链路数据记录；反之，如果所有微服务开关为true，则开启链路数据记录；
// 反之，则判断pod名称前缀是否在配置中，如在则开启链路数据记录。
func GetTraceEnabled(tc *CmTraceConfig) string {
	if tc.Enabled == "true" {
		if tc.EnabledAllPod == "true" {
			return "true"
		} else {
			for _, item := range strings.Split(tc.EnabledPods, ",") {
				if podName := os.Getenv("HOSTNAME"); getPodNamePrefix(podName) == item {
					return "true"
				}
			}
		}
	}
	return "false"
}

// GetLogEnabled 获取日志记录器开关配置。如果功能开关为false，则不开启日志记录；反之，如果所有微服务开关为true，则开启日志记录；
// 反之，则判断pod名称前缀是否在配置中，如在则开启日志记录。
func GetLogEnabled(lc *CmLogConfig) string {
	if lc.Enabled == "true" {
		if lc.EnabledAllPod == "true" {
			return "true"
		} else {
			for _, item := range strings.Split(lc.EnabledPods, ",") {
				if podName := os.Getenv("HOSTNAME"); getPodNamePrefix(podName) == item {
					return "true"
				}
			}
		}
	}
	return "false"
}

// getPodNamePrefix 获取pod名称前面不变的部分
func getPodNamePrefix(podName string) string {
	if lastIndex := strings.LastIndex(podName, "-"); lastIndex != -1 {
		if isAllDigits(podName[lastIndex+1:]) {
			return podName[:lastIndex]
		} else {
			if last2Index := strings.LastIndex(podName[:lastIndex], "-"); last2Index != -1 {
				return podName[:last2Index]
			} else {
				return podName[:lastIndex]
			}
		}
	} else {
		return podName
	}
}

// isAllDigits 检查字符串是否全部由数字组成。
func isAllDigits(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
