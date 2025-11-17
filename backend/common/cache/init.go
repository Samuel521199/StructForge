package cache

import (
	"context"
	"os"
	"sync"
	"time"

	memoryAdapter "StructForge/backend/common/cache/adapters/memory"
	redisAdapter "StructForge/backend/common/cache/adapters/redis"
)

var (
	globalCache Cache
	globalMu    sync.RWMutex
)

// InitOptions 初始化选项（函数式选项模式）
type InitOptions struct {
	// 适配器类型（必填）
	AdapterType AdapterType

	// 键前缀（可选）
	KeyPrefix string

	// 默认TTL（可选）
	DefaultTTL time.Duration

	// 序列化器类型（可选）
	SerializerType SerializerType

	// Redis配置（可选）
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// 内存缓存配置（可选）
	MemoryMaxSize  int64
	MemoryMaxItems int

	// 高级配置（可选，用于覆盖默认配置）
	Config *Config
}

// InitOption 初始化选项函数类型
type InitOption func(*InitOptions)

// WithKeyPrefix 设置键前缀
func WithKeyPrefix(prefix string) InitOption {
	return func(opts *InitOptions) {
		opts.KeyPrefix = prefix
	}
}

// WithDefaultTTL 设置默认TTL
func WithDefaultTTL(ttl time.Duration) InitOption {
	return func(opts *InitOptions) {
		opts.DefaultTTL = ttl
	}
}

// WithSerializer 设置序列化器类型
func WithSerializer(serializerType SerializerType) InitOption {
	return func(opts *InitOptions) {
		opts.SerializerType = serializerType
	}
}

// WithRedisAddr 设置Redis地址
func WithRedisAddr(addr string) InitOption {
	return func(opts *InitOptions) {
		opts.RedisAddr = addr
	}
}

// WithRedisPassword 设置Redis密码
func WithRedisPassword(password string) InitOption {
	return func(opts *InitOptions) {
		opts.RedisPassword = password
	}
}

// WithRedisDB 设置Redis数据库编号
func WithRedisDB(db int) InitOption {
	return func(opts *InitOptions) {
		opts.RedisDB = db
	}
}

// WithMemoryMaxSize 设置内存缓存最大大小
func WithMemoryMaxSize(maxSize int64) InitOption {
	return func(opts *InitOptions) {
		opts.MemoryMaxSize = maxSize
	}
}

// WithMemoryMaxItems 设置内存缓存最大条目数
func WithMemoryMaxItems(maxItems int) InitOption {
	return func(opts *InitOptions) {
		opts.MemoryMaxItems = maxItems
	}
}

// WithConfig 设置高级配置（会覆盖其他选项）
func WithConfig(config Config) InitOption {
	return func(opts *InitOptions) {
		opts.Config = &config
	}
}

// InitCache 初始化缓存系统
// adapterType: 适配器类型（"redis" 或 "memory"）
// options: 可选配置项
func InitCache(adapterType string, options ...InitOption) (Cache, error) {
	opts := &InitOptions{
		AdapterType:    AdapterType(adapterType),
		KeyPrefix:      "",
		DefaultTTL:     5 * time.Minute,
		SerializerType: SerializerJSON,
		RedisAddr:      getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RedisDB:        0,
		MemoryMaxSize:  100 * 1024 * 1024, // 100MB
		MemoryMaxItems: 1000,
	}

	// 应用选项
	for _, opt := range options {
		opt(opts)
	}

	// 构建配置
	var config Config
	if opts.Config != nil {
		config = *opts.Config
	} else {
		config = DefaultConfig()
		config.AdapterType = opts.AdapterType
		config.KeyPrefix = opts.KeyPrefix
		config.DefaultTTL = opts.DefaultTTL
		config.SerializerType = opts.SerializerType

		// 更新Redis配置
		if config.Redis != nil {
			if opts.RedisAddr != "" {
				config.Redis.Addr = opts.RedisAddr
			}
			if opts.RedisPassword != "" {
				config.Redis.Password = opts.RedisPassword
			}
			config.Redis.DB = opts.RedisDB
		}

		// 更新内存配置
		if config.Memory != nil {
			if opts.MemoryMaxSize > 0 {
				config.Memory.MaxSize = opts.MemoryMaxSize
			}
			if opts.MemoryMaxItems > 0 {
				config.Memory.MaxItems = opts.MemoryMaxItems
			}
		}
	}

	// 创建适配器
	cache, err := NewCache(config)
	if err != nil {
		return nil, err
	}

	// 设置为全局缓存
	SetGlobalCache(cache)

	return cache, nil
}

// InitCacheWithShutdown 初始化缓存系统并返回关闭函数
// 用法: defer cache.InitCacheWithShutdown("redis")()
func InitCacheWithShutdown(adapterType string, options ...InitOption) func() {
	cache, err := InitCache(adapterType, options...)
	if err != nil {
		// 如果初始化失败，返回空函数
		// 实际使用时应该处理错误
		return func() {}
	}

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = cache.Close()
		_ = Shutdown(ctx)
	}
}

// NewCache 根据配置创建缓存实例
func NewCache(config Config) (Cache, error) {
	switch config.AdapterType {
	case AdapterRedis:
		if config.Redis == nil {
			config.Redis = &RedisConfig{
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
			}
		}
		return newRedisAdapter(config)
	case AdapterMemory:
		if config.Memory == nil {
			config.Memory = &MemoryConfig{
				MaxSize:         100 * 1024 * 1024,
				MaxItems:        1000,
				Strategy:        "lru",
				CleanupInterval: 10 * time.Minute,
			}
		}
		return newMemoryAdapter(config)
	case AdapterMulti:
		return newMultiAdapter(config)
	default:
		return nil, ErrNotSupported
	}
}

// newRedisAdapter 创建Redis适配器（内部函数）
func newRedisAdapter(config Config) (Cache, error) {
	redisCfg := &redisAdapter.Config{
		Addr:         config.Redis.Addr,
		Password:     config.Redis.Password,
		DB:           config.Redis.DB,
		PoolSize:     config.Redis.PoolSize,
		MinIdleConns: config.Redis.MinIdleConns,
		MaxRetries:   config.Redis.MaxRetries,
		DialTimeout:  config.Redis.DialTimeout,
		ReadTimeout:  config.Redis.ReadTimeout,
		WriteTimeout: config.Redis.WriteTimeout,
		PoolTimeout:  config.Redis.PoolTimeout,
	}
	adapter, err := redisAdapter.NewAdapter(redisAdapter.AdapterConfig{
		RedisConfig: redisCfg,
		KeyPrefix:   config.KeyPrefix,
	})
	if err != nil {
		return nil, err
	}
	return &adapterWrapper{adapter: adapter}, nil
}

// newMemoryAdapter 创建内存适配器（内部函数）
func newMemoryAdapter(config Config) (Cache, error) {
	memoryCfg := &memoryAdapter.Config{
		MaxSize:         config.Memory.MaxSize,
		MaxItems:        config.Memory.MaxItems,
		Strategy:        config.Memory.Strategy,
		CleanupInterval: config.Memory.CleanupInterval,
	}
	adapter, err := memoryAdapter.NewAdapter(memoryAdapter.AdapterConfig{
		MemoryConfig: memoryCfg,
		KeyPrefix:    config.KeyPrefix,
	})
	if err != nil {
		return nil, err
	}
	return &adapterWrapper{adapter: adapter}, nil
}

// adapterWrapper 适配器包装器，将适配器接口转换为Cache接口
type adapterWrapper struct {
	adapter interface {
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
	}
}

func (w *adapterWrapper) Get(ctx context.Context, key string) ([]byte, error) {
	return w.adapter.Get(ctx, key)
}

func (w *adapterWrapper) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return w.adapter.Set(ctx, key, value, ttl)
}

func (w *adapterWrapper) Delete(ctx context.Context, key string) error {
	return w.adapter.Delete(ctx, key)
}

func (w *adapterWrapper) Exists(ctx context.Context, key string) (bool, error) {
	return w.adapter.Exists(ctx, key)
}

func (w *adapterWrapper) MGet(ctx context.Context, keys ...string) (map[string][]byte, error) {
	return w.adapter.MGet(ctx, keys...)
}

func (w *adapterWrapper) MSet(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	return w.adapter.MSet(ctx, items, ttl)
}

func (w *adapterWrapper) MDelete(ctx context.Context, keys ...string) error {
	return w.adapter.MDelete(ctx, keys...)
}

func (w *adapterWrapper) GetOrSet(ctx context.Context, key string, fn func() ([]byte, error), ttl time.Duration) ([]byte, error) {
	return w.adapter.GetOrSet(ctx, key, fn, ttl)
}

func (w *adapterWrapper) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	return w.adapter.Increment(ctx, key, delta)
}

func (w *adapterWrapper) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return w.adapter.Decrement(ctx, key, delta)
}

func (w *adapterWrapper) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return w.adapter.Expire(ctx, key, ttl)
}

func (w *adapterWrapper) TTL(ctx context.Context, key string) (time.Duration, error) {
	return w.adapter.TTL(ctx, key)
}

func (w *adapterWrapper) Keys(ctx context.Context, pattern string) ([]string, error) {
	return w.adapter.Keys(ctx, pattern)
}

func (w *adapterWrapper) Scan(ctx context.Context, pattern string, count int) ([]string, error) {
	return w.adapter.Scan(ctx, pattern, count)
}

func (w *adapterWrapper) Clear(ctx context.Context, pattern string) error {
	return w.adapter.Clear(ctx, pattern)
}

func (w *adapterWrapper) Close() error {
	return w.adapter.Close()
}

func (w *adapterWrapper) Ping(ctx context.Context) error {
	return w.adapter.Ping(ctx)
}

// newMultiAdapter 创建多级缓存适配器（内部函数，暂未实现）
func newMultiAdapter(config Config) (Cache, error) {
	return nil, ErrNotSupported
}

// SetGlobalCache 设置全局缓存实例
func SetGlobalCache(cache Cache) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalCache = cache
}

// GetGlobalCache 获取全局缓存实例
func GetGlobalCache() Cache {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalCache
}

// Shutdown 优雅关闭全局缓存
func Shutdown(ctx context.Context) error {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalCache != nil {
		err := globalCache.Close()
		globalCache = nil
		return err
	}
	return nil
}

// 全局便捷方法（使用全局缓存实例）

// Get 获取缓存值
func Get(ctx context.Context, key string) ([]byte, error) {
	cache := GetGlobalCache()
	if cache == nil {
		return nil, ErrAdapterNotInitialized
	}
	return cache.Get(ctx, key)
}

// Set 设置缓存值
func Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	cache := GetGlobalCache()
	if cache == nil {
		return ErrAdapterNotInitialized
	}
	return cache.Set(ctx, key, value, ttl)
}

// Delete 删除缓存键
func Delete(ctx context.Context, key string) error {
	cache := GetGlobalCache()
	if cache == nil {
		return ErrAdapterNotInitialized
	}
	return cache.Delete(ctx, key)
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	cache := GetGlobalCache()
	if cache == nil {
		return false, ErrAdapterNotInitialized
	}
	return cache.Exists(ctx, key)
}

// GetObject 获取对象（自动反序列化）
func GetObject(ctx context.Context, key string, v interface{}) error {
	cache := GetGlobalCache()
	if cache == nil {
		return ErrAdapterNotInitialized
	}

	data, err := cache.Get(ctx, key)
	if err != nil {
		return err
	}

	// 这里需要从配置中获取序列化器，简化实现
	serializer := NewSerializer(SerializerJSON)
	return serializer.Deserialize(data, v)
}

// SetObject 设置对象（自动序列化）
func SetObject(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	cache := GetGlobalCache()
	if cache == nil {
		return ErrAdapterNotInitialized
	}

	// 这里需要从配置中获取序列化器，简化实现
	serializer := NewSerializer(SerializerJSON)
	data, err := serializer.Serialize(v)
	if err != nil {
		return err
	}

	return cache.Set(ctx, key, data, ttl)
}
