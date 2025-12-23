package eventsdk

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func TestNewEvent(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		args args
		want Event
	}{
		{
			"",
			args{""},
			NewEvent(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEvent(); !reflect.DeepEqual(got.GetEventType(), tt.want.GetEventType()) {
				t.Errorf("NewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalEvents(t *testing.T) {
	myEvent := NewEvent()
	myEvents := []Event{myEvent}
	array, _ := json.Marshal(myEvents)
	bety, _ := json.Marshal(make([]Event, 1))
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []Event
		wantErr bool
	}{
		{
			"",
			args{array},
			make([]Event, 1),
			false,
		},
		{
			"",
			args{bety},
			make([]Event, 0),
			true,
		},
		{
			"",
			args{[]byte{}},
			make([]Event, 0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalEvents(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("UnmarshalEvents() got = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func TestEventGetData(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetEventID(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				EventID:   "12345",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			"12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetEventID(); got != tt.want {
				t.Errorf("GetEventID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetEventMap(t *testing.T) {
	myEvent := &event{
		EventID:   "",
		EventType: "",
		Time:      time.Time{},
		Level:     WARN,
		Resource:  nil,
		Subject:   "",
		Link:      &link{},
		Data:      nil,
	}
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     WARN,
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			myEvent.GetEventMap(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetEventMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEventMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetEventType(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "service.call",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			"service.call",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetEventType(); got != tt.want {
				t.Errorf("GetEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetLevel(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Level
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "ERROR",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			level("ERROR"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetLevel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetLink(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Link
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      nil,
				Data:      nil,
			},
			&link{
				TraceID: "",
				SpanID:  "",
			},
		},
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			&link{
				TraceID: "",
				SpanID:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetLink(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetResource(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   Resource
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  newResource(),
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			newResource(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetResource(); !reflect.DeepEqual(got.GetSchemaURL(), tt.want.GetSchemaURL()) {
				t.Errorf("GetResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetSubject(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "operating.obj",
				Link:      &link{},
				Data:      nil,
			},
			"operating.obj",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetSubject(); got != tt.want {
				t.Errorf("GetSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventGetTime(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.GetTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventSend(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				sent:      &sync.Once{},
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.Send()
		})
	}
}

func TestEventSetAttributes(t *testing.T) {
	type fields struct {
		EventID    string
		EventType  string
		Time       time.Time
		Level      level
		Attributes map[string]interface{}
		Resource   *resource
		Subject    string
		Link       *link
		Data       interface{}
	}
	type args struct {
		kvs []Attribute
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				EventID:    "",
				EventType:  "",
				Time:       time.Time{},
				Level:      "",
				Attributes: make(map[string]interface{}, 0),
				Resource:   newResource(),
				Subject:    "",
				Link:       &link{},
				Data:       nil,
			},
			args{[]Attribute{NewAttribute("key", "anything")}},
		},
		{
			"",
			fields{
				EventID:    "",
				EventType:  "",
				Time:       time.Time{},
				Level:      "",
				Attributes: make(map[string]interface{}, 0),
				Resource:   newResource(),
				Subject:    "",
				Link:       &link{},
				Data:       nil,
			},
			args{[]Attribute{NewAttribute("", "anything")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:    tt.fields.EventID,
				EventType:  tt.fields.EventType,
				Time:       tt.fields.Time,
				Level:      tt.fields.Level,
				Attributes: tt.fields.Attributes,
				Resource:   tt.fields.Resource,
				Subject:    tt.fields.Subject,
				Link:       tt.fields.Link,
				Data:       tt.fields.Data,
			}
			e.SetAttributes(tt.args.kvs...)
		})
	}
}

func TestEventSetData(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{struct{}{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetData(tt.args.data)
		})
	}
}

func TestEventSetEventType(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	type args struct {
		eventType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{"type"},
		},
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetEventType(tt.args.eventType)
		})
	}
}

func TestEventSetLink(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	type args struct {
		link trace.SpanContext
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{trace.SpanContext{}},
		},
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: [16]byte{0x1},
				SpanID:  [8]byte{0x1},
			})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetLink(tt.args.link)
		})
	}
}

func TestEventSetSubject(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	type args struct {
		subject string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{"subject"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetSubject(tt.args.subject)
		})
	}
}

func TestEventSetTime(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	type args struct {
		time time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{time.Time{}},
		},
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{time.Now()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.SetTime(tt.args.time)
		})
	}
}

func TestEventPrivate(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			e.private()
		})
	}
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"",
			generateID(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateID(); len(got) != len(tt.want) {
				t.Errorf("generateID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventValid(t *testing.T) {
	type fields struct {
		EventID   string
		EventType string
		Time      time.Time
		Level     level
		Resource  *resource
		Subject   string
		Link      *link
		Data      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"",
			fields{
				EventID:   "01GJWE6272JXT700WY7EG8V52S",
				EventType: "EventExporter/multiply",
				Time:      time.Now(),
				Level:     "INFO",
				Resource:  &resource{"", getDefaultAttributes()},
				Subject:   "",
				Link: &link{
					"5bd0ecc145cc8639007721df27ecda50",
					"4cbf3e2c1e8517e7",
				},
				Data: nil,
			},
			true,
		},
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				EventID:   tt.fields.EventID,
				EventType: tt.fields.EventType,
				Time:      tt.fields.Time,
				Level:     tt.fields.Level,
				Resource:  tt.fields.Resource,
				Subject:   tt.fields.Subject,
				Link:      tt.fields.Link,
				Data:      tt.fields.Data,
			}
			if got := e.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		data interface{}
		opts []EventStartOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.data, tt.args.opts...)
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		data interface{}
		opts []EventStartOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warn(tt.args.data, tt.args.opts...)
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		data interface{}
		opts []EventStartOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Error(tt.args.data, tt.args.opts...)
		})
	}
}

func TestSetServiceInfo(t *testing.T) {
	type args struct {
		ServiceName     string
		ServiceVersion  string
		ServiceInstance string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{
				"MYServiceName",
				"MYServiceVersion",
				"MYServiceInstance",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetServiceInfo(tt.args.ServiceName, tt.args.ServiceVersion, tt.args.ServiceInstance)
		})
	}
}

func TestEventSetLevel(t *testing.T) {
	type fields struct {
		sent       *sync.Once
		EventID    string
		EventType  string
		Time       time.Time
		Level      level
		Attributes map[string]interface{}
		Resource   *resource
		Subject    string
		Link       *link
		Data       interface{}
	}
	type args struct {
		level Level
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"",
			fields{
				EventID:   "",
				EventType: "",
				Time:      time.Time{},
				Level:     "",
				Resource:  nil,
				Subject:   "",
				Link:      &link{},
				Data:      nil,
			},
			args{level: newLevel("WARN")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				sent:       tt.fields.sent,
				EventID:    tt.fields.EventID,
				EventType:  tt.fields.EventType,
				Time:       tt.fields.Time,
				Level:      tt.fields.Level,
				Attributes: tt.fields.Attributes,
				Resource:   tt.fields.Resource,
				Subject:    tt.fields.Subject,
				Link:       tt.fields.Link,
				Data:       tt.fields.Data,
			}
			e.SetLevel(tt.args.level)
		})
	}
}

func TestEventGetAttributes(t *testing.T) {
	type fields struct {
		sent       *sync.Once
		EventID    string
		EventType  string
		Time       time.Time
		Level      level
		Attributes map[string]interface{}
		Resource   *resource
		Subject    string
		Link       *link
		Data       interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			"",
			fields{
				sent:       nil,
				EventID:    "",
				EventType:  "",
				Time:       time.Time{},
				Level:      "",
				Attributes: make(map[string]interface{}),
				Resource:   nil,
				Subject:    "",
				Link:       nil,
				Data:       nil,
			},
			map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &event{
				sent:       tt.fields.sent,
				EventID:    tt.fields.EventID,
				EventType:  tt.fields.EventType,
				Time:       tt.fields.Time,
				Level:      tt.fields.Level,
				Attributes: tt.fields.Attributes,
				Resource:   tt.fields.Resource,
				Subject:    tt.fields.Subject,
				Link:       tt.fields.Link,
				Data:       tt.fields.Data,
			}
			if got := e.GetAttributes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}
