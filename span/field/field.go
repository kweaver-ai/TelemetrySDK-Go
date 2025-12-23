package field

import "time"

type IntField int
type Float64Field float64
type StringField string
type TimeField time.Time
type FieldType int
type JsonFiled struct {
	Data interface{}
}
type MapField map[string]interface{}

const (
	IntType = iota
	Float64Type
	StringType
	TimeType
	ArrayType
	StructType
	MetricType
	JsonType
	MapType
)

type Field interface {
	Type() FieldType
	protect()
}

func (f IntField) Type() FieldType {
	return IntType
}

func (f IntField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f Float64Field) Type() FieldType {
	return Float64Type
}

// Avoiding irrelevant personnel to implement Field interface
func (f Float64Field) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f StringField) Type() FieldType {
	return StringType
}

// Avoiding irrelevant personnel to implement Field interface
func (f StringField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f TimeField) Type() FieldType {
	return TimeType
}

// Avoiding irrelevant personnel to implement Field interface
func (f TimeField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f *ArrayField) Type() FieldType {
	return ArrayType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *ArrayField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f *StructField) Type() FieldType {
	return StructType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *StructField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f *JsonFiled) Type() FieldType {
	return JsonType
}

// Avoiding irrelevant personnel to implement Field interface
func (f *JsonFiled) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func MallocJsonField(data interface{}) *JsonFiled {
	return &JsonFiled{
		Data: data,
	}
}

func (f MapField) Type() FieldType {
	return MapType
}

func (f MapField) protect() {
	// Avoiding irrelevant personnel to implement Field interface
}

func (f MapField) Append(key string, value Field) {
	if f == nil {
		return
	}
	f[key] = value
}

func MallocMapField() MapField {
	return make(map[string]interface{})
}
