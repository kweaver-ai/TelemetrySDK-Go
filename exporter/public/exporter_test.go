package public

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func channelWithStop() chan struct{} {
	stopCh := make(chan struct{})
	close(stopCh)
	return stopCh
}

func TestExporterExportData(t *testing.T) {
	type fields struct {
		name     string
		client   Client
		stopCh   chan struct{}
		stopOnce *sync.Once
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
			"空数据发本地",
			fields{
				"",
				NewStdoutClient(""),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		},
		{
			"测试数据发本地",
			fields{
				"",
				NewStdoutClient(""),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				byteData(),
			},
			false,
		},
		{
			"已停止的Exporter不能发数据",
			fields{
				"",
				NewStdoutClient(""),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				contextWithDone(),
				nil,
			},
			true,
		},
		{
			"已关闭的Exporter不能发数据",
			fields{
				"",
				NewStdoutClient(""),
				channelWithStop(),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		},
		{
			"空数据发HTTP",
			fields{
				"",
				NewHTTPClient(),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		},
		{
			"测试数据发HTTP",
			fields{
				"",
				NewHTTPClient(),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				context.Background(),
				byteData(),
			},
			// 目标地址正确时不报错。
			true,
		},
		{
			"已停止的Exporter不能发数据",
			fields{
				"",
				NewHTTPClient(),
				make(chan struct{}),
				&sync.Once{},
			},
			args{
				contextWithDone(),
				nil,
			},
			true,
		},
		{
			"已关闭的Exporter不能发数据",
			fields{
				"",
				NewHTTPClient(),
				channelWithStop(),
				&sync.Once{},
			},
			args{
				context.Background(),
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				name:     tt.fields.name,
				client:   tt.fields.client,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if err := e.ExportData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ExportData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExporterName(t *testing.T) {
	type fields struct {
		name     string
		client   Client
		stopCh   chan struct{}
		stopOnce *sync.Once
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"读取Exporter名字",
			fields{
				name:     "NAME",
				client:   nil,
				stopCh:   nil,
				stopOnce: &sync.Once{},
			},
			"NAME",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				name:     tt.fields.name,
				client:   tt.fields.client,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if got := e.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	type fields struct {
		name     string
		client   Client
		stopCh   chan struct{}
		stopOnce *sync.Once
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
			"关闭本地Exporter",
			fields{
				name:     "",
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				context.Background(),
			},
			false,
		},
		{
			"重复关闭本地Exporter",
			fields{
				name:     "",
				client:   NewStdoutClient(""),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				contextWithDone(),
			},
			true,
		},
		{
			"关闭HTTPExporter",
			fields{
				name:     "",
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				context.Background(),
			},
			false,
		},
		{
			"重复关闭HTTPExporter",
			fields{
				name:     "",
				client:   NewHTTPClient(),
				stopCh:   make(chan struct{}),
				stopOnce: &sync.Once{},
			},
			args{
				contextWithDone(),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exporter{
				name:     tt.fields.name,
				client:   tt.fields.client,
				stopCh:   tt.fields.stopCh,
				stopOnce: sync.Once{},
			}
			if err := e.Shutdown(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewExporter(t *testing.T) {
	type args struct {
		client Client
	}
	tests := []struct {
		name string
		args args
		want *Exporter
	}{
		{
			"创建本地Exporter",
			args{NewStdoutClient("")},
			NewExporter(NewStdoutClient("")),
		},
		{
			"创建HTTPExporter",
			args{NewHTTPClient()},
			NewExporter(NewHTTPClient()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExporter(tt.args.client); !reflect.DeepEqual(got.Name(), tt.want.Name()) {
				t.Errorf("NewExporter() = %v, want %v", got, tt.want)
			}
		})
	}
}
