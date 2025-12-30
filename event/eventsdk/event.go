package eventsdk

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/trace"

	"github.com/kweaver-ai/TelemetrySDK-Go/event/v2/custom_errors"
)

// Event 对外暴露的 event 接口。
type Event interface {
	// SetEventType 设置非空 EventType 。
	SetEventType(eventType string)
	// SetTime 设置 Time 。
	SetTime(time time.Time)
	//SetLevel 设置 level 。
	SetLevel(level Level)
	// SetAttributes 设置 Attributes 。
	SetAttributes(kvs ...Attribute)
	// SetSubject 设置 Subject 。
	SetSubject(subject string)
	// SetLink 设置 link 。
	SetLink(link trace.SpanContext)
	// SetData 设置 Data 。
	SetData(data interface{})

	// GetEventID 返回 EventID 。
	GetEventID() string
	// GetEventType 返回 EventType 。
	GetEventType() string
	// GetTime 返回 Time 。
	GetTime() time.Time
	// GetLevel 返回 level 。
	GetLevel() Level
	// GetAttributes 返回 attributes 。
	GetAttributes() map[string]interface{}
	// GetResource 返回 resource 。
	GetResource() Resource
	// GetSubject 返回 Subject 。
	GetSubject() string
	// GetLink 返回 link 。
	GetLink() Link
	// GetData 返回 Data 。
	GetData() interface{}
	// GetEventMap 返回 map[string]interface{} 形式的 event 。
	GetEventMap() map[string]interface{}
	// Send 上报 Event 到 AnyRobot Event 数据接收器。
	Send()

	// private 禁止用户自己实现接口。
	private()
	//SetEventID 当前不允许修改 EventID 。
	//SetEventID(eventID string)

}

// event 自定义 event 统一数据模型。
type event struct {
	sent       *sync.Once
	EventID    string                 `json:"EventID"`
	EventType  string                 `json:"EventType"`
	Time       time.Time              `json:"Time"`
	Level      level                  `json:"Level"`
	Attributes map[string]interface{} `json:"Attributes"`
	Resource   *resource              `json:"Resource"`
	Subject    string                 `json:"Subject"`
	Link       *link                  `json:"Link,omitempty"`
	Data       interface{}            `json:"Data"`
}

// Info 设置 Info 级别的事件并立即发送。
func Info(data interface{}, opts ...EventStartOption) {
	opts = append(opts, withData(data))
	opts = append(opts, withLevel(INFO))
	NewEvent(opts...).Send()
}

// Warn 设置 Warn 级别的事件并立即发送。
func Warn(data interface{}, opts ...EventStartOption) {
	opts = append(opts, withData(data))
	opts = append(opts, withLevel(WARN))
	NewEvent(opts...).Send()
}

// Error 设置 Error 级别的事件并立即发送。
func Error(data interface{}, opts ...EventStartOption) {
	opts = append(opts, withData(data))
	opts = append(opts, withLevel(ERROR))
	NewEvent(opts...).Send()
}

// NewEvent 创建新的 event ，默认填充ID、时间、事件级别、资源信息、事件类型。
func NewEvent(opts ...EventStartOption) Event {
	cfg := newEventStartConfig(opts...)
	e := &event{
		sent:       &sync.Once{},
		EventID:    cfg.EventID,
		EventType:  cfg.EventType,
		Time:       cfg.Time,
		Level:      newLevel(cfg.Level.Self()),
		Attributes: make(map[string]interface{}),
		Resource:   newResource(),
		Subject:    cfg.Subject,
		Link:       newLink(cfg.SpanContext),
		Data:       cfg.Data,
	}
	e.SetAttributes(cfg.Attributes...)
	return e
}

// generateID 生成全球唯一ULID。
func generateID() string {
	return ulid.Make().String()
}

func (e *event) SetEventType(eventType string) {
	if strings.TrimSpace(eventType) == "" {
		log.Println(errors.New(custom_errors.EmptyEventType))
		return
	}
	e.EventType = eventType
}

func (e *event) SetTime(t time.Time) {
	if t.Equal(time.Time{}) {
		log.Println(errors.New(custom_errors.ZeroTime))
		return
	}
	e.Time = t
}

func (e *event) SetLevel(level Level) {
	e.Level = newLevel(level.Self())
}

func (e *event) SetAttributes(kvs ...Attribute) {
	for _, kv := range kvs {
		// 校验 attribute 是否合法，合法的才放进map去重。
		if !kv.Valid() {
			log.Println(custom_errors.EmptyKey)
			continue
		}
		e.Attributes[kv.GetKey()] = kv.GetValue()
	}
}

func (e *event) SetSubject(subject string) {
	e.Subject = subject
}

func (e *event) SetLink(spanContext trace.SpanContext) {
	if !spanContext.IsValid() {
		log.Println(errors.New(custom_errors.InvalidLink))
		return
	}
	e.Link = newLink(spanContext)
}

func (e *event) SetData(data interface{}) {
	e.Data = data
}

func (e *event) GetEventID() string {
	return e.EventID
}

func (e *event) GetEventType() string {
	return e.EventType
}

func (e *event) GetTime() time.Time {
	return e.Time
}

func (e *event) GetLevel() Level {
	return e.Level
}

func (e *event) GetAttributes() map[string]interface{} {
	return e.Attributes
}

func (e *event) GetResource() Resource {
	return e.Resource
}

func (e *event) GetSubject() string {
	return e.Subject
}

func (e *event) GetLink() Link {
	if e.Link == nil {
		return &link{
			TraceID: "",
			SpanID:  "",
		}

	}
	return e.Link
}

func (e *event) GetData() interface{} {
	return e.Data
}

func (e *event) GetEventMap() map[string]interface{} {
	result := make(map[string]interface{}, 9)
	result["EventID"] = e.EventID
	result["EventType"] = e.EventType
	result["Time"] = e.Time
	result["Level"] = e.Level
	result["Attributes"] = e.Attributes
	result["Resource"] = e.Resource
	result["Subject"] = e.Subject
	if e.Link != nil {
		result["Link"] = e.Link
	}
	result["Data"] = e.Data

	return result
}

func (e *event) Send() {
	e.sent.Do(func() {
		GetEventProvider().(*eventProvider).loadEvent(e)
	})
}

// SetServiceInfo 设置服务信息，包括服务名、版本号、实例ID。
func SetServiceInfo(ServiceName string, ServiceVersion string, ServiceInstance string) {
	if strings.TrimSpace(ServiceName) != "" {
		globalServiceName = ServiceName
	}
	if strings.TrimSpace(ServiceVersion) != "" {
		globalServiceVersion = ServiceVersion
	}
	if strings.TrimSpace(ServiceInstance) != "" {
		globalServiceInstance = ServiceInstance
	}
}

func (e *event) private() {
	// private 禁止用户自己实现接口。
}

// UnmarshalEvents 将JSON解析成[]Event。
func UnmarshalEvents(b []byte) ([]Event, error) {
	// 使用实际类型 []*event 来接收数据。
	events := make([]*event, 0)
	err := json.Unmarshal(b, &events)

	// 返回接口类型 Event 。
	result := make([]Event, 0, len(events))
	for _, e := range events {
		// 校验格式
		if e != nil && e.Valid() {
			result = append(result, e)
		}
	}
	// 如果返回空切片说明传入的JSON格式错误。
	if len(result) == 0 {
		err = errors.New(custom_errors.InvalidJSON)
	}
	return result, err
}

func (e *event) Valid() bool {
	return len(e.GetEventID()) == 26 &&
		strings.TrimSpace(e.GetEventType()) != "" &&
		e.GetTime().After(time.Time{}) &&
		e.GetLevel().Valid() &&
		e.GetResource().Valid()
}
