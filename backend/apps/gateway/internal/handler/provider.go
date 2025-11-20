package handler

import (
	"github.com/google/wire"

	metricsMiddleware "StructForge/backend/apps/gateway/internal/middleware/metrics"
)

// ProviderSet Handler 模块依赖注入
var ProviderSet = wire.NewSet(
	NewGatewayHandler,
	metricsMiddleware.NewMetrics,
	metricsMiddleware.NewMetricsMiddleware,
	// 注意：router.ProviderSet 在 wire.go 中已经包含，这里不需要重复引入
)
