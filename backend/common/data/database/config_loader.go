package database

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	dbconfig "StructForge/backend/common/data/database/config"
)

// ConfigLoader 配置加载器
type ConfigLoader struct{}

// LoadFromFile 从 YAML 文件加载数据库配置
func LoadFromFile(configPath string) (*dbconfig.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return LoadFromYAML(data)
}

// LoadFromYAML 从 YAML 数据加载数据库配置
func LoadFromYAML(data []byte) (*dbconfig.Config, error) {
	var yamlConfig struct {
		Database struct {
			// 适配器类型
			AdapterType string `yaml:"adapter_type"`

			// 通用配置
			MaxOpenConns    int    `yaml:"max_open_conns"`
			MaxIdleConns    int    `yaml:"max_idle_conns"`
			ConnMaxLifetime string `yaml:"conn_max_lifetime"`  // 如: "1h", "30m"
			ConnMaxIdleTime string `yaml:"conn_max_idle_time"` // 如: "10m"
			LogLevel        string `yaml:"log_level"`
			SlowThreshold   string `yaml:"slow_threshold"` // 如: "200ms"

			// PostgreSQL 配置
			PostgreSQL struct {
				DSN      string `yaml:"dsn"`
				Host     string `yaml:"host"`
				Port     int    `yaml:"port"`
				User     string `yaml:"user"`
				Password string `yaml:"password"`
				DBName   string `yaml:"dbname"`
				SSLMode  string `yaml:"sslmode"`
				TimeZone string `yaml:"timezone"`
			} `yaml:"postgres"`

			// MySQL 配置
			MySQL struct {
				DSN       string `yaml:"dsn"`
				Host      string `yaml:"host"`
				Port      int    `yaml:"port"`
				User      string `yaml:"user"`
				Password  string `yaml:"password"`
				DBName    string `yaml:"dbname"`
				Charset   string `yaml:"charset"`
				ParseTime bool   `yaml:"parse_time"`
				Loc       string `yaml:"loc"`
			} `yaml:"mysql"`

			// SQLite 配置
			SQLite struct {
				DSN         string `yaml:"dsn"`
				Path        string `yaml:"path"`
				ForeignKeys bool   `yaml:"foreign_keys"`
				BusyTimeout string `yaml:"busy_timeout"` // 如: "5s"
			} `yaml:"sqlite"`
		} `yaml:"database"`
	}

	if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// 创建配置对象
	config := DefaultConfig()

	// 设置适配器类型
	if yamlConfig.Database.AdapterType != "" {
		config.AdapterType = dbconfig.AdapterType(yamlConfig.Database.AdapterType)
	}

	// 解析通用配置
	if yamlConfig.Database.MaxOpenConns > 0 {
		config.MaxOpenConns = yamlConfig.Database.MaxOpenConns
	}
	if yamlConfig.Database.MaxIdleConns > 0 {
		config.MaxIdleConns = yamlConfig.Database.MaxIdleConns
	}
	if yamlConfig.Database.ConnMaxLifetime != "" {
		if duration, err := time.ParseDuration(yamlConfig.Database.ConnMaxLifetime); err == nil {
			config.ConnMaxLifetime = duration
		}
	}
	if yamlConfig.Database.ConnMaxIdleTime != "" {
		if duration, err := time.ParseDuration(yamlConfig.Database.ConnMaxIdleTime); err == nil {
			config.ConnMaxIdleTime = duration
		}
	}
	if yamlConfig.Database.LogLevel != "" {
		config.LogLevel = yamlConfig.Database.LogLevel
	}
	if yamlConfig.Database.SlowThreshold != "" {
		if duration, err := time.ParseDuration(yamlConfig.Database.SlowThreshold); err == nil {
			config.SlowThreshold = duration
		}
	}

	// 解析 PostgreSQL 配置
	if yamlConfig.Database.PostgreSQL.Host != "" || yamlConfig.Database.PostgreSQL.DSN != "" {
		config.PostgreSQL = &dbconfig.PostgreSQLConfig{
			DSN:      yamlConfig.Database.PostgreSQL.DSN,
			Host:     getValueOrDefault(yamlConfig.Database.PostgreSQL.Host, "localhost"),
			Port:     getIntOrDefault(yamlConfig.Database.PostgreSQL.Port, 5432),
			User:     getValueOrDefault(yamlConfig.Database.PostgreSQL.User, "postgres"),
			Password: yamlConfig.Database.PostgreSQL.Password,
			DBName:   getValueOrDefault(yamlConfig.Database.PostgreSQL.DBName, "structforge"),
			SSLMode:  getValueOrDefault(yamlConfig.Database.PostgreSQL.SSLMode, "disable"),
			TimeZone: getValueOrDefault(yamlConfig.Database.PostgreSQL.TimeZone, "Asia/Shanghai"),
		}
	}

	// 解析 MySQL 配置
	if yamlConfig.Database.MySQL.Host != "" || yamlConfig.Database.MySQL.DSN != "" {
		config.MySQL = &dbconfig.MySQLConfig{
			DSN:       yamlConfig.Database.MySQL.DSN,
			Host:      getValueOrDefault(yamlConfig.Database.MySQL.Host, "localhost"),
			Port:      getIntOrDefault(yamlConfig.Database.MySQL.Port, 3306),
			User:      getValueOrDefault(yamlConfig.Database.MySQL.User, "root"),
			Password:  yamlConfig.Database.MySQL.Password,
			DBName:    getValueOrDefault(yamlConfig.Database.MySQL.DBName, "structforge"),
			Charset:   getValueOrDefault(yamlConfig.Database.MySQL.Charset, "utf8mb4"),
			ParseTime: yamlConfig.Database.MySQL.ParseTime,
			Loc:       getValueOrDefault(yamlConfig.Database.MySQL.Loc, "Local"),
		}
	}

	// 解析 SQLite 配置
	if yamlConfig.Database.SQLite.Path != "" || yamlConfig.Database.SQLite.DSN != "" {
		busyTimeout := 5 * time.Second
		if yamlConfig.Database.SQLite.BusyTimeout != "" {
			if duration, err := time.ParseDuration(yamlConfig.Database.SQLite.BusyTimeout); err == nil {
				busyTimeout = duration
			}
		}
		config.SQLite = &dbconfig.SQLiteConfig{
			DSN:         yamlConfig.Database.SQLite.DSN,
			Path:        getValueOrDefault(yamlConfig.Database.SQLite.Path, "data.db"),
			ForeignKeys: yamlConfig.Database.SQLite.ForeignKeys,
			BusyTimeout: busyTimeout,
		}
	}

	return config, nil
}

// getValueOrDefault 获取值，如果为空则返回默认值
func getValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// getIntOrDefault 获取整数值，如果为0则返回默认值
func getIntOrDefault(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}
