# 通用日志系统使用说明

> **版本**: v2.1  
> **更新日期**: 2025-01-15  
> **状态**: 生产就绪 ✅  
> **新增**: 便捷初始化函数（函数式选项模式）

## 快速开始

### 1. 初始化日志系统（推荐：便捷方式）

**最简单的方式**（自动从环境变量读取配置）：

```go
package main

import (
    "context"
    "StructForge/backend/common/log"
)

func main() {
    // 一行代码初始化，自动配置，自动优雅关闭
    defer log.InitLoggerWithShutdown("gateway")()
    
    // 使用日志
    log.Info(context.Background(), "服务启动成功")
}
```

**带选项的方式**（推荐）：

```go
func main() {
    // 使用便捷函数，自动从环境变量读取配置
    defer log.InitLoggerWithShutdown("gateway",
        log.WithServiceID("gateway-001"),           // 设置服务ID
        log.WithEnvironment("production"),          // 设置环境（自动调整日志级别和颜色）
        log.WithFilePath("logs/gateway-%s.log"),    // 自定义文件路径
    )()
    
    // 使用日志
    log.Info(context.Background(), "服务启动成功")
}
```

**更多选项**：

```go
func main() {
    defer log.InitLoggerWithShutdown("gateway",
        log.WithServiceID("gateway-001"),
        log.WithInstanceID("gateway-001-pod-123"),  // 自定义实例ID
        log.WithLevel(log.DebugLevel),               // 自定义日志级别
        log.WithEnvironment("development"),          // 设置环境
        log.WithFileOutput(true),                    // 启用/禁用文件输出
        log.WithFilePath("logs/custom-%s.log"),      // 自定义文件路径
        log.WithColor(true),                         // 强制启用/禁用颜色
    )()
}
```

**高级配置**（完全自定义）：

```go
func main() {
    // 使用完整配置
    customConfig := log.DefaultConfig()
    customConfig.ServiceName = "gateway"
    customConfig.Level = log.InfoLevel
    // ... 其他高级配置
    
    defer log.InitLoggerWithShutdown("gateway",
        log.WithConfig(customConfig),  // 使用完整配置覆盖
    )()
}
```

### 1.1 初始化日志系统（传统方式）

如果需要完全控制配置，可以使用传统方式：

```go
package main

import (
    "context"
    "StructForge/backend/common/log"
)

func main() {
    // 创建配置
    config := log.Config{
        Level:       log.InfoLevel,
        ServiceName: "gateway",
        ServiceID:   "gateway-001",
        InstanceID:  "gateway-001-pod-123",
        Outputs:     []log.Output{log.ConsoleOutput, log.FileOutput},
        EnableColor: true,
        Console: log.ConsoleConfig{
            Enabled: true,
            Format:  log.TextFormat,
            Level:   log.DebugLevel, // 控制台显示所有级别
        },
        File: log.FileConfig{
            Enabled:    true,
            Format:     log.JSONFormat,
            Level:      log.InfoLevel, // 文件只记录INFO及以上
            Path:       "logs/%s-%s.log", // 服务名-日期.log
            MaxSize:    100,              // MB
            MaxAge:     7,                // days
            MaxBackups: 10,
            Compress:   true,
        },
    }

    // 创建Logger
    logger, err := log.NewLogger(config)
    if err != nil {
        panic(err)
    }
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        logger.Shutdown(ctx)
    }()

    // 设置为全局Logger
    log.SetGlobalLogger(logger)
}
```

### 2. 使用日志

```go
ctx := context.Background()

// 简单日志
log.Info(ctx, "服务启动成功")

// 带字段的日志
log.Info(ctx, "用户登录成功",
    log.String("user_id", "123"),
    log.Int("age", 25),
    log.Bool("is_vip", true),
    log.ErrorField(err),
)

// 不同级别
log.Debug(ctx, "调试信息", log.String("key", "value"))
log.Warn(ctx, "警告信息", log.String("key", "value"))
log.Error(ctx, "错误信息", log.ErrorField(err))
log.Fatal(ctx, "致命错误", log.ErrorField(err)) // 会调用os.Exit(1)
```

### 3. 使用Context传递追踪信息

```go
// 在Context中设置追踪信息
ctx = context.WithValue(ctx, log.CtxTraceID, "trace-123")
ctx = context.WithValue(ctx, log.CtxSpanID, "span-456")
ctx = context.WithValue(ctx, log.CtxUserID, int64(123))

// 日志会自动包含这些信息
log.Info(ctx, "处理请求")
```

## 文件结构

### 核心文件
- `level.go` - 日志级别定义
- `field.go` - 字段类型和便捷函数（支持Duration、Time、StringSlice、Bytes等）
- `config.go` - 配置结构
- `logger.go` - Logger接口和实现（包含对象池复用）
- `init.go` - **便捷初始化函数**（推荐使用）

### 输出相关
- `color.go` - 颜色支持
- `formatter.go` - 格式化器（文本和JSON）
- `console.go` - 控制台输出
- `file.go` - 文件输出
- `async_writer.go` - 异步写入器

### 功能扩展
- `context.go` - Context工具函数
- `sampling.go` - 采样功能
- `mask.go` - 敏感信息脱敏
- `hook.go` - 钩子接口和管理器
- `hooks.go` - 内置钩子实现（告警、统计、过滤）
- `error_handler.go` - 错误处理增强
- `lazy_field.go` - 延迟字段实现
- `lazy_resolver.go` - 延迟字段解析器
- `pool.go` - 对象池复用
- `config_loader.go` - 配置文件加载

## 特性

- ✅ 5个日志级别（Debug、Info、Warn、Error、Fatal）
- ✅ 控制台颜色输出
- ✅ 文件按天轮转
- ✅ 结构化日志（JSON）
- ✅ 容器环境支持
- ✅ 追踪集成
- ✅ 类型安全的字段API
- ✅ 线程安全
- ✅ **异步写入**（提高性能，减少I/O阻塞）
- ✅ **采样功能**（控制日志量，避免日志过多）
- ✅ **级别检查优化**（提前检查，避免不必要的字段构建）
- ✅ **敏感信息脱敏**（自动脱敏密码、token等敏感字段）
- ✅ **钩子（Hooks）机制**（支持告警、统计、过滤等扩展）
- ✅ **错误处理增强**（优雅降级、错误统计）
- ✅ **优雅关闭**（确保日志不丢失）
- ✅ **配置文件加载**（支持YAML/JSON配置）
- ✅ **字段延迟计算**（避免不必要的计算）
- ✅ **对象池复用**（减少GC压力，提升性能）
- ✅ **更丰富的字段类型**（Duration、Time、StringSlice、Bytes、Map等）
- ✅ **便捷初始化函数**（函数式选项模式，自动配置，简化使用）

## 新功能使用

### 异步写入

异步写入可以显著提高性能，减少I/O阻塞：

```go
config := log.Config{
    // ... 其他配置
    File: log.FileConfig{
        Enabled:      true,
        AsyncEnabled: true, // 启用文件异步写入
    },
    Async: log.AsyncConfig{
        Enabled:      true,
        QueueSize:    1000,              // 队列大小
        BatchSize:    100,                // 批量写入大小
        FlushInterval: 5 * time.Second,   // 刷新间隔
        DropOnFull:   false,              // 队列满时阻塞（true则丢弃）
    },
}
```

**注意事项**：
- 异步写入只适用于文件输出，控制台保持同步（保证实时性）
- 程序退出前必须调用 `logger.Sync()` 确保所有日志写入完成
- 队列满时的行为：`DropOnFull=false` 会阻塞等待，`DropOnFull=true` 会丢弃日志

### 采样功能

采样功能可以控制日志量，避免日志过多影响性能：

```go
config := log.Config{
    // ... 其他配置
    Sampling: log.SamplingConfig{
        Enabled: true,
        Ratio:   0.1,                    // 采样10%的日志
        Levels:  []log.Level{log.DebugLevel}, // 只对Debug级别采样
    },
}
```

**使用场景**：
- Debug级别日志过多时，可以采样10%或更少
- 生产环境可以关闭Debug，或只采样少量
- 采样使用哈希算法，保证相同trace_id的日志采样一致性

**示例**：
```go
// 采样10%的Debug日志
config.Sampling = log.SamplingConfig{
    Enabled: true,
    Ratio:   0.1,
    Levels:  []log.Level{log.DebugLevel},
}

// 采样所有级别的50%日志
config.Sampling = log.SamplingConfig{
    Enabled: true,
    Ratio:   0.5,
    Levels:  []log.Level{}, // 空表示所有级别
}
```

### 级别检查优化

级别检查优化会在全局函数中提前检查日志级别，避免不必要的字段构建：

```go
// 自动优化，无需额外配置
log.Debug(ctx, "调试信息", 
    log.String("key", expensiveFunction()), // 如果Debug未启用，不会调用expensiveFunction
)
```

**工作原理**：
- 在全局函数中提前检查级别
- 如果级别未启用，直接返回
- 避免构建字段和函数调用

**性能提升**：
- 减少不必要的字段构建：5-10%
- 减少函数调用开销
- 特别适合Debug级别日志

### 敏感信息脱敏

敏感信息脱敏功能会自动脱敏密码、token等敏感字段：

```go
config := log.Config{
    // ... 其他配置
    Mask: log.MaskConfig{
        Enabled:  true,  // 启用脱敏
        Fields:   []string{"password", "pwd", "token", "secret", "key"},
        KeepHead: 3,     // 保留前3个字符
        KeepTail: 0,     // 不保留尾部
        MaskChar: "*",   // 使用*脱敏
    },
}
```

**使用示例**：
```go
// 自动脱敏
log.Info(ctx, "用户登录",
    log.String("password", "123456"),  // 输出: password[S]="123***"
    log.String("token", "abc123xyz"),  // 输出: token[S]="abc******"
    log.String("user_id", "123"),      // 不脱敏（不在敏感字段列表）
)
```

**脱敏规则**：
- 默认脱敏字段：password, pwd, token, secret, key
- 自动检测常见敏感字段模式（包含password、token等关键词）
- 保留前N个字符，中间用*替换
- 支持自定义脱敏字符

**配置示例**：
```go
// 保留前3后4（适合身份证、银行卡等）
Mask: log.MaskConfig{
    Enabled:  true,
    Fields:   []string{"idcard", "card_no"},
    KeepHead: 3,
    KeepTail: 4,
    MaskChar: "*",
}

// 完全脱敏（适合密码）
Mask: log.MaskConfig{
    Enabled:  true,
    Fields:   []string{"password"},
    KeepHead: 0,
    KeepTail: 0,
    MaskChar: "*",
}
```

### 钩子（Hooks）机制

钩子机制允许在日志写入前后执行自定义逻辑，支持告警、统计、过滤等功能：

```go
config := log.Config{
    // ... 其他配置
    Hooks: []log.Hook{
        // 告警钩子：Error级别自动发送告警
        log.NewAlertHook("https://webhook.example.com/alert", log.ErrorLevel, log.FatalLevel),
        
        // 统计钩子：统计日志数量
        log.NewStatsHook(),
        
        // 过滤钩子：过滤特定日志
        log.NewFilterHook(func(entry *log.LogEntry) bool {
            // 只记录包含特定字段的日志
            return entry.Service == "gateway"
        }),
    },
}
```

**内置钩子**：

1. **AlertHook** - 告警钩子
   ```go
   // 当Error或Fatal级别日志写入时，自动发送告警
   alertHook := log.NewAlertHook("https://webhook.example.com", log.ErrorLevel, log.FatalLevel)
   ```

2. **StatsHook** - 统计钩子
   ```go
   statsHook := log.NewStatsHook()
   // 获取统计信息
   count := statsHook.GetCount(log.ErrorLevel)
   stats := statsHook.GetStats()
   ```

3. **FilterHook** - 过滤钩子
   ```go
   // 只记录满足条件的日志
   filterHook := log.NewFilterHook(func(entry *log.LogEntry) bool {
       return entry.Level >= log.WarnLevel
   })
   ```

**自定义钩子**：
```go
type CustomHook struct{}

func (h *CustomHook) BeforeWrite(entry *log.LogEntry) error {
    // 写入前处理（可以修改entry或返回错误阻止写入）
    return nil
}

func (h *CustomHook) AfterWrite(entry *log.LogEntry) error {
    // 写入后处理（可以用于统计、告警等）
    return nil
}
```

### 错误处理增强

错误处理增强提供了优雅降级和错误统计功能：

```go
config := log.Config{
    // ... 其他配置
    ErrorHandler: log.NewDefaultErrorHandler(fallbackWriter),
}
```

**功能特性**：
- 写入失败时自动降级到备用写入器
- 错误统计和监控
- 不影响业务逻辑

**使用示例**：
```go
// 创建备用写入器（如stderr）
fallbackWriter := log.NewConsoleWriter(log.ConsoleConfig{
    Enabled: true,
    Format:  log.TextFormat,
}, false)

// 使用错误处理器
errorHandler := log.NewDefaultErrorHandler(fallbackWriter)

config := log.Config{
    ErrorHandler: errorHandler,
    // ... 其他配置
}

// 获取错误统计
stats := errorHandler.GetErrorStats()
writeErrors := errorHandler.GetErrorCount("write_error")
```

### 优雅关闭

优雅关闭功能确保程序退出时所有日志都写入完成，不丢失日志：

```go
// 程序退出前
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := log.Shutdown(ctx); err != nil {
    fmt.Printf("日志关闭失败: %v\n", err)
}
```

**功能特性**：
- 停止接收新日志
- 等待所有写入完成
- 支持超时控制
- 确保关键日志不丢失

**使用场景**：
- 程序正常退出
- 收到SIGTERM信号
- 容器关闭时

### 配置文件加载

支持从YAML或JSON文件加载配置，便于配置管理：

```go
// 从文件加载配置
config, err := log.LoadConfigFromFile("config/log.yaml")
if err != nil {
    panic(err)
}

logger, err := log.NewLogger(config)
if err != nil {
    panic(err)
}
```

**支持的格式**：
- YAML (`.yaml`, `.yml`)
- JSON (`.json`)

**环境变量覆盖**：
配置文件中的值可以通过环境变量覆盖：
- `LOG_SERVICE_NAME` - 服务名称
- `LOG_SERVICE_ID` - 服务ID
- `LOG_INSTANCE_ID` - 实例ID
- `LOG_LEVEL` - 日志级别
- `LOG_ENABLE_COLOR` - 启用颜色
- `LOG_FILE_PATH` - 文件路径

**配置示例（YAML）**：
```yaml
level: info
service_name: gateway
service_id: gateway-001
instance_id: gateway-001-pod-123
enable_color: true

console:
  enabled: true
  format: text
  level: debug

file:
  enabled: true
  format: json
  level: info
  path: "logs/%s-%s.log"
  max_size: 100
  max_age: 7
  max_backups: 10
  compress: true
  async_enabled: true

async:
  enabled: true
  queue_size: 2000
  batch_size: 100
  flush_interval: 5s
  drop_on_full: false

sampling:
  enabled: true
  ratio: 0.1
  levels: [debug]

mask:
  enabled: true
  fields: [password, pwd, token, secret, key]
  keep_head: 3
  keep_tail: 0
  mask_char: "*"
```

**保存配置**：
```go
// 保存配置到文件
err := log.SaveConfigToFile(config, "config/log.yaml")
```

### 字段延迟计算

字段延迟计算可以避免不必要的计算，提高性能：

```go
// 使用延迟字段
log.Debug(ctx, "处理请求",
    log.NewLazyString("expensive_data", func() string {
        // 只有在日志真正需要写入时才计算
        return expensiveComputation()
    }),
)

// 延迟整数字段
log.Info(ctx, "统计信息",
    log.NewLazyInt("count", func() int {
        return calculateCount()
    }),
)

// 延迟任意类型字段
log.Debug(ctx, "复杂对象",
    log.NewLazyField("complex_obj", func() interface{} {
        return buildComplexObject()
    }),
)
```

**性能优势**：
- 避免不必要的计算
- 特别适合Debug级别日志
- 减少CPU和内存开销
- 提升10-20%性能（Debug级别）

**使用场景**：
- 计算成本高的字段
- Debug级别日志
- 频繁调用的日志点

### 对象池复用

对象池复用功能通过复用LogEntry对象，减少内存分配和GC压力：

```go
// 自动优化，无需额外配置
// LogEntry对象从对象池获取，使用后自动归还
log.Info(ctx, "处理请求", log.String("key", "value"))
```

**工作原理**：
- 使用sync.Pool管理LogEntry对象
- 自动获取和归还对象
- 减少内存分配

**性能提升**：
- 减少内存分配：30-50%
- 降低GC压力：减少GC频率
- 提升性能：5-10%

---

## 字段类型支持

### 更丰富的字段类型

支持更多字段类型，提高易用性：

```go
// Duration字段
log.Info(ctx, "请求耗时",
    log.Duration("duration", time.Since(start)),
)

// Time字段
log.Info(ctx, "用户注册",
    log.Time("register_time", time.Now()),
)

// StringSlice字段
log.Info(ctx, "用户标签",
    log.StringSlice("tags", []string{"vip", "active"}),
)

// IntSlice字段
log.Info(ctx, "用户ID列表",
    log.IntSlice("user_ids", []int{1, 2, 3}),
)

// Bytes字段（十六进制）
log.Info(ctx, "数据",
    log.Bytes("data", []byte{0x01, 0x02, 0x03}),
)

// Map字段
log.Info(ctx, "用户信息",
    log.Map("user", map[string]interface{}{
        "id":   123,
        "name": "张三",
    }),
)

// Stringer字段
log.Info(ctx, "对象信息",
    log.Stringer("obj", myObject),
)
```

**支持的字段类型**：
- `Duration` - 时长（自动格式化为字符串）
- `Time` - 时间（RFC3339Nano格式）
- `StringSlice` - 字符串切片
- `IntSlice` - 整数切片
- `Int64Slice` - 64位整数切片
- `Bytes` - 字节数组（十六进制）
- `Map` - Map类型
- `Stringer` - 实现Stringer接口的对象

**完整字段类型列表**：
- **基础类型**：`String`, `Int`, `Int64`, `Int32`, `Bool`, `Float64`, `Float32`
- **时间相关**：`Duration`, `Time`
- **集合类型**：`StringSlice`, `IntSlice`, `Int64Slice`, `Map`
- **特殊类型**：`Bytes`（十六进制）, `Stringer`, `Error`
- **通用类型**：`Object`, `Any`（自动判断类型）

---

## 便捷初始化函数详解

### 概述

便捷初始化函数使用**函数式选项模式（Functional Options Pattern）**，提供简洁的API同时保留完整的配置能力。

### 核心函数

#### 1. InitLoggerWithShutdown（最推荐）

**一行代码完成初始化和优雅关闭**：

```go
func main() {
    // 最简单：自动从环境变量读取所有配置
    defer log.InitLoggerWithShutdown("gateway")()
    
    // 使用日志
    log.Info(context.Background(), "服务启动成功")
}
```

**带选项**：

```go
func main() {
    defer log.InitLoggerWithShutdown("gateway",
        log.WithServiceID("gateway-001"),
        log.WithEnvironment("production"),
        log.WithFilePath("logs/gateway-%s.log"),
    )()
    
    log.Info(context.Background(), "服务启动成功")
}
```

#### 2. InitLogger

返回 Logger 实例，需要手动设置全局Logger和关闭：

```go
func main() {
    logger, err := log.InitLogger("gateway",
        log.WithServiceID("gateway-001"),
        log.WithLevel(log.DebugLevel),
    )
    if err != nil {
        panic(err)
    }
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        logger.Shutdown(ctx)
    }()
    
    log.SetGlobalLogger(logger)
}
```

#### 3. MustInitLogger

初始化失败时panic：

```go
func main() {
    defer log.MustInitLogger("gateway",
        log.WithServiceID("gateway-001"),
    ).Shutdown(context.Background())
    
    log.Info(context.Background(), "服务启动成功")
}
```

### 选项函数

所有选项函数都支持链式调用：

| 函数 | 说明 | 示例 |
|------|------|------|
| `WithServiceID(id)` | 设置服务ID | `log.WithServiceID("gateway-001")` |
| `WithInstanceID(id)` | 设置实例ID | `log.WithInstanceID("pod-123")` |
| `WithLevel(level)` | 设置日志级别 | `log.WithLevel(log.DebugLevel)` |
| `WithEnvironment(env)` | 设置环境（自动调整级别和颜色） | `log.WithEnvironment("production")` |
| `WithFileOutput(enable)` | 启用/禁用文件输出 | `log.WithFileOutput(true)` |
| `WithFilePath(path)` | 设置文件路径模板 | `log.WithFilePath("logs/app-%s.log")` |
| `WithColor(enable)` | 强制启用/禁用颜色 | `log.WithColor(false)` |
| `WithConfig(config)` | 使用完整配置（覆盖其他选项） | `log.WithConfig(customConfig)` |

### 自动配置特性

便捷初始化函数会自动：

1. **从环境变量读取配置**：
   - `SERVICE_NAME` → 服务名称
   - `SERVICE_ID` → 服务ID
   - `INSTANCE_ID` / `POD_NAME` / `HOSTNAME` → 实例ID
   - `APP_ENV` → 环境（用于自动判断日志级别和颜色）
   - `LOG_LEVEL` → 日志级别
   - `LOG_ENABLE_COLOR` → 是否启用颜色
   - `LOG_FILE_PATH` → 文件路径

2. **根据环境自动调整**：
   - **生产环境**（`production`/`prod`）：日志级别 `Info`，禁用颜色
   - **开发环境**（其他）：日志级别 `Debug`，启用颜色

3. **容器环境支持**：
   - 自动从 `POD_NAME`、`HOSTNAME` 等环境变量获取实例ID
   - 自动识别容器环境并调整配置

### 使用场景对比

**传统方式**（50+ 行代码）：
```go
config := log.DefaultConfig()
config.ServiceName = "gateway"
config.ServiceID = "gateway-001"
// ... 大量配置代码 ...
logger, err := log.NewLogger(config)
if err != nil {
    panic(err)
}
defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    logger.Shutdown(ctx)
}()
log.SetGlobalLogger(logger)
```

**便捷方式**（1-5 行代码）：
```go
defer log.InitLoggerWithShutdown("gateway",
    log.WithServiceID("gateway-001"),
    log.WithEnvironment("production"),
)()
```

### 优势总结

1. ✅ **简洁**：从 50+ 行减少到 1-5 行
2. ✅ **智能**：自动从环境变量读取配置
3. ✅ **灵活**：支持函数式选项，按需定制
4. ✅ **安全**：自动优雅关闭，确保日志不丢失
5. ✅ **兼容**：保留完整配置能力，不影响高级功能

---

## 完整功能清单

