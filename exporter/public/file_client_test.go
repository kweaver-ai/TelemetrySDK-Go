package public

import (
	"context"
	"reflect"
	"testing"
)

func TestNewFileClient(t *testing.T) {
	type args struct {
		stdoutPath string
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			"创建未指定输出文件名的FileClient",
			args{stdoutPath: ""},
			NewFileClient(""),
		},
		{
			"创建指定输出文件名的FileClient",
			args{stdoutPath: "./simple.rst"},
			NewFileClient("./simple.rst"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileClient(tt.args.stdoutPath); !reflect.DeepEqual(got.Path(), tt.want.Path()) {
				t.Errorf("NewFileClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileClientPath(t *testing.T) {
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
			d := &FileClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if got := d.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileClientStop(t *testing.T) {
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
			"关闭运行中的FileClient",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{context.Background()},
			false,
		},
		{
			"重复关闭FileClient",
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
			d := &FileClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileClientUploadData(t *testing.T) {
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
			"FileClient数据写本地",
			fields{
				filepath: "./ObservableData.json",
				stopCh:   make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			false,
		},
		{
			"已关闭的FileClient写不了数据",
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
			d := &FileClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.UploadData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UploadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
