package encoder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/exporter"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
)

func getTestJson() field.Field {

	type A struct {
		Name string
		Age  int
	}
	var a = &A{Name: "123", Age: 12}

	return field.MallocJsonField(a)
}

func GetTestFieds() []field.Field {
	r0 := fakeTestStructField()

	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))

	r3 := getTestJson()
	return []field.Field{r0, r1, r3}
}

func fakeTestStructField() field.Field {
	_, msg, line, _ := runtime.Caller(1)

	res := field.MallocStructField(4)
	res.Set("Level", field.IntField(0))
	res.Set("Model", field.StringField("test Eacape \\\"\b\f\\t\t{}\r\n"))
	res.Set("Message", field.StringField(fmt.Sprintf("%s: %d", msg, line)))
	res.Set("eventNum", field.IntField(1))
	res.Set("float", field.Float64Field(1.3))

	return res
}

func TestJsonEncoder(t *testing.T) {
	check := map[string]interface{}{}
	var err error

	// test log encoder
	fields := GetTestFieds()

	b := bytes.NewBuffer(nil)
	enc := NewJsonEncoder(b)
	for _, i := range fields {
		if err := enc.Write(i); err != nil {
			t.Error(err)
		}
	}

	// fmt.Println("--------log-------------:")
	left := 0
	i := 0
	betty := b.Bytes()
	for ; i < len(betty); i += 1 {
		if betty[i] == '\n' {
			if err = json.Unmarshal(betty[left:i], &check); err != nil {
				t.Error(err)
			}
			left = i + 1
		}
	}
	fmt.Print(b.String())

}

func TestArrayField(t *testing.T) {
	capacity := 10
	length := 11
	a := field.MallocArrayField(capacity)

	for i := 0; i < length; i += 1 {
		a.Append(field.IntField(i))
	}

	for i := 0; i < length; i += 1 {
		assert.Equal(t, field.IntField(i), (*a)[i])
	}

	b := bytes.NewBuffer(nil)
	enc := NewJsonEncoder(b)
	_ = enc.Write(a)

}

func TestNewJsonEncoderBench(t *testing.T) {
	b := bytes.NewBuffer(nil)
	en := NewJsonEncoderBench(b)
	err := en.Write(field.TimeField(time.Now()))
	if err != nil {
		panic(err)
	}
	_ = en.Close()

}

func TestNewJsonEncoderWithExporters(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	type args struct {
		exporters []exporter.LogExporter
	}
	tests := []struct {
		name string
		args args
		want Encoder
	}{
		{
			"TestNewJsonEncoderWithExporters",
			args{[]exporter.LogExporter{exporter.GetRealTimeExporter()}},
			&JsonEncoder{
				w:            nil,
				buf:          nil,
				bufReal:      bytes.NewBuffer(make([]byte, 0, 4096)),
				End:          _lineFeed,
				logExporters: make(map[string]exporter.LogExporter),
				ctx:          ctx,
				cancelFunc:   cancel,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsonEncoderWithExporters(tt.args.exporters...); !reflect.DeepEqual(got.Close(), tt.want.Close()) {
				t.Errorf("NewJsonEncoderWithExporters(%v), want %v", got, tt.want)
			}
		})
	}
}
