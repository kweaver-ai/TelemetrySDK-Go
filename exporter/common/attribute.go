package common

import "go.opentelemetry.io/otel/attribute"

// Attribute 自定义Attribute统一Type。
type Attribute struct {
	Key   string `json:"Key"`
	Value Value  `json:"Value"`
}

// Value 自定义Value统一Type。
type Value struct {
	Type  string      `json:"Type"`
	Value interface{} `json:"Value"`
}

// AnyRobotAttributeFromKeyValue 单条KeyValue转换为*Attribute。
func AnyRobotAttributeFromKeyValue(keyValue attribute.KeyValue) *Attribute {
	return &Attribute{
		Key: string(keyValue.Key),
		Value: Value{
			Type:  standardizeValueType(keyValue.Value.Type().String()),
			Value: keyValue.Value.AsInterface(),
		},
	}
}

// AnyRobotAttributesFromKeyValues 批量KeyValue转换为[]*Attribute。
func AnyRobotAttributesFromKeyValues(keyValues []attribute.KeyValue) []*Attribute {
	if keyValues == nil {
		return make([]*Attribute, 0)
	}
	arAttributes := make([]*Attribute, 0, len(keyValues))
	for i := 0; i < len(keyValues); i++ {
		arAttributes = append(arAttributes, AnyRobotAttributeFromKeyValue(keyValues[i]))
	}
	return arAttributes
}

// standardizeValueType 标准化统一 ValueType 为各语言统一格式。
func standardizeValueType(valueType string) string {
	switch valueType {
	//case "BOOL":
	//	return "BOOL"
	case "BOOLSLICE":
		return "BOOLARRAY"
	case "INT64":
		return "INT"
	case "INT64SLICE":
		return "INTARRAY"
	case "FLOAT64":
		return "FLOAT"
	case "FLOAT64SLICE":
		return "FLOATARRAY"
	//case "STRING":
	//	return "STRING"
	case "STRINGSLICE":
		return "STRINGARRAY"
	default:
		return valueType
	}
}

// AnyRobotAttributesFromSet 批量Set转换为[]*Attribute。
func AnyRobotAttributesFromSet(set attribute.Set) []*Attribute {
	return AnyRobotAttributesFromKeyValues(set.ToSlice())
}
