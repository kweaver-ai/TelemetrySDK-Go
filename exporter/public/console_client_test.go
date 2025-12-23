package public

import (
	"context"
	"reflect"
	"testing"
)

func TestNewConsoleClient(t *testing.T) {
	tests := []struct {
		name string
		want Client
	}{
		{
			"创建ConsoleClient",
			NewConsoleClient(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConsoleClient(); !reflect.DeepEqual(got.Path(), tt.want.Path()) {
				t.Errorf("NewConsoleClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsoleClientPath(t *testing.T) {
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
				filepath: "CONSOLE",
				stopCh:   nil,
			},
			"CONSOLE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &ConsoleClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if got := d.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsoleClientStop(t *testing.T) {
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
			"关闭运行中的ConsoleClient",
			fields{
				filepath: "",
				stopCh:   make(chan struct{}),
			},
			args{context.Background()},
			false,
		},
		{
			"重复关闭ConsoleClient",
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
			d := &ConsoleClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConsoleClientUploadData(t *testing.T) {
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
			"ConsoleClient数据写控制台",
			fields{
				filepath: "CONSOLE",
				stopCh:   make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			false,
		},
		{
			"已关闭的ConsoleClient写不了数据",
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
			d := &ConsoleClient{
				filepath: tt.fields.filepath,
				stopCh:   tt.fields.stopCh,
			}
			if err := d.UploadData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UploadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
