package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"go.uber.org/automaxprocs/maxprocs"

	"StructForge/backend/apps/user/internal/conf"
	"StructForge/backend/common/data/database"
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
	flag.StringVar(&flagconf, "conf", "", "config path, eg: -conf user.yaml")
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
		configPath = "../../../../configs/local/user.yaml"
	case "test":
		configPath = "../../../../configs/test/user.yaml"
	case "prod":
		configPath = "../../../../configs/prod/user.yaml"
	default:
		configPath = "../../../../configs/local/user.yaml"
	}

	return configPath
}

func main() {
	flag.Parse()

	ctx := context.Background()

	// ========== 第一步：设置最大 CPU 核心数 ==========
	_, _ = maxprocs.Set()

	// ========== 第二步：初始化日志系统 ==========
	logConfig := log.DefaultConfig()
	logConfig.ServiceName = "user"
	logConfig.Level = log.InfoLevel

	logger, err := log.NewLogger(logConfig)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	log.SetGlobalLogger(logger)

	log.Info(ctx, "user 服务启动中")

	// ========== 第三步：加载配置文件 ==========
	configPath := getConfigPath()
	log.Info(ctx, "开始加载配置文件",
		log.String("config_path", configPath),
		log.String("environment", flagenv),
	)

	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		log.Error(ctx, "配置文件加载失败",
			log.ErrorField(err),
			log.String("config_path", configPath),
		)
		os.Exit(1)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.Error(ctx, "配置文件解析失败",
			log.ErrorField(err),
		)
		os.Exit(1)
	}

	log.Info(ctx, "启动配置文件加载成功",
		log.String("config_path", configPath),
	)

	// ========== 第四步：初始化 Nacos 配置中心（可选）==========
	if bc.Nacos != nil && bc.Nacos.ConfigCenter != nil && bc.Nacos.ConfigCenter.Enabled {
		log.Info(ctx, "正在创建 Nacos 配置客户端")

		// 转换服务器配置
		serverConfigs := make([]nacosClient.ServerConfig, 0, len(bc.Nacos.ServerConfigs))
		for _, sc := range bc.Nacos.ServerConfigs {
			serverConfigs = append(serverConfigs, nacosClient.ServerConfig{
				IpAddr:      sc.IpAddr,
				Port:        sc.Port,
				ContextPath: sc.ContextPath,
				Scheme:      sc.Scheme,
			})
		}

		// 转换客户端配置
		clientConfig := nacosClient.ClientConfig{
			NamespaceId:         bc.Nacos.ClientConfig.NamespaceId,
			TimeoutMs:           bc.Nacos.ClientConfig.TimeoutMs,
			NotLoadCacheAtStart: bc.Nacos.ClientConfig.NotLoadCacheAtStart,
			LogDir:              bc.Nacos.ClientConfig.LogDir,
			CacheDir:            bc.Nacos.ClientConfig.CacheDir,
			LogLevel:            bc.Nacos.ClientConfig.LogLevel,
			Username:            bc.Nacos.ClientConfig.Username,
			Password:            bc.Nacos.ClientConfig.Password,
			AccessKey:           bc.Nacos.ClientConfig.AccessKey,
			SecretKey:           bc.Nacos.ClientConfig.SecretKey,
		}

		// 转换配置中心配置
		configCenter := nacosClient.ConfigCenter{
			Enabled:   bc.Nacos.ConfigCenter.Enabled,
			DataId:    bc.Nacos.ConfigCenter.DataId,
			Group:     bc.Nacos.ConfigCenter.Group,
			Namespace: bc.Nacos.ConfigCenter.Namespace,
		}

		// 构建 StartupConfig（Nacos 客户端需要的配置格式）
		startupConfig := nacosClient.StartupConfig{
			Nacos: nacosClient.NacosConfig{
				ServerConfigs: serverConfigs,
				ClientConfig:  clientConfig,
				ConfigCenter:  configCenter,
			},
		}

		// 创建 Nacos 配置客户端
		configClient, err := nacosClient.NewNacosConfigClient(&startupConfig)
		if err != nil {
			log.Error(ctx, "Nacos 配置客户端创建失败",
				log.ErrorField(err),
			)
			os.Exit(1)
		}

		log.Info(ctx, "Nacos 配置客户端创建成功")

		// 从 Nacos 获取完整配置
		log.Info(ctx, "正在从 Nacos 获取完整配置")
		nacosConfigData, err := configClient.GetConfig(
			bc.Nacos.ConfigCenter.DataId,
			bc.Nacos.ConfigCenter.Group,
		)
		if err != nil {
			log.Error(ctx, "从 Nacos 获取配置失败",
				log.ErrorField(err),
			)
			os.Exit(1)
		}

		log.Info(ctx, "从 Nacos 获取配置成功",
			log.String("data_id", bc.Nacos.ConfigCenter.DataId),
			log.String("group", bc.Nacos.ConfigCenter.Group),
		)

		// 合并 Nacos 配置到本地配置
		if nacosConfigData != "" {
			// 这里可以解析 Nacos 配置并合并到 bc
			// 暂时跳过，直接使用本地配置
		}
	}

	log.Info(ctx, "配置加载完成",
		log.String("config_path", configPath),
	)

	// ========== 第五步：初始化数据库 ==========
	log.Info(ctx, "正在初始化数据库系统")

	// 始终从配置文件加载数据库配置（配置文件包含完整的数据库配置）
	db, err := database.InitDatabaseFromFile(configPath)
	if err != nil {
		log.Error(ctx, "数据库初始化失败",
			log.ErrorField(err),
			log.String("config_path", configPath),
		)
		os.Exit(1)
	}

	// 创建关闭函数
	dbCleanup := func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Close(); err != nil {
			log.Warn(cleanupCtx, "关闭数据库连接失败",
				log.ErrorField(err),
			)
		} else {
			log.Info(cleanupCtx, "数据库连接已关闭")
		}
	}
	defer dbCleanup()

	// ========== 第六步：使用 Wire 进行依赖注入，创建应用实例 ==========
	log.Info(ctx, "正在初始化应用实例")

	app, cleanup, err := wireApp(&bc)
	if err != nil {
		log.Error(ctx, "user 服务初始化失败",
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

	log.Info(ctx, "user 服务开始运行",
		log.String("server_id", serverID),
	)

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动应用
	if err := app.Run(); err != nil {
		log.Error(ctx, "user 服务运行失败",
			log.ErrorField(err),
		)
		os.Exit(1)
	}

	// 等待中断信号
	<-quit

	log.Info(ctx, "收到停止信号，正在关闭服务...")

	// 优雅关闭（Kratos v2 的 Stop 不接受参数）
	if err := app.Stop(); err != nil {
		log.Error(ctx, "服务关闭失败",
			log.ErrorField(err),
		)
	} else {
		log.Info(ctx, "服务已优雅关闭")
	}
}
