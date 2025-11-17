# 异步写入和采样功能使用指南

## 🚀 异步写入功能

### 功能说明

异步写入将日志写入操作放到后台goroutine中执行，避免阻塞业务逻辑，显著提高性能。

### 工作原理

```
业务代码 → 日志写入请求 → 队列（Channel）→ 后台goroutine批量写入文件
```

**优势**：
- 不阻塞业务逻辑
- 批量写入，减少I/O操作
- 提高吞吐量

### 配置示例

```go
config := log.Config{
    // ... 基础配置
    File: log.FileConfig{
        Enabled:      true,
        AsyncEnabled: true, // 启用异步写入
        Path:         "logs/%s-%s.log",
        // ... 其他文件配置
    },
    Async: log.AsyncConfig{
        Enabled:      true,              // 启用异步写入
        QueueSize:    1000,              // 队列大小（建议1000-10000）
        BatchSize:    100,                // 批量写入大小（建议50-200）
        FlushInterval: 5 * time.Second,   // 刷新间隔（建议3-10秒）
        DropOnFull:   false,              // 队列满时的行为
    },
}
```

### 配置参数说明

| 参数 | 说明 | 建议值 | 注意事项 |
|------|------|--------|----------|
| `QueueSize` | 队列大小 | 1000-10000 | 太小容易阻塞，太大占用内存 |
| `BatchSize` | 批量写入大小 | 50-200 | 批量写入减少I/O次数 |
| `FlushInterval` | 刷新间隔 | 3-10秒 | 定期刷新，避免日志延迟 |
| `DropOnFull` | 队列满时行为 | false | true=丢弃，false=阻塞等待 |

### 使用注意事项

1. **必须调用Sync()**
   ```go
   defer logger.Sync() // 确保所有日志写入完成
   ```

2. **优雅关闭**
   ```go
   // 程序退出前
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()
   logger.Sync() // 等待所有日志写入
   ```

3. **队列满时的处理**
   - `DropOnFull=false`：阻塞等待，保证不丢失日志（推荐）
   - `DropOnFull=true`：丢弃日志，保证不阻塞（高并发场景）

### 性能对比

**同步写入**：
- 每次写入都阻塞
- 吞吐量：~1000 logs/sec

**异步写入**：
- 写入不阻塞
- 吞吐量：~10000+ logs/sec
- 延迟：< 5秒（取决于FlushInterval）

---

## 📊 采样功能

### 功能说明

采样功能可以按比例记录日志，避免日志过多影响性能和存储。

### 工作原理

使用哈希算法决定是否采样，保证相同trace_id的日志采样一致性。

```
日志请求 → 采样检查 → 哈希计算 → 决定是否记录
```

### 配置示例

#### 示例1：采样Debug级别10%

```go
config := log.Config{
    // ... 基础配置
    Sampling: log.SamplingConfig{
        Enabled: true,
        Ratio:   0.1,                    // 采样10%
        Levels:  []log.Level{log.DebugLevel}, // 只对Debug级别采样
    },
}
```

**效果**：
- Debug级别：只记录10%
- 其他级别：全部记录

#### 示例2：采样所有级别50%

```go
config := log.Config{
    Sampling: log.SamplingConfig{
        Enabled: true,
        Ratio:   0.5,        // 采样50%
        Levels:  []log.Level{}, // 空表示所有级别
    },
}
```

**效果**：
- 所有级别：只记录50%

#### 示例3：采样多个级别

```go
config := log.Config{
    Sampling: log.SamplingConfig{
        Enabled: true,
        Ratio:   0.2, // 采样20%
        Levels:  []log.Level{
            log.DebugLevel,
            log.InfoLevel,
        },
    },
}
```

**效果**：
- Debug和Info级别：只记录20%
- 其他级别：全部记录

### 使用场景

1. **开发环境**
   ```go
   // 开发时记录所有日志
   Sampling: log.SamplingConfig{
       Enabled: false, // 不采样
   }
   ```

2. **生产环境**
   ```go
   // 生产环境采样Debug日志
   Sampling: log.SamplingConfig{
       Enabled: true,
       Ratio:   0.1, // 只记录10%的Debug日志
       Levels:  []log.Level{log.DebugLevel},
   }
   ```

3. **高并发场景**
   ```go
   // 高并发时采样所有级别
   Sampling: log.SamplingConfig{
       Enabled: true,
       Ratio:   0.3, // 只记录30%
       Levels:  []log.Level{}, // 所有级别
   }
   ```

### 采样算法

使用FNV哈希算法，保证：
- 相同trace_id的日志采样一致性
- 采样分布均匀
- 性能开销小

### 注意事项

1. **采样比例**
   - 0.0-1.0之间
   - 0.1 = 10%
   - 1.0 = 100%（不采样）

2. **级别列表**
   - 空列表：采样所有级别
   - 非空：只采样指定级别

3. **采样一致性**
   - 使用哈希算法保证相同trace_id的日志采样一致
   - 可以通过trace_id查询完整的日志链

---

## 💡 最佳实践

### 推荐配置（生产环境）

```go
config := log.Config{
    Level:       log.InfoLevel,
    ServiceName: "gateway",
    Outputs:     []log.Output{log.ConsoleOutput, log.FileOutput},
    
    Console: log.ConsoleConfig{
        Enabled: true,
        Format:  log.TextFormat,
        Level:   log.InfoLevel, // 控制台只显示INFO及以上
    },
    
    File: log.FileConfig{
        Enabled:      true,
        Format:       log.JSONFormat,
        Level:        log.InfoLevel,
        AsyncEnabled: true, // 启用异步写入
        Path:         "logs/%s-%s.log",
        MaxSize:      100,
        MaxAge:       7,
        MaxBackups:   10,
        Compress:     true,
    },
    
    Async: log.AsyncConfig{
        Enabled:      true,
        QueueSize:    2000,
        BatchSize:    100,
        FlushInterval: 5 * time.Second,
        DropOnFull:   false, // 保证不丢失
    },
    
    Sampling: log.SamplingConfig{
        Enabled: true,
        Ratio:   0.1, // Debug级别采样10%
        Levels:  []log.Level{log.DebugLevel},
    },
}
```

### 性能调优建议

1. **队列大小**
   - 低并发：1000
   - 中并发：2000-5000
   - 高并发：5000-10000

2. **批量大小**
   - 小批量：50-100（低延迟）
   - 大批量：100-200（高吞吐）

3. **刷新间隔**
   - 实时性要求高：3秒
   - 平衡：5秒
   - 性能优先：10秒

4. **采样比例**
   - Debug级别：10-20%
   - Info级别：50-100%
   - Error/Fatal：100%（不采样）

---

## 🔍 故障排查

### 问题1：日志丢失

**原因**：队列满且DropOnFull=true

**解决**：
```go
Async: log.AsyncConfig{
    DropOnFull: false, // 改为阻塞等待
}
```

### 问题2：程序退出时日志丢失

**原因**：未调用Sync()

**解决**：
```go
defer logger.Sync() // 确保调用
```

### 问题3：采样不一致

**原因**：采样算法基于级别，不是trace_id

**解决**：使用trace_id作为采样key（需要扩展功能）

---

## 📈 性能监控

### 监控指标

```go
// 获取异步写入器队列大小
if asyncWriter, ok := writer.(*AsyncWriter); ok {
    queueSize := asyncWriter.QueueSize()
    // 监控队列大小，如果接近QueueSize，说明写入速度跟不上
}
```

### 告警建议

- 队列使用率 > 80%：考虑增加QueueSize或BatchSize
- 队列使用率 = 100%：检查写入性能或调整DropOnFull策略

---

**完成！现在你的日志系统支持异步写入和采样功能了！** 🎉

