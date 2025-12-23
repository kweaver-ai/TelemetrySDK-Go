package common

import (
	"encoding/json"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/sdk/resource"
)

func TestAnyRobotResourceFromResource(t *testing.T) {
	type args struct {
		resource *resource.Resource
	}
	tests := []struct {
		name string
		args args
		want Resource
	}{
		{
			"转换空资源信息",
			args{},
			*AnyRobotResourceFromResource(nil),
		},
		{
			"转换非空资源信息",
			args{resource.Default()},
			Resource{AnyRobotAttributesFromKeyValues(resource.Default().Set().ToSlice()),
				resource.Default().SchemaURL()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := *AnyRobotResourceFromResource(tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotResourceFromResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

var resourceByte, _ = json.Marshal(Resource{}.Attributes)

func TestResourceMarshalJSON(t *testing.T) {
	type fields struct {
		Attributes []*Attribute
		SchemaURL  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			"JSON格式输出Resource",
			fields{
				Attributes: nil,
				SchemaURL:  "",
			},
			resourceByte,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Resource{
				Attributes: tt.fields.Attributes,
				SchemaURL:  tt.fields.SchemaURL,
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
