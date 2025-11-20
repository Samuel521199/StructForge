package router

import (
	"context"
	"fmt"
	"strings"

	"StructForge/backend/apps/gateway/internal/conf"
	"StructForge/backend/common/log"
)

// ValidateGatewayConfig 验证网关配置
func ValidateGatewayConfig(config *conf.GatewayConfig) error {
	if config == nil {
		return fmt.Errorf("网关配置为空")
	}

	// 验证路由配置
	if config.Routes != nil {
		for i, route := range config.Routes.Routes {
			if err := validateRoute(route, i); err != nil {
				return fmt.Errorf("路由配置错误 [索引 %d]: %w", i, err)
			}
		}
	}

	// 验证服务配置
	if config.Services != nil {
		for serviceName, instances := range config.Services.Services {
			if err := validateService(serviceName, instances); err != nil {
				return fmt.Errorf("服务配置错误 [%s]: %w", serviceName, err)
			}
		}
	}

	// 验证JWT配置
	if config.JWT != nil {
		if err := validateJWT(config.JWT); err != nil {
			return fmt.Errorf("JWT配置错误: %w", err)
		}
	}

	// 验证CORS配置
	if config.CORS != nil {
		if err := validateCORS(config.CORS); err != nil {
			return fmt.Errorf("CORS配置错误: %w", err)
		}
	}

	return nil
}

// validateRoute 验证路由配置
func validateRoute(route conf.RouteRule, index int) error {
	// 验证路径
	if route.Path == "" {
		return fmt.Errorf("路径不能为空")
	}
	if !strings.HasPrefix(route.Path, "/") {
		return fmt.Errorf("路径必须以 / 开头")
	}

	// 验证匹配类型
	validMatchTypes := map[string]bool{
		"exact":  true,
		"prefix": true,
		"regex":  true,
	}
	if route.MatchType != "" && !validMatchTypes[route.MatchType] {
		return fmt.Errorf("无效的匹配类型: %s (支持: exact, prefix, regex)", route.MatchType)
	}

	// 验证服务名称
	if route.Service == "" {
		return fmt.Errorf("服务名称不能为空")
	}

	// 验证超时时间
	if route.Timeout < 0 {
		return fmt.Errorf("超时时间不能为负数")
	}
	if route.Timeout == 0 {
		log.Warn(context.TODO(), "路由超时时间为0，将使用默认值30秒",
			log.String("path", route.Path),
			log.Int("index", index),
		)
	}

	// 验证重试次数
	if route.Retries < 0 {
		return fmt.Errorf("重试次数不能为负数")
	}
	if route.Retries > 5 {
		log.Warn(context.TODO(), "重试次数过大，建议不超过5次",
			log.String("path", route.Path),
			log.Int("retries", route.Retries),
		)
	}

	// 验证负载均衡策略
	validStrategies := map[string]bool{
		"round_robin":       true,
		"random":            true,
		"least_connections": true,
	}
	if route.LoadBalanceStrategy != "" && !validStrategies[route.LoadBalanceStrategy] {
		return fmt.Errorf("无效的负载均衡策略: %s (支持: round_robin, random, least_connections)", route.LoadBalanceStrategy)
	}

	// 验证限流配置
	if route.RateLimit != nil {
		if route.RateLimit.QPS < 0 {
			return fmt.Errorf("QPS不能为负数")
		}
		if route.RateLimit.Burst < 0 {
			return fmt.Errorf("Burst不能为负数")
		}
		if route.RateLimit.Burst < route.RateLimit.QPS {
			log.Warn(context.TODO(), "Burst小于QPS，建议Burst >= QPS",
				log.String("path", route.Path),
				log.Int("qps", route.RateLimit.QPS),
				log.Int("burst", route.RateLimit.Burst),
			)
		}
	}

	// 验证熔断器配置
	if route.CircuitBreaker != nil && route.CircuitBreaker.Enabled {
		if route.CircuitBreaker.FailureThreshold < 0 || route.CircuitBreaker.FailureThreshold > 1 {
			return fmt.Errorf("失败率阈值必须在0-1之间")
		}
		if route.CircuitBreaker.MinRequests < 0 {
			return fmt.Errorf("最小请求数不能为负数")
		}
		if route.CircuitBreaker.WindowSize < 0 {
			return fmt.Errorf("时间窗口不能为负数")
		}
		if route.CircuitBreaker.OpenDuration < 0 {
			return fmt.Errorf("打开状态持续时间不能为负数")
		}
		if route.CircuitBreaker.HalfOpenRequests < 0 {
			return fmt.Errorf("半开状态请求数不能为负数")
		}
		if route.CircuitBreaker.Timeout < 0 {
			return fmt.Errorf("超时时间不能为负数")
		}
	}

	// 验证缓存配置
	if route.Cache != nil && route.Cache.Enabled {
		if route.Cache.TTL < 0 {
			return fmt.Errorf("缓存TTL不能为负数")
		}
		if route.Cache.TTL == 0 {
			log.Warn(context.TODO(), "缓存TTL为0，将使用默认值300秒",
				log.String("path", route.Path),
			)
		}
		if len(route.Cache.Methods) == 0 {
			log.Warn(context.TODO(), "缓存方法列表为空，将默认只缓存GET请求",
				log.String("path", route.Path),
			)
		}
	}

	return nil
}

// validateService 验证服务配置
func validateService(serviceName string, instances []conf.ServiceInstance) error {
	if serviceName == "" {
		return fmt.Errorf("服务名称不能为空")
	}

	if len(instances) == 0 {
		return fmt.Errorf("服务实例列表不能为空")
	}

	for i, instance := range instances {
		// 验证实例ID
		if instance.ID == "" {
			return fmt.Errorf("实例ID不能为空 [索引 %d]", i)
		}

		// 验证主机地址
		if instance.Host == "" {
			return fmt.Errorf("主机地址不能为空 [实例 %s]", instance.ID)
		}

		// 验证端口
		if instance.Port <= 0 || instance.Port > 65535 {
			return fmt.Errorf("端口号无效: %d (范围: 1-65535) [实例 %s]", instance.Port, instance.ID)
		}

		// 验证权重
		if instance.Weight < 0 {
			return fmt.Errorf("权重不能为负数 [实例 %s]", instance.ID)
		}
		if instance.Weight > 1000 {
			log.Warn(context.TODO(), "权重过大，建议不超过1000",
				log.String("service", serviceName),
				log.String("instance", instance.ID),
				log.Int("weight", instance.Weight),
			)
		}
	}

	return nil
}

// validateJWT 验证JWT配置
func validateJWT(jwt *conf.JWTConfig) error {
	if jwt.SecretKey == "" {
		return fmt.Errorf("JWT密钥不能为空")
	}
	if len(jwt.SecretKey) < 32 {
		log.Warn(context.TODO(), "JWT密钥长度过短，建议至少32个字符",
			log.Int("length", len(jwt.SecretKey)),
		)
	}
	if jwt.SecretKey == "your-secret-key-change-in-production" {
		log.Warn(context.TODO(), "使用默认JWT密钥，生产环境必须修改")
	}

	if jwt.TokenDuration != "" {
		// 验证时间格式（由调用方解析，这里只做基本检查）
		if !strings.Contains(jwt.TokenDuration, "h") && !strings.Contains(jwt.TokenDuration, "d") {
			log.Warn(context.TODO(), "Token有效期格式可能不正确，建议使用如 '24h' 或 '7d' 格式",
				log.String("duration", jwt.TokenDuration),
			)
		}
	}

	return nil
}

// validateCORS 验证CORS配置
func validateCORS(cors *conf.CORSConfig) error {
	// 验证允许的源
	if len(cors.AllowedOrigins) == 0 {
		log.Warn(context.TODO(), "CORS允许的源列表为空，将使用默认值 '*'")
	}

	// 验证允许的方法
	validMethods := map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"OPTIONS": true,
		"HEAD":    true,
	}
	for _, method := range cors.AllowedMethods {
		if !validMethods[strings.ToUpper(method)] {
			log.Warn(context.TODO(), "CORS允许的方法可能无效",
				log.String("method", method),
			)
		}
	}

	// 验证MaxAge
	if cors.MaxAge < 0 {
		return fmt.Errorf("CORS MaxAge不能为负数")
	}
	if cors.MaxAge > 86400*7 {
		log.Warn(context.TODO(), "CORS MaxAge过大，建议不超过7天",
			log.Int("max_age", cors.MaxAge),
		)
	}

	return nil
}
