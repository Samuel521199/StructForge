package database

import (
	"fmt"
	"os"
	"time"

	dbconfig "StructForge/backend/common/data/database/config"
)

// 类型别名，保持向后兼容
type (
	AdapterType      = dbconfig.AdapterType
	Config           = dbconfig.Config
	PostgreSQLConfig = dbconfig.PostgreSQLConfig
	MySQLConfig      = dbconfig.MySQLConfig
	SQLiteConfig     = dbconfig.SQLiteConfig
)

const (
	AdapterPostgreSQL = dbconfig.AdapterPostgreSQL
	AdapterMySQL      = dbconfig.AdapterMySQL
	AdapterSQLite     = dbconfig.AdapterSQLite
)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		AdapterType:     AdapterPostgreSQL,
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
		LogLevel:        "warn",
		SlowThreshold:   200 * time.Millisecond,
		PostgreSQL: &PostgreSQLConfig{
			Host:            getEnvOrDefault("DB_HOST", "localhost"),
			Port:            getEnvIntOrDefault("DB_PORT", 5432),
			User:            getEnvOrDefault("DB_USER", "postgres"),
			Password:        getEnvOrDefault("DB_PASSWORD", ""),
			DBName:          getEnvOrDefault("DB_NAME", "structforge"),
			SSLMode:         getEnvOrDefault("DB_SSLMODE", "disable"),
			TimeZone:        getEnvOrDefault("DB_TIMEZONE", "Asia/Shanghai"),
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 10 * time.Minute,
		},
		MySQL: &dbconfig.MySQLConfig{
			Host:            getEnvOrDefault("DB_HOST", "localhost"),
			Port:            getEnvIntOrDefault("DB_PORT", 3306),
			User:            getEnvOrDefault("DB_USER", "root"),
			Password:        getEnvOrDefault("DB_PASSWORD", ""),
			DBName:          getEnvOrDefault("DB_NAME", "structforge"),
			Charset:         "utf8mb4",
			ParseTime:       true,
			Loc:             "Local",
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 10 * time.Minute,
		},
		SQLite: &dbconfig.SQLiteConfig{
			Path:        getEnvOrDefault("DB_PATH", "data.db"),
			ForeignKeys: true,
			BusyTimeout: 5 * time.Second,
		},
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault 获取环境变量并转换为整数，如果不存在则返回默认值
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}
