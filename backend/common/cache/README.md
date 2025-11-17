# 通用缓存系统

> **版本**: v1.0  
> **状态**: 生产就绪 ✅

## 概述

通用缓存系统提供了统一的缓存接口，支持多种缓存后端（Redis、内存缓存等）。设计参考了 `common/log` 系统的架构模式，保持代码风格和架构一致性。

## 快速开始

### 1. 初始化缓存系统（最简单方式）

```go
package main

import (
    "context"
    "StructForge/backend/common/cache"
)

func main() {
    // 一行代码初始化，自动配置，自动优雅关闭
    defer cache.InitCacheWithShutdown("memory")()
    
    // 使用缓存
    ctx := context.Background()
    cache.Set(ctx, "key", []byte("value"), 10*time.Minute)
    data, err := cache.Get(ctx, "key")
}
```

### 2. 使用Redis缓存

```go
defer cache.InitCacheWithShutdown("redis",
    cache.WithKeyPrefix("app:"),
    cache.WithDefaultTTL(5 * time.Minute),
    cache.WithRedisAddr("localhost:6379"),
)()
```

### 3. 使用内存缓存

```go
defer cache.InitCacheWithShutdown("memory",
    cache.WithKeyPrefix("app:"),
    cache.WithMemoryMaxSize(100 * 1024 * 1024), // 100MB
    cache.WithMemoryMaxItems(1000),
)()
```

## 核心功能

### 基础操作

- **Get**: 获取缓存值
- **Set**: 设置缓存值
- **Delete**: 删除缓存键
- **Exists**: 检查键是否存在

### 批量操作

- **MGet**: 批量获取
- **MSet**: 批量设置
- **MDelete**: 批量删除

### 高级操作

- **GetOrSet**: 获取或设置（缓存穿透保护）
- **Increment/Decrement**: 数值增减
- **Expire/TTL**: 过期时间管理
- **Keys/Scan**: 键扫描
- **Clear**: 清除匹配的键

### 对象操作

```go
// 自动序列化/反序列化
var user User
err := cache.GetObject(ctx, "user:123", &user)
err := cache.SetObject(ctx, "user:123", user, 10*time.Minute)
```

## 适配器

### Redis适配器

**注意**: 需要先安装依赖：
```bash
go get github.com/redis/go-redis/v9
```

然后取消注释 `adapters/redis/adapter.go` 中的实现代码。

**配置**:
- `Addr`: Redis地址
- `Password`: 密码
- `DB`: 数据库编号
- `PoolSize`: 连接池大小
- `MinIdleConns`: 最小空闲连接数

### 内存适配器

**特性**:
- 支持LRU、LFU、FIFO淘汰策略
- 内存限制和条目数限制
- 自动清理过期项
- 线程安全

**配置**:
- `MaxSize`: 最大内存使用（字节）
- `MaxItems`: 最大条目数
- `Strategy`: 淘汰策略（"lru", "lfu", "fifo"）
- `CleanupInterval`: 清理间隔

## 配置选项

### 函数式选项

```go
cache.InitCache("redis",
    cache.WithKeyPrefix("app:"),              // 键前缀
    cache.WithDefaultTTL(5 * time.Minute),    // 默认TTL
    cache.WithSerializer(cache.SerializerJSON), // 序列化器
    cache.WithRedisAddr("localhost:6379"),    // Redis地址
    cache.WithRedisPassword("password"),      // Redis密码
    cache.WithRedisDB(0),                     // Redis数据库
    cache.WithMemoryMaxSize(100 * 1024 * 1024), // 内存最大大小
    cache.WithMemoryMaxItems(1000),           // 内存最大条目数
)
```

### 环境变量

- `REDIS_ADDR`: Redis地址（默认: localhost:6379）
- `REDIS_PASSWORD`: Redis密码

## 错误处理

标准错误定义：
- `ErrNotFound`: 键不存在
- `ErrExpired`: 键已过期
- `ErrSerialization`: 序列化错误
- `ErrConnection`: 连接错误
- `ErrTimeout`: 操作超时

## 最佳实践

1. **键命名**: 使用有意义的键前缀，如 `user:123`, `workflow:456`
2. **TTL设置**: 根据数据特性设置合理的过期时间
3. **错误处理**: 缓存错误不应影响主流程
4. **批量操作**: 优先使用批量操作提高性能
5. **内存缓存**: 适合热点数据，Redis适合持久化数据

## 架构设计

### 目录结构

```
backend/common/cache/
├── cache.go              # 核心接口定义
├── adapter.go            # 适配器接口
├── config.go             # 配置结构
├── errors.go             # 错误定义
├── serialization.go      # 序列化接口
├── init.go               # 初始化函数
├── adapters/             # 适配器实现
│   ├── redis/            # Redis适配器
│   └── memory/           # 内存适配器
└── README.md             # 使用文档
```

### 设计模式

- **适配器模式**: 统一接口，多种实现
- **函数式选项**: 灵活的配置方式
- **全局单例**: 便捷的全局API

## 注意事项

1. **循环导入**: 适配器包不直接导入cache包，避免循环依赖
2. **Redis依赖**: Redis适配器需要安装 `github.com/redis/go-redis/v9`
3. **线程安全**: 所有适配器实现都是线程安全的
4. **资源清理**: 使用 `InitCacheWithShutdown` 确保资源正确清理

## 扩展

### 添加新适配器

1. 在 `adapters/` 目录下创建新适配器包
2. 实现 `CacheAdapter` 接口
3. 在 `init.go` 的 `NewCache` 函数中添加新适配器类型
4. 实现 `newXxxAdapter` 函数

## 示例

### 完整示例

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "StructForge/backend/common/cache"
)

func main() {
    // 初始化缓存
    defer cache.InitCacheWithShutdown("memory",
        cache.WithKeyPrefix("app:"),
        cache.WithDefaultTTL(5 * time.Minute),
    )()
    
    ctx := context.Background()
    
    // 设置缓存
    err := cache.Set(ctx, "user:123", []byte("user data"), 10*time.Minute)
    if err != nil {
        fmt.Printf("设置缓存失败: %v\n", err)
        return
    }
    
    // 获取缓存
    data, err := cache.Get(ctx, "user:123")
    if err != nil {
        fmt.Printf("获取缓存失败: %v\n", err)
        return
    }
    
    fmt.Printf("缓存值: %s\n", string(data))
    
    // GetOrSet示例
    val, err := cache.GetOrSet(ctx, "key", func() ([]byte, error) {
        // 如果缓存不存在，调用此函数获取值
        return []byte("computed value"), nil
    }, 5*time.Minute)
    
    fmt.Printf("值: %s\n", string(val))
}
```

---

**最后更新**: 2025-01-15  
**版本**: v1.0

