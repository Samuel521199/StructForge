package conf

// RouteConfig 路由配置（从配置文件加载）
type RouteConfig struct {
	Routes []RouteRule `yaml:"routes" json:"routes"`
}

// RouteRule 路由规则配置
type RouteRule struct {
	Path                string                `yaml:"path" json:"path"`
	MatchType           string                `yaml:"match_type" json:"match_type"`
	Service             string                `yaml:"service" json:"service"`
	TargetPath          string                `yaml:"target_path" json:"target_path"`
	RequireAuth         bool                  `yaml:"require_auth" json:"require_auth"`
	Timeout             int                   `yaml:"timeout" json:"timeout"`
	Retries             int                   `yaml:"retries" json:"retries"`
	LoadBalanceStrategy string                `yaml:"load_balance_strategy" json:"load_balance_strategy"`
	RateLimit           *RateLimitConfig      `yaml:"rate_limit" json:"rate_limit"`
	CircuitBreaker      *CircuitBreakerConfig `yaml:"circuit_breaker" json:"circuit_breaker"`
	Cache               *CacheConfig          `yaml:"cache" json:"cache"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	QPS   int `yaml:"qps" json:"qps"`
	Burst int `yaml:"burst" json:"burst"`
}

// CircuitBreakerConfig 熔断器配置
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

// ServiceConfig 服务配置（静态服务发现）
type ServiceConfig struct {
	Services map[string][]ServiceInstance `yaml:"services" json:"services"`
}

// ServiceInstance 服务实例配置
type ServiceInstance struct {
	ID       string            `yaml:"id" json:"id"`
	Host     string            `yaml:"host" json:"host"`
	Port     int               `yaml:"port" json:"port"`
	Weight   int               `yaml:"weight" json:"weight"`
	Healthy  bool              `yaml:"healthy" json:"healthy"`
	Metadata map[string]string `yaml:"metadata" json:"metadata"`
}

// GatewayConfig 网关完整配置
type GatewayConfig struct {
	// JWT 配置
	JWT *JWTConfig `yaml:"jwt" json:"jwt"`
	// 路由配置
	Routes *RouteConfig `yaml:"routes" json:"routes"`
	// 服务配置（静态服务发现）
	Services *ServiceConfig `yaml:"services" json:"services"`
	// 前端配置
	Frontend *FrontendConfig `yaml:"frontend" json:"frontend"`
	// CORS 配置
	CORS *CORSConfig `yaml:"cors" json:"cors"`
}

// FrontendConfig 前端配置
type FrontendConfig struct {
	// 前端地址（用于CORS等配置）
	URL string `yaml:"url" json:"url"`
	// 允许的前端地址列表（支持多个）
	AllowedURLs []string `yaml:"allowed_urls" json:"allowed_urls"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	// 允许的源（支持通配符 *）
	AllowedOrigins []string `yaml:"allowed_origins" json:"allowed_origins"`
	// 允许的方法
	AllowedMethods []string `yaml:"allowed_methods" json:"allowed_methods"`
	// 允许的请求头
	AllowedHeaders []string `yaml:"allowed_headers" json:"allowed_headers"`
	// 暴露的响应头
	ExposedHeaders []string `yaml:"exposed_headers" json:"exposed_headers"`
	// 是否允许携带凭证
	AllowCredentials bool `yaml:"allow_credentials" json:"allow_credentials"`
	// 预检请求的缓存时间（秒）
	MaxAge int `yaml:"max_age" json:"max_age"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	SecretKey     string `yaml:"secret_key" json:"secret_key"`
	TokenDuration string `yaml:"token_duration" json:"token_duration"` // 如: "24h", "7d"
}
