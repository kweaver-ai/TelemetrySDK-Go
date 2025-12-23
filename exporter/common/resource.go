package common

import (
	"encoding/json"

	"go.opentelemetry.io/otel/sdk/resource"
)

// Resource 自定义 Resource 统一attribute。
type Resource struct {
	Attributes []*Attribute `json:"Attributes"`
	SchemaURL  string       `json:"SchemaURL"`
}

// AnyRobotResourceFromResource resource.Resource转换为 Resource。
func AnyRobotResourceFromResource(resource *resource.Resource) *Resource {
	if resource == nil {
		return &Resource{}
	}
	return &Resource{
		Attributes: AnyRobotAttributesFromKeyValues(resource.Set().ToSlice()),
		SchemaURL:  resource.SchemaURL(),
	}

}

// MarshalJSON 只输出 Attributes 不输出 SchemaURL。
func (r Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Attributes)
}
