package common

import (
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

var keyValue = attribute.String("123", "123")

func TestAnyRobotAttributeFromKeyValue(t *testing.T) {
	type args struct {
		attribute attribute.KeyValue
	}
	tests := []struct {
		name string
		args args
		want *Attribute
	}{
		{
			"转换空Attribute",
			args{},
			AnyRobotAttributeFromKeyValue(attribute.KeyValue{}),
		},
		{
			"转换非空Attribute",
			args{keyValue},
			AnyRobotAttributeFromKeyValue(keyValue),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotAttributeFromKeyValue(tt.args.attribute); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotAttributeFromKeyValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyRobotAttributesFromKeyValues(t *testing.T) {
	type args struct {
		attributes []attribute.KeyValue
	}
	tests := []struct {
		name string
		args args
		want []*Attribute
	}{
		{
			"转换空Attributes",
			args{},
			[]*Attribute{},
		},
		{
			"转换非空Attributes",
			args{[]attribute.KeyValue{keyValue}},
			AnyRobotAttributesFromKeyValues([]attribute.KeyValue{keyValue}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyRobotAttributesFromKeyValues(tt.args.attributes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnyRobotAttributesFromKeyValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStandardizeValueType(t *testing.T) {
	type args struct {
		valueType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"BOOL转换成BOOL",
			args{valueType: "BOOL"},
			"BOOL",
		},
		{
			"BOOLSLICE转换成BOOLARRAY",
			args{valueType: "BOOLSLICE"},
			"BOOLARRAY",
		},
		{
			"INT64转换成INT",
			args{valueType: "INT64"},
			"INT",
		},
		{
			"INT64SLICE转换成INTARRAY",
			args{valueType: "INT64SLICE"},
			"INTARRAY",
		},
		{
			"FLOAT64转换成FLOAT",
			args{valueType: "FLOAT64"},
			"FLOAT",
		},
		{
			"FLOAT64SLICE转换成FLOATARRAY",
			args{valueType: "FLOAT64SLICE"},
			"FLOATARRAY",
		},
		{
			"STRING转换成STRING",
			args{valueType: "STRING"},
			"STRING",
		},
		{
			"STRINGSLICE转换成STRINGARRAY",
			args{valueType: "STRINGSLICE"},
			"STRINGARRAY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := standardizeValueType(tt.args.valueType); got != tt.want {
				t.Errorf("standardizeValueType() = %v, want %v", got, tt.want)
			}
		})
	}
}
