package main

import (
	"context"
	"strings"

	"StructForge/backend/apps/gateway/internal/conf"
	"StructForge/backend/common/log"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

// loadConfig 加载配置文件
func loadConfig(configPath string) (*conf.Bootstrap, error) {
	ctx := context.Background()

	// 创建文件配置源
	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)

	// 加载配置
	if err := c.Load(); err != nil {
		return nil, err
	}

	// 创建 Bootstrap 结构
	var bc conf.Bootstrap

	// 扫描配置到结构体
	if err := c.Scan(&bc); err != nil {
		log.Error(ctx, "扫描配置失败",
			log.ErrorField(err),
		)
		return nil, err
	}

	// 添加调试日志
	log.Info(ctx, "配置扫描完成",
		log.Bool("has_server", bc.Server != nil),
		log.Bool("has_redis", bc.Redis != nil),
		log.Bool("has_gateway", bc.Gateway != nil),
	)

	if bc.Gateway != nil {
		log.Info(ctx, "Gateway 配置详情",
			log.Bool("has_jwt", bc.Gateway.JWT != nil),
			log.Bool("has_routes", bc.Gateway.Routes != nil),
			log.Bool("has_services", bc.Gateway.Services != nil),
			log.Bool("has_frontend", bc.Gateway.Frontend != nil),
			log.Bool("has_cors", bc.Gateway.CORS != nil),
		)
		if bc.Gateway.Routes != nil {
			log.Info(ctx, "路由数量",
				log.Int("count", len(bc.Gateway.Routes.Routes)),
			)
		}
		if bc.Gateway.CORS != nil {
			originsStr := ""
			if len(bc.Gateway.CORS.AllowedOrigins) > 0 {
				originsStr = strings.Join(bc.Gateway.CORS.AllowedOrigins, ", ")
			}
			log.Info(ctx, "CORS 配置",
				log.String("allowed_origins", originsStr),
				log.Int("origins_count", len(bc.Gateway.CORS.AllowedOrigins)),
			)
		}
	} else {
		log.Warn(ctx, "Gateway 配置为 nil，检查 YAML 文件结构")
	}

	// 关闭配置源
	defer func() {
		if err := c.Close(); err != nil {
			log.Warn(ctx, "关闭配置源失败",
				log.ErrorField(err),
			)
		}
	}()

	return &bc, nil
}
