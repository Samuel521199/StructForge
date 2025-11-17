package handler

import (
	"github.com/go-kratos/kratos/v2/transport/http"
)

// GatewayHandler Gateway处理器
type GatewayHandler struct {
}

// NewGatewayHandler 创建Gateway处理器
func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

// RegisterRoutes 注册路由
func (h *GatewayHandler) RegisterRoutes(srv *http.Server) {
	// 注册健康检查路由
	srv.Route("/").GET("/health", h.Health)
	srv.Route("/api/v1").GET("/health", h.Health)
}

// Health 健康检查接口
func (h *GatewayHandler) Health(ctx http.Context) error {
	return ctx.JSON(200, map[string]interface{}{
		"status":  "ok",
		"service": "gateway",
	})
}
