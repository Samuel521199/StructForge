package handler

import (
	"context"
	"net/http"

	"StructForge/backend/common/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// RegisterMetricsRoutes 注册指标路由
func (h *GatewayHandler) RegisterMetricsRoutes(srv *kratosHttp.Server) {
	// 注册 Prometheus 指标端点
	// 注意：Kratos 的 http.Context 不直接支持 promhttp.Handler
	// 这里提供一个简单的 JSON 响应，实际生产环境建议使用独立的 HTTP 服务器
	srv.Route("/").GET("/metrics", h.Metrics)
}

// Metrics Prometheus 指标端点
// 注意：由于 Kratos 的 http.Context 限制，这里返回 JSON 格式
// 生产环境建议使用独立的 HTTP 服务器在单独端口暴露 /metrics 端点
func (h *GatewayHandler) Metrics(ctx kratosHttp.Context) error {
	// 返回提示信息
	// 实际指标可以通过 promhttp.Handler() 在独立服务器上暴露
	return ctx.JSON(200, map[string]interface{}{
		"message": "Prometheus metrics endpoint",
		"note":    "Metrics are collected. For Prometheus scraping, use a separate HTTP server on port 9090",
		"endpoint": "/metrics",
	})
}

// NewMetricsHandler 创建指标处理器（用于独立的指标服务器）
func NewMetricsHandler() http.Handler {
	return promhttp.Handler()
}

// StartMetricsServer 启动独立的指标服务器（可选）
// 建议在生产环境使用此方法在单独端口（如 9090）暴露指标
func StartMetricsServer(addr string) error {
	http.Handle("/metrics", promhttp.Handler())
	
	log.Info(context.Background(), "启动 Prometheus 指标服务器",
		log.String("addr", addr),
		log.String("path", "/metrics"),
	)
	
	return http.ListenAndServe(addr, nil)
}

