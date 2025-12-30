package public

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kweaver-ai/TelemetrySDK-Go/event/v2/eventsdk"
)

func TestNewStdoutClient(t *testing.T) {
	type args struct {
		stdoutPath string
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			"创建未指定输出文件名的StdoutClient",
			args{stdoutPath: ""},
			NewStdoutClient(""),
		},
		{
			"创建指定输出文件名的StdoutClient",
			args{stdoutPath: "./simple.rst"},
			NewStdoutClient("./simple.rst"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStdoutClient(tt.args.stdoutPath); !reflect.DeepEqual(got.Path(), tt.want.Path()) {
				t.Errorf("NewStdoutClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdoutClientPath(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"获取上报地址",
			fields{
				filepath: "/path",
				stopCh:   nil,
			},
			"/path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &StdoutClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if got := d.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdoutClientStop(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"关闭运行中的StdoutClient",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{context.Background()},
			false,
		},
		{
			"重复关闭StdoutClient",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{contextWithDone()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &StdoutClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func byteData() []byte {
	data, _ := json.Marshal(eventsdk.NewEvent())
	return data
}

func TestStdoutClientUploadData(t *testing.T) {
	type fields struct {
		filepath string
		stopCh   chan struct{}
	}
	type args struct {
		ctx  context.Context
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"StdoutClient数据写本地",
			fields{
				filepath: "./AnyRobotData.json",
				stopCh:   make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			false,
		},
		{
			"已关闭的StdoutClient写不了数据",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{
				contextWithDone(),
				byteData(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &StdoutClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.UploadData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UploadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
