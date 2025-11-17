//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"StructForge/backend/apps/gateway/internal/conf"
	"StructForge/backend/apps/gateway/internal/handler"
	"StructForge/backend/apps/gateway/internal/server"
)

// wireApp 初始化应用（由 Wire 生成）
// 注意：此文件只在 wireinject 构建标签下编译
// 运行 wire 命令后会生成 wire_gen.go 文件
func wireApp(*conf.Bootstrap, *conf.Redis) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		handler.ProviderSet,
		newLogger,
		newApp,
	))
}
