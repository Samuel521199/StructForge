package database

import (
	"context"
	"fmt"
	"time"

	dbconfig "StructForge/backend/common/data/database/config"
	"StructForge/backend/common/log"
)

// InitOption 初始化选项
type InitOption func(*Config)

// WithAdapterType 设置适配器类型
func WithAdapterType(adapterType dbconfig.AdapterType) InitOption {
	return func(c *dbconfig.Config) {
		c.AdapterType = adapterType
	}
}

// WithPostgreSQL 设置 PostgreSQL 配置
func WithPostgreSQL(cfg *dbconfig.PostgreSQLConfig) InitOption {
	return func(c *dbconfig.Config) {
		c.PostgreSQL = cfg
	}
}

// WithMySQL 设置 MySQL 配置
func WithMySQL(cfg *dbconfig.MySQLConfig) InitOption {
	return func(c *dbconfig.Config) {
		c.MySQL = cfg
	}
}

// WithSQLite 设置 SQLite 配置
func WithSQLite(cfg *dbconfig.SQLiteConfig) InitOption {
	return func(c *dbconfig.Config) {
		c.SQLite = cfg
	}
}

// WithMaxOpenConns 设置最大打开连接数
func WithMaxOpenConns(maxOpenConns int) InitOption {
	return func(c *dbconfig.Config) {
		c.MaxOpenConns = maxOpenConns
	}
}

// WithMaxIdleConns 设置最大空闲连接数
func WithMaxIdleConns(maxIdleConns int) InitOption {
	return func(c *dbconfig.Config) {
		c.MaxIdleConns = maxIdleConns
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(level string) InitOption {
	return func(c *dbconfig.Config) {
		c.LogLevel = level
	}
}

// WithSlowThreshold 设置慢查询阈值
func WithSlowThreshold(threshold time.Duration) InitOption {
	return func(c *dbconfig.Config) {
		c.SlowThreshold = threshold
	}
}

// InitDatabase 初始化数据库系统
func InitDatabase(adapterType string, options ...InitOption) (Database, error) {
	ctx := context.Background()

	// 创建默认配置
	config := DefaultConfig()

	// 应用选项
	for _, option := range options {
		option(config)
	}

	// 设置适配器类型
	if adapterType != "" {
		config.AdapterType = dbconfig.AdapterType(adapterType)
	}

	// 创建适配器
	db, err := NewDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 设置为全局数据库
	SetGlobalDB(db)

	// 测试连接
	if err := db.Ping(ctx); err != nil {
		log.Warn(ctx, "数据库连接测试失败（不影响启动）",
			log.ErrorField(err),
		)
	} else {
		log.Info(ctx, "数据库系统初始化成功",
			log.String("adapter", string(config.AdapterType)),
		)
	}

	return db, nil
}

// InitDatabaseFromFile 从配置文件初始化数据库系统
// configPath: 配置文件路径（YAML 格式）
// 配置文件格式示例：
//
//	database:
//	  adapter_type: postgres
//	  postgres:
//	    host: localhost
//	    port: 5432
//	    user: postgres
//	    password: password
//	    dbname: structforge
func InitDatabaseFromFile(configPath string, options ...InitOption) (Database, error) {
	ctx := context.Background()

	// 从文件加载配置
	config, err := LoadFromFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load database config from file: %w", err)
	}

	// 应用额外选项（选项优先级高于配置文件）
	for _, option := range options {
		option(config)
	}

	// 创建适配器
	db, err := NewDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 设置为全局数据库
	SetGlobalDB(db)

	// 测试连接
	if err := db.Ping(ctx); err != nil {
		log.Warn(ctx, "数据库连接测试失败（不影响启动）",
			log.ErrorField(err),
		)
	} else {
		log.Info(ctx, "数据库系统初始化成功（从配置文件加载）",
			log.String("adapter", string(config.AdapterType)),
			log.String("config_path", configPath),
		)
	}

	return db, nil
}

// InitDatabaseFromYAML 从 YAML 数据初始化数据库系统
func InitDatabaseFromYAML(yamlData []byte, options ...InitOption) (Database, error) {
	ctx := context.Background()

	// 从 YAML 数据加载配置
	config, err := LoadFromYAML(yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to load database config from YAML: %w", err)
	}

	// 应用额外选项（选项优先级高于配置文件）
	for _, option := range options {
		option(config)
	}

	// 创建适配器
	db, err := NewDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 设置为全局数据库
	SetGlobalDB(db)

	// 测试连接
	if err := db.Ping(ctx); err != nil {
		log.Warn(ctx, "数据库连接测试失败（不影响启动）",
			log.ErrorField(err),
		)
	} else {
		log.Info(ctx, "数据库系统初始化成功（从 YAML 加载）",
			log.String("adapter", string(config.AdapterType)),
		)
	}

	return db, nil
}

// InitDatabaseWithShutdown 初始化数据库系统并返回关闭函数
// 用法: defer database.InitDatabaseWithShutdown("postgres")()
func InitDatabaseWithShutdown(adapterType string, options ...InitOption) func() {
	db, err := InitDatabase(adapterType, options...)
	if err != nil {
		// 如果初始化失败，返回空函数
		// 实际使用时应该处理错误
		return func() {}
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = db.Close()
		log.Info(ctx, "数据库连接已关闭")
	}
}

// InitDatabaseFromFileWithShutdown 从配置文件初始化数据库系统并返回关闭函数
// 用法: defer database.InitDatabaseFromFileWithShutdown("config.yaml")()
func InitDatabaseFromFileWithShutdown(configPath string, options ...InitOption) func() {
	db, err := InitDatabaseFromFile(configPath, options...)
	if err != nil {
		// 如果初始化失败，返回空函数
		// 实际使用时应该处理错误
		return func() {}
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = db.Close()
		log.Info(ctx, "数据库连接已关闭")
	}
}

// NewDatabase 函数已移至 factory.go 以避免导入循环

// Shutdown 优雅关闭数据库系统
func Shutdown(ctx context.Context) error {
	if globalDB == nil {
		return nil
	}
	return globalDB.Close()
}
