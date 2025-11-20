package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"StructForge/backend/common/log"
)

// Limiter 限流器接口
type Limiter interface {
	// Allow 检查是否允许请求
	Allow(ctx context.Context, key string) (bool, error)
	// Reset 重置限流器
	Reset(key string)
}

// TokenBucketLimiter Token Bucket 限流器
type TokenBucketLimiter struct {
	// QPS: 每秒允许的请求数
	qps int
	// Burst: 突发请求数（桶容量）
	burst int
	// buckets: 每个key的令牌桶
	buckets map[string]*tokenBucket
	mu      sync.RWMutex
	// 清理间隔
	cleanupInterval time.Duration
	// 最后清理时间
	lastCleanup time.Time
}

// tokenBucket 令牌桶
type tokenBucket struct {
	// 当前令牌数
	tokens float64
	// 最后更新时间
	lastUpdate time.Time
	// 每秒添加的令牌数（等于QPS）
	rate float64
	// 桶容量（等于Burst）
	capacity float64
	mu       sync.Mutex
}

// NewTokenBucketLimiter 创建 Token Bucket 限流器
func NewTokenBucketLimiter(qps, burst int) *TokenBucketLimiter {
	if qps <= 0 {
		qps = 100 // 默认值
	}
	if burst <= 0 {
		burst = 200 // 默认值
	}

	limiter := &TokenBucketLimiter{
		qps:             qps,
		burst:           burst,
		buckets:         make(map[string]*tokenBucket),
		cleanupInterval: 5 * time.Minute, // 5分钟清理一次
		lastCleanup:     time.Now(),
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter
}

// Allow 检查是否允许请求
func (l *TokenBucketLimiter) Allow(ctx context.Context, key string) (bool, error) {
	l.mu.RLock()
	bucket, exists := l.buckets[key]
	l.mu.RUnlock()

	now := time.Now()

	// 如果不存在，创建新的令牌桶
	if !exists {
		l.mu.Lock()
		// 双重检查
		if bucket, exists = l.buckets[key]; !exists {
			bucket = &tokenBucket{
				tokens:     float64(l.burst), // 初始时桶是满的
				lastUpdate: now,
				rate:       float64(l.qps),
				capacity:   float64(l.burst),
			}
			l.buckets[key] = bucket
		}
		l.mu.Unlock()
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	// 计算应该添加的令牌数（基于时间差）
	elapsed := now.Sub(bucket.lastUpdate).Seconds()
	tokensToAdd := elapsed * bucket.rate

	// 添加令牌（不超过容量）
	bucket.tokens = min(bucket.capacity, bucket.tokens+tokensToAdd)
	bucket.lastUpdate = now

	// 检查是否有足够的令牌
	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true, nil
	}

	// 被限流
	return false, nil
}

// Reset 重置限流器（清除指定key的令牌桶）
func (l *TokenBucketLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// cleanup 定期清理不活跃的令牌桶
func (l *TokenBucketLimiter) cleanup() {
	ticker := time.NewTicker(l.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		// 清理超过10分钟未使用的桶
		for key, bucket := range l.buckets {
			bucket.mu.Lock()
			if now.Sub(bucket.lastUpdate) > 10*time.Minute {
				delete(l.buckets, key)
			}
			bucket.mu.Unlock()
		}
		l.mu.Unlock()
	}
}

// min 返回两个浮点数中的较小值
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// RateLimitManager 限流管理器（管理多个限流器）
type RateLimitManager struct {
	limiters map[string]Limiter // key: route path, value: limiter
	mu       sync.RWMutex
}

// NewRateLimitManager 创建限流管理器
func NewRateLimitManager() *RateLimitManager {
	return &RateLimitManager{
		limiters: make(map[string]Limiter),
	}
}

// GetLimiter 获取或创建限流器
func (m *RateLimitManager) GetLimiter(path string, qps, burst int) Limiter {
	// 生成限流器key
	key := fmt.Sprintf("%s:qps:%d:burst:%d", path, qps, burst)

	m.mu.RLock()
	limiter, exists := m.limiters[key]
	m.mu.RUnlock()

	if exists {
		return limiter
	}

	// 创建新的限流器
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if limiter, exists = m.limiters[key]; !exists {
		limiter = NewTokenBucketLimiter(qps, burst)
		m.limiters[key] = limiter
	}

	return limiter
}

// ExtractKey 提取限流key（支持IP、用户、API级别）
func ExtractKey(ctx context.Context, path string, requireAuth bool) string {
	// 默认使用路径作为key（API级别限流）
	key := path

	// 如果支持IP级别限流，可以从请求中提取IP
	// 如果支持用户级别限流，可以从JWT中提取用户ID
	// 这里先实现API级别，后续可以扩展

	return key
}

// CheckRateLimit 检查限流（在handler中调用）
func CheckRateLimit(ctx context.Context, manager *RateLimitManager, path string, qps, burst int, requireAuth bool) (bool, error) {
	if qps <= 0 || burst <= 0 {
		// 如果没有配置限流，直接通过
		return true, nil
	}

	// 提取限流key
	key := ExtractKey(ctx, path, requireAuth)

	// 获取限流器
	limiter := manager.GetLimiter(path, qps, burst)

	// 检查是否允许
	allowed, err := limiter.Allow(ctx, key)
	if err != nil {
		log.Error(ctx, "限流检查失败",
			log.ErrorField(err),
			log.String("path", path),
		)
		// 出错时允许通过，避免影响正常请求
		return true, nil
	}

	if !allowed {
		log.Warn(ctx, "请求被限流",
			log.String("path", path),
			log.String("key", key),
			log.Int("qps", qps),
			log.Int("burst", burst),
		)
		return false, fmt.Errorf("请求过于频繁，请稍后再试")
	}

	return true, nil
}
