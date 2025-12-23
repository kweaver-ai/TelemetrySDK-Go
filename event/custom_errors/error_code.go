package custom_errors

// custom_errors 定义错误码和错误描述。
const (
	ModuleName = "TelemetrySDK-Event(Go).Error: "

	EmptyKey        = ModuleName + "Attribute设置了无意义的空键或与默认值冲突"
	InvalidJSON     = ModuleName + "传入了非法的JSON，应该传入[]Event类型"
	AlreadyShutdown = ModuleName + "已经关闭了Event Exporter，不能再发送"
	EmptyEventType  = ModuleName + "设置了无意义的空EventType"
	ZeroTime        = ModuleName + "设置了无意义的Event时间"
	InvalidLink     = ModuleName + "设置了无效的Trace关联记录"
)
