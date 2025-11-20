package main

import (
	"context"

	"StructForge/backend/apps/gateway/internal/conf"
	"StructForge/backend/common/log"
)

// getGatewayConfig 从 Bootstrap 获取 Gateway 配置（Wire provider）
func getGatewayConfig(bc *conf.Bootstrap) *conf.GatewayConfig {
	ctx := context.Background()

	// 添加调试日志
	if bc == nil {
		log.Error(ctx, "Bootstrap 配置为空")
		return &conf.GatewayConfig{}
	}

	if bc.Gateway == nil {
		log.Warn(ctx, "Gateway 配置未提供，使用默认配置",
			log.String("bootstrap_server_id", bc.Server.Id),
			log.Bool("has_redis", bc.Redis != nil),
		)
		return &conf.GatewayConfig{}
	}

	log.Info(ctx, "Gateway 配置已加载",
		log.Bool("has_jwt", bc.Gateway.JWT != nil),
		log.Bool("has_routes", bc.Gateway.Routes != nil),
		log.Bool("has_services", bc.Gateway.Services != nil),
		log.Bool("has_frontend", bc.Gateway.Frontend != nil),
		log.Bool("has_cors", bc.Gateway.CORS != nil),
	)

	if bc.Gateway.Routes != nil {
		log.Info(ctx, "路由配置详情",
			log.Int("routes_count", len(bc.Gateway.Routes.Routes)),
		)
	}

	return bc.Gateway
}
