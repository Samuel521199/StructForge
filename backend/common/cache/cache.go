package cache

import (
	"context"
	"time"
)

// Cache 缓存接口
// 提供统一的缓存操作API
type Cache interface {
	// Get 获取缓存值
	// key: 缓存键
	// 返回: 缓存值（字节数组）和错误
	Get(ctx context.Context, key string) ([]byte, error)

	// Set 设置缓存值
	// key: 缓存键
	// value: 缓存值（字节数组）
	// ttl: 过期时间，0表示不过期
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete 删除缓存键
	// key: 缓存键
	Delete(ctx context.Context, key string) error

	// Exists 检查键是否存在
	// key: 缓存键
	// 返回: 是否存在和错误
	Exists(ctx context.Context, key string) (bool, error)

	// MGet 批量获取缓存值
	// keys: 缓存键列表
	// 返回: 键值对映射（不存在的键不会出现在结果中）
	MGet(ctx context.Context, keys ...string) (map[string][]byte, error)

	// MSet 批量设置缓存值
	// items: 键值对映射
	// ttl: 过期时间，0表示不过期
	MSet(ctx context.Context, items map[string][]byte, ttl time.Duration) error

	// MDelete 批量删除缓存键
	// keys: 缓存键列表
	MDelete(ctx context.Context, keys ...string) error

	// GetOrSet 获取缓存值，如果不存在则调用函数获取并缓存
	// key: 缓存键
	// fn: 获取值的函数
	// ttl: 过期时间
	// 返回: 缓存值（字节数组）和错误
	GetOrSet(ctx context.Context, key string, fn func() ([]byte, error), ttl time.Duration) ([]byte, error)

	// Increment 增加数值
	// key: 缓存键
	// delta: 增量
	// 返回: 增加后的值和错误
	Increment(ctx context.Context, key string, delta int64) (int64, error)

	// Decrement 减少数值
	// key: 缓存键
	// delta: 减量
	// 返回: 减少后的值和错误
	Decrement(ctx context.Context, key string, delta int64) (int64, error)

	// Expire 设置键的过期时间
	// key: 缓存键
	// ttl: 过期时间
	Expire(ctx context.Context, key string, ttl time.Duration) error

	// TTL 获取键的剩余过期时间
	// key: 缓存键
	// 返回: 剩余过期时间，-1表示不过期，-2表示键不存在
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Keys 获取匹配模式的所有键（谨慎使用，可能影响性能）
	// pattern: 匹配模式，支持通配符（如 "user:*"）
	// 返回: 匹配的键列表
	Keys(ctx context.Context, pattern string) ([]string, error)

	// Scan 扫描匹配模式的键（推荐使用，性能更好）
	// pattern: 匹配模式
	// count: 每次扫描的数量
	// 返回: 匹配的键列表
	Scan(ctx context.Context, pattern string, count int) ([]string, error)

	// Clear 清除匹配模式的所有键
	// pattern: 匹配模式
	Clear(ctx context.Context, pattern string) error

	// Close 关闭缓存连接
	Close() error

	// Ping 检查连接是否正常
	Ping(ctx context.Context) error
}
