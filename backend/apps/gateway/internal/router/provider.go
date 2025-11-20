package router

import (
	"time"

	"StructForge/backend/apps/gateway/internal/conf"
	corsMiddleware "StructForge/backend/apps/gateway/internal/middleware/cors"
	jwtMiddleware "StructForge/backend/apps/gateway/internal/middleware/jwt"
	"StructForge/backend/apps/gateway/internal/router/discovery"

	"github.com/google/wire"
)

// ProviderSet 路由模块依赖注入
var ProviderSet = wire.NewSet(
	NewStaticDiscovery,
	NewJWTManagerFromConfig,
	NewCORSHandlerFromConfig,
	LoadRouterFromConfig, // LoadRouterFromConfig 内部会调用 NewRouter
)

// NewStaticDiscovery 创建静态服务发现
func NewStaticDiscovery() *discovery.StaticDiscovery {
	return discovery.NewStaticDiscovery()
}

// NewJWTManagerFromConfig 从配置创建 JWT 管理器（Wire provider）
func NewJWTManagerFromConfig(config *conf.GatewayConfig) *jwtMiddleware.Manager {
	if config != nil && config.JWT != nil {
		secretKey := config.JWT.SecretKey
		if secretKey == "" {
			secretKey = "your-secret-key-change-in-production"
		}

		tokenDuration := 24 * time.Hour
		if config.JWT.TokenDuration != "" {
			if duration, err := time.ParseDuration(config.JWT.TokenDuration); err == nil {
				tokenDuration = duration
			}
		}

		return jwtMiddleware.NewManager(secretKey, tokenDuration)
	}

	// 使用默认配置
	return jwtMiddleware.NewManager("your-secret-key-change-in-production", 24*time.Hour)
}

// NewCORSHandlerFromConfig 从配置创建 CORS 处理器（Wire provider）
func NewCORSHandlerFromConfig(config *conf.GatewayConfig) *corsMiddleware.CORSHandler {
	var options *corsMiddleware.CORSOptions

	if config != nil && config.CORS != nil {
		// 从配置创建 CORS 选项
		options = &corsMiddleware.CORSOptions{
			AllowedOrigins:   config.CORS.AllowedOrigins,
			AllowedMethods:   config.CORS.AllowedMethods,
			AllowedHeaders:   config.CORS.AllowedHeaders,
			ExposedHeaders:   config.CORS.ExposedHeaders,
			AllowCredentials: config.CORS.AllowCredentials,
			MaxAge:           config.CORS.MaxAge,
		}

		// 如果配置了前端地址，自动添加到允许的源
		if config.Frontend != nil {
			if config.Frontend.URL != "" {
				// 检查是否已存在
				exists := false
				for _, origin := range options.AllowedOrigins {
					if origin == config.Frontend.URL {
						exists = true
						break
					}
				}
				if !exists {
					options.AllowedOrigins = append(options.AllowedOrigins, config.Frontend.URL)
				}
			}

			// 添加多个前端地址
			if len(config.Frontend.AllowedURLs) > 0 {
				for _, url := range config.Frontend.AllowedURLs {
					exists := false
					for _, origin := range options.AllowedOrigins {
						if origin == url {
							exists = true
							break
						}
					}
					if !exists {
						options.AllowedOrigins = append(options.AllowedOrigins, url)
					}
				}
			}
		}

		// 设置默认值
		if len(options.AllowedOrigins) == 0 {
			options.AllowedOrigins = []string{"*"}
		}
		if len(options.AllowedMethods) == 0 {
			options.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
		}
		if len(options.AllowedHeaders) == 0 {
			options.AllowedHeaders = []string{"*"}
		}
		if options.MaxAge == 0 {
			options.MaxAge = 86400 // 24小时
		}
	} else {
		// 使用默认配置
		options = corsMiddleware.DefaultCORSOptions()
	}

	return corsMiddleware.NewCORSHandler(options)
}
