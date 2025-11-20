package cache

import (
	"context"
	"net/http"
	"testing"
	"time"

	"StructForge/backend/common/cache"
)

// mockCache 模拟缓存实现
type mockCache struct {
	data map[string][]byte
	ttl  map[string]time.Time
}

func newMockCache() *mockCache {
	return &mockCache{
		data: make(map[string][]byte),
		ttl:  make(map[string]time.Time),
	}
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	if data, ok := m.data[key]; ok {
		if ttl, ok := m.ttl[key]; ok && time.Now().Before(ttl) {
			return data, nil
		}
		delete(m.data, key)
		delete(m.ttl, key)
	}
	return nil, cache.ErrNotFound
}

func (m *mockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.data[key] = value
	m.ttl[key] = time.Now().Add(ttl)
	return nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	delete(m.ttl, key)
	return nil
}

func (m *mockCache) Clear(ctx context.Context, pattern string) error {
	// 简单实现：清除所有匹配的键
	if pattern == "*" {
		m.data = make(map[string][]byte)
		m.ttl = make(map[string]time.Time)
	}
	return nil
}

func (m *mockCache) Exists(ctx context.Context, key string) (bool, error) {
	_, ok := m.data[key]
	return ok, nil
}

func (m *mockCache) MGet(ctx context.Context, keys ...string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	for _, key := range keys {
		if data, ok := m.data[key]; ok {
			if ttl, ok := m.ttl[key]; ok && time.Now().Before(ttl) {
				result[key] = data
			}
		}
	}
	return result, nil
}

func (m *mockCache) MSet(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	for key, value := range items {
		m.data[key] = value
		m.ttl[key] = time.Now().Add(ttl)
	}
	return nil
}

func (m *mockCache) MDelete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(m.data, key)
		delete(m.ttl, key)
	}
	return nil
}

func (m *mockCache) GetOrSet(ctx context.Context, key string, fn func() ([]byte, error), ttl time.Duration) ([]byte, error) {
	if data, err := m.Get(ctx, key); err == nil && data != nil {
		return data, nil
	}
	data, err := fn()
	if err != nil {
		return nil, err
	}
	return data, m.Set(ctx, key, data, ttl)
}

func (m *mockCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	// 简化实现
	return 0, nil
}

func (m *mockCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	// 简化实现
	return 0, nil
}

func (m *mockCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if _, ok := m.data[key]; ok {
		m.ttl[key] = time.Now().Add(ttl)
	}
	return nil
}

func (m *mockCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	if ttl, ok := m.ttl[key]; ok {
		remaining := time.Until(ttl)
		if remaining < 0 {
			return -2, nil // 键不存在
		}
		return remaining, nil
	}
	return -2, nil // 键不存在
}

func (m *mockCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys, nil
}

func (m *mockCache) Scan(ctx context.Context, pattern string, count int) ([]string, error) {
	return m.Keys(ctx, pattern)
}

func (m *mockCache) Close() error {
	return nil
}

func (m *mockCache) Ping(ctx context.Context) error {
	return nil
}

// TestShouldCache 测试缓存判断逻辑
func TestShouldCache(t *testing.T) {
	// 设置全局缓存
	originalCache := cache.GetGlobalCache()
	defer cache.SetGlobalCache(originalCache)

	mock := newMockCache()
	cache.SetGlobalCache(mock)

	config := DefaultCacheConfig()
	config.Methods = []string{"GET"}
	config.Paths = []string{"/api/v1/users*"}

	middleware, err := NewCacheMiddleware(config)
	if err != nil {
		t.Fatalf("创建缓存中间件失败: %v", err)
	}

	testCases := []struct {
		method string
		path   string
		should bool
	}{
		{"GET", "/api/v1/users", true},
		{"GET", "/api/v1/users/123", true},
		{"POST", "/api/v1/users", false},
		{"GET", "/api/v1/orders", false},
	}

	for _, tc := range testCases {
		result := middleware.ShouldCache(tc.method, tc.path)
		if result != tc.should {
			t.Errorf("方法 %s, 路径 %s: 期望 %v, 实际 %v", tc.method, tc.path, tc.should, result)
		}
	}
}

// TestCacheKeyGeneration 测试缓存键生成
func TestCacheKeyGeneration(t *testing.T) {
	originalCache := cache.GetGlobalCache()
	defer cache.SetGlobalCache(originalCache)

	mock := newMockCache()
	cache.SetGlobalCache(mock)

	config := DefaultCacheConfig()
	config.IncludeQueryParams = true

	middleware, err := NewCacheMiddleware(config)
	if err != nil {
		t.Fatalf("创建缓存中间件失败: %v", err)
	}

	// 测试相同路径和查询参数生成相同键
	query1 := map[string][]string{"id": {"123"}}
	query2 := map[string][]string{"id": {"123"}}
	headers := map[string][]string{}

	key1 := middleware.GenerateCacheKey("GET", "/api/v1/users", query1, headers)
	key2 := middleware.GenerateCacheKey("GET", "/api/v1/users", query2, headers)

	if key1 != key2 {
		t.Error("相同参数应该生成相同的缓存键")
	}

	// 测试不同查询参数生成不同键
	query3 := map[string][]string{"id": {"456"}}
	key3 := middleware.GenerateCacheKey("GET", "/api/v1/users", query3, headers)

	if key1 == key3 {
		t.Error("不同参数应该生成不同的缓存键")
	}
}

// TestCacheHandler 测试缓存处理器
func TestCacheHandler(t *testing.T) {
	originalCache := cache.GetGlobalCache()
	defer cache.SetGlobalCache(originalCache)

	mock := newMockCache()
	cache.SetGlobalCache(mock)

	config := DefaultCacheConfig()
	middleware, err := NewCacheMiddleware(config)
	if err != nil {
		t.Fatalf("创建缓存中间件失败: %v", err)
	}

	handler := NewCacheHandler(middleware)

	// 创建测试请求
	req, _ := http.NewRequest("GET", "/api/v1/users?id=123", nil)
	ctx := context.Background()

	// 第一次请求，应该未命中
	cachedResp, hit := handler.HandleRequest(ctx, req)
	if hit {
		t.Error("第一次请求不应该命中缓存")
	}

	// 写入缓存
	headers := map[string][]string{"Content-Type": {"application/json"}}
	body := []byte(`{"id": 123, "name": "test"}`)
	handler.HandleResponse(ctx, req, 200, headers, body)

	// 第二次请求，应该命中
	cachedResp, hit = handler.HandleRequest(ctx, req)
	if !hit {
		t.Error("第二次请求应该命中缓存")
	}
	if cachedResp == nil {
		t.Fatal("缓存响应不应该为空")
	}
	if cachedResp.StatusCode != 200 {
		t.Errorf("期望状态码 200，实际 %d", cachedResp.StatusCode)
	}
	if string(cachedResp.Body) != string(body) {
		t.Error("缓存响应体不匹配")
	}
}

// TestCacheExpiration 测试缓存过期
func TestCacheExpiration(t *testing.T) {
	originalCache := cache.GetGlobalCache()
	defer cache.SetGlobalCache(originalCache)

	mock := newMockCache()
	cache.SetGlobalCache(mock)

	config := DefaultCacheConfig()
	config.TTL = 1 // 1秒过期
	middleware, err := NewCacheMiddleware(config)
	if err != nil {
		t.Fatalf("创建缓存中间件失败: %v", err)
	}

	handler := NewCacheHandler(middleware)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	ctx := context.Background()

	// 写入缓存
	handler.HandleResponse(ctx, req, 200, map[string][]string{}, []byte("test"))

	// 立即读取，应该命中
	_, hit := handler.HandleRequest(ctx, req)
	if !hit {
		t.Error("应该命中缓存")
	}

	// 等待过期
	time.Sleep(1100 * time.Millisecond)

	// 再次读取，应该未命中
	_, hit = handler.HandleRequest(ctx, req)
	if hit {
		t.Error("缓存过期后不应该命中")
	}
}

// TestCacheOnlySuccess 测试只缓存成功响应
func TestCacheOnlySuccess(t *testing.T) {
	originalCache := cache.GetGlobalCache()
	defer cache.SetGlobalCache(originalCache)

	mock := newMockCache()
	cache.SetGlobalCache(mock)

	config := DefaultCacheConfig()
	middleware, err := NewCacheMiddleware(config)
	if err != nil {
		t.Fatalf("创建缓存中间件失败: %v", err)
	}

	handler := NewCacheHandler(middleware)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	ctx := context.Background()

	// 写入错误响应（4xx），不应该缓存
	handler.HandleResponse(ctx, req, 404, map[string][]string{}, []byte("not found"))

	// 读取，应该未命中
	_, hit := handler.HandleRequest(ctx, req)
	if hit {
		t.Error("错误响应不应该被缓存")
	}

	// 写入成功响应（2xx），应该缓存
	handler.HandleResponse(ctx, req, 200, map[string][]string{}, []byte("success"))

	// 读取，应该命中
	_, hit = handler.HandleRequest(ctx, req)
	if !hit {
		t.Error("成功响应应该被缓存")
	}
}
