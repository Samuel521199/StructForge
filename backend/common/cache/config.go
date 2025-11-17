package cache

import (
	"os"
	"time"
)

// AdapterType 适配器类型
type AdapterType string

const (
	// AdapterRedis Redis适配器
	AdapterRedis AdapterType = "redis"

	// AdapterMemory 内存适配器
	AdapterMemory AdapterType = "memory"

	// AdapterMulti 多级缓存适配器
	AdapterMulti AdapterType = "multi"
)

// SerializerType 序列化器类型
type SerializerType string

const (
	// SerializerJSON JSON序列化器
	SerializerJSON SerializerType = "json"

	// SerializerMsgPack MsgPack序列化器
	SerializerMsgPack SerializerType = "msgpack"

	// SerializerGob Gob序列化器
	SerializerGob SerializerType = "gob"
)

// Config 缓存系统配置
type Config struct {
	// 适配器类型
	AdapterType AdapterType

	// 默认TTL
	DefaultTTL time.Duration

	// 键前缀（用于多租户隔离）
	KeyPrefix string

	// 序列化器类型
	SerializerType SerializerType

	// 适配器配置（根据类型选择）
	Redis  *RedisConfig
	Memory *MemoryConfig
	Multi  *MultiConfig

	// 性能配置
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration

	// 日志配置
	EnableLogging bool
}

// RedisConfig Redis适配器配置
type RedisConfig struct {
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

// MemoryConfig 内存缓存配置
type MemoryConfig struct {
	// MaxSize 最大内存使用（字节），0表示不限制
	MaxSize int64

	// MaxItems 最大条目数，0表示不限制
	MaxItems int

	// Strategy 淘汰策略："lru", "lfu", "fifo"
	Strategy string

	// CleanupInterval 清理间隔
	CleanupInterval time.Duration
}

// MultiConfig 多级缓存配置
type MultiConfig struct {
	// L1Config L1缓存配置（内存）
	L1Config *MemoryConfig

	// L2Config L2缓存配置（Redis）
	L2Config *RedisConfig

	// L1TTL L1缓存TTL（通常比L2短）
	L1TTL time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		AdapterType:    AdapterRedis,
		DefaultTTL:     5 * time.Minute,
		KeyPrefix:      "",
		SerializerType: SerializerJSON,
		Redis: &RedisConfig{
			Addr:         getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 5,
			MaxRetries:   3,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,
		},
		Memory: &MemoryConfig{
			MaxSize:         100 * 1024 * 1024, // 100MB
			MaxItems:        1000,
			Strategy:        "lru",
			CleanupInterval: 10 * time.Minute,
		},
		MaxRetries:    3,
		RetryDelay:    100 * time.Millisecond,
		Timeout:       5 * time.Second,
		EnableLogging: true,
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
