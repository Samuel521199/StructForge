package redis

import (
	"time"
)

// Config Redis适配器配置
type Config struct {
	// Addr Redis地址，如 "localhost:6379"
	Addr string

	// Password 密码
	Password string

	// DB 数据库编号
	DB int

	// PoolSize 连接池大小
	PoolSize int

	// MinIdleConns 最小空闲连接数
	MinIdleConns int

	// MaxRetries 最大重试次数
	MaxRetries int

	// DialTimeout 连接超时
	DialTimeout time.Duration

	// ReadTimeout 读取超时
	ReadTimeout time.Duration

	// WriteTimeout 写入超时
	WriteTimeout time.Duration

	// PoolTimeout 连接池超时
	PoolTimeout time.Duration
}
