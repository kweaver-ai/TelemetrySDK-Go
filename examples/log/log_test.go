package examplelog

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/ar_log"
)

func Benchmark_InfoWithProton(b *testing.B) {
	ar_log.InitLogger("yaml", "log-sdk-config-with-proton", "my-service-2")

	for n := 0; n < b.N; n++ {
		ar_log.Info(context.Background(), "this is log")
	}
}

func Benchmark_InfoWithFile(b *testing.B) {
	ar_log.InitLogger("yaml", "log-sdk-config-with-file", "my-service-2")

	for n := 0; n < b.N; n++ {
		ar_log.Info(context.Background(), "this is log")
	}
}

func Test_send(t *testing.T) {
	ar_log.InitLogger("yaml", "log-sdk-config-with-file", "my-service-2")
	for i := 0; i < 1; i++ {
		ar_log.Info(context.Background(), "this is log")
	}
}

func Benchmark_SendTest(b *testing.B) {

	w := &kafka.Writer{
		Addr:  kafka.TCP("10.4.110.244:31000"),
		Topic: "topic_benchmark",
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{Username: "anyrobot", Password: "eisoo.com123"},
		},
		BatchSize:              1,
		BatchTimeout:           time.Second * 10,
		AllowAutoTopicCreation: true,
	}
	defer w.Close()

	for n := 0; n < b.N; n++ {
		err := w.WriteMessages(context.Background(), kafka.Message{Value: []byte("this is log")})
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}

}
