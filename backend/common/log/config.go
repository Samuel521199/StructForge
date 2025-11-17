package log

import (
	"os"
	"time"
)

// Output 输出目标类型
type Output string

const (
	ConsoleOutput Output = "console" // 控制台输出
	FileOutput    Output = "file"    // 文件输出
)

// Format 输出格式
type Format string

const (
	TextFormat Format = "text" // 文本格式
	JSONFormat Format = "json" // JSON格式
)

// Config 日志系统配置
type Config struct {
	// 基础配置
	Level       Level  // 全局日志级别
	ServiceName string // 服务名称
	ServiceID   string // 服务ID
	InstanceID  string // 实例ID

	// 输出配置
	Outputs     []Output // 输出目标列表
	EnableColor bool     // 是否启用颜色（控制台）

	// 控制台配置
	Console ConsoleConfig

	// 文件配置
	File FileConfig

	// 采样配置
	Sampling SamplingConfig

	// 异步写入配置
	Async AsyncConfig

	// 脱敏配置
	Mask MaskConfig

	// 钩子配置
	Hooks []Hook // 日志钩子列表

	// 错误处理配置
	ErrorHandler ErrorHandler // 错误处理器（可选）
}

// ConsoleConfig 控制台配置
type ConsoleConfig struct {
	Enabled bool   // 是否启用
	Format  Format // 输出格式（text/json）
	Level   Level  // 控制台日志级别（可独立设置）
}

// FileConfig 文件配置
type FileConfig struct {
	Enabled         bool   // 是否启用
	Format          Format // 输出格式（通常为json）
	Level           Level  // 文件日志级别
	Path            string // 文件路径模板，支持 %s 占位符（日期）
	MaxSize         int    // 单个文件最大大小（MB）
	MaxAge          int    // 保留天数
	MaxBackups      int    // 最大备份文件数
	Compress        bool   // 是否压缩旧文件
	SeparateByLevel bool   // 是否按级别分离文件
	AsyncEnabled    bool   // 是否启用异步写入
}

// SamplingConfig 采样配置
type SamplingConfig struct {
	Enabled bool    // 是否启用采样
	Ratio   float64 // 采样比例，0.1表示10%，1.0表示100%
	Levels  []Level // 需要采样的级别（为空则采样所有级别）
}

// AsyncConfig 异步写入配置
type AsyncConfig struct {
	Enabled       bool          // 是否启用异步写入
	QueueSize     int           // 队列大小
	BatchSize     int           // 批量写入大小
	FlushInterval time.Duration // 刷新间隔
	DropOnFull    bool          // 队列满时是否丢弃（false则阻塞）
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "app"
	}

	return Config{
		Level:       InfoLevel,
		ServiceName: serviceName,
		ServiceID:   os.Getenv("SERVICE_ID"),
		InstanceID:  os.Getenv("INSTANCE_ID"),
		Outputs:     []Output{ConsoleOutput},
		EnableColor: true,
		Console: ConsoleConfig{
			Enabled: true,
			Format:  TextFormat,
			Level:   DebugLevel,
		},
		File: FileConfig{
			Enabled:         false,
			Format:          JSONFormat,
			Level:           InfoLevel,
			Path:            "logs/%s-%s.log", // 服务名-日期.log
			MaxSize:         100,              // MB
			MaxAge:          7,                // days
			MaxBackups:      10,
			Compress:        true,
			SeparateByLevel: false,
			AsyncEnabled:    false,
		},
		Sampling: SamplingConfig{
			Enabled: false,
			Ratio:   1.0,                 // 默认不采样（100%）
			Levels:  []Level{DebugLevel}, // 默认只对Debug级别采样
		},
		Async: AsyncConfig{
			Enabled:       false,
			QueueSize:     1000,
			BatchSize:     100,
			FlushInterval: 5 * time.Second,
			DropOnFull:    false, // 队列满时阻塞
		},
		Mask: MaskConfig{
			Enabled:  true, // 默认启用脱敏
			Fields:   []string{"password", "pwd", "token", "secret", "key"},
			KeepHead: 3,   // 保留前3个字符
			KeepTail: 0,   // 不保留尾部
			MaskChar: "*", // 使用*脱敏
		},
	}
}
