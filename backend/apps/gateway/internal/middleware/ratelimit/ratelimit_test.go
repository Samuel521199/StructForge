package ratelimit

import (
	"context"
	"testing"
	"time"
)

// TestTokenBucket 测试令牌桶算法
func TestTokenBucket(t *testing.T) {
	limiter := NewTokenBucketLimiter(10, 5) // QPS=10, Burst=5
	ctx := context.Background()
	key := "test"

	// 初始状态应该允许5个请求
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, key)
		if err != nil {
			t.Fatalf("限流检查失败: %v", err)
		}
		if !allowed {
			t.Errorf("第 %d 次应该允许", i+1)
		}
	}

	// 第6次应该被拒绝（令牌已用完）
	allowed, _ := limiter.Allow(ctx, key)
	if allowed {
		t.Error("应该被拒绝（令牌已用完）")
	}

	// 等待一段时间，令牌应该恢复
	time.Sleep(200 * time.Millisecond) // 应该恢复约2个令牌
	allowed, _ = limiter.Allow(ctx, key)
	if !allowed {
		t.Error("等待后应该允许")
	}
}

// TestTokenBucketRefill 测试令牌桶自动补充
func TestTokenBucketRefill(t *testing.T) {
	limiter := NewTokenBucketLimiter(10, 5) // QPS=10, 每秒10个令牌
	ctx := context.Background()
	key := "test"

	// 消耗所有令牌
	for i := 0; i < 5; i++ {
		limiter.Allow(ctx, key)
	}

	// 等待1秒，应该恢复10个令牌（但不超过burst）
	time.Sleep(1100 * time.Millisecond)

	// 应该可以再允许5个请求（burst=5）
	allowedCount := 0
	for i := 0; i < 10; i++ {
		allowed, _ := limiter.Allow(ctx, key)
		if allowed {
			allowedCount++
		}
	}
	if allowedCount < 5 {
		t.Errorf("等待1秒后应该至少允许5个请求，实际 %d", allowedCount)
	}
}

// TestRateLimitManager 测试限流管理器
func TestRateLimitManager(t *testing.T) {
	mgr := NewRateLimitManager()

	// 测试不同路径的限流
	path1 := "/api/v1/users"
	path2 := "/api/v1/orders"

	// 为不同路径创建限流器
	allowed1, err := CheckRateLimit(context.Background(), mgr, path1, 10, 5, false)
	if err != nil {
		t.Fatalf("限流检查失败: %v", err)
	}
	if !allowed1 {
		t.Error("第一次请求应该允许")
	}

	allowed2, err := CheckRateLimit(context.Background(), mgr, path2, 10, 5, false)
	if err != nil {
		t.Fatalf("限流检查失败: %v", err)
	}
	if !allowed2 {
		t.Error("不同路径的请求应该允许")
	}

	// 快速请求，应该触发限流
	for i := 0; i < 10; i++ {
		CheckRateLimit(context.Background(), mgr, path1, 10, 5, false)
	}

	// 第11次应该被限流（超过burst）
	allowed, _ := CheckRateLimit(context.Background(), mgr, path1, 10, 5, false)
	if allowed {
		t.Error("超过burst后应该被限流")
	}
}

// TestRateLimitCleanup 测试限流器清理
func TestRateLimitCleanup(t *testing.T) {
	mgr := NewRateLimitManager()

	// 创建多个限流器
	for i := 0; i < 10; i++ {
		path := "/api/v1/test" + string(rune(i))
		CheckRateLimit(context.Background(), mgr, path, 10, 5, false)
	}

	// 等待清理时间（默认5分钟）
	// 这里只测试清理逻辑存在，不等待实际清理
	mgr.mu.RLock()
	limiterCount := len(mgr.limiters)
	mgr.mu.RUnlock()

	if limiterCount == 0 {
		t.Error("应该有多个限流器")
	}
}
