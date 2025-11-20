package metrics

import (
	"context"
	"strconv"
	"time"

	"StructForge/backend/common/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 指标收集器
type Metrics struct {
	// HTTP 请求总数（按方法、路径、状态码）
	httpRequestsTotal *prometheus.CounterVec
	// HTTP 请求持续时间（按方法、路径）
	httpRequestDuration *prometheus.HistogramVec
	// HTTP 请求大小（按方法、路径）
	httpRequestSize *prometheus.HistogramVec
	// HTTP 响应大小（按方法、路径）
	httpResponseSize *prometheus.HistogramVec
	// 活跃请求数（按方法、路径）
	httpRequestsInFlight prometheus.Gauge
	// 限流拒绝的请求数（按路径）
	rateLimitRejected *prometheus.CounterVec
	// 熔断器打开次数（按服务）
	circuitBreakerOpened *prometheus.CounterVec
	// 熔断器状态（按服务）
	circuitBreakerState *prometheus.GaugeVec
	// 下游服务请求总数（按服务、状态码）
	downstreamRequestsTotal *prometheus.CounterVec
	// 下游服务请求持续时间（按服务）
	downstreamRequestDuration *prometheus.HistogramVec
	// 缓存命中数（按路径）
	cacheHits *prometheus.CounterVec
	// 缓存未命中数（按路径）
	cacheMisses *prometheus.CounterVec
}

// NewMetrics 创建指标收集器
func NewMetrics() *Metrics {
	return &Metrics{
		// HTTP 请求总数
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status_code"},
		),
		// HTTP 请求持续时间
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),
		// HTTP 请求大小
		httpRequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_request_size_bytes",
				Help:    "HTTP request size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B to 100MB
			},
			[]string{"method", "path"},
		),
		// HTTP 响应大小
		httpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 7), // 100B to 100MB
			},
			[]string{"method", "path"},
		),
		// 活跃请求数
		httpRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gateway_http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		),
		// 限流拒绝的请求数
		rateLimitRejected: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_rate_limit_rejected_total",
				Help: "Total number of requests rejected by rate limiter",
			},
			[]string{"path"},
		),
		// 熔断器打开次数
		circuitBreakerOpened: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_circuit_breaker_opened_total",
				Help: "Total number of times circuit breaker opened",
			},
			[]string{"service"},
		),
		// 熔断器状态（0=closed, 1=open, 2=half-open）
		circuitBreakerState: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "gateway_circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
			},
			[]string{"service"},
		),
		// 下游服务请求总数
		downstreamRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_downstream_requests_total",
				Help: "Total number of requests to downstream services",
			},
			[]string{"service", "status_code"},
		),
		// 下游服务请求持续时间
		downstreamRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gateway_downstream_request_duration_seconds",
				Help:    "Downstream service request duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"service"},
		),
		// 缓存命中数
		cacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"path"},
		),
		// 缓存未命中数
		cacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gateway_cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"path"},
		),
	}
}

// RecordHTTPRequest 记录 HTTP 请求
func (m *Metrics) RecordHTTPRequest(method, path string, statusCode int, duration time.Duration, requestSize, responseSize int64) {
	statusCodeStr := strconv.Itoa(statusCode)
	
	// 记录请求总数
	m.httpRequestsTotal.WithLabelValues(method, path, statusCodeStr).Inc()
	
	// 记录请求持续时间
	m.httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
	
	// 记录请求大小
	if requestSize > 0 {
		m.httpRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	}
	
	// 记录响应大小
	if responseSize > 0 {
		m.httpResponseSize.WithLabelValues(method, path).Observe(float64(responseSize))
	}
}

// IncRequestsInFlight 增加活跃请求数
func (m *Metrics) IncRequestsInFlight() {
	m.httpRequestsInFlight.Inc()
}

// DecRequestsInFlight 减少活跃请求数
func (m *Metrics) DecRequestsInFlight() {
	m.httpRequestsInFlight.Dec()
}

// RecordRateLimitRejected 记录限流拒绝的请求
func (m *Metrics) RecordRateLimitRejected(path string) {
	m.rateLimitRejected.WithLabelValues(path).Inc()
}

// RecordCircuitBreakerOpened 记录熔断器打开
func (m *Metrics) RecordCircuitBreakerOpened(service string) {
	m.circuitBreakerOpened.WithLabelValues(service).Inc()
}

// SetCircuitBreakerState 设置熔断器状态
func (m *Metrics) SetCircuitBreakerState(service string, state int) {
	m.circuitBreakerState.WithLabelValues(service).Set(float64(state))
}

// RecordDownstreamRequest 记录下游服务请求
func (m *Metrics) RecordDownstreamRequest(service string, statusCode int, duration time.Duration) {
	statusCodeStr := strconv.Itoa(statusCode)
	m.downstreamRequestsTotal.WithLabelValues(service, statusCodeStr).Inc()
	m.downstreamRequestDuration.WithLabelValues(service).Observe(duration.Seconds())
}

// RecordCacheHit 记录缓存命中
func (m *Metrics) RecordCacheHit(path string) {
	m.cacheHits.WithLabelValues(path).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (m *Metrics) RecordCacheMiss(path string) {
	m.cacheMisses.WithLabelValues(path).Inc()
}

// GetRegistry 获取 Prometheus 注册表（用于暴露指标）
// 注意：promauto 使用默认注册表，这里返回 nil 表示使用默认注册表
func (m *Metrics) GetRegistry() *prometheus.Registry {
	// promauto 使用 prometheus.DefaultRegisterer
	// 如果需要自定义注册表，需要修改 NewMetrics 使用 prometheus.NewRegistry()
	return nil // 返回 nil 表示使用默认注册表
}

// MetricsMiddleware HTTP 指标中间件
type MetricsMiddleware struct {
	metrics *Metrics
}

// NewMetricsMiddleware 创建指标中间件
func NewMetricsMiddleware(metrics *Metrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
	}
}

// IncRequestsInFlight 增加活跃请求数
func (m *MetricsMiddleware) IncRequestsInFlight() {
	m.metrics.IncRequestsInFlight()
}

// DecRequestsInFlight 减少活跃请求数
func (m *MetricsMiddleware) DecRequestsInFlight() {
	m.metrics.DecRequestsInFlight()
}

// RecordRequest 记录请求（在请求处理前后调用）
func (m *MetricsMiddleware) RecordRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, requestSize, responseSize int64) {
	m.metrics.RecordHTTPRequest(method, path, statusCode, duration, requestSize, responseSize)
}

// RecordRateLimit 记录限流事件
func (m *MetricsMiddleware) RecordRateLimit(ctx context.Context, path string) {
	m.metrics.RecordRateLimitRejected(path)
	log.Warn(ctx, "请求被限流",
		log.String("path", path),
	)
}

// RecordCircuitBreaker 记录熔断器事件
func (m *MetricsMiddleware) RecordCircuitBreaker(ctx context.Context, service string, state string) {
	// 状态映射：closed=0, open=1, half-open=2
	stateValue := 0
	switch state {
	case "open":
		stateValue = 1
		m.metrics.RecordCircuitBreakerOpened(service)
	case "half-open":
		stateValue = 2
	}
	m.metrics.SetCircuitBreakerState(service, stateValue)
}

// RecordDownstream 记录下游服务请求
func (m *MetricsMiddleware) RecordDownstream(ctx context.Context, service string, statusCode int, duration time.Duration) {
	m.metrics.RecordDownstreamRequest(service, statusCode, duration)
}

// RecordCacheHit 记录缓存命中
func (m *MetricsMiddleware) RecordCacheHit(ctx context.Context, path string) {
	m.metrics.RecordCacheHit(path)
}

// RecordCacheMiss 记录缓存未命中
func (m *MetricsMiddleware) RecordCacheMiss(ctx context.Context, path string) {
	m.metrics.RecordCacheMiss(path)
}

