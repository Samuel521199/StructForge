# 数据库抽象层

> **版本**: v1.0  
> **状态**: 生产就绪 ✅

## 概述

数据库抽象层提供了统一的数据库接口，支持多种数据库后端（PostgreSQL、MySQL、SQLite等）。设计参考了 `common/cache` 和 `common/log` 系统的架构模式，保持代码风格和架构一致性。

**核心特性**：
- ✅ 支持多种数据库（PostgreSQL、MySQL、SQLite）
- ✅ 基于 GORM，提供强大的 ORM 功能
- ✅ 统一的接口，业务代码与具体数据库无关
- ✅ 支持从配置文件（YAML）加载配置
- ✅ 支持环境变量配置
- ✅ 连接池管理
- ✅ 慢查询监控
- ✅ 统一日志集成
- ✅ 优雅关闭

## 快速开始

### 1. 安装依赖

```bash
cd backend
go get gorm.io/gorm
go get gorm.io/driver/postgres    # PostgreSQL
go get gorm.io/driver/mysql       # MySQL
go get gorm.io/driver/sqlite       # SQLite
```

### 2. 方式一：从配置文件初始化（推荐）

#### 2.1 创建配置文件

在 `backend/configs/local/` 目录下创建 `database.yaml`：

```yaml
database:
  # 适配器类型：postgres, mysql, sqlite
  adapter_type: postgres

  # 通用配置
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h
  conn_max_idle_time: 10m
  log_level: warn          # silent, error, warn, info
  slow_threshold: 200ms    # 慢查询阈值

  # PostgreSQL 配置
  postgres:
    host: localhost
    port: 5432
    user: postgres
    password: your_password
    dbname: structforge
    sslmode: disable
    timezone: Asia/Shanghai
    # 或者直接使用 DSN
    # dsn: "host=localhost user=postgres password=password dbname=structforge port=5432 sslmode=disable TimeZone=Asia/Shanghai"

  # MySQL 配置（示例）
  # mysql:
  #   host: localhost
  #   port: 3306
  #   user: root
  #   password: your_password
  #   dbname: structforge
  #   charset: utf8mb4
  #   parse_time: true
  #   loc: Local

  # SQLite 配置（示例）
  # sqlite:
  #   path: data.db
  #   foreign_keys: true
  #   busy_timeout: 5s
```

#### 2.2 在代码中使用

```go
package main

import (
    "context"
    "StructForge/backend/common/data/database"
    "StructForge/backend/common/log"
)

func main() {
    ctx := context.Background()
    
    // 初始化日志系统（必需）
    logConfig := log.DefaultConfig()
    logConfig.ServiceName = "my-service"
    logger, _ := log.NewLogger(logConfig)
    log.SetGlobalLogger(logger)
    defer logger.Sync()
    
    // 从配置文件初始化数据库（推荐方式）
    defer database.InitDatabaseFromFileWithShutdown("configs/local/database.yaml")()
    
    // 使用数据库
    db := database.GetDB()
    if db != nil {
        // 执行数据库操作
        var users []User
        db.Find(&users)
    }
}
```

### 3. 方式二：使用代码配置

```go
package main

import (
    "context"
    "time"
    "StructForge/backend/common/data/database"
    "StructForge/backend/common/log"
)

func main() {
    ctx := context.Background()
    
    // 初始化日志系统
    logConfig := log.DefaultConfig()
    logConfig.ServiceName = "my-service"
    logger, _ := log.NewLogger(logConfig)
    log.SetGlobalLogger(logger)
    defer logger.Sync()
    
    // 使用代码配置初始化数据库
    defer database.InitDatabaseWithShutdown("postgres",
        database.WithPostgreSQL(&database.PostgreSQLConfig{
            Host:     "localhost",
            Port:     5432,
            User:     "postgres",
            Password: "password",
            DBName:   "structforge",
            SSLMode:  "disable",
            TimeZone: "Asia/Shanghai",
        }),
        database.WithMaxOpenConns(100),
        database.WithMaxIdleConns(10),
        database.WithLogLevel("info"),
    )()
    
    // 使用数据库
    db := database.GetDB()
    // ...
}
```

### 4. 方式三：使用环境变量

```bash
# 设置环境变量
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=structforge
```

```go
// 使用默认配置（从环境变量读取）
defer database.InitDatabaseWithShutdown("postgres")()
```

## 核心功能

### 基础操作

```go
// 获取 GORM 实例
db := database.GetDB()

// 查询
var user User
db.First(&user, 1)

// 创建
db.Create(&User{Name: "John"})

// 更新
db.Model(&user).Update("name", "Jane")

// 删除
db.Delete(&user)
```

### 事务操作

```go
ctx := context.Background()
db := database.GetGlobalDB()

err := db.Transaction(ctx, func(tx *gorm.DB) error {
    // 在事务中执行操作
    if err := tx.Create(&user1).Error; err != nil {
        return err
    }
    if err := tx.Create(&user2).Error; err != nil {
        return err
    }
    return nil
})
```

### 数据库迁移

```go
ctx := context.Background()
db := database.GetGlobalDB()

// 自动迁移（根据模型创建/更新表结构）
type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string
}

err := db.AutoMigrate(ctx, &User{})
```

### 健康检查

```go
ctx := context.Background()
db := database.GetGlobalDB()

if err := db.Ping(ctx); err != nil {
    log.Error(ctx, "数据库连接失败", log.ErrorField(err))
} else {
    log.Info(ctx, "数据库连接正常")
}
```

## 配置选项

### 适配器类型

- `postgres` - PostgreSQL
- `mysql` - MySQL
- `sqlite` - SQLite

### 通用配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `max_open_conns` | int | 100 | 最大打开连接数 |
| `max_idle_conns` | int | 10 | 最大空闲连接数 |
| `conn_max_lifetime` | duration | 1h | 连接最大生存时间 |
| `conn_max_idle_time` | duration | 10m | 连接最大空闲时间 |
| `log_level` | string | warn | 日志级别：silent, error, warn, info |
| `slow_threshold` | duration | 200ms | 慢查询阈值 |

### PostgreSQL 配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `host` | string | localhost | 数据库主机 |
| `port` | int | 5432 | 数据库端口 |
| `user` | string | postgres | 用户名 |
| `password` | string | - | 密码 |
| `dbname` | string | structforge | 数据库名 |
| `sslmode` | string | disable | SSL 模式 |
| `timezone` | string | Asia/Shanghai | 时区 |
| `dsn` | string | - | 完整 DSN（可选，覆盖其他配置） |

### MySQL 配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `host` | string | localhost | 数据库主机 |
| `port` | int | 3306 | 数据库端口 |
| `user` | string | root | 用户名 |
| `password` | string | - | 密码 |
| `dbname` | string | structforge | 数据库名 |
| `charset` | string | utf8mb4 | 字符集 |
| `parse_time` | bool | true | 解析时间 |
| `loc` | string | Local | 时区 |
| `dsn` | string | - | 完整 DSN（可选，覆盖其他配置） |

### SQLite 配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `path` | string | data.db | 数据库文件路径 |
| `foreign_keys` | bool | true | 启用外键约束 |
| `busy_timeout` | duration | 5s | 忙等待超时 |
| `dsn` | string | - | 完整 DSN（可选，覆盖其他配置） |

## 接口说明

### Database 接口

所有数据库适配器都实现了 `Database` 接口：

```go
type Database interface {
    // GetDB 获取 GORM 数据库实例
    GetDB() *gorm.DB

    // Ping 检查数据库连接
    Ping(ctx context.Context) error

    // Close 关闭数据库连接
    Close() error

    // Health 健康检查
    Health(ctx context.Context) error

    // Migrate 执行数据库迁移
    Migrate(ctx context.Context, models ...interface{}) error

    // AutoMigrate 自动迁移（根据模型创建/更新表结构）
    AutoMigrate(ctx context.Context, models ...interface{}) error

    // Transaction 执行事务
    Transaction(ctx context.Context, fn func(*gorm.DB) error) error
}
```

## 最佳实践

### 1. 使用配置文件

优先使用配置文件方式，便于管理和部署：

```yaml
# configs/local/database.yaml
database:
  adapter_type: postgres
  postgres:
    host: localhost
    # ...
```

### 2. 优雅关闭

使用 `InitDatabaseFromFileWithShutdown` 确保数据库连接正确关闭：

```go
defer database.InitDatabaseFromFileWithShutdown("configs/local/database.yaml")()
```

### 3. 使用事务

对于需要原子性的操作，使用事务：

```go
err := db.Transaction(ctx, func(tx *gorm.DB) error {
    // 多个操作
    return nil
})
```

### 4. 连接池配置

根据实际负载调整连接池大小：

```yaml
database:
  max_open_conns: 100    # 高并发场景可以增大
  max_idle_conns: 10     # 保持合理的空闲连接数
```

### 5. 慢查询监控

启用慢查询监控，及时发现性能问题：

```yaml
database:
  log_level: info
  slow_threshold: 200ms
```

## 与项目配置系统集成

### 在服务配置文件中添加数据库配置

在服务的配置文件中（如 `backend/configs/local/user.yaml`）：

```yaml
# 服务配置
server:
  id: user-service
  name: user
  version: v1.0.0

# 数据库配置
database:
  adapter_type: postgres
  postgres:
    host: localhost
    port: 5432
    user: postgres
    password: password
    dbname: structforge
    sslmode: disable
    timezone: Asia/Shanghai
```

### 在服务启动代码中使用

```go
package main

import (
    "context"
    "StructForge/backend/common/data/database"
    "StructForge/backend/common/log"
)

func main() {
    ctx := context.Background()
    
    // 初始化日志
    // ...
    
    // 从配置文件初始化数据库
    defer database.InitDatabaseFromFileWithShutdown("configs/local/user.yaml")()
    
    // 启动服务
    // ...
}
```

## 故障排查

### 1. 连接失败

检查：
- 数据库服务是否运行
- 配置的主机、端口是否正确
- 用户名、密码是否正确
- 防火墙是否允许连接

### 2. 慢查询

启用 `log_level: info` 查看详细日志，调整 `slow_threshold` 阈值。

### 3. 连接池耗尽

增大 `max_open_conns` 配置，或检查是否有连接泄漏。

## 扩展

### 添加新的数据库适配器

1. 在 `adapters/` 目录下创建新适配器目录
2. 实现 `Database` 接口
3. 在 `NewDatabase` 函数中添加新的 case

示例：

```go
// adapters/mongodb/adapter.go
package mongodb

import (
    "StructForge/backend/common/data/database"
)

type Adapter struct {
    // ...
}

func NewAdapter(config *database.Config) (*Adapter, error) {
    // 实现适配器
}
```

然后在 `init.go` 的 `NewDatabase` 函数中添加：

```go
case AdapterMongoDB:
    return mongodbAdapter.NewAdapter(&config)
```

## 相关文档

- [GORM 官方文档](https://gorm.io/docs/)
- [PostgreSQL 驱动](https://github.com/go-gorm/postgres)
- [MySQL 驱动](https://github.com/go-gorm/mysql)
- [SQLite 驱动](https://github.com/go-gorm/sqlite)

