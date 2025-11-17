package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// AdapterConfig 内存适配器配置
type AdapterConfig struct {
	MemoryConfig *Config
	KeyPrefix    string
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

// item 缓存项
type item struct {
	value       []byte
	expiresAt   time.Time
	accessTime  time.Time
	accessCount int64
}

// isExpired 检查是否过期
func (i *item) isExpired() bool {
	if i.expiresAt.IsZero() {
		return false
	}
	return time.Now().After(i.expiresAt)
}

// adapter 内存缓存适配器
type adapter struct {
	mu            sync.RWMutex
	items         map[string]*item
	config        Config
	prefix        string
	currentSize   int64
	maxSize       int64
	maxItems      int
	strategy      string
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

// NewAdapter 创建内存缓存适配器
func NewAdapter(config AdapterConfig) (CacheAdapter, error) {
	if config.MemoryConfig == nil {
		config.MemoryConfig = &Config{
			MaxSize:         100 * 1024 * 1024, // 100MB
			MaxItems:        1000,
			Strategy:        "lru",
			CleanupInterval: 10 * time.Minute,
		}
	}

	a := &adapter{
		items:       make(map[string]*item),
		config:      *config.MemoryConfig,
		prefix:      config.KeyPrefix,
		maxSize:     config.MemoryConfig.MaxSize,
		maxItems:    config.MemoryConfig.MaxItems,
		strategy:    config.MemoryConfig.Strategy,
		stopCleanup: make(chan struct{}),
	}

	// 启动清理协程
	if config.MemoryConfig.CleanupInterval > 0 {
		a.cleanupTicker = time.NewTicker(config.MemoryConfig.CleanupInterval)
		go a.cleanupLoop()
	}

	return a, nil
}

// cleanupLoop 清理循环
func (a *adapter) cleanupLoop() {
	for {
		select {
		case <-a.cleanupTicker.C:
			a.cleanup()
		case <-a.stopCleanup:
			return
		}
	}
}

// cleanup 清理过期项
func (a *adapter) cleanup() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for key, item := range a.items {
		if item.isExpired() {
			delete(a.items, key)
			a.currentSize -= int64(len(item.value))
		}
	}

	// 如果超过限制，执行淘汰
	if a.shouldEvict() {
		a.evict()
	}
}

// shouldEvict 检查是否需要淘汰
func (a *adapter) shouldEvict() bool {
	if a.maxSize > 0 && a.currentSize > a.maxSize {
		return true
	}
	if a.maxItems > 0 && len(a.items) > a.maxItems {
		return true
	}
	return false
}

// evict 执行淘汰
func (a *adapter) evict() {
	switch a.strategy {
	case "lru":
		a.evictLRU()
	case "lfu":
		a.evictLFU()
	case "fifo":
		a.evictFIFO()
	default:
		a.evictLRU() // 默认使用LRU
	}
}

// evictLRU 最近最少使用淘汰
func (a *adapter) evictLRU() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, item := range a.items {
		if first || item.accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.accessTime
			first = false
		}
	}

	if oldestKey != "" {
		if item, ok := a.items[oldestKey]; ok {
			a.currentSize -= int64(len(item.value))
			delete(a.items, oldestKey)
		}
	}
}

// evictLFU 最少使用淘汰
func (a *adapter) evictLFU() {
	var leastKey string
	var leastCount int64 = -1

	for key, item := range a.items {
		if leastCount == -1 || item.accessCount < leastCount {
			leastKey = key
			leastCount = item.accessCount
		}
	}

	if leastKey != "" {
		if item, ok := a.items[leastKey]; ok {
			a.currentSize -= int64(len(item.value))
			delete(a.items, leastKey)
		}
	}
}

// evictFIFO 先进先出淘汰
func (a *adapter) evictFIFO() {
	// FIFO需要记录插入顺序，这里简化实现，使用任意一个键
	for key := range a.items {
		if item, ok := a.items[key]; ok {
			a.currentSize -= int64(len(item.value))
			delete(a.items, key)
		}
		break
	}
}

// Name 返回适配器名称
func (a *adapter) Name() string {
	return "memory"
}

// HealthCheck 健康检查
func (a *adapter) HealthCheck(ctx context.Context) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// 内存缓存总是健康的
	return nil
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
	a.mu.RLock()
	item, exists := a.items[a.buildKey(key)]
	a.mu.RUnlock()

	if !exists || item.isExpired() {
		if exists {
			// 删除过期项
			a.mu.Lock()
			delete(a.items, a.buildKey(key))
			if exists {
				a.currentSize -= int64(len(item.value))
			}
			a.mu.Unlock()
		}
		return nil, errors.New("key not found")
	}

	// 更新访问时间和访问次数
	a.mu.Lock()
	item.accessTime = time.Now()
	item.accessCount++
	a.mu.Unlock()

	// 返回值的副本
	result := make([]byte, len(item.value))
	copy(result, item.value)
	return result, nil
}

// Set 设置缓存值
func (a *adapter) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	key = a.buildKey(key)

	// 计算新项大小
	newSize := int64(len(value))

	a.mu.Lock()
	defer a.mu.Unlock()

	// 如果键已存在，先减去旧值大小
	if oldItem, exists := a.items[key]; exists {
		a.currentSize -= int64(len(oldItem.value))
	}

	// 检查是否需要淘汰
	for a.shouldEvict() && len(a.items) > 0 {
		a.evict()
	}

	// 创建新项
	now := time.Now()
	newItem := &item{
		value:       make([]byte, len(value)),
		accessTime:  now,
		accessCount: 1,
	}
	copy(newItem.value, value)

	if ttl > 0 {
		newItem.expiresAt = now.Add(ttl)
	}

	a.items[key] = newItem
	a.currentSize += newSize

	return nil
}

// Delete 删除缓存键
func (a *adapter) Delete(ctx context.Context, key string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	key = a.buildKey(key)
	if item, exists := a.items[key]; exists {
		a.currentSize -= int64(len(item.value))
		delete(a.items, key)
	}

	return nil
}

// Exists 检查键是否存在
func (a *adapter) Exists(ctx context.Context, key string) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	item, exists := a.items[a.buildKey(key)]
	if !exists {
		return false, nil
	}

	if item.isExpired() {
		return false, nil
	}

	return true, nil
}

// MGet 批量获取缓存值
func (a *adapter) MGet(ctx context.Context, keys ...string) (map[string][]byte, error) {
	result := make(map[string][]byte)

	a.mu.RLock()
	for _, key := range keys {
		item, exists := a.items[a.buildKey(key)]
		if exists && !item.isExpired() {
			// 返回值的副本
			value := make([]byte, len(item.value))
			copy(value, item.value)
			result[key] = value

			// 更新访问信息
			item.accessTime = time.Now()
			item.accessCount++
		}
	}
	a.mu.RUnlock()

	return result, nil
}

// MSet 批量设置缓存值
func (a *adapter) MSet(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	for key, value := range items {
		if err := a.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}
	return nil
}

// MDelete 批量删除缓存键
func (a *adapter) MDelete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		if err := a.Delete(ctx, key); err != nil {
			return err
		}
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
	a.mu.Lock()
	defer a.mu.Unlock()

	key = a.buildKey(key)
	it, exists := a.items[key]

	var currentValue int64
	if exists && !it.isExpired() {
		// 尝试解析现有值
		if len(it.value) > 0 {
			// 简化实现，假设值是数字的字符串表示
			// 实际应该使用更健壮的解析
			fmt.Sscanf(string(it.value), "%d", &currentValue)
		}
	}

	newValue := currentValue + delta
	newValueBytes := []byte(fmt.Sprintf("%d", newValue))

	// 使用默认TTL
	ttl := 5 * time.Minute
	now := time.Now()
	newIt := &item{
		value:       newValueBytes,
		accessTime:  now,
		accessCount: 1,
		expiresAt:   now.Add(ttl),
	}

	if exists {
		a.currentSize -= int64(len(it.value))
	}

	a.items[key] = newIt
	a.currentSize += int64(len(newValueBytes))

	return newValue, nil
}

// Decrement 减少数值
func (a *adapter) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return a.Increment(ctx, key, -delta)
}

// Expire 设置键的过期时间
func (a *adapter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	key = a.buildKey(key)
	item, exists := a.items[key]
	if !exists || item.isExpired() {
		return errors.New("key not found")
	}

	item.expiresAt = time.Now().Add(ttl)
	return nil
}

// TTL 获取键的剩余过期时间
func (a *adapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	item, exists := a.items[a.buildKey(key)]
	if !exists {
		return -2 * time.Second, nil // -2表示键不存在
	}

	if item.isExpired() {
		return -2 * time.Second, nil
	}

	if item.expiresAt.IsZero() {
		return -1 * time.Second, nil // -1表示不过期
	}

	remaining := time.Until(item.expiresAt)
	if remaining < 0 {
		return 0, nil
	}

	return remaining, nil
}

// Keys 获取匹配模式的所有键
func (a *adapter) Keys(ctx context.Context, pattern string) ([]string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// 简化实现，不支持通配符匹配
	// 实际应该使用正则表达式或字符串匹配
	keys := make([]string, 0, len(a.items))

	for key := range a.items {
		// 移除前缀
		originalKey := key
		if a.prefix != "" {
			keyPrefix := a.prefix + ":"
			if len(key) > len(keyPrefix) && key[:len(keyPrefix)] == keyPrefix {
				originalKey = key[len(keyPrefix):]
			}
		}

		// 简单的字符串匹配（实际应该支持通配符）
		if pattern == "*" || originalKey == pattern {
			keys = append(keys, originalKey)
		}
	}

	return keys, nil
}

// Scan 扫描匹配模式的键
func (a *adapter) Scan(ctx context.Context, pattern string, count int) ([]string, error) {
	// 内存缓存中，Scan和Keys实现相同
	return a.Keys(ctx, pattern)
}

// Clear 清除匹配模式的所有键
func (a *adapter) Clear(ctx context.Context, pattern string) error {
	keys, err := a.Keys(ctx, pattern)
	if err != nil {
		return err
	}

	return a.MDelete(ctx, keys...)
}

// Close 关闭缓存连接
func (a *adapter) Close() error {
	if a.cleanupTicker != nil {
		a.cleanupTicker.Stop()
		close(a.stopCleanup)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.items = make(map[string]*item)
	a.currentSize = 0

	return nil
}

// Ping 检查连接是否正常
func (a *adapter) Ping(ctx context.Context) error {
	return a.HealthCheck(ctx)
}
