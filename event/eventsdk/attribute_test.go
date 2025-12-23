package eventsdk

import (
	"reflect"
	"testing"

	"github.com/shirou/gopsutil/v3/host"
)

func TestNewAttribute(t *testing.T) {
	type args struct {
		key string
		v   interface{}
	}
	tests := []struct {
		name string
		args args
		want Attribute
	}{
		{
			"",
			args{
				key: "key",
				v:   "STRING",
			},
			&attribute{
				Key:   "key",
				Value: "STRING",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAttribute(tt.args.key, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeGetKey(t *testing.T) {
	type fields struct {
		Key   string
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				Key:   "key",
				Value: 123,
			},
			"key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.GetKey(); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeGetValue(t *testing.T) {
	type fields struct {
		Key   string
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			"",
			fields{
				Key:   "key",
				Value: 123,
			},
			123,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeValid(t *testing.T) {
	type fields struct {
		Key   string
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				Key:   "key",
				Value: 123,
			},
			true,
		},
		{
			"",
			fields{
				Key:   "",
				Value: 123,
			},
			false,
		},
		{
			"",
			fields{
				Key:   "os",
				Value: 123,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributeKeyNotNil(t *testing.T) {
	type fields struct {
		Key   string
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				Key:   "key",
				Value: 123,
			},
			true,
		},
		{
			"",
			fields{
				Key:   "",
				Value: 123,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			if got := a.keyNotNil(); got != tt.want {
				t.Errorf("keyNotNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributePrivate(t *testing.T) {
	type fields struct {
		Key   string
		Value interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				Key:   "key",
				Value: 123,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &attribute{
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
			}
			a.private()
		})
	}
}

func TestGetDefaultAttributes(t *testing.T) {
	tests := []struct {
		name string
		want map[string]interface{}
	}{
		{
			"",
			getDefaultAttributes(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDefaultAttributes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostIP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"",
			getHostIP(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostIP(); got != tt.want {
				t.Errorf("getHostIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHostInfo(t *testing.T) {
	tests := []struct {
		name string
		want *host.InfoStat
	}{
		{
			"",
			getHostInfo(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHostInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHostInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
