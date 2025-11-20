package database

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

var (
	// ErrNotSupported 不支持的数据库类型
	ErrNotSupported = errors.New("database adapter not supported")
	// ErrNotInitialized 数据库未初始化
	ErrNotInitialized = errors.New("database not initialized")
	// ErrConnectionFailed 数据库连接失败
	ErrConnectionFailed = errors.New("database connection failed")
)

// Database 数据库抽象接口
// 所有数据库适配器必须实现此接口
type Database interface {
	// GetDB 获取 GORM 数据库实例
	GetDB() *gorm.DB

	// Ping 检查数据库连接
	Ping(ctx context.Context) error

	// Close 关闭数据库连接
	Close() error

	// Health 健康检查
	Health(ctx context.Context) error

	// Migrate 执行数据库迁移
	Migrate(ctx context.Context, models ...interface{}) error

	// AutoMigrate 自动迁移（根据模型创建/更新表结构）
	AutoMigrate(ctx context.Context, models ...interface{}) error

	// Transaction 执行事务
	Transaction(ctx context.Context, fn func(*gorm.DB) error) error
}

// globalDB 全局数据库实例
var globalDB Database

// SetGlobalDB 设置全局数据库实例
func SetGlobalDB(db Database) {
	globalDB = db
}

// GetGlobalDB 获取全局数据库实例
func GetGlobalDB() Database {
	return globalDB
}

// GetDB 获取全局 GORM 数据库实例（便捷方法）
func GetDB() *gorm.DB {
	if globalDB == nil {
		return nil
	}
	return globalDB.GetDB()
}

