package main

import (
	"StructForge/backend/apps/gateway/internal/conf"
	"StructForge/backend/apps/gateway/internal/handler"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// newApp 创建Kratos应用实例
func newApp(
	bc *conf.Bootstrap,
	httpSrv *http.Server,
	gatewayHandler *handler.GatewayHandler,
	logger log.Logger,
) *kratos.App {
	// 注册路由
	gatewayHandler.RegisterRoutes(httpSrv)

	// 创建应用实例
	opts := []kratos.Option{
		kratos.Logger(logger),
		kratos.Server(httpSrv),
	}

	// 如果配置中有服务信息，添加到应用实例
	if bc.Server != nil {
		if bc.Server.Id != "" {
			opts = append(opts, kratos.ID(bc.Server.Id))
		}
		if bc.Server.Name != "" {
			opts = append(opts, kratos.Name(bc.Server.Name))
		}
		if bc.Server.Version != "" {
			opts = append(opts, kratos.Version(bc.Server.Version))
		}
	}

	return kratos.New(opts...)
}
