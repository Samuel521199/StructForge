package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	cacheMiddleware "StructForge/backend/apps/gateway/internal/middleware/cache"
	corsMiddleware "StructForge/backend/apps/gateway/internal/middleware/cors"
	jwtMiddleware "StructForge/backend/apps/gateway/internal/middleware/jwt"
	loggingMiddleware "StructForge/backend/apps/gateway/internal/middleware/logging"
	metricsMiddleware "StructForge/backend/apps/gateway/internal/middleware/metrics"
	ratelimit "StructForge/backend/apps/gateway/internal/middleware/ratelimit"
	"StructForge/backend/apps/gateway/internal/router"
	"StructForge/backend/common/log"

	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// GatewayHandler Gateway处理器
type GatewayHandler struct {
	router        *router.Router
	jwtManager    *jwtMiddleware.Manager
	rateLimitMgr  *ratelimit.RateLimitManager
	corsHandler   *corsMiddleware.CORSHandler
	requestLogger *loggingMiddleware.RequestLogger
	metrics       *metricsMiddleware.MetricsMiddleware
	cacheHandlers map[string]*cacheMiddleware.CacheHandler // 按路由路径存储缓存处理器
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status          string                            `json:"status"`                     // 服务状态：ok, degraded, down
	Service         string                            `json:"service"`                    // 服务名称
	Timestamp       string                            `json:"timestamp"`                  // 时间戳
	Version         string                            `json:"version,omitempty"`          // 版本号（可选）
	Uptime          string                            `json:"uptime,omitempty"`           // 运行时间（可选）
	Services        map[string]ServiceHealth          `json:"services,omitempty"`         // 下游服务健康状态（可选）
	CircuitBreakers map[string]map[string]interface{} `json:"circuit_breakers,omitempty"` // 熔断器状态（可选）
}

// ServiceHealth 服务健康状态
type ServiceHealth struct {
	Status        string `json:"status"`         // 状态：ok, degraded, down
	InstanceCount int    `json:"instance_count"` // 可用实例数
	LastCheck     string `json:"last_check"`     // 最后检查时间
}

// NewGatewayHandler 创建Gateway处理器
func NewGatewayHandler(router *router.Router, jwtManager *jwtMiddleware.Manager, corsHandler *corsMiddleware.CORSHandler, metrics *metricsMiddleware.MetricsMiddleware) *GatewayHandler {
	return &GatewayHandler{
		router:        router,
		jwtManager:    jwtManager,
		rateLimitMgr:  ratelimit.NewRateLimitManager(),
		corsHandler:   corsHandler,
		requestLogger: loggingMiddleware.NewRequestLogger(),
		metrics:       metrics,
		cacheHandlers: make(map[string]*cacheMiddleware.CacheHandler),
	}
}

// RegisterRoutes 注册路由
func (h *GatewayHandler) RegisterRoutes(srv *kratosHttp.Server) {
	// 注册健康检查路由（不需要认证）
	srv.Route("/").GET("/health", h.Health)
	srv.Route("/api/v1").GET("/health", h.Health)

	// 注册Kubernetes探针路由
	srv.Route("/").GET("/ready", h.Readiness)
	srv.Route("/").GET("/live", h.Liveness)

	// 注册通用路由转发（使用通配符匹配所有 /api/v1/* 请求）
	// 注意：Kratos 的路由匹配机制
	// 由于 {path} 可能无法匹配多级路径（如 /api/v1/users/register），
	// 我们需要注册更具体的路由，或者使用不同的路由语法
	apiRoute := srv.Route("/api/v1")

	// 注册 users 服务的路由
	// 注意：Kratos 的 {path} 参数可能只能匹配单级路径
	// 对于多级路径，需要注册更具体的路由
	usersRoute := srv.Route("/api/v1/users")
	usersRoute.GET("/{path}", h.Proxy)
	usersRoute.POST("/{path}", h.Proxy)
	usersRoute.PUT("/{path}", h.Proxy)
	usersRoute.DELETE("/{path}", h.Proxy)
	usersRoute.PATCH("/{path}", h.Proxy)
	usersRoute.HEAD("/{path}", h.Proxy)
	usersRoute.OPTIONS("/{path}", h.Proxy)

	// 同时注册根路径的通配符路由，作为后备
	apiRoute.GET("/{path}", h.Proxy)
	apiRoute.POST("/{path}", h.Proxy)
	apiRoute.PUT("/{path}", h.Proxy)
	apiRoute.DELETE("/{path}", h.Proxy)
	apiRoute.PATCH("/{path}", h.Proxy)
	apiRoute.HEAD("/{path}", h.Proxy)
	apiRoute.OPTIONS("/{path}", h.Proxy) // OPTIONS 也由 Proxy 处理
}

// Health 健康检查接口
func (h *GatewayHandler) Health(ctx kratosHttp.Context) error {
	response := HealthResponse{
		Status:    "ok",
		Service:   "gateway",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0", // 可以从配置或环境变量读取
	}

	// 检查下游服务健康状态
	services := h.checkDownstreamServices(ctx.Request().Context())
	if len(services) > 0 {
		response.Services = services

		// 如果有服务不健康，将网关状态设置为 degraded
		for _, serviceHealth := range services {
			if serviceHealth.Status != "ok" {
				response.Status = "degraded"
				break
			}
		}
	}

	// 获取熔断器状态
	cbStats := h.router.GetCircuitBreakerStats()
	if len(cbStats) > 0 {
		response.CircuitBreakers = cbStats
		// 检查是否有熔断器处于打开状态
		for _, stats := range cbStats {
			if state, ok := stats["state"].(string); ok && state == "open" {
				if response.Status == "ok" {
					response.Status = "degraded"
				}
			}
		}
	}

	return ctx.JSON(200, response)
}

// Readiness 就绪检查接口（用于Kubernetes readiness probe）
func (h *GatewayHandler) Readiness(ctx kratosHttp.Context) error {
	// 检查关键依赖是否就绪
	// 例如：数据库连接、Redis连接等

	// 目前简单返回ok，后续可以添加更详细的检查
	return ctx.JSON(200, map[string]interface{}{
		"status":  "ready",
		"service": "gateway",
	})
}

// Liveness 存活检查接口（用于Kubernetes liveness probe）
func (h *GatewayHandler) Liveness(ctx kratosHttp.Context) error {
	// 简单的存活检查，只要服务在运行就返回ok
	return ctx.JSON(200, map[string]interface{}{
		"status":  "alive",
		"service": "gateway",
	})
}

// checkDownstreamServices 检查下游服务健康状态
func (h *GatewayHandler) checkDownstreamServices(ctx context.Context) map[string]ServiceHealth {
	services := make(map[string]ServiceHealth)

	// 获取所有已注册的服务名称
	serviceNames := h.router.GetAllServiceNames()

	// 检查每个服务的健康状态
	for _, serviceName := range serviceNames {
		instances, err := h.router.GetServiceInstances(ctx, serviceName)
		if err != nil {
			// 如果获取实例失败，标记为down
			services[serviceName] = ServiceHealth{
				Status:        "down",
				InstanceCount: 0,
				LastCheck:     time.Now().Format(time.RFC3339),
			}
			continue
		}

		// 统计健康实例数
		healthyCount := 0
		for _, instance := range instances {
			if instance.Healthy {
				healthyCount++
			}
		}

		// 确定服务状态
		status := "ok"
		if healthyCount == 0 {
			status = "down"
		} else if healthyCount < len(instances) {
			status = "degraded"
		}

		services[serviceName] = ServiceHealth{
			Status:        status,
			InstanceCount: healthyCount,
			LastCheck:     time.Now().Format(time.RFC3339),
		}
	}

	return services
}

// CORS 处理CORS预检请求
func (h *GatewayHandler) CORS(ctx kratosHttp.Context) error {
	// 添加调试日志
	log.Info(ctx.Request().Context(), "处理 CORS 预检请求",
		log.String("method", ctx.Request().Method),
		log.String("path", ctx.Request().URL.Path),
		log.String("origin", ctx.Request().Header.Get("Origin")),
	)
	return h.corsHandler.HandleCORS(ctx)
}

// Proxy 代理处理器（处理所有 /api/v1/* 请求）
func (h *GatewayHandler) Proxy(ctx kratosHttp.Context) error {
	startTime := time.Now()
	method := ctx.Request().Method
	path := ctx.Request().URL.Path

	// 优先处理 OPTIONS 预检请求
	if method == "OPTIONS" {
		log.Info(ctx.Request().Context(), "Proxy 方法收到 OPTIONS 请求，转发到 CORS 处理器",
			log.String("path", path),
			log.String("origin", ctx.Request().Header.Get("Origin")),
		)
		return h.corsHandler.HandleCORS(ctx)
	}

	// 生成或获取 TraceID
	requestCtx := ctx.Request().Context()
	traceID := getTraceIDFromRequest(ctx.Request())
	if traceID == "" {
		traceID = generateTraceID()
	}
	// 将 TraceID 注入到 Context 中
	requestCtx = context.WithValue(requestCtx, log.CtxTraceID, traceID)
	requestCtx = context.WithValue(requestCtx, log.CtxRequestID, traceID)
	// 注意：Kratos 的 Context 不支持直接更新 Request 的 Context
	// 我们通过传递 requestCtx 来使用带 TraceID 的 Context

	// 在响应头中添加 TraceID
	ctx.Response().Header().Set("X-Trace-ID", traceID)
	ctx.Response().Header().Set("X-Request-ID", traceID)

	// 增加活跃请求数
	if h.metrics != nil {
		h.metrics.IncRequestsInFlight()
		defer h.metrics.DecRequestsInFlight()
	}

	// 记录请求大小
	var requestSize int64
	if ctx.Request().Body != nil {
		requestSize = ctx.Request().ContentLength
	}

	// 处理CORS（在代理之前，对于非 OPTIONS 请求）
	// 注意：OPTIONS 请求已经在上面处理了
	if err := h.corsHandler.HandleCORS(ctx); err != nil {
		h.requestLogger.LogError(ctx, err, startTime)
		if h.metrics != nil {
			duration := time.Since(startTime)
			h.metrics.RecordRequest(requestCtx, method, path, 500, duration, requestSize, 0)
		}
		return ctx.JSON(500, NewErrorResponse(requestCtx, 500, "CORS 处理失败", err))
	}

	// 查找匹配的路由
	route := h.router.FindRoute(path)
	if route == nil {
		log.Warn(requestCtx, "未找到匹配的路由",
			log.String("path", path),
		)
		if h.metrics != nil {
			duration := time.Since(startTime)
			h.metrics.RecordRequest(requestCtx, method, path, 404, duration, requestSize, 0)
		}
		return ctx.JSON(404, ErrNotFound(requestCtx))
	}

	// 检查缓存（仅在 GET 请求且配置了缓存时）
	var cacheHandler *cacheMiddleware.CacheHandler
	if route.Cache != nil && route.Cache.Enabled {
		// 获取或创建缓存处理器
		cacheHandlerKey := route.Path
		if handler, exists := h.cacheHandlers[cacheHandlerKey]; exists {
			cacheHandler = handler
		} else {
			// 转换配置类型
			cacheConfig := &cacheMiddleware.CacheConfig{
				Enabled:            route.Cache.Enabled,
				TTL:                route.Cache.TTL,
				KeyPrefix:          route.Cache.KeyPrefix,
				Methods:            route.Cache.Methods,
				Paths:              route.Cache.Paths,
				ExcludePaths:       route.Cache.ExcludePaths,
				IncludeQueryParams: route.Cache.IncludeQueryParams,
				IncludeHeaders:     route.Cache.IncludeHeaders,
			}
			// 创建新的缓存中间件和处理器
			cacheMid, err := cacheMiddleware.NewCacheMiddleware(cacheConfig)
			if err != nil {
				log.Warn(ctx, "创建缓存中间件失败",
					log.ErrorField(err),
					log.String("path", path),
				)
			} else {
				cacheHandler = cacheMiddleware.NewCacheHandler(cacheMid)
				h.cacheHandlers[cacheHandlerKey] = cacheHandler
			}
		}

		// 检查缓存
		if cacheHandler != nil {
			cachedResp, hit := cacheHandler.HandleRequest(requestCtx, ctx.Request())
			if hit && cachedResp != nil {
				// 缓存命中，直接返回
				// 复制响应头
				for key, values := range cachedResp.Headers {
					for _, value := range values {
						ctx.Response().Header().Set(key, value)
					}
				}

				// 添加缓存相关的响应头
				ctx.Response().Header().Set("X-Cache", "HIT")
				ctx.Response().Header().Set("X-Cache-Age", fmt.Sprintf("%.0f", time.Since(cachedResp.CachedAt).Seconds()))

				// 设置状态码
				ctx.Response().WriteHeader(cachedResp.StatusCode)

				// 写入响应体
				_, err := ctx.Response().Write(cachedResp.Body)
				if err != nil {
					log.Warn(requestCtx, "写入缓存响应失败",
						log.ErrorField(err),
					)
				}

				// 记录缓存命中指标
				if h.metrics != nil {
					h.metrics.RecordCacheHit(requestCtx, path)
					duration := time.Since(startTime)
					h.metrics.RecordRequest(requestCtx, method, path, cachedResp.StatusCode, duration, requestSize, int64(len(cachedResp.Body)))
				}

				return nil
			} else {
				// 缓存未命中，记录指标
				if h.metrics != nil {
					h.metrics.RecordCacheMiss(requestCtx, path)
				}
			}
		}
	}

	// 检查限流（在认证之前检查，避免浪费资源）
	if route.RateLimit != nil {
		allowed, err := ratelimit.CheckRateLimit(
			requestCtx,
			h.rateLimitMgr,
			path,
			route.RateLimit.QPS,
			route.RateLimit.Burst,
			route.RequireAuth,
		)
		if err != nil {
			log.Warn(requestCtx, "限流检查失败",
				log.ErrorField(err),
				log.String("path", path),
			)
		}
		if !allowed {
			if h.metrics != nil {
				h.metrics.RecordRateLimit(requestCtx, path)
				duration := time.Since(startTime)
				h.metrics.RecordRequest(requestCtx, method, path, 429, duration, requestSize, 0)
			}
			return ctx.JSON(429, ErrRateLimit(requestCtx))
		}
	}

	// 检查是否需要认证
	if route.RequireAuth {
		// 获取 Authorization 头
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			if h.metrics != nil {
				duration := time.Since(startTime)
				h.metrics.RecordRequest(requestCtx, method, path, 401, duration, requestSize, 0)
			}
			return ctx.JSON(401, ErrUnauthorized(requestCtx))
		}

		// 解析 Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			if h.metrics != nil {
				duration := time.Since(startTime)
				h.metrics.RecordRequest(requestCtx, method, path, 401, duration, requestSize, 0)
			}
			return ctx.JSON(401, ErrInvalidAuth(requestCtx))
		}

		token := parts[1]

		// 验证 Token
		_, err := h.jwtManager.ValidateToken(token)
		if err != nil {
			log.Warn(requestCtx, "Token 验证失败",
				log.ErrorField(err),
			)
			if h.metrics != nil {
				duration := time.Since(startTime)
				h.metrics.RecordRequest(requestCtx, method, path, 401, duration, requestSize, 0)
			}
			return ctx.JSON(401, ErrInvalidToken(requestCtx, err))
		}
	}

	// 转发请求
	downstreamStartTime := time.Now()

	// 定义缓存回调函数
	var cacheCallback router.CacheCallback
	if cacheHandler != nil {
		cacheCallback = func(statusCode int, headers map[string][]string, body []byte) {
			// 调用缓存处理器写入缓存
			cacheHandler.HandleResponse(requestCtx, ctx.Request(), statusCode, headers, body)
		}
	}

	// 转发请求（传递缓存回调）
	err := h.router.Forward(ctx, route, cacheCallback)
	downstreamDuration := time.Since(downstreamStartTime)

	if err != nil {
		// 根据错误类型确定状态码和错误响应
		statusCode := 500
		var errorResp *StandardResponse

		// 分析错误类型
		if isTimeoutError(err) {
			statusCode = 504
			errorResp = ErrRequestTimeout(requestCtx, time.Duration(route.Timeout)*time.Second)
		} else if strings.Contains(err.Error(), "熔断器") {
			statusCode = 503
			errorResp = ErrCircuitBreakerOpen(requestCtx, route.Service)
		} else if strings.Contains(err.Error(), "没有可用实例") {
			statusCode = 503
			errorResp = ErrNoServiceInstance(requestCtx, route.Service)
		} else if isNetworkError(err) {
			statusCode = 502
			errorResp = NewErrorResponse(requestCtx, statusCode, "下游服务网络错误", err)
		} else {
			statusCode = 500
			errorResp = ErrServiceUnavailable(requestCtx, err)
		}

		log.Error(requestCtx, "转发请求失败",
			log.ErrorField(err),
			log.String("path", path),
			log.String("service", route.Service),
		)
		h.requestLogger.LogError(ctx, err, startTime)

		// 记录下游服务请求（失败）
		if h.metrics != nil {
			h.metrics.RecordDownstream(requestCtx, route.Service, statusCode, downstreamDuration)
		}

		// 记录总请求指标
		if h.metrics != nil {
			totalDuration := time.Since(startTime)
			h.metrics.RecordRequest(requestCtx, method, path, statusCode, totalDuration, requestSize, 0)
		}

		return ctx.JSON(statusCode, errorResp)
	}

	// 记录下游服务请求（成功）
	responseStatusCode := 200 // Forward 方法已经写入了响应，假设成功
	if h.metrics != nil {
		h.metrics.RecordDownstream(requestCtx, route.Service, responseStatusCode, downstreamDuration)
	}

	// 记录请求日志
	h.requestLogger.LogRequest(ctx, startTime)

	// 记录总请求指标
	if h.metrics != nil {
		totalDuration := time.Since(startTime)
		// 注意：由于响应已经写入，我们无法获取响应体大小
		// 如果需要精确的响应大小，需要在 Forward 方法中返回
		h.metrics.RecordRequest(requestCtx, method, path, responseStatusCode, totalDuration, requestSize, 0)
	}

	return nil
}
