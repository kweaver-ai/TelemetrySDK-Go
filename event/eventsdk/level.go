package eventsdk

// Level 定义了 level 的3种级别。
type Level interface {
	// Self 返回事件级别。
	Self() string
	// Valid 校验是否合法。
	Valid() bool
	// private 禁止用户自己实现接口。
	private()
	// ERROR 事件级别：错误。
	// WARN 事件级别：警告。
	// INFO 事件级别：一般。
}

// level 实际为 string 类型。
type level string

// newLevel 创建新的 level 。
func newLevel(l string) level {
	return level(l)
}

func (l level) Self() string {
	return string(l)
}

func (l level) Valid() bool {
	switch l.Self() {
	case "ERROR":
		return true
	case "WARN":
		return true
	case "INFO":
		return true
	default:
		return false
	}
}

func (l level) private() {
	// private 禁止用户自己实现接口。
}

const (
	ERROR level = "ERROR"
	WARN  level = "WARN"
	INFO  level = "INFO"
)
