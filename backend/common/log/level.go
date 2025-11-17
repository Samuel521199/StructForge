package log

// Level 日志级别
type Level int

const (
	// DebugLevel 调试级别，最详细的日志
	DebugLevel Level = iota
	// InfoLevel 信息级别，一般信息
	InfoLevel
	// WarnLevel 警告级别，需要注意但不影响运行
	WarnLevel
	// ErrorLevel 错误级别，错误信息
	ErrorLevel
	// FatalLevel 致命级别，致命错误，程序会退出
	FatalLevel
)

// String 返回级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel 从字符串解析日志级别
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return DebugLevel
	case "info", "INFO":
		return InfoLevel
	case "warn", "WARN", "warning", "WARNING":
		return WarnLevel
	case "error", "ERROR":
		return ErrorLevel
	case "fatal", "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}
