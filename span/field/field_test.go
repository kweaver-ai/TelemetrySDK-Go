package field

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArrayField(t *testing.T) {
	capacity := 10
	length := 11
	f := MallocArrayField(capacity)

	for i := 0; i < length; i += 1 {
		f.Append(IntField(i))
	}

	for i := 0; i < length; i += 1 {
		v := IntField(i)
		assert.Equal(t, v, (*f)[i])
		v1, err := f.At(i)
		assert.Equal(t, v, v1)
		assert.Equal(t, err, nil)
	}

	fs := MallocStructField(1)
	f.Append(fs)
	assert.Equal(t, (*f)[length], fs)

	fa := MallocArrayField(0)
	f.Append(fa)
	assert.Equal(t, (*f)[length+1], fa)

	assert.Equal(t, f.Length(), length+2)

	v, err := f.At(length + 3)
	assert.Equal(t, err, OverIndexError)
	assert.Equal(t, v, nil)

}

func TestStructField(t *testing.T) {
	capacity := 10
	length := 11
	f := MallocStructField(capacity)
	for i := 0; i < length; i += 1 {
		f.Set(strconv.Itoa(i), IntField(i))
	}

	for i := 0; i < length; i += 1 {
		k := strconv.Itoa(i)
		v := IntField(i)
		assert.Equal(t, v, f.values[i])
		assert.Equal(t, k, f.keys[i])
		k1, v1, err := f.At(i)
		assert.Equal(t, k1, k)
		assert.Equal(t, v1, v)
		assert.Equal(t, err, nil)
	}

	fs := MallocStructField(1)
	fs.Set("test", StringField("sting"))
	f.Set("struct", fs)
	assert.Equal(t, f.keys[length], "struct")
	assert.Equal(t, f.values[length], fs)

	fa := MallocArrayField(0)
	f.Set("array", fa)
	assert.Equal(t, f.keys[length+1], "array")
	assert.Equal(t, f.values[length+1], fa)

	assert.Equal(t, f.Length(), length+2)

	k, v, err := f.At(length + 3)
	assert.Equal(t, k, "")
	assert.Equal(t, v, nil)
	assert.Equal(t, err, OverIndexError)
}

func TestFieldType(t *testing.T) {
	assert.Equal(t, FieldType(IntType), IntField(0).Type())
	assert.Equal(t, FieldType(Float64Type), Float64Field(0).Type())
	assert.Equal(t, FieldType(StringType), StringField("").Type())
	assert.Equal(t, FieldType(TimeType), TimeField(time.Now()).Type())
	assert.Equal(t, FieldType(ArrayType), MallocArrayField(0).Type())
	assert.Equal(t, FieldType(StructType), MallocStructField(0).Type())

	assert.Equal(t, FieldType(JsonType), (&JsonFiled{}).Type())
}

func TestMallocJsonField(t *testing.T) {
	type People struct {
	}
	var p = &People{}
	j := MallocJsonField(p)
	assert.Equal(t, FieldType(JsonType), j.Type())
}

func TestAllProtect(t *testing.T) {
	var (
		intField     IntField
		float64Field Float64Field
		stringField  StringField
		timeField    TimeField
		arrayField   ArrayField
		structField  StructField
		jsonFiled    JsonFiled
		mapField     MapField
	)
	intField.protect()
	float64Field.protect()
	stringField.protect()
	timeField.protect()
	arrayField.protect()
	structField.protect()
	jsonFiled.protect()
	mapField.protect()
}

func TestMapFieldAppend(t *testing.T) {
	type args struct {
		key   string
		value Field
	}
	tests := []struct {
		name string
		f    MapField
		args args
	}{
		{
			"",
			nil,
			args{
				key:   "123",
				value: StringField("456"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Append(tt.args.key, tt.args.value)
		})
	}
}
