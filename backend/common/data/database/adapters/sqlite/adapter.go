package sqlite

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	dbconfig "StructForge/backend/common/data/database/config"
	"StructForge/backend/common/log"
)

// Adapter SQLite 适配器
type Adapter struct {
	db     *gorm.DB
	config *dbconfig.SQLiteConfig
}

// NewAdapter 创建 SQLite 适配器
func NewAdapter(config *dbconfig.Config) (*Adapter, error) {
	sqliteConfig := config.SQLite
	if sqliteConfig == nil {
		return nil, fmt.Errorf("SQLite config is required")
	}

	// 构建 DSN（文件路径）
	dsn := sqliteConfig.DSN
	if dsn == "" {
		dsn = sqliteConfig.Path
	}
	if dsn == "" {
		dsn = "data.db"
	}

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: newGormLogger(config.LogLevel, config.SlowThreshold),
	}

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	// SQLite 特定配置
	if sqliteConfig.ForeignKeys {
		// 启用外键约束
		if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
			return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
		}
	}

	if sqliteConfig.BusyTimeout > 0 {
		// 设置忙等待超时
		timeoutMs := int(sqliteConfig.BusyTimeout.Milliseconds())
		if err := db.Exec(fmt.Sprintf("PRAGMA busy_timeout = %d", timeoutMs)).Error; err != nil {
			return nil, fmt.Errorf("failed to set busy timeout: %w", err)
		}
	}

	adapter := &Adapter{
		db:     db,
		config: sqliteConfig,
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := adapter.Ping(ctx); err != nil {
		return nil, fmt.Errorf("SQLite connection test failed: %w", err)
	}

	return adapter, nil
}

// GetDB 获取 GORM 数据库实例
func (a *Adapter) GetDB() *gorm.DB {
	return a.db
}

// Ping 检查数据库连接
func (a *Adapter) Ping(ctx context.Context) error {
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Close 关闭数据库连接
func (a *Adapter) Close() error {
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health 健康检查
func (a *Adapter) Health(ctx context.Context) error {
	return a.Ping(ctx)
}

// Migrate 执行数据库迁移
func (a *Adapter) Migrate(ctx context.Context, models ...interface{}) error {
	return a.db.WithContext(ctx).AutoMigrate(models...)
}

// AutoMigrate 自动迁移（根据模型创建/更新表结构）
func (a *Adapter) AutoMigrate(ctx context.Context, models ...interface{}) error {
	return a.db.WithContext(ctx).AutoMigrate(models...)
}

// Transaction 执行事务
func (a *Adapter) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return a.db.WithContext(ctx).Transaction(fn)
}

// newGormLogger 创建 GORM 日志记录器
func newGormLogger(level string, slowThreshold time.Duration) logger.Interface {
	var logLevel logger.LogLevel
	switch level {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}

	return logger.New(
		&gormLogWriter{},
		logger.Config{
			SlowThreshold:             slowThreshold,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// gormLogWriter GORM 日志写入器（使用统一日志系统）
type gormLogWriter struct{}

func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	ctx := context.Background()
	msg := fmt.Sprintf(format, args...)
	log.Info(ctx, "[GORM] "+msg)
}
