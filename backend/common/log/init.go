package log

import (
	"context"
	"os"
	"strings"
	"time"
)

// InitOptions 初始化选项（函数式选项模式）
type InitOptions struct {
	// 基础配置
	ServiceName string // 服务名称（必填）
	ServiceID   string // 服务ID（可选，默认从环境变量或ServiceName获取）
	InstanceID  string // 实例ID（可选，默认从环境变量获取）

	// 日志级别（可选，默认根据环境自动判断）
	Level Level

	// 环境（可选，用于自动判断日志级别和颜色）
	// 如果为空，从 APP_ENV 环境变量获取
	Environment string

	// 是否启用文件输出（可选，默认true）
	EnableFile bool

	// 文件路径模板（可选，默认 "logs/{serviceName}-%s.log"）
	FilePath string

	// 是否启用控制台颜色（可选，默认根据环境判断）
	EnableColor *bool

	// 高级配置（可选，用于覆盖默认配置）
	Config *Config
}

// InitOption 初始化选项函数类型
type InitOption func(*InitOptions)

// WithServiceID 设置服务ID
func WithServiceID(serviceID string) InitOption {
	return func(opts *InitOptions) {
		opts.ServiceID = serviceID
	}
}

// WithInstanceID 设置实例ID
func WithInstanceID(instanceID string) InitOption {
	return func(opts *InitOptions) {
		opts.InstanceID = instanceID
	}
}

// WithLevel 设置日志级别
func WithLevel(level Level) InitOption {
	return func(opts *InitOptions) {
		opts.Level = level
	}
}

// WithEnvironment 设置环境
func WithEnvironment(env string) InitOption {
	return func(opts *InitOptions) {
		opts.Environment = env
	}
}

// WithFileOutput 设置是否启用文件输出
func WithFileOutput(enable bool) InitOption {
	return func(opts *InitOptions) {
		opts.EnableFile = enable
	}
}

// WithFilePath 设置文件路径模板
func WithFilePath(path string) InitOption {
	return func(opts *InitOptions) {
		opts.FilePath = path
	}
}

// WithColor 设置是否启用颜色
func WithColor(enable bool) InitOption {
	return func(opts *InitOptions) {
		opts.EnableColor = &enable
	}
}

// WithConfig 设置高级配置（会覆盖其他选项）
func WithConfig(config Config) InitOption {
	return func(opts *InitOptions) {
		opts.Config = &config
	}
}

// InitLogger 便捷初始化日志系统
// serviceName: 服务名称（必填）
// options: 可选配置项
//
// 示例：
//
//	// 最简单的方式
//	logger, err := log.InitLogger("gateway")
//
//	// 带服务ID
//	logger, err := log.InitLogger("gateway", log.WithServiceID("gateway-001"))
//
//	// 完整配置
//	logger, err := log.InitLogger("gateway",
//	    log.WithServiceID("gateway-001"),
//	    log.WithLevel(log.DebugLevel),
//	    log.WithEnvironment("development"),
//	)
func InitLogger(serviceName string, options ...InitOption) (Logger, error) {
	if serviceName == "" {
		serviceName = os.Getenv("SERVICE_NAME")
		if serviceName == "" {
			serviceName = "app"
		}
	}

	// 解析选项
	opts := &InitOptions{
		ServiceName: serviceName,
		EnableFile:  true,
		// FilePath 默认为空，会在后面根据 ServiceName 自动生成
	}

	for _, opt := range options {
		opt(opts)
	}

	// 如果提供了完整配置，直接使用
	if opts.Config != nil {
		logger, err := NewLogger(*opts.Config)
		if err != nil {
			return nil, err
		}
		SetGlobalLogger(logger)
		return logger, nil
	}

	// 构建配置
	config := DefaultConfig()
	config.ServiceName = opts.ServiceName

	// 设置服务ID
	if opts.ServiceID != "" {
		config.ServiceID = opts.ServiceID
	} else if serviceID := os.Getenv("SERVICE_ID"); serviceID != "" {
		config.ServiceID = serviceID
	} else {
		config.ServiceID = opts.ServiceName // 默认使用服务名
	}

	// 设置实例ID（优先从环境变量获取）
	if opts.InstanceID != "" {
		config.InstanceID = opts.InstanceID
	} else if instanceID := os.Getenv("POD_NAME"); instanceID != "" {
		config.InstanceID = instanceID
	} else if instanceID := os.Getenv("HOSTNAME"); instanceID != "" {
		config.InstanceID = instanceID
	} else if instanceID := os.Getenv("INSTANCE_ID"); instanceID != "" {
		config.InstanceID = instanceID
	} else {
		config.InstanceID = config.ServiceID
	}

	// 确定环境
	env := opts.Environment
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		env = "development" // 默认开发环境
	}

	// 设置日志级别
	if opts.Level != 0 {
		config.Level = opts.Level
	} else {
		// 根据环境自动判断
		envLower := strings.ToLower(env)
		if envLower == "production" || envLower == "prod" {
			config.Level = InfoLevel
		} else {
			config.Level = DebugLevel
		}
	}

	// 配置控制台输出
	config.Console.Enabled = true
	config.Console.Format = TextFormat
	config.Console.Level = DebugLevel // 控制台显示所有级别

	// 设置颜色（根据环境自动判断）
	if opts.EnableColor != nil {
		config.EnableColor = *opts.EnableColor
	} else {
		envLower := strings.ToLower(env)
		config.EnableColor = envLower != "production" && envLower != "prod"
	}

	// 配置文件输出
	if opts.EnableFile {
		config.File.Enabled = true
		config.File.Format = JSONFormat
		config.File.Level = InfoLevel // 文件只记录INFO及以上
		if opts.FilePath != "" {
			config.File.Path = opts.FilePath
		} else {
			// 使用默认路径模板：logs/{serviceName}-%s.log
			// 注意：%s 会被替换为日期，服务名需要手动替换
			config.File.Path = "logs/" + opts.ServiceName + "-%s.log"
		}
		config.File.MaxSize = 100 // MB
		config.File.MaxAge = 7    // days
		config.File.MaxBackups = 10
		config.File.Compress = true
	}

	// 创建Logger
	logger, err := NewLogger(config)
	if err != nil {
		return nil, err
	}

	// 设置为全局Logger
	SetGlobalLogger(logger)

	return logger, nil
}

// InitLoggerAndSetGlobal 初始化日志系统并设置为全局Logger（更简洁的API）
// 这是 InitLogger 的便捷包装，自动设置全局Logger
//
// 示例：
//
//	// 最简单的方式
//	if err := log.InitLoggerAndSetGlobal("gateway"); err != nil {
//	    panic(err)
//	}
//
//	// 带选项
//	if err := log.InitLoggerAndSetGlobal("gateway",
//	    log.WithServiceID("gateway-001"),
//	    log.WithLevel(log.DebugLevel),
//	); err != nil {
//	    panic(err)
//	}
func InitLoggerAndSetGlobal(serviceName string, options ...InitOption) error {
	_, err := InitLogger(serviceName, options...)
	return err
}

// MustInitLogger 初始化日志系统（失败时panic）
// 适用于main函数中的初始化
//
// 示例：
//
//	func main() {
//	    defer log.MustInitLogger("gateway").Shutdown(context.Background())
//	    // ... 其他代码
//	}
func MustInitLogger(serviceName string, options ...InitOption) Logger {
	logger, err := InitLogger(serviceName, options...)
	if err != nil {
		panic("初始化日志系统失败: " + err.Error())
	}
	return logger
}

// InitLoggerWithShutdown 初始化日志系统并返回带优雅关闭的defer函数
// 返回的defer函数应该在main函数中使用 defer 调用
// 如果初始化失败，会panic（符合启动阶段的错误处理）
//
// 示例：
//
//	func main() {
//	    defer log.InitLoggerWithShutdown("gateway")()
//	    // ... 其他代码
//	}
//
//	// 带选项
//	func main() {
//	    defer log.InitLoggerWithShutdown("gateway",
//	        log.WithServiceID("gateway-001"),
//	        log.WithEnvironment("production"),
//	    )()
//	    // ... 其他代码
//	}
func InitLoggerWithShutdown(serviceName string, options ...InitOption) func() {
	logger, err := InitLogger(serviceName, options...)
	if err != nil {
		// 日志系统初始化失败，panic（启动阶段的关键错误）
		panic("初始化日志系统失败: " + err.Error())
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := logger.Shutdown(ctx); err != nil {
			// 关闭失败，输出到stderr（符合规则例外）
			os.Stderr.WriteString("警告: 日志系统关闭失败: " + err.Error() + "\n")
		}
	}
}
