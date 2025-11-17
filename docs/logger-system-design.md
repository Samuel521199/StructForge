# 通用日志系统设计文档

## 📋 需求分析

### 1. 控制台打印需求
- **需求**：能够在控制台打印，区分 info、debug、warn、error、fatal 等级别的日志
- **分析**：
  - 需要支持5个日志级别：Debug、Info、Warn、Error、Fatal
  - 控制台输出应该清晰易读
  - 需要支持结构化输出和普通文本输出两种模式
  - 开发环境建议使用文本模式，生产环境建议使用结构化模式

### 2. 颜色区分需求
- **需求**：控制台打印时不同的日志级别可以用不同的颜色来区分
- **颜色方案建议**：
  - **Debug**: 灰色/浅蓝色（开发调试信息，不重要）
  - **Info**: 绿色/白色（正常信息，默认色）
  - **Warn**: 黄色/橙色（警告信息，需要注意）
  - **Error**: 红色（错误信息，需要立即关注）
  - **Fatal**: 红色+粗体/背景色（致命错误，程序即将退出）
- **技术实现**：
  - 使用 ANSI 转义码实现颜色
  - 检测终端是否支持颜色（避免在非终端环境出错）
  - 支持通过配置关闭颜色（适合日志重定向到文件）

### 3. 容器追踪需求
- **需求**：如果使用日志的进程（或微服务）在容器中可以被追踪或者查看
- **分析**：
  - 需要在日志中包含容器/服务标识信息
  - 支持从环境变量获取容器信息（POD_NAME, HOSTNAME, NAMESPACE等）
  - 支持服务名称、服务ID、实例ID等标识
  - 支持追踪ID（TraceID）和SpanID（用于分布式追踪）
  - 日志格式应该便于日志收集系统（如ELK、Loki）解析

### 4. 文件输出需求
- **需求**：日志系统可以打印成为文件，每天一个文件，每条记录都有简短但准确的服务器名和准确日期和时间
- **文件命名规则**：
  - 格式：`{服务名}-{日期}.log` 或 `{服务名}-{日期}-{级别}.log`
  - 示例：`gateway-2025-01-15.log` 或 `gateway-2025-01-15-error.log`
  - 支持按级别分离文件（可选）
- **时间格式**：
  - 建议使用 ISO8601 格式：`2025-01-15T14:30:45.123Z`
  - 或带时区：`2025-01-15T14:30:45.123+08:00`
  - 时间精度到毫秒级别
- **服务器名**：
  - 每条日志包含服务名称
  - 格式：`[服务名]` 或作为字段 `service: "gateway"`
- **文件管理**：
  - 自动按天轮转
  - 支持文件大小限制
  - 支持保留天数配置
  - 支持自动压缩旧文件

### 5. 合并功能需求
- **需求**：日志可以通过添加符号和类型代号来传达参数显示
- **分析**：
  - 需要支持结构化日志（JSON格式）
  - 需要支持键值对形式的参数
  - 需要支持类型标识（string, int, bool, error等）
  - 需要支持嵌套结构
  - 需要支持数组/切片
- **参数格式建议**：
  - 使用符号标识：`key=value` 或 `key: value`
  - 使用类型代号：`[S]`字符串、`[I]`整数、`[B]`布尔、`[E]`错误、`[O]`对象
  - 示例：`user_id[I]=123 name[S]="张三" error[E]=err`

---

## 🏗️ 系统架构设计

### 1. 核心组件

```
┌─────────────────────────────────────────┐
│          Logger Interface               │
│  (统一的日志接口，供业务代码调用)         │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│       Logger Implementation             │
│  (日志实现层，处理日志级别、格式化等)     │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         Output Manager                  │
│  (输出管理器，管理多个输出目标)          │
└─────────────────────────────────────────┘
        ↓                    ↓
┌──────────────┐    ┌──────────────┐
│ Console Core │    │  File Core   │
│ (控制台输出) │    │ (文件输出)   │
└──────────────┘    └──────────────┘
```

### 2. 模块划分

#### 2.1 配置模块（Config）
- **LoggerConfig**: 日志系统全局配置
  - 日志级别（Level）
  - 输出目标（Outputs: console, file）
  - 格式（Format: text, json）
  - 颜色支持（EnableColor）
  - 服务信息（ServiceName, ServiceID, InstanceID）

#### 2.2 控制台输出模块（Console）
- **ConsoleWriter**: 控制台写入器
  - 支持颜色输出
  - 支持文本格式和JSON格式
  - 自动检测终端能力
  - 支持禁用颜色（用于重定向）

#### 2.3 文件输出模块（File）
- **FileWriter**: 文件写入器
  - 按天自动轮转
  - 支持文件大小限制
  - 支持保留策略
  - 支持压缩
  - 线程安全的文件写入

#### 2.4 格式化模块（Formatter）
- **TextFormatter**: 文本格式化器
  - 可读性强的文本格式
  - 支持颜色
  - 支持自定义格式模板

- **JSONFormatter**: JSON格式化器
  - 标准JSON格式
  - 便于日志收集系统解析
  - 支持嵌套结构

#### 2.5 字段管理模块（Fields）
- **FieldBuilder**: 字段构建器
  - 类型安全的字段添加
  - 支持多种数据类型
  - 支持嵌套字段
  - 自动序列化复杂类型

---

## 📐 接口设计

### 1. 核心接口

```go
// Logger 日志接口
type Logger interface {
    Debug(ctx context.Context, msg string, fields ...Field)
    Info(ctx context.Context, msg string, fields ...Field)
    Warn(ctx context.Context, msg string, fields ...Field)
    Error(ctx context.Context, msg string, fields ...Field)
    Fatal(ctx context.Context, msg string, fields ...Field)
    
    // 带级别的通用方法
    Log(ctx context.Context, level Level, msg string, fields ...Field)
    
    // 同步刷新
    Sync() error
    
    // 创建子Logger（带默认字段）
    With(fields ...Field) Logger
}
```

### 2. 配置接口

```go
// Config 日志配置
type Config struct {
    // 基础配置
    Level       Level    // 日志级别
    ServiceName string   // 服务名称
    ServiceID   string   // 服务ID
    InstanceID  string   // 实例ID
    
    // 输出配置
    Outputs     []Output // 输出目标列表
    EnableColor bool     // 是否启用颜色
    
    // 控制台配置
    Console ConsoleConfig
    
    // 文件配置
    File   FileConfig
}

// ConsoleConfig 控制台配置
type ConsoleConfig struct {
    Enabled bool   // 是否启用
    Format  Format // 输出格式（text/json）
    Level   Level  // 控制台日志级别（可独立设置）
}

// FileConfig 文件配置
type FileConfig struct {
    Enabled     bool   // 是否启用
    Format      Format // 输出格式（通常为json）
    Level       Level  // 文件日志级别
    Path        string // 文件路径模板
    MaxSize     int    // 单个文件最大大小（MB）
    MaxAge      int    // 保留天数
    MaxBackups  int    // 最大备份文件数
    Compress    bool   // 是否压缩旧文件
    SeparateByLevel bool // 是否按级别分离文件
}
```

### 3. 字段接口

```go
// Field 日志字段
type Field interface {
    Key() string
    Value() interface{}
    Type() FieldType
}

// FieldType 字段类型
type FieldType string

const (
    FieldTypeString FieldType = "S"  // 字符串
    FieldTypeInt    FieldType = "I"  // 整数
    FieldTypeBool   FieldType = "B"  // 布尔
    FieldTypeFloat  FieldType = "F"  // 浮点数
    FieldTypeError  FieldType = "E"  // 错误
    FieldTypeObject FieldType = "O"  // 对象
    FieldTypeArray  FieldType = "A"  // 数组
)

// 便捷函数
func String(key, value string) Field
func Int(key string, value int) Field
func Int64(key string, value int64) Field
func Bool(key string, value bool) Field
func Float64(key string, value float64) Field
func Error(err error) Field
func Object(key string, value interface{}) Field
func Any(key string, value interface{}) Field
```

---

## 🎨 输出格式设计

### 1. 控制台文本格式（开发环境）

```
[2025-01-15 14:30:45.123] [INFO] [gateway] 用户登录成功
  user_id[I]=123
  username[S]="zhangsan"
  ip[S]="192.168.1.1"
  duration[F]=0.123
```

**带颜色版本**：
- 时间：灰色
- 级别：对应颜色（INFO绿色、ERROR红色等）
- 服务名：蓝色
- 消息：白色
- 字段：浅灰色

### 2. 控制台JSON格式（可选）

```json
{
  "timestamp": "2025-01-15T14:30:45.123Z",
  "level": "INFO",
  "service": "gateway",
  "service_id": "gateway-001",
  "instance_id": "gateway-001-pod-123",
  "message": "用户登录成功",
  "fields": {
    "user_id": 123,
    "username": "zhangsan",
    "ip": "192.168.1.1",
    "duration": 0.123
  }
}
```

### 3. 文件JSON格式（生产环境）

```json
{
  "timestamp": "2025-01-15T14:30:45.123+08:00",
  "level": "INFO",
  "service": "gateway",
  "service_id": "gateway-001",
  "instance_id": "gateway-001-pod-123",
  "pod_name": "gateway-pod-123",
  "hostname": "gateway-host-1",
  "namespace": "default",
  "trace_id": "abc123",
  "span_id": "def456",
  "message": "用户登录成功",
  "fields": {
    "user_id": 123,
    "username": "zhangsan",
    "ip": "192.168.1.1",
    "duration": 0.123
  },
  "caller": "handler/user.go:45"
}
```

---

## 🔧 功能特性

### 1. 多输出目标
- 支持同时输出到控制台和文件
- 每个输出目标可以独立配置级别和格式
- 例如：控制台显示所有级别（开发用），文件只记录INFO及以上（生产用）

### 2. 上下文支持
- 从Context中提取追踪信息（TraceID、SpanID）
- 从Context中提取用户信息（UserID）
- 从Context中提取请求信息（RequestID）

### 3. 自动字段注入
- 自动注入服务信息（服务名、ID、实例ID）
- 自动注入容器信息（POD_NAME、HOSTNAME等）
- 自动注入调用者信息（文件名、行号、函数名）

### 4. 性能优化
- 异步写入文件（可选）
- 缓冲写入
- 级别检查提前返回（避免不必要的格式化）
- 字段延迟计算（只在需要时序列化）

### 5. 线程安全
- 所有操作线程安全
- 支持并发写入
- 文件写入使用锁保护

---

## 📦 使用方式设计

### 1. 初始化

```go
// 方式1：使用配置结构
config := log.Config{
    Level:       log.InfoLevel,
    ServiceName: "gateway",
    ServiceID:   "gateway-001",
    Outputs:     []log.Output{log.ConsoleOutput, log.FileOutput},
    Console: log.ConsoleConfig{
        Enabled: true,
        Format:  log.TextFormat,
        Level:   log.DebugLevel,
    },
    File: log.FileConfig{
        Enabled:    true,
        Format:     log.JSONFormat,
        Level:      log.InfoLevel,
        Path:       "logs/gateway-%s.log",
        MaxSize:    100, // MB
        MaxAge:     7,   // days
        MaxBackups: 10,
        Compress:   true,
    },
}

logger, err := log.NewLogger(config)
if err != nil {
    panic(err)
}
defer logger.Sync()

// 设置为全局Logger
log.SetGlobalLogger(logger)
```

### 2. 基本使用

```go
// 简单日志
log.Info(ctx, "服务启动成功")

// 带字段的日志
log.Info(ctx, "用户登录成功",
    log.String("user_id", "123"),
    log.Int("age", 25),
    log.Bool("is_vip", true),
    log.Error(err),
)

// 不同级别
log.Debug(ctx, "调试信息", log.String("key", "value"))
log.Warn(ctx, "警告信息", log.String("key", "value"))
log.Error(ctx, "错误信息", log.Error(err))
log.Fatal(ctx, "致命错误", log.Error(err)) // 会调用os.Exit(1)
```

### 3. 创建子Logger

```go
// 创建带默认字段的Logger
userLogger := logger.With(
    log.String("module", "user"),
    log.String("version", "v1.0"),
)

// 使用子Logger，所有日志都会自动包含默认字段
userLogger.Info(ctx, "用户操作", log.String("action", "login"))
```

---

## 🎯 扩展性设计

### 1. 自定义格式化器
- 实现 `Formatter` 接口
- 可以自定义输出格式
- 支持插件化

### 2. 自定义输出目标
- 实现 `Writer` 接口
- 可以输出到远程日志系统（如ELK、Loki）
- 可以输出到消息队列
- 可以输出到数据库

### 3. 钩子（Hooks）
- 支持日志钩子，在日志写入前后执行自定义逻辑
- 例如：发送告警、统计日志数量等

### 4. 采样（Sampling）
- 支持日志采样，避免日志过多
- 例如：Debug级别日志只记录10%

---

## 🔍 容器环境支持

### 1. 环境变量自动检测
- `POD_NAME`: Pod名称
- `HOSTNAME`: 主机名
- `POD_NAMESPACE`: 命名空间
- `POD_IP`: Pod IP
- `NODE_NAME`: 节点名称

### 2. 服务发现集成
- 自动从服务注册中心获取服务信息
- 支持Nacos、Consul等服务发现

### 3. 追踪集成
- 自动从Context提取TraceID和SpanID
- 支持OpenTelemetry、Jaeger等追踪系统

---

## 📊 性能考虑

### 1. 零分配设计
- 字段对象复用
- 字符串拼接优化
- 避免不必要的内存分配

### 2. 异步写入
- 文件写入使用异步队列
- 控制台写入同步（保证实时性）

### 3. 级别检查
- 在格式化前检查级别
- 避免不必要的字符串操作

### 4. 缓冲管理
- 合理的缓冲区大小
- 定期刷新机制

---

## 🛡️ 错误处理

### 1. 初始化错误
- 配置验证
- 文件权限检查
- 目录创建失败处理

### 2. 运行时错误
- 文件写入失败降级（只输出到控制台）
- 格式化错误容错
- 不因日志错误影响业务

### 3. 资源清理
- 优雅关闭
- 确保缓冲区刷新
- 文件句柄正确关闭

---

## 📝 总结

### 设计原则
1. **简单易用**：API设计简洁，学习成本低
2. **高性能**：零分配、异步写入、级别检查
3. **可扩展**：支持自定义格式化器、输出目标、钩子
4. **生产就绪**：支持容器环境、追踪、文件管理
5. **灵活配置**：支持多输出、多格式、多级别

### 核心特性
- ✅ 5个日志级别（Debug、Info、Warn、Error、Fatal）
- ✅ 控制台颜色输出
- ✅ 文件按天轮转
- ✅ 结构化日志（JSON）
- ✅ 容器环境支持
- ✅ 追踪集成
- ✅ 类型安全的字段API
- ✅ 线程安全
- ✅ 高性能

---

**下一步**：根据此设计文档实现具体的代码。

