# 阶段3：依赖注入（Wire）

## 📚 学习目标
- 理解依赖注入的概念和好处
- 掌握Wire的使用方法
- 学习Provider和ProviderSet
- 理解代码生成的原理

---

## 🎯 什么是依赖注入？

### 传统方式（手动创建依赖）

```go
// 不好的方式：手动创建依赖
func main() {
    // 创建Redis客户端
    redisClient := redis.NewClient(...)
    
    // 创建Handler，需要传入Redis
    handler := NewHandler(redisClient)
    
    // 创建Server，需要传入Handler
    server := NewServer(handler)
    
    // 启动服务
    server.Start()
}
```

**问题：**
- 依赖关系复杂时，代码难以维护
- 测试时难以替换依赖（如用Mock替换Redis）
- 创建顺序必须正确

### 依赖注入方式（Wire自动管理）

```go
// 好的方式：Wire自动管理依赖
func main() {
    // Wire自动创建所有依赖，按正确顺序
    app, cleanup, err := wireApp(config, redisConfig)
    if err != nil {
        panic(err)
    }
    defer cleanup()
    
    app.Run()
}
```

**优点：**
- 依赖关系清晰
- 易于测试（可以轻松替换依赖）
- 自动处理创建顺序

---

## 🔧 Wire核心概念

### 1. Provider（提供者）

**Provider是一个函数，用于创建某个对象：**

```go
// Provider函数：创建Redis客户端
func NewRedisClient(config *RedisConfig) (redis.UniversalClient, func(), error) {
    client := redis.NewClient(&redis.Options{
        Addr: config.Addr,
    })
    
    cleanup := func() {
        client.Close()
    }
    
    return client, cleanup, nil
}
```

**Provider函数的特点：**
- 函数名通常以 `New` 开头
- 返回值可以是：`(对象, error)` 或 `(对象, cleanup函数, error)`
- cleanup函数用于资源清理

### 2. ProviderSet（提供者集合）

**将多个Provider组合在一起：**

```go
// 数据层的ProviderSet
var ProviderSet = wire.NewSet(
    redis.NewRedisClient,  // Redis客户端
    NewDatabaseClient,     // 数据库客户端
    NewCacheClient,        // 缓存客户端
)
```

**使用ProviderSet的好处：**
- 模块化：相关Provider放在一起
- 复用：可以在多个地方使用同一个ProviderSet

---

## 📝 核心代码：wire.go

### 完整代码（手写练习）

```go
//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
    "HeronGame/apps/gateway/internal/conf"
    "HeronGame/apps/gateway/internal/data"
    "HeronGame/apps/gateway/internal/handler"
    "HeronGame/apps/gateway/internal/manager"
    "HeronGame/apps/gateway/internal/remote"
    "HeronGame/apps/gateway/internal/server"
    redisconf "HeronGame/common/data/redis/conf"
    glog "HeronGame/common/log"

    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
    server.ProviderSet,    // 服务器Provider
    data.ProviderSet,      // 数据层Provider
    handler.ProviderSet,   // 处理器Provider
    manager.ProviderSet,    // 管理器Provider
    remote.ProviderSet,     // 远程调用Provider
    NewLogger,             // 日志Provider
)

// NewLogger creates a new logger instance
func NewLogger() log.Logger {
    return glog.ZapToKratosLogger()
}

// wireApp init kratos application.
// 这个函数会被Wire自动生成实现
func wireApp(*conf.Bootstrap, *redisconf.Redis) (*kratos.App, func(), error) {
    panic(wire.Build(ProviderSet, newApp))
}
```

---

## 📖 代码详解

### 1. Build Tag（构建标签）

```go
//go:build wireinject
// +build wireinject
```

**作用：**
- 这个文件只在生成Wire代码时编译
- 正常编译时不会包含这个文件
- 避免编译错误（因为函数体是 `panic`）

### 2. ProviderSet组合

```go
var ProviderSet = wire.NewSet(
    server.ProviderSet,    // 来自server包的ProviderSet
    data.ProviderSet,       // 来自data包的ProviderSet
    handler.ProviderSet,    // 来自handler包的ProviderSet
    manager.ProviderSet,    // 来自manager包的ProviderSet
    remote.ProviderSet,     // 来自remote包的ProviderSet
    NewLogger,              // 当前文件的Provider
)
```

**知识点：**
- `wire.NewSet()`：创建ProviderSet
- ProviderSet可以包含其他ProviderSet（嵌套）
- 可以包含单个Provider函数

### 3. wire.Build

```go
func wireApp(*conf.Bootstrap, *redisconf.Redis) (*kratos.App, func(), error) {
    panic(wire.Build(ProviderSet, newApp))
}
```

**作用：**
- `wire.Build()`：告诉Wire需要生成什么代码
- 参数：所有需要的Provider和最终目标函数（`newApp`）
- Wire会分析依赖关系，生成 `wire_gen.go`

**函数签名说明：**
- 参数：`*conf.Bootstrap, *redisconf.Redis` - 这些是外部提供的（从main传入）
- 返回值：`*kratos.App, func(), error` - Wire会生成代码来创建这些

---

## 🔍 子模块的ProviderSet示例

### server包的ProviderSet

```go
// apps/gateway/internal/server/server.go
package server

import (
    "github.com/google/wire"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
    selector.Provider,      // 选择器Provider
    NewHTTPServer,         // HTTP服务器Provider
    NewGRPCServer,         // gRPC服务器Provider
    NewWebSocketServer,    // WebSocket服务器Provider
)
```

### data包的ProviderSet

```go
// apps/gateway/internal/data/provider.go
package data

import (
    "HeronGame/common/data"
    "github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
    data.ProviderSet,           // 公共数据层ProviderSet
    client.ProviderSet,        // 客户端ProviderSet
    ratelimit.NewRateLimitManager, // 限流管理器Provider
)
```

### handler包的ProviderSet

```go
// apps/gateway/internal/handler/handler.go
package handler

import (
    "github.com/google/wire"
)

// ProviderSet is handler providers.
var ProviderSet = wire.NewSet(
    NewMessageHandler,    // 消息处理器Provider
    NewCommonHandler,     // 通用处理器Provider
    NewMessageRouter,     // 消息路由Provider
    NewGatewayHandler,    // 网关处理器Provider
    NewChatHandler,       // 聊天处理器Provider
    // ... 其他Handler
)
```

---

## 🛠️ Wire代码生成

### 生成命令

```bash
# 在wire.go所在目录执行
go generate ./...

# 或者直接运行wire命令
go run github.com/google/wire/cmd/wire
```

### 生成的代码（wire_gen.go）

Wire会自动分析依赖关系，生成类似这样的代码：

```go
// Code generated by Wire. DO NOT EDIT.

func wireApp(bootstrap *conf.Bootstrap, confRedis *redisconf.Redis) (*kratos.App, func(), error) {
    // 1. 创建Redis客户端
    universalClient, cleanup, err := redis.NewRedisClient(confRedis)
    if err != nil {
        return nil, nil, err
    }
    
    // 2. 创建Nacos客户端
    iNamingClient, err := client.NewNacosClient(bootstrap)
    if err != nil {
        cleanup()  // 清理已创建的资源
        return nil, nil, err
    }
    
    // 3. 创建连接管理器
    connectionManager, cleanup2 := manager.NewConnectionManager(bootstrap)
    
    // 4. 创建各种Handler
    chatHandler := handler.NewChatHandler(...)
    
    // 5. 创建服务器
    grpcServer := server.NewGRPCServer(...)
    wsServer := server.NewWebSocketServer(...)
    
    // 6. 创建Kratos应用
    app := newApp(bootstrap, grpcServer, wsServer, registry)
    
    // 7. 返回应用和清理函数
    cleanupFunc := func() {
        cleanup2()
        cleanup()
    }
    
    return app, cleanupFunc, nil
}
```

**Wire的智能之处：**
- 自动分析依赖关系
- 按正确顺序创建对象
- 自动处理错误和资源清理
- 如果依赖缺失，编译时就会报错

---

## 💡 Provider函数编写规范

### 标准Provider函数

```go
// 方式1：只有对象和错误
func NewUserService(repo UserRepository) (*UserService, error) {
    return &UserService{repo: repo}, nil
}

// 方式2：对象、清理函数和错误（推荐）
func NewRedisClient(config *RedisConfig) (redis.UniversalClient, func(), error) {
    client := redis.NewClient(&redis.Options{
        Addr: config.Addr,
    })
    
    cleanup := func() {
        if err := client.Close(); err != nil {
            log.Error("关闭Redis失败", err)
        }
    }
    
    return client, cleanup, nil
}
```

### Provider函数参数

```go
// Wire会自动注入参数
func NewHandler(
    redisClient redis.UniversalClient,  // Wire会自动找到Redis的Provider
    dbClient *sql.DB,                    // Wire会自动找到DB的Provider
    config *conf.Bootstrap,              // 需要从外部传入
) (*Handler, error) {
    return &Handler{
        redis: redisClient,
        db:    dbClient,
        config: config,
    }, nil
}
```

**规则：**
- 如果参数类型有对应的Provider，Wire会自动注入
- 如果参数类型没有Provider，需要从外部传入（如config）

---

## 🎓 依赖注入的优势

### 1. 解耦合

```go
// 传统方式：Handler直接依赖Redis
type Handler struct {
    redis *redis.Client  // 强依赖
}

// 依赖注入：Handler依赖接口
type Handler struct {
    cache Cache  // 依赖接口，可以替换实现
}
```

### 2. 易于测试

```go
// 测试时可以用Mock替换真实依赖
func TestHandler(t *testing.T) {
    mockCache := &MockCache{}  // Mock对象
    handler := NewHandler(mockCache)  // 注入Mock
    // 测试...
}
```

### 3. 清晰的依赖关系

```go
// 从ProviderSet就能看出依赖关系
var ProviderSet = wire.NewSet(
    NewRedisClient,      // 需要RedisConfig
    NewHandler,           // 需要RedisClient
    NewServer,            // 需要Handler
)
```

---

## 🔧 实践练习

### 练习1：理解依赖关系

画出以下依赖关系图：
```
Bootstrap → RedisClient
Bootstrap → NacosClient
RedisClient + NacosClient → ConnectionManager
ConnectionManager → Handlers
Handlers + Servers → KratosApp
```

### 练习2：编写Provider函数

尝试编写一个简单的Provider函数：

```go
// 创建一个UserService的Provider
func NewUserService(repo UserRepository) (*UserService, error) {
    // TODO: 实现创建逻辑
}
```

### 练习3：理解Wire生成过程

1. 查看 `wire.go` 文件
2. 运行 `go generate` 生成 `wire_gen.go`
3. 对比生成前后的代码，理解Wire做了什么

---

## 📌 下一阶段预告

**阶段4：服务器配置（HTTP/gRPC/WebSocket）**
- 学习如何创建HTTP服务器
- 理解gRPC服务器的配置
- 掌握WebSocket服务器的实现

---

## ❓ 思考题

1. 为什么要在 `wire.go` 中使用 `panic(wire.Build(...))`？
2. Provider函数的返回值中，cleanup函数的作用是什么？
3. 如果两个Provider返回相同类型，Wire如何区分？
4. 依赖注入和工厂模式有什么区别？

---

**完成本阶段后，请继续学习阶段4！** 🚀

