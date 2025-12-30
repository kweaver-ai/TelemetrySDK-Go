package public

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"

	"github.com/kweaver-ai/TelemetrySDK-Go/exporter/v2/config"
)

func TestHTTPClientPath(t *testing.T) {
	type fields struct {
		cfg       *config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"获取上报地址",
			fields{
				cfg:       config.DefaultHTTPConfig(),
				retryFunc: nil,
				client:    nil,
				stopCh:    nil,
			},
			"localhost:5678/api/feed_ingester/v1/jobs/{jobid}/data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			if got := d.Path(); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPClientStop(t *testing.T) {
	type fields struct {
		cfg       *config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
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
			"关闭HTTP Client",
			fields{
				config.DefaultHTTPConfig(),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{context.Background()},
			false,
		},
		{
			"重复关闭HTTP Client",
			fields{
				config.DefaultHTTPConfig(),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{contextWithDone()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			if err := d.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func configWithStatusCode(statusCode string) *config.HTTPConfig {
	cfg := config.DefaultHTTPConfig()
	cfg.Path = statusCode
	cfg.Insecure = false
	cfg.Headers = map[string]string{"self": "defined"}
	cfg.Compression = 0
	return cfg
}

func TestHTTPClientUploadData(t *testing.T) {
	sth := gomonkey.ApplyFunc(send, func(d *HttpClient, req *http.Request) (*http.Response, error) {
		if d.cfg.Path == "500" {
			code := 500
			return &http.Response{StatusCode: code,
				Body:   io.NopCloser(bytes.NewReader([]byte{})),
				Header: map[string][]string{"Retry-After": {"12"}}}, config.RetryableError{Throttle: 12}
		}
		if len(d.cfg.Path) == 3 {
			code, _ := strconv.Atoi(d.cfg.Path)
			return &http.Response{StatusCode: code,
				Body: io.NopCloser(bytes.NewReader([]byte{}))}, nil
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader([]byte{}))}, nil
	})
	defer sth.Reset()

	type fields struct {
		cfg       *config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
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
			"HTTPClient发送测试数据",
			fields{
				config.DefaultHTTPConfig(),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			false,
		},
		{
			"200_StatusOK",
			fields{
				configWithStatusCode("200"),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			false,
		},
		{
			"400_StatusBadRequest",
			fields{
				configWithStatusCode("400"),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			true,
		},
		{
			"404_StatusNotFound",
			fields{
				configWithStatusCode("404"),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			true,
		},
		{
			"413_StatusRequestEntityTooLarge",
			fields{
				configWithStatusCode("413"),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			true,
		},
		{
			"429_StatusTooManyRequests",
			fields{
				configWithStatusCode("429"),
				config.RetryConfig{
					Enabled:         true,
					InitialInterval: 5 * time.Second,
					MaxInterval:     30 * time.Second,
					MaxElapsedTime:  time.Second,
				}.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			true,
		},
		{
			"500_StatusInternalServerError",
			fields{
				configWithStatusCode("500"),
				config.RetryConfig{
					Enabled:         true,
					InitialInterval: 5 * time.Second,
					MaxInterval:     30 * time.Second,
					MaxElapsedTime:  time.Second,
				}.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			true,
		},
		{
			"511_StatusNetworkAuthenticationRequired",
			fields{
				configWithStatusCode("511"),
				config.RetryConfig{
					Enabled:         true,
					InitialInterval: 5 * time.Second,
					MaxInterval:     30 * time.Second,
					MaxElapsedTime:  time.Second,
				}.RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{
				context.Background(),
				byteData(),
			},
			true,
		},
		{
			"已关闭的HTTPClient发送数据",
			fields{
				config.DefaultHTTPConfig(),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
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
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			if err := d.UploadData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UploadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func cancelFuncWithContext() context.CancelFunc {
	_, cancel := context.WithCancel(context.Background())
	return cancel
}

func contextWithCancelFunc() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel
	return ctx
}

func TestHTTPClientContextWithStop(t *testing.T) {
	type fields struct {
		cfg       *config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Context
		want1  context.CancelFunc
	}{
		{
			"未关闭的HTTPClient",
			fields{
				config.DefaultHTTPConfig(),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{ctx: context.Background()},
			contextWithCancelFunc(),
			cancelFuncWithContext(),
		},
		{
			"已关闭的HTTPClient",
			fields{
				config.DefaultHTTPConfig(),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				make(chan struct{}),
			},
			args{ctx: contextWithDone()},
			contextWithDone(),
			cancelFuncWithContext(),
		},
		{
			"已关闭的HTTPClient",
			fields{
				config.DefaultHTTPConfig(),
				config.DefaultRetryConfig().RetryFunc(),
				http.DefaultClient,
				channelWithStop(),
			},
			args{ctx: context.Background()},
			contextWithCancelFunc(),
			cancelFuncWithContext(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			got, got1 := d.contextWithStop(tt.args.ctx)
			if got.Err() != nil && tt.want.Err() != nil {
				assert.Equal(t, tt.want.Err().Error(), got.Err().Error())
			} else {
				assert.Equal(t, tt.want.Err(), got.Err())
			}
			// if !reflect.DeepEqual(got.Err(), tt.want.Err()) {
			// 	t.Errorf("contextWithStop() got = %v, want %v", got, tt.want)
			// }
			if got1 == nil || tt.want1 == nil {
				t.Errorf("contextWithStop() got1 = %v, want %v", &got1, &tt.want1)
			}
		})
	}
}

func TestHTTPClientGetScheme(t *testing.T) {
	type fields struct {
		cfg       *config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"获取发送方式",
			fields{
				cfg:       &config.HTTPConfig{},
				retryFunc: nil,
				client:    nil,
				stopCh:    nil,
			},
			"https",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			if got := d.getScheme(); got != tt.want {
				t.Errorf("getScheme() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPClientNewRequest(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
	type fields struct {
		cfg       *config.HTTPConfig
		retryFunc config.RetryFunc
		client    *http.Client
		stopCh    chan struct{}
	}
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    arRequest
		wantErr bool
	}{
		{
			"创建HTTP请求",
			fields{
				cfg:       &config.HTTPConfig{},
				retryFunc: nil,
				client:    nil,
				stopCh:    nil,
			},
			args{
				[]byte{},
			},
			arRequest{
				Request:    r,
				bodyReader: bodyReader([]byte{}),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &HttpClient{
				cfg:       tt.fields.cfg,
				retryFunc: tt.fields.retryFunc,
				client:    tt.fields.client,
				stopCh:    tt.fields.stopCh,
			}
			got, err := d.newRequest(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Method, tt.want.Method) {
				t.Errorf("newRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHTTPClient(t *testing.T) {
	type args struct {
		opts []config.Option
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			"创建HTTPClient",
			args{nil},
			NewHTTPClient(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHTTPClient(tt.args.opts...); !reflect.DeepEqual(got.Path(), tt.want.Path()) {
				t.Errorf("NewHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARRequestReset(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
	type fields struct {
		Request    *http.Request
		bodyReader func() io.ReadCloser
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"重置arRequest",
			fields{
				Request:    r,
				bodyReader: bodyReader([]byte{}),
			},
			args{
				context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &arRequest{
				Request:    tt.fields.Request,
				bodyReader: tt.fields.bodyReader,
			}
			r.reset(tt.args.ctx)
		})
	}
}

func TestBodyReader(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want func() io.ReadCloser
	}{
		{
			"返回Reader",
			args{
				[]byte{},
			},
			bodyReader([]byte{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bodyReader(tt.args.buf); !reflect.DeepEqual(got(), tt.want()) {
				t.Errorf("bodyReader() = %v, want %v", got(), tt.want())
			}
		})
	}
}

func TestNewResponseError(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"可重发错误",
			args{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := newResponseError(tt.args.header); (err != nil) != tt.wantErr {
				t.Errorf("newResponseError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSend(t *testing.T) {
	req, _ := http.NewRequest(http.MethodHead, "http://255.255.255.255", nil)
	res, _ := http.DefaultClient.Do(req)
	type args struct {
		d   *HttpClient
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			"发送HTTP请求",
			args{
				d: &HttpClient{
					cfg:       config.DefaultConfig().HTTPConfig,
					retryFunc: config.DefaultConfig().RetryConfig.RetryFunc(),
					stopCh:    make(chan struct{}),
					client:    http.DefaultClient,
				},
				req: req,
			},
			res,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := send(tt.args.d, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("send() got = %v, want %v", got, tt.want)
			}
		})
	}
}
