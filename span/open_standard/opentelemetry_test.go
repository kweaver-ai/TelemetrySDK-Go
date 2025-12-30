package open_standard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/encoder"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/exporter"
	"github.com/kweaver-ai/TelemetrySDK-Go/span/v2/field"
)

func setTestSpance(s field.LogSpan) {
	r1 := field.MallocStructField(2)
	r1.Set("Level", field.IntField(1))
	r1.Set("eventNum", field.IntField(2))
	s.SetRecord(r1)
}

func TestOpenTelemetryWrite(t *testing.T) {

	rootSpan := field.NewSpanFromPool(nil, nil)

	setTestSpance(rootSpan)
	var rootSpans []field.LogSpan
	rootSpans = append(rootSpans, rootSpan)
	buf := bytes.NewBuffer(nil)
	open := OpenTelemetryWriter(
		encoder.NewJsonEncoder(buf),
		field.IntField(0))

	defer func(open Writer) {
		_ = open.Close()
	}(open)

	err := open.Write(rootSpans)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// check result
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
		}
	}

	fmt.Print(buf.String())
}

func TestNewSyncWriter(t *testing.T) {
	w := NewSyncWriter(encoder.NewSyncEncoder(exporter.SyncRealTimeExporter()), nil)
	_ = w.Close()
}
