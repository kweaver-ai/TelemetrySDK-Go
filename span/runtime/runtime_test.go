package runtime

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/AISHU-Technology/TelemetrySDK-Go/span/v2/encoder"
	"github.com/AISHU-Technology/TelemetrySDK-Go/span/v2/field"
	"github.com/AISHU-Technology/TelemetrySDK-Go/span/v2/open_standard"
)

func setTestSpance(s field.LogSpan) {
	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))
	s.SetRecord(r1)
}

func TestRecord(t *testing.T) {
	// runtimeSpan := NewRuntime(NewSpanFromPool)
	buf := bytes.NewBuffer(nil)
	writer := open_standard.OpenTelemetryWriter(
		encoder.NewJsonEncoder(buf),
		field.IntField(0))
	runtimeSpan := NewRuntime(writer, field.NewSpanFromPool)
	// go runtimeSpan.Run()
	runtimeSpan.SetUploadInternalAndMaxLog(3*time.Second, 10)
	// task thread0, log quick
	go func() {
		// a new web request task
		root := runtimeSpan.Children(nil) //nolint
		// web request complete, wait children span
		// finally send web's span control to loger runtime
		defer root.Signal()

		// do some task and log
		setTestSpance(root)
	}()
	// task thread1, live 1s
	go func() {
		root := runtimeSpan.Children(nil) //nolint
		defer root.Signal()
		setTestSpance(root)

		time.Sleep(1 * time.Second)
	}()

	// stop runtime after 2s
	go func() {
		time.Sleep(4 * time.Second)
		runtimeSpan.Signal()
	}()

	runtimeSpan.Run()

	// test runtime stop
	afterStop := runtimeSpan.Children(nil) //nolint
	assert.Equal(t, nil, afterStop)

	// check result
	var err error
	capacity := map[string]interface{}{}
	betty := buf.Bytes()
	left := 0
	i := 0
	n := 0
	for ; i < len(betty); i += 1 {
		if betty[i] == '\n' {
			if err = json.Unmarshal(betty[left:i], &capacity); err != nil {
				t.Error(err)
				t.FailNow()
			} else {
				n += 1
			}
			left = i + 1
		}
	}
	if left < len(betty) {
		if err = json.Unmarshal(betty[left:i], &capacity); err != nil {
			t.Error(err)
			t.FailNow()
		} else {
			n += 1
		}
	}

	assert.Equal(t, 2, n)

}

func TestSignal(t *testing.T) {
	writer := &open_standard.OpenTelemetry{
		Encoder: encoder.NewJsonEncoder(io.Discard),
	}
	runtimeSpan := NewRuntime(writer, field.NewSpanFromPool)

	runtimeSpan.Signal()
	runtimeSpan.Signal()
}
