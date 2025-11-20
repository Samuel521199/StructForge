package router

import (
	"context"
	"fmt"
	"io"
	stdHttp "net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"StructForge/backend/apps/gateway/internal/conf"
	circuitbreaker "StructForge/backend/apps/gateway/internal/middleware/circuitbreaker"
	"StructForge/backend/apps/gateway/internal/router/discovery"
	"StructForge/backend/apps/gateway/internal/router/loadbalancer"
	"StructForge/backend/common/log"

	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// Route 路由规则
type Route struct {
	// 路径匹配规则（支持前缀匹配和精确匹配）
	Path string `yaml:"path" json:"path"`
	// 路径匹配类型：prefix（前缀匹配）、exact（精确匹配）、regex（正则匹配）
	MatchType string `yaml:"match_type" json:"match_type"`
	// 目标服务名称
	Service string `yaml:"service" json:"service"`
	// 目标服务路径（可选，如果不指定则使用原始路径）
	TargetPath string `yaml:"target_path" json:"target_path"`
	// 是否需要认证
	RequireAuth bool `yaml:"require_auth" json:"require_auth"`
	// 超时时间（秒）
	Timeout int `yaml:"timeout" json:"timeout"`
	// 重试次数
	Retries int `yaml:"retries" json:"retries"`
	// 负载均衡策略：round_robin, random, least_connections
	LoadBalanceStrategy string `yaml:"load_balance_strategy" json:"load_balance_strategy"`
	// 限流配置
	RateLimit *RateLimitConfig `yaml:"rate_limit" json:"rate_limit"`
	// 熔断器配置
	CircuitBreaker *conf.CircuitBreakerConfig `yaml:"circuit_breaker" json:"circuit_breaker"`
	// 缓存配置
	Cache *conf.CacheConfig `yaml:"cache" json:"cache"`
}

// CircuitBreakerConfig 熔断器配置（与 conf.CircuitBreakerConfig 相同，避免循环依赖）
type CircuitBreakerConfig struct {
	// 是否启用熔断器
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 失败率阈值（0-1），超过此值将打开熔断器
	FailureThreshold float64 `yaml:"failure_threshold" json:"failure_threshold"`
	// 最小请求数，达到此数量后才开始计算失败率
	MinRequests int `yaml:"min_requests" json:"min_requests"`
	// 时间窗口（秒），在此时间窗口内统计失败率
	WindowSize int `yaml:"window_size" json:"window_size"`
	// 打开状态持续时间（秒），熔断器打开后等待此时间后进入半开状态
	OpenDuration int `yaml:"open_duration" json:"open_duration"`
	// 半开状态允许的请求数
	HalfOpenRequests int `yaml:"half_open_requests" json:"half_open_requests"`
	// 超时时间（秒），请求超过此时间视为失败
	Timeout int `yaml:"timeout" json:"timeout"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// 每秒允许的请求数
	QPS int `yaml:"qps" json:"qps"`
	// 突发请求数
	Burst int `yaml:"burst" json:"burst"`
}

// Router 路由管理器
type Router struct {
	routes          []*Route
	discovery       discovery.ServiceDiscovery
	loadBalancers   map[string]loadbalancer.LoadBalancer
	circuitBreakers *circuitbreaker.CircuitBreakerManager
	httpClient      *stdHttp.Client
	mu              sync.RWMutex
}

// NewRouter 创建路由管理器
func NewRouter(discovery discovery.ServiceDiscovery) *Router {
	return &Router{
		routes:          make([]*Route, 0),
		discovery:       discovery,
		loadBalancers:   make(map[string]loadbalancer.LoadBalancer),
		circuitBreakers: circuitbreaker.NewCircuitBreakerManager(),
		httpClient: &stdHttp.Client{
			Timeout: 30 * time.Second,
			Transport: &stdHttp.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// AddRoute 添加路由规则
func (r *Router) AddRoute(route *Route) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 设置默认值
	if route.MatchType == "" {
		route.MatchType = "prefix"
	}
	if route.Timeout == 0 {
		route.Timeout = 30
	}
	if route.Retries == 0 {
		route.Retries = 0
	}
	if route.LoadBalanceStrategy == "" {
		route.LoadBalanceStrategy = "round_robin"
	}

	r.routes = append(r.routes, route)

	// 为每个服务创建负载均衡器
	if _, exists := r.loadBalancers[route.Service]; !exists {
		r.loadBalancers[route.Service] = loadbalancer.NewLoadBalancer(route.LoadBalanceStrategy)
	}

	log.Info(context.Background(), "路由规则已添加",
		log.String("path", route.Path),
		log.String("service", route.Service),
		log.String("match_type", route.MatchType),
	)
}

// AddRoutes 批量添加路由规则
func (r *Router) AddRoutes(routes []*Route) {
	for _, route := range routes {
		r.AddRoute(route)
	}
}

// FindRoute 查找匹配的路由
func (r *Router) FindRoute(path string) *Route {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 按添加顺序匹配，第一个匹配的路由生效
	for _, route := range r.routes {
		if r.matchPath(path, route) {
			return route
		}
	}

	return nil
}

// matchPath 匹配路径
func (r *Router) matchPath(path string, route *Route) bool {
	switch route.MatchType {
	case "exact":
		return path == route.Path
	case "prefix":
		return strings.HasPrefix(path, route.Path)
	case "regex":
		// 实现正则匹配
		return r.matchRegex(path, route.Path)
	default:
		return strings.HasPrefix(path, route.Path)
	}
}

// Forward 转发请求到目标服务
// CacheCallback 缓存回调函数类型
type CacheCallback func(statusCode int, headers map[string][]string, body []byte)

func (r *Router) Forward(ctx kratosHttp.Context, route *Route, cacheCallback CacheCallback) error {
	// 从 HTTP 请求中获取标准 context
	requestCtx := ctx.Request().Context()

	// 获取服务实例
	instances, err := r.discovery.GetInstances(requestCtx, route.Service)
	if err != nil {
		log.Error(ctx, "获取服务实例失败",
			log.ErrorField(err),
			log.String("service", route.Service),
		)
		return fmt.Errorf("服务不可用: %s", route.Service)
	}

	if len(instances) == 0 {
		log.Error(ctx, "服务实例为空",
			log.String("service", route.Service),
		)
		return fmt.Errorf("服务 %s 没有可用实例", route.Service)
	}

	// 使用负载均衡选择实例
	lb := r.loadBalancers[route.Service]
	instance := lb.Select(instances)

	if instance == nil {
		return fmt.Errorf("无法选择服务实例: %s", route.Service)
	}

	// 构建目标URL
	targetPath := route.TargetPath
	if targetPath == "" {
		// 如果未指定目标路径，使用原始路径
		// 对于前缀匹配，需要保留完整路径；对于精确匹配，直接使用路径
		if route.MatchType == "exact" {
			targetPath = route.Path
		} else {
			// 前缀匹配：保留完整路径
			targetPath = ctx.Request().URL.Path
		}
		// 确保路径以 / 开头
		if !strings.HasPrefix(targetPath, "/") {
			targetPath = "/" + targetPath
		}
	}

	targetURL := fmt.Sprintf("http://%s:%d%s", instance.Host, instance.Port, targetPath)
	if ctx.Request().URL.RawQuery != "" {
		targetURL += "?" + ctx.Request().URL.RawQuery
	}

	log.Info(ctx, "转发请求",
		log.String("from", ctx.Request().URL.Path),
		log.String("to", targetURL),
		log.String("service", route.Service),
		log.String("instance", fmt.Sprintf("%s:%d", instance.Host, instance.Port)),
	)

	// 设置超时
	timeout := time.Duration(route.Timeout) * time.Second
	if timeout > 0 {
		var cancel context.CancelFunc
		requestCtx, cancel = context.WithTimeout(requestCtx, timeout)
		defer cancel()
	}

	// 创建代理请求
	req, err := stdHttp.NewRequestWithContext(requestCtx, ctx.Request().Method, targetURL, ctx.Request().Body)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 复制请求头
	for key, values := range ctx.Request().Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 检查熔断器配置
	var cbConfig *circuitbreaker.Config
	if route.CircuitBreaker != nil && route.CircuitBreaker.Enabled {
		cbConfig = &circuitbreaker.Config{
			FailureThreshold: route.CircuitBreaker.FailureThreshold,
			MinRequests:      route.CircuitBreaker.MinRequests,
			WindowSize:       route.CircuitBreaker.WindowSize,
			OpenDuration:     route.CircuitBreaker.OpenDuration,
			HalfOpenRequests: route.CircuitBreaker.HalfOpenRequests,
			Timeout:          route.CircuitBreaker.Timeout,
		}
		// 设置默认值
		if cbConfig.FailureThreshold == 0 {
			cbConfig.FailureThreshold = 0.5
		}
		if cbConfig.MinRequests == 0 {
			cbConfig.MinRequests = 10
		}
		if cbConfig.WindowSize == 0 {
			cbConfig.WindowSize = 60
		}
		if cbConfig.OpenDuration == 0 {
			cbConfig.OpenDuration = 30
		}
		if cbConfig.HalfOpenRequests == 0 {
			cbConfig.HalfOpenRequests = 3
		}
		if cbConfig.Timeout == 0 {
			cbConfig.Timeout = route.Timeout
			if cbConfig.Timeout == 0 {
				cbConfig.Timeout = 30
			}
		}
	}

	// 发送请求（支持重试和熔断保护）
	var resp *stdHttp.Response
	var httpErr error

	maxRetries := route.Retries
	if maxRetries < 0 {
		maxRetries = 0
	}

	// 执行请求的函数
	executeRequest := func() error {
		var attemptErr error
		for attempt := 0; attempt <= maxRetries; attempt++ {
			if attempt > 0 {
				// 重试前等待（指数退避）
				backoff := time.Duration(attempt) * 100 * time.Millisecond
				if backoff > 1*time.Second {
					backoff = 1 * time.Second
				}
				time.Sleep(backoff)

				log.Info(ctx, "重试请求",
					log.Int("attempt", attempt),
					log.Int("max_retries", maxRetries),
					log.String("target", targetURL),
				)
			}

			// 发送请求
			resp, attemptErr = r.httpClient.Do(req)
			if attemptErr == nil {
				// 请求成功，检查状态码
				if resp.StatusCode < 500 {
					// 4xx 错误不重试（客户端错误），但视为失败
					if resp.StatusCode >= 400 {
						return fmt.Errorf("client error: %d", resp.StatusCode)
					}
					// 2xx/3xx 成功
					return nil
				}
				// 5xx 错误需要重试
				if attempt < maxRetries {
					resp.Body.Close()
					resp = nil
					continue
				}
				// 最后一次尝试仍然失败
				return fmt.Errorf("server error: %d", resp.StatusCode)
			} else {
				// 网络错误，需要重试
				if attempt < maxRetries {
					// 检查是否是可重试的错误
					if isRetryableError(attemptErr) {
						continue
					}
					// 不可重试的错误，直接返回
					return attemptErr
				}
				// 最后一次尝试仍然失败
				return attemptErr
			}
		}
		return attemptErr
	}

	// 使用熔断器执行请求（如果启用）
	if cbConfig != nil {
		httpErr = r.circuitBreakers.Execute(requestCtx, route.Service, cbConfig, executeRequest)
		if httpErr != nil {
			// 检查是否是熔断器打开错误
			if circuitbreaker.IsCircuitBreakerError(httpErr) {
				log.Warn(ctx, "熔断器已打开，拒绝请求",
					log.String("service", route.Service),
					log.String("target", targetURL),
				)
				return fmt.Errorf("服务暂时不可用（熔断器已打开）")
			}
			// 其他错误
			log.Error(ctx, "转发请求失败",
				log.ErrorField(httpErr),
				log.String("target", targetURL),
				log.Int("retries", maxRetries),
			)
			return fmt.Errorf("转发请求失败: %w", httpErr)
		}
	} else {
		// 不使用熔断器，直接执行
		httpErr = executeRequest()
		if httpErr != nil {
			log.Error(ctx, "转发请求失败",
				log.ErrorField(httpErr),
				log.String("target", targetURL),
				log.Int("retries", maxRetries),
			)
			return fmt.Errorf("转发请求失败: %w", httpErr)
		}
	}

	if resp == nil {
		return fmt.Errorf("响应为空")
	}

	defer resp.Body.Close()

	// 读取响应体（用于缓存）
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %w", err)
	}

	// 复制响应头
	responseHeaders := make(map[string][]string)
	for key, values := range resp.Header {
		responseHeaders[key] = values
		for _, value := range values {
			ctx.Response().Header().Set(key, value)
		}
	}

	// 设置状态码
	ctx.Response().WriteHeader(resp.StatusCode)

	// 写入响应体
	_, err = ctx.Response().Write(responseBody)
	if err != nil {
		return fmt.Errorf("写入响应失败: %w", err)
	}

	// 如果提供了缓存回调，调用它
	if cacheCallback != nil {
		cacheCallback(resp.StatusCode, responseHeaders, responseBody)
	}

	return nil
}

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// 网络超时、连接错误等可重试
	retryableErrors := []string{
		"timeout",
		"deadline exceeded",
		"connection refused",
		"connection reset",
		"no such host",
		"network is unreachable",
		"temporary failure",
	}

	for _, retryableErr := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), retryableErr) {
			return true
		}
	}

	return false
}

// regexCache 正则表达式缓存（避免重复编译）
var (
	regexCache = make(map[string]*regexp.Regexp)
	regexMu    sync.RWMutex
)

// matchRegex 正则匹配
func (r *Router) matchRegex(path, pattern string) bool {
	// 从缓存获取或编译正则表达式
	regexMu.RLock()
	re, exists := regexCache[pattern]
	regexMu.RUnlock()

	if !exists {
		// 编译正则表达式
		var err error
		re, err = regexp.Compile(pattern)
		if err != nil {
			// 正则表达式编译失败，记录日志并返回false
			log.Warn(context.Background(), "正则表达式编译失败",
				log.String("pattern", pattern),
				log.ErrorField(err),
			)
			return false
		}

		// 存入缓存
		regexMu.Lock()
		regexCache[pattern] = re
		regexMu.Unlock()
	}

	return re.MatchString(path)
}

// GetServiceInstances 获取服务实例（用于健康检查）
func (r *Router) GetServiceInstances(ctx context.Context, serviceName string) ([]discovery.Instance, error) {
	return r.discovery.GetInstances(ctx, serviceName)
}

// GetAllServiceNames 获取所有已注册的服务名称
func (r *Router) GetAllServiceNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	serviceNames := make(map[string]bool)
	for _, route := range r.routes {
		serviceNames[route.Service] = true
	}

	names := make([]string, 0, len(serviceNames))
	for name := range serviceNames {
		names = append(names, name)
	}

	return names
}

// UpdateServiceInstances 更新服务实例（由服务发现调用）
func (r *Router) UpdateServiceInstances(serviceName string, instances []discovery.Instance) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if lb, exists := r.loadBalancers[serviceName]; exists {
		lb.UpdateInstances(instances)
		log.Info(context.Background(), "服务实例已更新",
			log.String("service", serviceName),
			log.Int("instances", len(instances)),
		)
	}
}

// GetCircuitBreakerStats 获取所有熔断器的统计信息
func (r *Router) GetCircuitBreakerStats() map[string]map[string]interface{} {
	return r.circuitBreakers.GetBreakerStats()
}
