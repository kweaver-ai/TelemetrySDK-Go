package eventsdk

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewResource(t *testing.T) {
	tests := []struct {
		name string
		want *resource
	}{
		{
			"",
			newResource(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newResource(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceGetAttributes(t *testing.T) {
	type fields struct {
		SchemaURL     string
		AttributesMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			"",
			fields{
				SchemaURL:     "",
				AttributesMap: make(map[string]interface{}),
			},
			make(map[string]interface{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resource{
				SchemaURL:     tt.fields.SchemaURL,
				AttributesMap: tt.fields.AttributesMap,
			}
			if got := r.GetAttributes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceGetSchemaURL(t *testing.T) {
	type fields struct {
		SchemaURL     string
		AttributesMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				SchemaURL:     "123",
				AttributesMap: nil,
			},
			"123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resource{
				SchemaURL:     tt.fields.SchemaURL,
				AttributesMap: tt.fields.AttributesMap,
			}
			if got := r.GetSchemaURL(); got != tt.want {
				t.Errorf("GetSchemaURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceMarshalJSON(t *testing.T) {
	bety, _ := json.Marshal(make(map[string]interface{}))
	type fields struct {
		SchemaURL     string
		AttributesMap map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"",
			fields{
				SchemaURL:     "1",
				AttributesMap: make(map[string]interface{}),
			},
			bety,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resource{
				SchemaURL:     tt.fields.SchemaURL,
				AttributesMap: tt.fields.AttributesMap,
			}
			got, err := r.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceUnmarshalJSON(t *testing.T) {
	bety, _ := json.Marshal(make(map[string]interface{}))
	type fields struct {
		SchemaURL     string
		AttributesMap map[string]interface{}
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{
				SchemaURL:     "",
				AttributesMap: nil,
			},
			args{bety},
			false,
		},
		{
			"",
			fields{
				SchemaURL:     "",
				AttributesMap: nil,
			},
			args{[]byte{}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resource{
				SchemaURL:     tt.fields.SchemaURL,
				AttributesMap: tt.fields.AttributesMap,
			}
			if err := r.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourcePrivate(t *testing.T) {
	type fields struct {
		SchemaURL     string
		AttributesMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				SchemaURL:     "",
				AttributesMap: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resource{
				SchemaURL:     tt.fields.SchemaURL,
				AttributesMap: tt.fields.AttributesMap,
			}
			r.private()
		})
	}
}

func TestResourceValid(t *testing.T) {
	type fields struct {
		SchemaURL     string
		AttributesMap map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				SchemaURL:     "",
				AttributesMap: getDefaultAttributes(),
			},
			true,
		},
		{
			"",
			fields{
				SchemaURL:     "",
				AttributesMap: nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resource{
				SchemaURL:     tt.fields.SchemaURL,
				AttributesMap: tt.fields.AttributesMap,
			}
			if got := r.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
