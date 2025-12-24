package public

import (
	"testing"

	msqclient "github.com/AISHU-Technology/proton-mq-sdk-go"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/config"
)

var (
	client msqclient.ProtonMQClient
	err    error
)

func TestMain(m *testing.M) {
	config := ProtonMqConfig{
		SubType:    config.ProtonMqKafka,
		BrokerIp:   "localhost",
		BrokerPort: 31000,
		Topic:      "opentelemetry.log",
		UserName:   "test",
		PassWord:   "testpwd",
	}
	client, err = initProtonMqClient(config)
	if err != nil {
		return
	}
	m.Run()
}

//func TestNewProtonMqClient(t *testing.T) {
//	msg := `[{"Link": {"TraceId": "00000000000000000000000000000000","SpanId": "0000000000000000"},"Timestamp": "2024-03-28T06:28:20.172550047Z","SeverityText": "Info","Body": {"Message": "AnyRobot Logger init success"},"Attributes": {},"Resource": {"host":{"arch":"x86_64","ip":"localhost","name":"ubuntu"},"os":{"description":"ubuntu","type":"linux","version":"22.04"},"service":{"instance":{"id":""},"name":"my-service-2","version":"UnknownServiceVersion"},"telemetry":{"sdk":{"language":"go","name":"TelemetrySDK-Go/exporter/ar_log","version":"2.7.5"}}}},{"Link": {"TraceId": "00000000000000000000000000000000","SpanId": "0000000000000000"},"Timestamp": "2024-03-28T06:28:20.172740878Z","SeverityText": "Info","Body": {"Message": "/root/GoProject/TelemetrySDK-Go/exporter/log/main.go:13:main: this is log"},"Attributes": {},"Resource": {"host":{"arch":"x86_64","ip":"localhost","name":"ubuntu"},"os":{"description":"ubuntu","type":"linux","version":"22.04"},"service":{"instance":{"id":""},"name":"my-service-2","version":"UnknownServiceVersion"},"telemetry":{"sdk":{"language":"go","name":"TelemetrySDK-Go/exporter/ar_log","version":"2.7.5"}}}}]`
//
//	if err := client.Pub("opentelemetry.log", []byte(msg)); err != nil {
//
//		t.Error(err)
//		return
//	}
//	t.Log("done")
//}

func Test_randomSlice(t *testing.T) {
	slice := []string{"aaa", "bbb", "ccc"}
	r := randomSlice(slice)
	t.Log(r)
}
