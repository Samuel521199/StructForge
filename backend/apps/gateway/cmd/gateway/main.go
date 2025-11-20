package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"go.uber.org/automaxprocs/maxprocs"

	"StructForge/backend/apps/gateway/internal/conf"
	"StructForge/backend/common/cache"
	"StructForge/backend/common/log"
	nacosClient "StructForge/backend/common/middleware/nacos"
)

// 命令行参数定义
var (
	flagconf string // 配置文件路径
	flagenv  string // 环境标识（local/test/prod）
)

// init 函数：程序启动时自动执行
func init() {
	// 默认环境为 local
	flag.StringVar(&flagenv, "env", "local", "environment: local, test, prod")
	flag.StringVar(&flagconf, "conf", "", "config path, eg: -conf gateway.yaml")
}

// getConfigPath 根据环境获取配置文件路径
func getConfigPath() string {
	// 如果指定了具体的配置文件路径，直接使用
	if flagconf != "" {
		return flagconf
	}

	// 优先使用环境变量，然后使用命令行参数
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = flagenv
	}

	// 根据环境获取配置文件路径
	var configPath string
	switch env {
	case "local":
		configPath = "../../../../configs/local/gateway.yaml"
	case "test":
		configPath = "../../../../configs/test/gateway.yaml"
	case "prod":
		configPath = "../../../../configs/prod/gateway.yaml"
	default:
		os.Stderr.WriteString("Unknown environment: " + env + ", using local environment\n")
		configPath = "../../../../configs/local/gateway.yaml"
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.Stderr.WriteString("Config file not found: " + configPath + "\n")
		os.Stderr.WriteString("Available environments: local, test, prod\n")
		os.Exit(1)
	}

	return configPath
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 设置最大 CPU 核心数（自动适配容器环境）
	_, err := maxprocs.Set()
	if err != nil {
		panic(err)
	}

	// 获取环境变量
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = flagenv
	}

	// ========== 第一步：初始化基础日志系统（最前端）==========
	// 先初始化一个基础的日志系统，这样后续的配置加载和服务初始化日志都能被打印
	defer log.InitLoggerWithShutdown("gateway",
		log.WithEnvironment(env),
		log.WithFilePath("logs/gateway-%s.log"),
	)()

	// 创建启动上下文（包含追踪信息）
	ctx := context.Background()
	if podName := os.Getenv("POD_NAME"); podName != "" {
		ctx = context.WithValue(ctx, log.CtxTraceID, podName)
	}

	// 记录启动开始
	log.Info(ctx, "gateway 服务启动中",
		log.String("environment", env),
		log.String("pod_name", os.Getenv("POD_NAME")),
		log.String("hostname", os.Getenv("HOSTNAME")),
	)

	// 获取配置文件路径
	configPath := getConfigPath()
	if configPath == "" {
		log.Error(ctx, "配置文件路径为空")
		os.Exit(1)
	}

	log.Info(ctx, "开始加载配置文件",
		log.String("config_path", configPath),
		log.String("environment", env),
	)

	// ========== 第二步：加载启动配置文件（StartupConfig）==========
	c := config.New(config.WithSource(
		file.NewSource(configPath),
	))
	defer c.Close()

	if err := c.Load(); err != nil {
		log.Error(ctx, "加载启动配置文件失败",
			log.String("config_path", configPath),
			log.ErrorField(err),
		)
		panic(err)
	}

	var startupConfig nacosClient.StartupConfig
	if err := c.Scan(&startupConfig); err != nil {
		log.Error(ctx, "解析启动配置文件失败",
			log.String("config_path", configPath),
			log.ErrorField(err),
		)
		panic(err)
	}

	log.Info(ctx, "启动配置文件加载成功",
		log.String("config_path", configPath),
	)

	// ========== 第三步：创建Nacos配置客户端 ==========
	log.Info(ctx, "正在创建 Nacos 配置客户端")
	nacosConfigClient, err := nacosClient.NewNacosConfigClient(&startupConfig)
	if err != nil {
		log.Error(ctx, "创建 Nacos 配置客户端失败",
			log.ErrorField(err),
		)
		os.Exit(1)
	}
	log.Info(ctx, "Nacos 配置客户端创建成功")

	// ========== 第四步：加载完整配置（Bootstrap）==========
	// 如果启用了 Nacos 配置中心，从 Nacos 获取配置
	// 否则直接从本地文件加载配置
	var bc *conf.Bootstrap

	if startupConfig.Nacos.ConfigCenter.Enabled {
		// 从 Nacos 获取配置
		log.Info(ctx, "正在从 Nacos 获取完整配置")
		var bcTemp conf.Bootstrap
		dm, err := nacosClient.NewConfigManager(nacosConfigClient, &startupConfig, &bcTemp)
		if err != nil {
			log.Error(ctx, "从 Nacos 获取配置失败",
				log.ErrorField(err),
			)
			os.Exit(1)
		}
		defer dm.Close()
		bc = &bcTemp

		log.Info(ctx, "从 Nacos 获取配置成功",
			log.String("data_id", startupConfig.Nacos.ConfigCenter.DataId),
			log.String("group", startupConfig.Nacos.ConfigCenter.Group),
		)
	} else {
		// 直接从本地文件加载配置
		log.Info(ctx, "正在从本地文件加载配置")
		var err error
		bc, err = loadConfig(configPath)
		if err != nil {
			log.Error(ctx, "加载配置文件失败",
				log.ErrorField(err),
				log.String("config_path", configPath),
			)
			os.Exit(1)
		}
		log.Info(ctx, "从本地文件加载配置成功",
			log.String("config_path", configPath),
		)
	}

	// ========== 第五步：更新日志系统配置（如果需要）==========
	// 如果 Bootstrap 配置中有服务ID，可以更新日志系统
	if bc.Server != nil && bc.Server.Id != "" {
		log.Info(ctx, "配置加载完成",
			log.String("server_id", bc.Server.Id),
			log.String("server_name", bc.Server.Name),
			log.String("config_path", configPath),
			log.String("environment", env),
		)
	} else {
		log.Info(ctx, "配置加载完成",
			log.String("config_path", configPath),
			log.String("environment", env),
		)
	}

	// ========== 第五步（补充）：初始化缓存系统 ==========
	log.Info(ctx, "正在初始化缓存系统")

	// 根据配置选择缓存适配器（优先使用Redis，如果Redis配置不存在则使用内存缓存）
	adapterType := "memory" // 默认使用内存缓存
	var cacheOptions []cache.InitOption

	if bc.Redis != nil && bc.Redis.Addr != "" {
		// 使用Redis缓存
		adapterType = "redis"
		cacheOptions = append(cacheOptions,
			cache.WithKeyPrefix("gateway:"),
			cache.WithDefaultTTL(5*time.Minute),
			cache.WithRedisAddr(bc.Redis.Addr),
		)

		if bc.Redis.Password != "" {
			cacheOptions = append(cacheOptions, cache.WithRedisPassword(bc.Redis.Password))
		}

		if bc.Redis.Db > 0 {
			cacheOptions = append(cacheOptions, cache.WithRedisDB(int(bc.Redis.Db)))
		}

		log.Info(ctx, "使用Redis缓存",
			log.String("addr", bc.Redis.Addr),
			log.Int("db", int(bc.Redis.Db)),
		)
	} else {
		// 使用内存缓存
		cacheOptions = append(cacheOptions,
			cache.WithKeyPrefix("gateway:"),
			cache.WithDefaultTTL(5*time.Minute),
			cache.WithMemoryMaxSize(100*1024*1024), // 100MB
			cache.WithMemoryMaxItems(1000),
		)

		log.Info(ctx, "使用内存缓存（Redis配置未提供）")
	}

	// 初始化缓存并设置优雅关闭
	cacheCleanup := cache.InitCacheWithShutdown(adapterType, cacheOptions...)
	defer cacheCleanup()

	// 测试缓存连接
	cacheInstance := cache.GetGlobalCache()
	if cacheInstance != nil {
		if err := cacheInstance.Ping(ctx); err != nil {
			log.Warn(ctx, "缓存连接测试失败（不影响启动）",
				log.ErrorField(err),
			)
		} else {
			log.Info(ctx, "缓存系统初始化成功",
				log.String("adapter", adapterType),
			)
		}
	} else {
		log.Warn(ctx, "缓存实例未初始化")
	}

	// ========== 第六步：使用Wire进行依赖注入，创建应用实例 ==========
	log.Info(ctx, "正在初始化应用实例")

	app, cleanup, err := wireApp(bc, bc.Redis)
	if err != nil {
		log.Error(ctx, "gateway 服务初始化失败",
			log.ErrorField(err),
		)
		os.Exit(1)
	}
	defer cleanup()
	log.Info(ctx, "应用实例初始化成功")

	// ========== 第七步：启动服务 ==========
	serverID := ""
	if bc.Server != nil {
		serverID = bc.Server.Id
	}
	log.Info(ctx, "gateway 服务开始运行",
		log.String("server_id", serverID),
		log.String("pod_name", os.Getenv("POD_NAME")),
		log.String("hostname", os.Getenv("HOSTNAME")),
		log.String("env_namespace", os.Getenv("POD_NAMESPACE")),
	)

	// start and wait for stop signal
	// TODO: 需要实现完整的 wireApp 函数，返回 *kratos.App 类型
	// 当前 wireApp 返回 interface{}，需要类型断言后才能调用 Run()
	if app == nil {
		log.Error(ctx, "应用实例为空，请先实现完整的 wireApp 函数")
		os.Exit(1)
	}

	// 临时：等待实现完整的 wireApp 函数
	// 实际使用时需要：
	// 1. 创建 server、handler、data、biz 等模块
	// 2. 定义 ProviderSet
	// 3. 运行 wire 命令生成完整的依赖注入代码
	log.Warn(ctx, "应用实例已创建，但 wireApp 函数需要完整实现才能运行服务")
	log.Info(ctx, "请运行 wire 命令生成完整的依赖注入代码: cd backend/apps/gateway/cmd/gateway && wire")

	// 等待信号（临时实现）
	// select {} // 永久阻塞，等待中断信号

	// 启动服务（app已经是*kratos.App类型，不需要类型断言）
	if err := app.Run(); err != nil {
		log.Error(ctx, "gateway 服务运行失败", log.ErrorField(err))
		os.Exit(1)
	}

	log.Info(ctx, "gateway 服务结束运行",
		log.String("server_id", serverID),
		log.String("pod_name", os.Getenv("POD_NAME")),
		log.String("hostname", os.Getenv("HOSTNAME")),
	)
}
