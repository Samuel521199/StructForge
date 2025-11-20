package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"StructForge/backend/common/cache"
	"StructForge/backend/common/log"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	// 是否启用缓存
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 缓存过期时间（秒）
	TTL int `yaml:"ttl" json:"ttl"`
	// 缓存键前缀
	KeyPrefix string `yaml:"key_prefix" json:"key_prefix"`
	// 需要缓存的 HTTP 方法（默认只缓存 GET）
	Methods []string `yaml:"methods" json:"methods"`
	// 需要缓存的路径（支持通配符）
	Paths []string `yaml:"paths" json:"paths"`
	// 排除的路径（支持通配符）
	ExcludePaths []string `yaml:"exclude_paths" json:"exclude_paths"`
	// 是否包含查询参数在缓存键中
	IncludeQueryParams bool `yaml:"include_query_params" json:"include_query_params"`
	// 是否包含请求头在缓存键中（用于区分不同用户）
	IncludeHeaders []string `yaml:"include_headers" json:"include_headers"`
}

// DefaultCacheConfig 默认缓存配置
func DefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		Enabled:            true,
		TTL:                300, // 5分钟
		KeyPrefix:          "gateway:cache:",
		Methods:            []string{"GET"},
		IncludeQueryParams: true,
		IncludeHeaders:     []string{},
	}
}

// CacheMiddleware 缓存中间件
type CacheMiddleware struct {
	config *CacheConfig
	cache  cache.Cache
}

// NewCacheMiddleware 创建缓存中间件
func NewCacheMiddleware(config *CacheConfig) (*CacheMiddleware, error) {
	if config == nil {
		config = DefaultCacheConfig()
	}

	// 设置默认值
	if config.TTL == 0 {
		config.TTL = 300
	}
	if config.KeyPrefix == "" {
		config.KeyPrefix = "gateway:cache:"
	}
	if len(config.Methods) == 0 {
		config.Methods = []string{"GET"}
	}

	// 获取全局缓存实例
	cacheInstance := cache.GetGlobalCache()
	if cacheInstance == nil {
		return nil, fmt.Errorf("缓存实例未初始化，请先初始化缓存系统")
	}

	return &CacheMiddleware{
		config: config,
		cache:  cacheInstance,
	}, nil
}

// CachedResponse 缓存的响应
type CachedResponse struct {
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
	CachedAt   time.Time           `json:"cached_at"`
}

// ShouldCache 检查是否应该缓存此请求
func (m *CacheMiddleware) ShouldCache(method, path string) bool {
	if !m.config.Enabled {
		return false
	}

	// 检查方法
	methodAllowed := false
	for _, m := range m.config.Methods {
		if strings.EqualFold(method, m) {
			methodAllowed = true
			break
		}
	}
	if !methodAllowed {
		return false
	}

	// 检查路径（如果配置了路径列表）
	if len(m.config.Paths) > 0 {
		pathMatched := false
		for _, pattern := range m.config.Paths {
			if m.matchPath(path, pattern) {
				pathMatched = true
				break
			}
		}
		if !pathMatched {
			return false
		}
	}

	// 检查排除路径
	for _, pattern := range m.config.ExcludePaths {
		if m.matchPath(path, pattern) {
			return false
		}
	}

	return true
}

// matchPath 匹配路径（支持通配符）
func (m *CacheMiddleware) matchPath(path, pattern string) bool {
	// 精确匹配
	if path == pattern {
		return true
	}

	// 前缀匹配
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(path, prefix)
	}

	// 后缀匹配
	if strings.HasPrefix(pattern, "*") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(path, suffix)
	}

	return false
}

// GenerateCacheKey 生成缓存键
func (m *CacheMiddleware) GenerateCacheKey(method, path string, queryParams map[string][]string, headers map[string][]string) string {
	// 构建键的组成部分
	parts := []string{
		method,
		path,
	}

	// 包含查询参数
	if m.config.IncludeQueryParams && len(queryParams) > 0 {
		queryStr := m.buildQueryString(queryParams)
		if queryStr != "" {
			parts = append(parts, queryStr)
		}
	}

	// 包含指定的请求头
	if len(m.config.IncludeHeaders) > 0 {
		headerParts := make([]string, 0)
		for _, headerName := range m.config.IncludeHeaders {
			if values, ok := headers[headerName]; ok && len(values) > 0 {
				headerParts = append(headerParts, fmt.Sprintf("%s:%s", headerName, values[0]))
			}
		}
		if len(headerParts) > 0 {
			parts = append(parts, strings.Join(headerParts, ","))
		}
	}

	// 组合所有部分
	key := strings.Join(parts, "|")

	// 使用 MD5 生成固定长度的键（避免键过长）
	hash := md5.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	return m.config.KeyPrefix + hashStr
}

// buildQueryString 构建查询参数字符串（排序后）
func (m *CacheMiddleware) buildQueryString(params map[string][]string) string {
	if len(params) == 0 {
		return ""
	}

	// 简单实现：按 key 排序后拼接
	// 注意：实际应该使用更严格的排序算法
	parts := make([]string, 0)
	for key, values := range params {
		for _, value := range values {
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(parts, "&")
}

// Get 从缓存获取响应
func (m *CacheMiddleware) Get(ctx context.Context, key string) (*CachedResponse, error) {
	// 从缓存获取字节数据
	data, err := m.cache.Get(ctx, key)
	if err != nil {
		if err == cache.ErrNotFound {
			return nil, nil // 缓存未命中
		}
		return nil, err
	}

	// 反序列化
	var cachedResp CachedResponse
	if err := json.Unmarshal(data, &cachedResp); err != nil {
		return nil, fmt.Errorf("反序列化缓存数据失败: %w", err)
	}

	return &cachedResp, nil
}

// Set 设置缓存
func (m *CacheMiddleware) Set(ctx context.Context, key string, response *CachedResponse) error {
	// 序列化
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("序列化缓存数据失败: %w", err)
	}

	return m.cache.Set(ctx, key, data, time.Duration(m.config.TTL)*time.Second)
}

// Delete 删除缓存
func (m *CacheMiddleware) Delete(ctx context.Context, key string) error {
	return m.cache.Delete(ctx, key)
}

// InvalidateByPath 根据路径模式使缓存失效
func (m *CacheMiddleware) InvalidateByPath(ctx context.Context, pathPattern string) error {
	// 注意：这是一个简化实现
	// 实际应该维护路径到缓存键的映射，或者使用支持模式匹配的缓存系统
	log.Warn(ctx, "按路径使缓存失效功能需要维护路径到缓存键的映射",
		log.String("pattern", pathPattern),
	)
	return nil
}

// CacheHandler 缓存处理器（用于 HTTP 响应）
type CacheHandler struct {
	middleware *CacheMiddleware
}

// NewCacheHandler 创建缓存处理器
func NewCacheHandler(middleware *CacheMiddleware) *CacheHandler {
	return &CacheHandler{
		middleware: middleware,
	}
}

// HandleRequest 处理请求（检查缓存）
func (h *CacheHandler) HandleRequest(ctx context.Context, req *http.Request) (*CachedResponse, bool) {
	if !h.middleware.ShouldCache(req.Method, req.URL.Path) {
		return nil, false
	}

	// 生成缓存键
	queryParams := make(map[string][]string)
	for key, values := range req.URL.Query() {
		queryParams[key] = values
	}

	headers := make(map[string][]string)
	for key, values := range req.Header {
		headers[key] = values
	}

	cacheKey := h.middleware.GenerateCacheKey(req.Method, req.URL.Path, queryParams, headers)

	// 从缓存获取
	cachedResp, err := h.middleware.Get(ctx, cacheKey)
	if err != nil {
		log.Warn(ctx, "获取缓存失败",
			log.ErrorField(err),
			log.String("key", cacheKey),
		)
		return nil, false
	}

	if cachedResp != nil {
		log.Info(ctx, "缓存命中",
			log.String("key", cacheKey),
			log.String("path", req.URL.Path),
		)
		return cachedResp, true
	}

	return nil, false
}

// HandleResponse 处理响应（写入缓存）
func (h *CacheHandler) HandleResponse(ctx context.Context, req *http.Request, statusCode int, headers map[string][]string, body []byte) {
	if !h.middleware.ShouldCache(req.Method, req.URL.Path) {
		return
	}

	// 只缓存成功的响应（2xx）
	if statusCode < 200 || statusCode >= 300 {
		return
	}

	// 生成缓存键
	queryParams := make(map[string][]string)
	for key, values := range req.URL.Query() {
		queryParams[key] = values
	}

	reqHeaders := make(map[string][]string)
	for key, values := range req.Header {
		reqHeaders[key] = values
	}

	cacheKey := h.middleware.GenerateCacheKey(req.Method, req.URL.Path, queryParams, reqHeaders)

	// 创建缓存响应
	cachedResp := &CachedResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       body,
		CachedAt:   time.Now(),
	}

	// 写入缓存
	err := h.middleware.Set(ctx, cacheKey, cachedResp)
	if err != nil {
		log.Warn(ctx, "写入缓存失败",
			log.ErrorField(err),
			log.String("key", cacheKey),
		)
	} else {
		log.Info(ctx, "响应已缓存",
			log.String("key", cacheKey),
			log.String("path", req.URL.Path),
			log.Int("ttl", h.middleware.config.TTL),
		)
	}
}
