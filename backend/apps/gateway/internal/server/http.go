package server

import (
	"StructForge/backend/apps/gateway/internal/conf"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(c *conf.Bootstrap) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(), // 恢复中间件
		),
	}

	// 设置默认地址（如果配置中没有指定）
	addr := ":8000"
	if c.Server != nil && c.Server.Id != "" {
		// 可以根据服务ID设置不同的端口
		addr = ":8000"
	}

	opts = append(opts, http.Address(addr))

	return http.NewServer(opts...)
}
