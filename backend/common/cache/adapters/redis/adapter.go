package redis

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// AdapterConfig Redis适配器配置
type AdapterConfig struct {
	RedisConfig *Config
	KeyPrefix   string
}

// CacheAdapter 缓存适配器接口（避免循环导入）
type CacheAdapter interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	MGet(ctx context.Context, keys ...string) (map[string][]byte, error)
	MSet(ctx context.Context, items map[string][]byte, ttl time.Duration) error
	MDelete(ctx context.Context, keys ...string) error
	GetOrSet(ctx context.Context, key string, fn func() ([]byte, error), ttl time.Duration) ([]byte, error)
	Increment(ctx context.Context, key string, delta int64) (int64, error)
	Decrement(ctx context.Context, key string, delta int64) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, pattern string, count int) ([]string, error)
	Clear(ctx context.Context, pattern string) error
	Close() error
	Ping(ctx context.Context) error
	Name() string
	HealthCheck(ctx context.Context) error
}

// adapter Redis缓存适配器
type adapter struct {
	client interface{} // 使用interface{}避免直接依赖，实际类型为 *redis.Client
	config Config
	prefix string
}

// RedisClient Redis客户端接口（用于类型断言）
type RedisClient interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	Close() error
}

// NewAdapter 创建Redis适配器
func NewAdapter(config AdapterConfig) (CacheAdapter, error) {
	if config.RedisConfig == nil {
		return nil, errors.New("redis config is required")
	}

	// TODO: 实际实现需要取消注释以下代码并导入redis包
	/*
		import "github.com/redis/go-redis/v9"

		client := redis.NewClient(&redis.Options{
			Addr:         config.RedisConfig.Addr,
			Password:     config.RedisConfig.Password,
			DB:           config.RedisConfig.DB,
			PoolSize:     config.RedisConfig.PoolSize,
			MinIdleConns: config.RedisConfig.MinIdleConns,
			MaxRetries:   config.RedisConfig.MaxRetries,
			DialTimeout:  config.RedisConfig.DialTimeout,
			ReadTimeout:  config.RedisConfig.ReadTimeout,
			WriteTimeout: config.RedisConfig.WriteTimeout,
			PoolTimeout:  config.RedisConfig.PoolTimeout,
		})

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("redis connection failed: %w", err)
		}

		return &adapter{
			client: client,
			config: *config.RedisConfig,
			prefix: config.KeyPrefix,
		}, nil
	*/

	// 临时返回错误，提示需要实现
	return nil, fmt.Errorf("redis adapter requires github.com/redis/go-redis/v9, please install it and uncomment the implementation")
}

// Name 返回适配器名称
func (a *adapter) Name() string {
	return "redis"
}

// HealthCheck 健康检查
func (a *adapter) HealthCheck(ctx context.Context) error {
	if a.client == nil {
		return errors.New("redis client not initialized")
	}

	client, ok := a.client.(RedisClient)
	if !ok {
		return errors.New("invalid redis client type")
	}

	return client.Ping(ctx)
}

// buildKey 构建带前缀的键
func (a *adapter) buildKey(key string) string {
	if a.prefix == "" {
		return key
	}
	return a.prefix + ":" + key
}

// Get 获取缓存值
func (a *adapter) Get(ctx context.Context, key string) ([]byte, error) {
	client, ok := a.client.(RedisClient)
	if !ok {
		return nil, errors.New("invalid redis client type")
	}

	val, err := client.Get(ctx, a.buildKey(key))
	if err != nil {
		// 检查是否是键不存在的错误
		if err.Error() == "redis: nil" || err.Error() == "redis: key does not exist" {
			return nil, errors.New("key not found")
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}
	return []byte(val), nil
}

// Set 设置缓存值
func (a *adapter) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	client, ok := a.client.(RedisClient)
	if !ok {
		return errors.New("invalid redis client type")
	}

	err := client.Set(ctx, a.buildKey(key), value, ttl)
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

// Delete 删除缓存键
func (a *adapter) Delete(ctx context.Context, key string) error {
	client, ok := a.client.(RedisClient)
	if !ok {
		return errors.New("invalid redis client type")
	}

	_, err := client.Del(ctx, a.buildKey(key))
	if err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}

// Exists 检查键是否存在
func (a *adapter) Exists(ctx context.Context, key string) (bool, error) {
	client, ok := a.client.(RedisClient)
	if !ok {
		return false, errors.New("invalid redis client type")
	}

	count, err := client.Exists(ctx, a.buildKey(key))
	if err != nil {
		return false, fmt.Errorf("redis exists error: %w", err)
	}
	return count > 0, nil
}

// MGet 批量获取缓存值
func (a *adapter) MGet(ctx context.Context, keys ...string) (map[string][]byte, error) {
	if len(keys) == 0 {
		return make(map[string][]byte), nil
	}

	client, ok := a.client.(RedisClient)
	if !ok {
		return nil, errors.New("invalid redis client type")
	}

	// 构建带前缀的键
	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = a.buildKey(key)
	}

	vals, err := client.MGet(ctx, prefixedKeys...)
	if err != nil {
		return nil, fmt.Errorf("redis mget error: %w", err)
	}

	result := make(map[string][]byte)
	for i, val := range vals {
		if val != nil {
			if str, ok := val.(string); ok {
				result[keys[i]] = []byte(str)
			}
		}
	}

	return result, nil
}

// MSet 批量设置缓存值
func (a *adapter) MSet(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	if len(items) == 0 {
		return nil
	}

	// 简化实现：逐个设置
	// 实际Redis实现应该使用Pipeline
	for key, value := range items {
		if err := a.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}

	return nil
}

// MDelete 批量删除缓存键
func (a *adapter) MDelete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	client, ok := a.client.(RedisClient)
	if !ok {
		return errors.New("invalid redis client type")
	}

	prefixedKeys := make([]string, len(keys))
	for i, key := range keys {
		prefixedKeys[i] = a.buildKey(key)
	}

	_, err := client.Del(ctx, prefixedKeys...)
	if err != nil {
		return fmt.Errorf("redis mdelete error: %w", err)
	}

	return nil
}

// GetOrSet 获取缓存值，如果不存在则调用函数获取并缓存
func (a *adapter) GetOrSet(ctx context.Context, key string, fn func() ([]byte, error), ttl time.Duration) ([]byte, error) {
	// 先尝试获取
	val, err := a.Get(ctx, key)
	if err == nil {
		return val, nil
	}
	if err.Error() != "key not found" {
		return nil, err
	}

	// 不存在，调用函数获取
	val, err = fn()
	if err != nil {
		return nil, err
	}

	// 设置缓存（忽略错误）
	_ = a.Set(ctx, key, val, ttl)

	return val, nil
}

// Increment 增加数值
func (a *adapter) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	client, ok := a.client.(RedisClient)
	if !ok {
		return 0, errors.New("invalid redis client type")
	}

	val, err := client.IncrBy(ctx, a.buildKey(key), delta)
	if err != nil {
		return 0, fmt.Errorf("redis increment error: %w", err)
	}
	return val, nil
}

// Decrement 减少数值
func (a *adapter) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return a.Increment(ctx, key, -delta)
}

// Expire 设置键的过期时间
func (a *adapter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	client, ok := a.client.(RedisClient)
	if !ok {
		return errors.New("invalid redis client type")
	}

	_, err := client.Expire(ctx, a.buildKey(key), ttl)
	if err != nil {
		return fmt.Errorf("redis expire error: %w", err)
	}
	return nil
}

// TTL 获取键的剩余过期时间
func (a *adapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	client, ok := a.client.(RedisClient)
	if !ok {
		return 0, errors.New("invalid redis client type")
	}

	ttl, err := client.TTL(ctx, a.buildKey(key))
	if err != nil {
		return 0, fmt.Errorf("redis ttl error: %w", err)
	}
	return ttl, nil
}

// Keys 获取匹配模式的所有键
func (a *adapter) Keys(ctx context.Context, pattern string) ([]string, error) {
	client, ok := a.client.(RedisClient)
	if !ok {
		return nil, errors.New("invalid redis client type")
	}

	prefixedPattern := a.buildKey(pattern)
	keys, err := client.Keys(ctx, prefixedPattern)
	if err != nil {
		return nil, fmt.Errorf("redis keys error: %w", err)
	}

	// 移除前缀
	if a.prefix != "" {
		prefix := a.prefix + ":"
		for i, key := range keys {
			if len(key) > len(prefix) && key[:len(prefix)] == prefix {
				keys[i] = key[len(prefix):]
			}
		}
	}

	return keys, nil
}

// Scan 扫描匹配模式的键
func (a *adapter) Scan(ctx context.Context, pattern string, count int) ([]string, error) {
	client, ok := a.client.(RedisClient)
	if !ok {
		return nil, errors.New("invalid redis client type")
	}

	prefixedPattern := a.buildKey(pattern)
	var keys []string
	var cursor uint64 = 0

	for {
		batch, nextCursor, err := client.Scan(ctx, cursor, prefixedPattern, int64(count))
		if err != nil {
			return nil, fmt.Errorf("redis scan error: %w", err)
		}

		keys = append(keys, batch...)
		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	// 移除前缀
	if a.prefix != "" {
		prefix := a.prefix + ":"
		for i, key := range keys {
			if len(key) > len(prefix) && key[:len(prefix)] == prefix {
				keys[i] = key[len(prefix):]
			}
		}
	}

	return keys, nil
}

// Clear 清除匹配模式的所有键
func (a *adapter) Clear(ctx context.Context, pattern string) error {
	keys, err := a.Scan(ctx, pattern, 100)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return a.MDelete(ctx, keys...)
}

// Close 关闭缓存连接
func (a *adapter) Close() error {
	if a.client == nil {
		return nil
	}

	client, ok := a.client.(RedisClient)
	if !ok {
		return errors.New("invalid redis client type")
	}

	return client.Close()
}

// Ping 检查连接是否正常
func (a *adapter) Ping(ctx context.Context) error {
	return a.HealthCheck(ctx)
}
