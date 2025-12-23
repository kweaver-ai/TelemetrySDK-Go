package custom_errors

// custom_errors 定义TelemetrySDK-Exporter错误码和错误描述。
const (
	ModuleName = "TelemetrySDK-Exporter(Go).Error: "

	JobIdNotFound          = ModuleName + "接收器上报地址不正确"
	PayloadTooLarge        = ModuleName + "数据太大超过了5MB限制"
	InvalidFormat          = ModuleName + "格式校验不通过"
	InvalidURL             = ModuleName + "URL非法，请检查"
	InvalidCompression     = ModuleName + "压缩方式不存在"
	RetryTooLong           = ModuleName + "重发持续时间太长"
	SentFailed             = ModuleName + "发送失败，检查日志"
	ExceedRetryElapsedTime = ModuleName + "超过最大重发时间限制"
	DurationTooLong        = ModuleName + "超过最长连接时间限制"
	RetryFailure           = ModuleName + "数据正在重发"
)
