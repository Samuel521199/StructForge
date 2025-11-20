//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2"
	kratosLog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"StructForge/backend/apps/user/internal/biz"
	"StructForge/backend/apps/user/internal/conf"
	"StructForge/backend/apps/user/internal/data"
	"StructForge/backend/apps/user/internal/server"
	"StructForge/backend/apps/user/internal/service"
	"StructForge/backend/common/data/database"
	"StructForge/backend/common/email"
	"StructForge/backend/common/log"
)

// wireApp 初始化应用（Wire 会自动生成 wire_gen.go）
func wireApp(bc *conf.Bootstrap) (*kratos.App, func(), error) {
	panic(wire.Build(
		// 数据访问层
		data.ProviderSet,
		// JWT 配置提供者
		jwtSecretKeyProvider,
		jwtTokenDurationProvider,
		// 邮件服务
		emailProvider,
		// 业务逻辑层（包含 JWT Manager）
		biz.ProviderSet,
		// 服务层
		service.ProviderSet,
		// 服务器
		server.ProviderSet,
		// 数据库
		databaseProvider,
		// 日志
		logProvider,
		// 应用
		newApp,
	))
}

// emailProvider 提供邮件服务实例
func emailProvider() email.EmailService {
	// TODO: 从配置文件读取邮件配置
	config := email.DefaultConfig()
	// 可以从环境变量或配置文件读取
	// config.SMTPHost = os.Getenv("SMTP_HOST")
	// config.SMTPPort = ...
	// config.SMTPUser = os.Getenv("SMTP_USER")
	// config.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	// config.FromEmail = os.Getenv("FROM_EMAIL")
	// config.FromName = os.Getenv("FROM_NAME")
	return email.NewEmailService(config)
}

// jwtSecretKeyProvider 提供 JWT 密钥
func jwtSecretKeyProvider() string {
	// TODO: 从配置文件读取
	return "your-secret-key-change-in-production"
}

// jwtTokenDurationProvider 提供 JWT Token 有效期
func jwtTokenDurationProvider() time.Duration {
	// TODO: 从配置文件读取
	return 24 * time.Hour
}

// logProvider 提供日志实例（返回 Kratos 兼容的日志接口）
func logProvider() kratosLog.Logger {
	// 使用全局日志实例（通过包级别的函数）
	// 注意：如果全局日志未初始化，创建一个默认的
	config := log.DefaultConfig()
	config.ServiceName = "user"
	logger, err := log.NewLogger(config)
	if err != nil {
		// 如果创建失败，返回一个 no-op logger
		return kratosLog.NewStdLogger(os.Stderr)
	}
	// 设置为全局日志
	log.SetGlobalLogger(logger)
	// 返回适配器
	return newKratosLogAdapter(logger)
}

// kratosLogAdapter 适配器，将我们的日志接口适配到 Kratos 日志接口
type kratosLogAdapter struct {
	logger log.Logger
}

// newKratosLogAdapter 创建 Kratos 日志适配器
func newKratosLogAdapter(logger log.Logger) kratosLog.Logger {
	return &kratosLogAdapter{logger: logger}
}

// Log 实现 Kratos Logger 接口
func (a *kratosLogAdapter) Log(level kratosLog.Level, keyvals ...interface{}) error {
	// 将 Kratos 日志级别转换为我们的日志级别
	var logLevel log.Level
	switch level {
	case kratosLog.LevelDebug:
		logLevel = log.DebugLevel
	case kratosLog.LevelInfo:
		logLevel = log.InfoLevel
	case kratosLog.LevelWarn:
		logLevel = log.WarnLevel
	case kratosLog.LevelError:
		logLevel = log.ErrorLevel
	case kratosLog.LevelFatal:
		logLevel = log.FatalLevel
	default:
		logLevel = log.InfoLevel
	}

	// 解析 keyvals（Kratos 使用 key-value 对）
	msg := ""
	fields := make([]log.Field, 0)
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key := keyvals[i]
			val := keyvals[i+1]
			if keyStr, ok := key.(string); ok {
				if keyStr == "msg" {
					if msgStr, ok := val.(string); ok {
						msg = msgStr
					}
				} else {
					// 转换为我们的 Field
					switch v := val.(type) {
					case string:
						fields = append(fields, log.String(keyStr, v))
					case int:
						fields = append(fields, log.Int(keyStr, v))
					case int64:
						fields = append(fields, log.Int64(keyStr, v))
					case bool:
						fields = append(fields, log.Bool(keyStr, v))
					case error:
						fields = append(fields, log.ErrorField(v))
					default:
						fields = append(fields, log.Any(keyStr, v))
					}
				}
			}
		}
	}

	// 如果没有消息，使用默认消息
	if msg == "" {
		msg = "log message"
	}

	// 使用 context.Background()（Kratos 的 Log 方法不提供 context）
	ctx := context.Background()
	a.logger.Log(ctx, logLevel, msg, fields...)
	return nil
}

// databaseProvider 提供数据库实例
func databaseProvider() (database.Database, func(), error) {
	db := database.GetGlobalDB()
	if db == nil {
		// 如果全局数据库未初始化，返回错误
		// 数据库应该在 main.go 中已经初始化
		return nil, nil, fmt.Errorf("数据库未初始化，请在 main.go 中先调用 database.InitDatabaseFromFile")
	}
	return db, func() {}, nil
}

// newApp 创建应用实例
func newApp(
	logger kratosLog.Logger,
	grpcServer *server.GRPCServer,
	httpServer *server.HTTPServer,
) *kratos.App {
	return kratos.New(
		kratos.Name("user"),
		kratos.Version("v1.0.0"),
		kratos.Logger(logger),
		kratos.Server(
			grpcServer,
			httpServer,
		),
	)
}

// 注意：需要运行以下命令生成 wire_gen.go：
// cd backend/apps/user/cmd/user && wire
