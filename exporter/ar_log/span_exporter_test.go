package ar_log

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/config"
	"github.com/AISHU-Technology/TelemetrySDK-Go/exporter/v2/public"
)

func TestNewExporter(t *testing.T) {
	type args struct {
		c public.Client
	}
	tests := []struct {
		name string
		args args
		want *SpanExporter
	}{
		{
			"HTTPClient的LogExporter",
			args{c: public.NewHTTPClient()},
			NewExporter(public.NewHTTPClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.c); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestLogExporterExportLogs(t *testing.T) {
	type fields struct {
		Exporter *public.Exporter
	}
	type args struct {
		ctx context.Context
		log []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"发送空数据",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx: context.Background(),
				log: nil,
			},
			false,
		},
		{
			"发送Log",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx: context.Background(),
				log: []byte("test"),
			},
			false,
		},
		{
			"已关闭的Exporter不能发Log",
			fields{Exporter: public.NewExporter(public.NewStdoutClient(""))},
			args{
				ctx: contextWithDone(),
				log: []byte("test"),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &SpanExporter{
				Exporter: tt.fields.Exporter,
			}
			if err := e.ExportLogs(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("ExportLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitLogger(t *testing.T) {
	//file, err := ioutil.ReadFile("C:\\Users\\frank.liu01\\GolandProjects\\TelemetrySDK-Go\\examples\\api_service\\ob-app-config-log.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var data = config.YamlLogConfig{}
	//if err := yaml.Unmarshal(file, &data); err != nil {
	//	t.Error(err)
	//	return
	//}

	config.CfgFileNameLog = "ob-app-config-log"
	// 初始化配置
	config.NewLogConfig()
	config.LoadLogConfig()

}

func Test_loadConfigMapData(t *testing.T) {

	Convey("测试加载configMap", t, func() {
		var (
			nameSpace = "default"
			ctx       = context.Background()
		)
		client := fake.NewSimpleClientset()
		_, _ = client.CoreV1().ConfigMaps(nameSpace).Create(ctx, &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "cmConfig",
			},
			Data: map[string]string{
				"ob": "xxx",
			},
		}, metav1.CreateOptions{})

		Convey("正常加载", func() {
			value, err := loadConfigMapData(ctx, client, nameSpace, "cmConfig", "ob")
			So(err, ShouldBeNil)
			So(value, ShouldEqual, "xxx")
		})

		Convey("异常configMap", func() {
			_, err := loadConfigMapData(ctx, client, nameSpace, "xxx", "xxx")
			So(err, ShouldBeError)
		})

		Convey("异常key", func() {
			_, err := loadConfigMapData(ctx, client, nameSpace, "cmConfig", "xxx")
			So(err, ShouldBeError)
		})

	})
}

func Test_getNamespace(t *testing.T) {
	Convey("Test_getNamespace", t, func() {
		Convey("读取环境变量", func() {
			sth := gomonkey.ApplyFunc(os.Getenv, func(v string) string {
				return "TEST"
			})
			defer sth.Reset()
			res, err := getNameSpace()
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "TEST")
		})
		Convey("读取文件", func() {
			sth := gomonkey.ApplyFunc(os.ReadFile, func(v string) ([]byte, error) {
				return []byte("anyrobot"), nil
			})
			defer sth.Reset()
			res, err := getNameSpace()
			So(err, ShouldBeNil)
			So(res, ShouldEqual, "anyrobot")
		})
	})
}
