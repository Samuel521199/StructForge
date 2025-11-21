# Node Service 架构设计 - 快速参考

## 架构方案：单一 Node Service + 插件化

### 核心设计
- **统一服务**：所有节点通过一个 Node Service 管理
- **插件化**：支持自定义节点扩展
- **水平扩展**：可以部署多个实例

## 目录结构

```
backend/apps/node/
├── internal/
│   ├── biz/
│   │   ├── node.go          # 节点业务逻辑
│   │   ├── registry.go      # 节点注册表
│   │   ├── executor.go      # 执行器接口
│   │   └── factory.go        # 执行器工厂
│   ├── service/
│   │   └── node.go           # gRPC 服务实现
│   └── nodes/                # 节点执行器实现
│       ├── trigger/          # 触发节点
│       ├── ai/               # AI 节点（调用 AI Service）
│       ├── data/             # 数据处理节点
│       ├── integration/      # 集成节点
│       ├── control/          # 控制节点
│       └── tool/             # 工具节点（调用 Tool Service）
```

## 节点分类

### 触发节点（4种）
- `webhook` - Webhook 触发
- `timer` - 定时触发
- `manual` - 手动触发
- `event` - 事件触发

### AI 节点（4种）
- `text_generation` - 文本生成（调用 AI Service）
- `image_generation` - 图像生成（调用 AI Service）
- `chat` - 对话（调用 AI Service）
- `embedding` - 向量化（调用 AI Service）

### 数据处理节点（5种）
- `transform` - 数据转换
- `filter` - 数据过滤
- `aggregate` - 数据聚合
- `map` - 数据映射
- `sort` - 数据排序

### 集成节点（5种）
- `http` - HTTP 请求
- `database` - 数据库操作
- `message_queue` - 消息队列
- `email` - 邮件发送
- `webhook` - Webhook 发送

### 控制节点（5种）
- `condition` - 条件判断
- `loop` - 循环
- `parallel` - 并行执行
- `switch` - 分支选择
- `delay` - 延迟

### 工具节点（3种）
- `script` - 脚本执行（调用 Tool Service）
- `file_operation` - 文件操作
- `system_command` - 系统命令

**总计：26 种节点类型**

## 核心接口

```go
type NodeExecutor interface {
    Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error)
    Validate(ctx context.Context, config map[string]interface{}) error
    GetMetadata() *NodeMetadata
}
```

## 执行流程

```
Workflow Service → Node Service → Node Executor → 执行逻辑
                                      ↓
                            AI/Tool Service (如需要)
```

## 详细文档

- **第一部分**：架构概述、目录结构、触发节点和 AI 节点
- **第二部分**：集成节点、控制节点、工具节点、核心接口设计
- **第三部分**：插件系统、外部服务交互、API 设计、错误处理
- **第四部分**：实现示例、初始化流程、配置模式、测试策略

