package server

import (
	"time"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	v1 "StructForge/backend/api/user/v1"
	"StructForge/backend/apps/user/internal/conf"
	"StructForge/backend/apps/user/internal/service"
)

// GRPCServer gRPC 服务器类型别名
type GRPCServer = grpc.Server

// NewGRPCServer 创建 gRPC 服务器
func NewGRPCServer(c *conf.Bootstrap, userService *service.UserService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}

	// 从配置中读取 gRPC 服务器地址
	if c.Server != nil && c.Server.Grpc != nil {
		if c.Server.Grpc.Addr != "" {
			opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
		}
		if c.Server.Grpc.Timeout > 0 {
			opts = append(opts, grpc.Timeout(time.Duration(c.Server.Grpc.Timeout)*time.Second))
		}
	} else {
		// 默认地址
		opts = append(opts, grpc.Address(":9001"))
	}

	srv := grpc.NewServer(opts...)
	v1.RegisterUserServiceServer(srv, userService)

	return srv
}
