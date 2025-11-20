package router

import (
	"context"
	"fmt"

	"StructForge/backend/apps/gateway/internal/conf"
	"StructForge/backend/apps/gateway/internal/router/discovery"
	"StructForge/backend/common/log"
)

// LoadRouterFromConfig 从配置加载路由（Wire provider，只返回 Router）
func LoadRouterFromConfig(config *conf.GatewayConfig, staticDiscovery *discovery.StaticDiscovery) (*Router, error) {
	ctx := context.Background()

	// 验证配置
	if err := ValidateGatewayConfig(config); err != nil {
		log.Error(ctx, "网关配置验证失败",
			log.ErrorField(err),
		)
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 创建路由管理器
	router := NewRouter(staticDiscovery)

	// 加载路由规则
	if config != nil && config.Routes != nil {
		for _, routeConfig := range config.Routes.Routes {
			route := &Route{
				Path:                routeConfig.Path,
				MatchType:           routeConfig.MatchType,
				Service:             routeConfig.Service,
				TargetPath:          routeConfig.TargetPath,
				RequireAuth:         routeConfig.RequireAuth,
				Timeout:             routeConfig.Timeout,
				Retries:             routeConfig.Retries,
				LoadBalanceStrategy: routeConfig.LoadBalanceStrategy,
			}

			if routeConfig.RateLimit != nil {
				route.RateLimit = &RateLimitConfig{
					QPS:   routeConfig.RateLimit.QPS,
					Burst: routeConfig.RateLimit.Burst,
				}
			}

			if routeConfig.CircuitBreaker != nil {
				route.CircuitBreaker = routeConfig.CircuitBreaker
			}

			if routeConfig.Cache != nil {
				route.Cache = routeConfig.Cache
			}

			router.AddRoute(route)
		}
	}

	// 加载服务实例（静态服务发现）
	if config != nil && config.Services != nil {
		for serviceName, instances := range config.Services.Services {
			discoveryInstances := make([]discovery.Instance, 0, len(instances))
			for _, instanceConfig := range instances {
				discoveryInstances = append(discoveryInstances, discovery.Instance{
					ID:       instanceConfig.ID,
					Host:     instanceConfig.Host,
					Port:     instanceConfig.Port,
					Weight:   instanceConfig.Weight,
					Healthy:  instanceConfig.Healthy,
					Metadata: instanceConfig.Metadata,
				})
			}
			staticDiscovery.RegisterService(serviceName, discoveryInstances)
		}
	}

	log.Info(ctx, "路由配置加载完成",
		log.Int("routes", len(router.routes)),
	)

	return router, nil
}
