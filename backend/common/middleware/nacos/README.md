# Nacos 中间件

本包提供了 Nacos 服务发现和配置管理的功能。

## 功能特性

- **配置管理**: 从 Nacos 配置中心获取和管理配置
- **服务发现**: 服务注册、发现和订阅
- **配置监听**: 支持配置变化监听和自动更新

## 使用方法

### 1. 配置结构

在配置文件中定义 Nacos 配置：

```yaml
nacos:
  server_configs:
    - ip_addr: "127.0.0.1"
      port: 8848
      context_path: "/nacos"
      scheme: "http"
  client_config:
    namespace_id: "public"
    timeout_ms: 5000
    not_load_cache_at_start: true
    log_dir: "tmp/nacos/log"
    cache_dir: "tmp/nacos/cache"
    log_level: "info"
    username: ""  # 可选
    password: ""   # 可选
  config_center:
    enabled: true
    data_id: "gateway-config"
    group: "DEFAULT_GROUP"
    namespace: "public"
```

### 2. 创建配置客户端

```go
import (
    nacosClient "StructForge/backend/common/middleware/nacos"
)

// 加载启动配置
var startupConfig nacosClient.StartupConfig
// ... 从配置文件加载 startupConfig

// 创建 Nacos 配置客户端
configClient, err := nacosClient.NewNacosConfigClient(&startupConfig)
if err != nil {
    // 处理错误
}
```

### 3. 使用配置管理器

```go
// 定义 Bootstrap 配置结构
var bc conf.Bootstrap

// 创建配置管理器（会自动从 Nacos 加载配置并扫描到 bc）
dm, err := nacosClient.NewConfigManager(configClient, &startupConfig, &bc)
if err != nil {
    // 处理错误
}
defer dm.Close()

// bc 现在已经包含了从 Nacos 加载的配置
```

### 4. 服务发现

```go
// 创建服务发现客户端
namingClient, err := nacosClient.NewNacosNamingClient(&startupConfig)
if err != nil {
    // 处理错误
}

// 注册服务实例
success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
    Ip:          "127.0.0.1",
    Port:        8080,
    ServiceName: "gateway-service",
    GroupName:   "DEFAULT_GROUP",
    Weight:      1.0,
    Enable:      true,
    Healthy:     true,
    Ephemeral:   true,
})

// 获取服务实例
instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
    ServiceName: "gateway-service",
    GroupName:   "DEFAULT_GROUP",
    HealthyOnly: true,
})

// 选择一个健康的实例
instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
    ServiceName: "gateway-service",
    GroupName:   "DEFAULT_GROUP",
})
```

### 5. 配置监听

```go
// 监听配置变化
err := configClient.ListenConfig("gateway-config", "DEFAULT_GROUP", func(namespace, group, dataId, data string) {
    // 配置变化时的回调
    // 可以在这里重新加载配置
})
```

## API 参考

### NacosConfigClient

- `NewNacosConfigClient(startupConfig *StartupConfig) (*NacosConfigClient, error)`: 创建配置客户端
- `GetConfig(dataId, group string) (string, error)`: 获取配置
- `PublishConfig(dataId, group, content string) (bool, error)`: 发布配置
- `DeleteConfig(dataId, group string) (bool, error)`: 删除配置
- `ListenConfig(dataId, group string, onChange func(...)) error`: 监听配置变化

### NacosNamingClient

- `NewNacosNamingClient(startupConfig *StartupConfig) (*NacosNamingClient, error)`: 创建服务发现客户端
- `RegisterInstance(param vo.RegisterInstanceParam) (bool, error)`: 注册服务实例
- `DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error)`: 注销服务实例
- `GetService(param vo.GetServiceParam) (vo.Service, error)`: 获取服务
- `SelectInstances(param vo.SelectInstancesParam) ([]vo.Instance, error)`: 选择服务实例
- `SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*vo.Instance, error)`: 选择一个健康的实例
- `Subscribe(param *vo.SubscribeParam) error`: 订阅服务变化
- `Unsubscribe(param *vo.SubscribeParam) error`: 取消订阅

### ConfigManager

- `NewConfigManager(configClient *NacosConfigClient, startupConfig *StartupConfig, bootstrapPtr interface{}) (*ConfigManager, error)`: 创建配置管理器
- `GetBootstrap() interface{}`: 获取 Bootstrap 配置
- `GetConfig() config.Config`: 获取底层配置对象
- `Watch(key string, callback func(string, config.Value)) error`: 监听配置变化
- `Value(key string) (config.Value, error)`: 获取配置值
- `Load() error`: 重新加载配置
- `Close() error`: 关闭配置管理器

## 依赖

- `github.com/nacos-group/nacos-sdk-go`: Nacos Go SDK
- `github.com/go-kratos/kratos/v2/config`: Kratos 配置框架
- `github.com/go-kratos/kratos/v2/config/nacos`: Kratos Nacos 配置源

## 注意事项

1. 确保 Nacos 服务器正在运行
2. 配置中心的 `enabled` 必须为 `true` 才能使用配置管理功能
3. 服务发现和配置管理使用相同的客户端配置
4. 建议在生产环境中使用命名空间进行隔离

