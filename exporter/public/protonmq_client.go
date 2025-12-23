package public

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	msqclient "github.com/AISHU-Technology/proton-mq-sdk-go"
	"github.com/pkg/errors"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/cipters"
	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/config"
)

// ProtonMqClient 客户端结构体。
type ProtonMqClient struct {
	stopCh chan struct{}
	client msqclient.ProtonMQClient
	cfg    *config.ProtonMqExporterTyp
	Broker string
}

type ProtonMqConfig struct {
	SubType    config.ExportersSubTyp
	BrokerIp   string
	BrokerPort int
	Topic      string
	UserName   string
	PassWord   string
	opts       []msqclient.ClientOpt
}

// Path 获取上报地址。
func (c *ProtonMqClient) Path() string {
	return fmt.Sprintf("%s--%s", c.cfg.Config.SubType.String(), c.Broker)
}

// Stop 关闭发送器。
func (c *ProtonMqClient) Stop(ctx context.Context) error {
	if c.client != nil {
		c.client.Close()
	}
	if c.stopCh != nil {
		close(c.stopCh)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// UploadData 批量发送可观测性数据。
func (c *ProtonMqClient) UploadData(ctx context.Context, data []byte) error {
	// 退出逻辑关闭了发送。
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.stopCh:
		return nil
	default:

	}
	var (
		//topic 命名需要确定，谁来定义
		topic = c.cfg.Config.Topic
	)
	if len(topic) == 0 {
		return errors.New("未配置正确的Topic")
	}

	if err := c.Pub(topic, data); err != nil {
		return err
	}
	return nil
}

func (c *ProtonMqClient) Pub(topic string, data []byte) error {
	if err := c.client.Pub(topic, data); err != nil {
		return err
	}
	return nil
}

// NewProtonMqClient 创建Exporter的控制台+本地文件发送客户端。
func NewProtonMqClient(config *config.ProtonMqExporterTyp) (Client, error) {
	if config == nil || !config.Enable {
		return nil, errors.New("NewProtonMqClient 未找到配置文件或者未启用")
	}

	var (
		opts                 []msqclient.ClientOpt
		protonMqOutputConfig = config.Config
		username             = protonMqOutputConfig.UserName
		password             = protonMqOutputConfig.PassWord
		brokerList           = protonMqOutputConfig.BrokerList
	)
	if len(brokerList) <= 0 {
		return nil, errors.New("brokerList 为空")
	}
	//broker处理
	randomBroker := randomSlice(brokerList)
	brokerInfos := strings.Split(randomBroker, ":") //随机取出一个broker
	if len(brokerInfos) < 2 {
		return nil, errors.New("brokerList 格式有误")
	}

	var (
		pubServer  = brokerInfos[0]
		pubPort, _ = strconv.Atoi(brokerInfos[1])
	)
	if len(pubServer) <= 0 {
		return nil, errors.New("brokerList ip有误")
	}
	if pubPort <= 0 {
		return nil, errors.New("brokerList 端口有误")
	}

	//如果存在用户名密码 那么加入 验证
	//如果存在用户名密码 那么加入 验证
	if len(username) > 0 || len(password) > 0 {
		user, err := cipters.RsaDecryptBase64(username)
		if err != nil {
			return nil, errors.Wrap(err, "RsaDecryptBase64 username")
		}
		username = user

		passwd, err := cipters.RsaDecryptBase64(password)
		if err != nil {
			return nil, errors.Wrap(err, "RsaDecryptBase64 password")
		}
		password = passwd

		opts = []msqclient.ClientOpt{
			msqclient.UserInfo(username, password),
			msqclient.AuthMechanism("PLAIN"),
		}
	}

	client, err := initProtonMqClient(ProtonMqConfig{
		SubType:    config.Config.SubType,
		BrokerIp:   pubServer,
		BrokerPort: pubPort,
		UserName:   username,
		PassWord:   password,
		opts:       opts,
	})
	if err != nil {
		return nil, err
	}
	return &ProtonMqClient{cfg: config, client: client, Broker: randomBroker, stopCh: make(chan struct{})}, nil
}

func initProtonMqClient(config ProtonMqConfig) (msqclient.ProtonMQClient, error) {
	var (
		pubServer = config.BrokerIp
		pubPort   = config.BrokerPort
		subType   = config.SubType
		opts      = config.opts
	)

	client, err := msqclient.NewProtonMQClient(pubServer, pubPort, pubServer, pubPort, subType.String(), opts...)

	if err != nil {
		return nil, errors.Errorf("failed to create a proton mq client: %+v", err)
	}
	return client, nil
}

// randomSlice 随机获取一个slice
func randomSlice(slice []string) string {
	if len(slice) <= 0 {
		return ""
	}
	rand.New(rand.NewSource(time.Now().Unix()))
	// 随机获取数组元素
	randomIndex := rand.Intn(len(slice))
	return slice[randomIndex]
}
