package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

// YamlTraceConfig 链路数据记录器配置，结构体映射到YAML数据结构
type YamlTraceConfig struct {
	Enabled  string `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

// YamlLogConfig 程序日志记录器配置，结构体映射到YAML数据结构
type YamlLogConfig struct {
	Enabled   string              `mapstructure:"enabled"`
	Endpoint  string              `mapstructure:"endpoint"`
	Exporters *ExportersTypConfig `mapstructure:"exporters"` //新的输出配置
	Level     string              `mapstructure:"level"`
}

type ExportersTypConfig struct {
	FileExporters     *FileExporterTyp     `mapstructure:"file" yaml:"file"  `
	ConsoleExporter   *ConsoleExporterTyp  `mapstructure:"console" yaml:"console"`
	HttpExporters     *HttpExporterTyp     `mapstructure:"http" yaml:"http"`
	ProtonMqExporters *ProtonMqExporterTyp `mapstructure:"proton_mq" yaml:"proton_mq"`
}

// ExportersFileTyp
type FileExporterConfig struct {
	Path string `mapstructure:"path" yaml:"path"`
}
type FileExporterTyp struct {
	Enable bool               `mapstructure:"enable" yaml:"enable"`
	Config FileExporterConfig `mapstructure:"config" yaml:"config"`
}

type ConsoleExporterTyp struct {
	Enable bool `mapstructure:"enable" yaml:"enable"`
}

// HttpExportersConfig
type HttpExporterConfig struct {
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint"`
}
type HttpExporterTyp struct {
	Enable bool               `mapstructure:"enable" yaml:"enable"`
	Config HttpExporterConfig `mapstructure:"config" yaml:"config"`
}

// ProtonmqExportersConfig
type ProtonmqExporterConfig struct {
	SubType    ExportersSubTyp `mapstructure:"sub_type" yaml:"sub_type"`
	BrokerList []string        `mapstructure:"broker_list" yaml:"broker_list"`
	Topic      string          `mapstructure:"topic" yaml:"topic"` //Topic
	UserName   string          `mapstructure:"username" yaml:"username"`
	PassWord   string          `mapstructure:"password" yaml:"password"`
}
type ProtonMqExporterTyp struct {
	Enable bool                   `mapstructure:"enable" yaml:"enable"`
	Config ProtonmqExporterConfig `mapstructure:"config" yaml:"config"`
}

func (subTyp ExportersSubTyp) String() string {
	return string(subTyp)
}

var (
	// ConfigFile 配置文件信息
	cfgFilePath      = "./"
	CfgFileNameTrace = "ob-app-config-trace"
	CfgFileNameLog   = "ob-app-config-log"
	cfgFileType      = "yaml"
)

var YamlTraceCfg *YamlTraceConfig
var YamlLogCfg *YamlLogConfig
var TraceVP *viper.Viper
var LogVP *viper.Viper

var traceCfgOnce sync.Once
var logCfgOnce sync.Once

// NewTraceConfig 初始化配置
func NewTraceConfig() {
	traceCfgOnce.Do(func() {
		if TraceVP == nil || YamlTraceCfg == nil {
			YamlTraceCfg = &YamlTraceConfig{}
			TraceVP = viper.New()
			initTraceConfig()
		}
	})

}

// NewLogConfig 初始化配置
func NewLogConfig() {
	logCfgOnce.Do(func() {
		if LogVP == nil || YamlLogCfg == nil {
			YamlLogCfg = &YamlLogConfig{}
			LogVP = viper.New()
			initLogConfig()
		}
	})

}

// initTraceConfig 初始化配置
func initTraceConfig() {
	fmt.Printf("Init Trace Setting From File %s%s.%s\n", cfgFilePath, CfgFileNameTrace, cfgFileType)

	TraceVP.AddConfigPath(cfgFilePath)
	TraceVP.SetConfigName(CfgFileNameTrace)
	TraceVP.SetConfigType(cfgFileType)

	LoadTraceConfig()

	TraceVP.WatchConfig()
}

// initLogConfig 初始化配置
func initLogConfig() {
	fmt.Printf("Init Log Setting From File %s%s.%s\n", cfgFilePath, CfgFileNameLog, cfgFileType)

	LogVP.AddConfigPath(cfgFilePath)
	LogVP.SetConfigName(CfgFileNameLog)
	LogVP.SetConfigType(cfgFileType)

	LoadLogConfig()

	LogVP.WatchConfig()
}

// LoadTraceConfig 读取配置文件
func LoadTraceConfig() {
	fmt.Printf("Load Trace Config File %s%s.%s\n", cfgFilePath, CfgFileNameTrace, cfgFileType)

	if err := TraceVP.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	if err := TraceVP.Unmarshal(YamlTraceCfg); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	s, _ := json.Marshal(YamlTraceCfg)
	fmt.Printf("Trace Config Content: %s\n", string(s))
}

// LoadLogConfig 读取配置文件
func LoadLogConfig() {
	fmt.Printf("Load Log Config File %s%s.%s\n", cfgFilePath, CfgFileNameLog, cfgFileType)

	if err := LogVP.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	if err := LogVP.Unmarshal(YamlLogCfg); err != nil {
		fmt.Printf("err:%s\n", err)
	}

	s, _ := json.Marshal(YamlLogCfg)
	fmt.Printf("Log Config Content: %s\n", string(s))
}
