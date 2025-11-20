package config

import "time"

// AdapterType 数据库适配器类型
type AdapterType string

const (
	// AdapterPostgreSQL PostgreSQL 适配器
	AdapterPostgreSQL AdapterType = "postgres"
	// AdapterMySQL MySQL 适配器
	AdapterMySQL AdapterType = "mysql"
	// AdapterSQLite SQLite 适配器
	AdapterSQLite AdapterType = "sqlite"
)

// Config 数据库系统配置
type Config struct {
	// 适配器类型
	AdapterType AdapterType

	// 通用配置
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生存时间
	ConnMaxIdleTime time.Duration // 连接最大空闲时间

	// 日志配置
	LogLevel      string        // 日志级别：silent, error, warn, info
	SlowThreshold time.Duration // 慢查询阈值

	// 适配器特定配置
	PostgreSQL *PostgreSQLConfig
	MySQL      *MySQLConfig
	SQLite     *SQLiteConfig
}

// PostgreSQLConfig PostgreSQL 配置
type PostgreSQLConfig struct {
	// DSN 数据源名称，格式：host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	DSN string

	// 或者使用单独的配置项
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string // disable, require, verify-ca, verify-full
	TimeZone string // Asia/Shanghai

	// 连接池配置（如果使用 DSN，这些会被 DSN 中的参数覆盖）
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	// DSN 数据源名称，格式：user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	DSN string

	// 或者使用单独的配置项
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	Charset   string // utf8mb4
	ParseTime bool   // true
	Loc       string // Local

	// 连接池配置
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// SQLiteConfig SQLite 配置
type SQLiteConfig struct {
	// DSN 数据源名称，通常是文件路径，如：test.db
	DSN string

	// 或者使用文件路径
	Path string

	// 其他选项
	ForeignKeys bool          // 启用外键约束
	BusyTimeout time.Duration // 忙等待超时
}
