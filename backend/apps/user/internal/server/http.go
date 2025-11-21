package server

import (
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	v1 "StructForge/backend/api/user/v1"
	"StructForge/backend/apps/user/internal/conf"
	"StructForge/backend/apps/user/internal/handler"
	"StructForge/backend/apps/user/internal/service"
)

// HTTPServer HTTP 服务器类型别名
type HTTPServer = http.Server

// NewHTTPServer 创建 HTTP 服务器（用于 HTTP Gateway）
func NewHTTPServer(c *conf.Bootstrap, userService *service.UserService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}

	// 从配置中读取 HTTP 服务器地址
	if c.Server != nil && c.Server.Http != nil {
		if c.Server.Http.Addr != "" {
			opts = append(opts, http.Address(c.Server.Http.Addr))
		}
		if c.Server.Http.Timeout > 0 {
			opts = append(opts, http.Timeout(time.Duration(c.Server.Http.Timeout)*time.Second))
		}
	} else {
		// 默认地址
		opts = append(opts, http.Address(":8001"))
	}

	srv := http.NewServer(opts...)
	v1.RegisterUserServiceHTTPServer(srv, userService)

	// 注册自定义路由（头像上传）
	srv.Route("/api/v1/users").POST("/avatar", handler.UploadAvatar)

	return srv
}
