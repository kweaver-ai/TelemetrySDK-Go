package eventsdk

import (
	"encoding/json"

	"github.com/kweaver-ai/TelemetrySDK-Go/event/v2/version"
)

// Resource ，记录资源信息例如服务名、版本号、主机信息等。
type Resource interface {
	// GetSchemaURL 返回 SchemaURL 。
	GetSchemaURL() string
	// GetAttributes 返回 Attributes 。
	GetAttributes() map[string]interface{}
	// Valid 校验是否合法。
	Valid() bool
	// private 禁止用户自己实现接口。
	private()
}

// resource 自定义 event resource 和 Trace 输出格式一致。
type resource struct {
	SchemaURL     string                 `json:"SchemaURL"`
	AttributesMap map[string]interface{} `json:"Attributes"`
}

// newResource 创建新的 *resource 。
func newResource() *resource {
	r := &resource{
		SchemaURL:     version.EventInstrumentationURL,
		AttributesMap: copyDefaultAttributes(),
	}
	return r
}

func (r *resource) GetSchemaURL() string {
	return r.SchemaURL
}

func (r *resource) GetAttributes() map[string]interface{} {
	return r.AttributesMap
}

func (r *resource) Valid() bool {
	return r != nil && r.GetAttributes() != nil && len(r.GetAttributes()) > 0
}

func (r *resource) private() {
	// private 禁止用户自己实现接口。
}

// MarshalJSON 只输出 AttributesMap 。
func (r *resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.AttributesMap)
}

// UnmarshalJSON AttributesMap 变回 resource 。
func (r *resource) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &r.AttributesMap)
	return err
}
