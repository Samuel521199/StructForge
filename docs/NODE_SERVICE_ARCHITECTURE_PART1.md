# Node Service 架构设计文档 - 第一部分

## 一、架构方案概述

### 1.1 方案选择：单一 Node Service + 插件化架构

**核心设计理念：**
- 统一的服务入口，简化调用和运维
- 插件化的节点执行器，支持灵活扩展
- 清晰的职责划分，便于维护和开发

### 1.2 架构优势

**运维优势：**
- ✅ 服务数量少，只需维护一个 Node Service
- ✅ 部署简单，只需部署一个服务实例
- ✅ 监控集中，所有节点执行统一监控
- ✅ 日志统一，便于问题排查

**开发优势：**
- ✅ 统一接口，Workflow Service 只需调用一个服务
- ✅ 代码复用，节点间可以共享通用功能
- ✅ 技术栈统一，全部使用 Go 语言
- ✅ 开发效率高，新节点开发快速

**扩展优势：**
- ✅ 水平扩展，可以部署多个实例
- ✅ 插件化扩展，支持自定义节点
- ✅ 按需扩展，高负载节点可以单独优化
- ✅ 未来可拆分，如果某个类别节点变得复杂，可以拆分

### 1.3 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    Workflow Service                         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  执行引擎                                              │   │
│  │  - 解析工作流定义（DAG）                               │   │
│  │  - 按拓扑顺序执行节点                                   │   │
│  │  - 管理执行状态                                         │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │ gRPC 调用
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Node Service (统一服务)                    │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Node Registry (节点注册表)                            │   │
│  │  - 节点类型注册                                        │   │
│  │  - 节点元数据管理                                      │   │
│  │  - 节点执行器映射                                      │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Executor Factory (执行器工厂)                        │   │
│  │  - 根据节点类型创建执行器                              │   │
│  │  - 执行器生命周期管理                                  │   │
│  │  - 执行器池管理                                        │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Executor Implementations (执行器实现)                │   │
│  │  ├── Trigger Executors (触发节点)                     │   │
│  │  ├── AI Executors (AI 节点，调用 AI Service)          │   │
│  │  ├── Data Executors (数据处理节点)                     │   │
│  │  ├── Integration Executors (集成节点)                  │   │
│  │  ├── Control Executors (控制节点)                      │   │
│  │  └── Tool Executors (工具节点，调用 Tool Service)      │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Plugin System (插件系统)                             │   │
│  │  - 自定义节点插件加载                                  │   │
│  │  - 插件隔离执行                                        │   │
│  │  - 插件生命周期管理                                    │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
        │                          │                          │
        ▼                          ▼                          ▼
┌──────────────┐        ┌──────────────┐        ┌──────────────┐
│ AI Service   │        │ Tool Service │        │ 外部服务      │
│ AI 模型管理   │        │ 工具集成管理  │        │ (HTTP/DB等)  │
└──────────────┘        └──────────────┘        └──────────────┘
```

---

## 二、目录结构设计

### 2.1 完整目录结构

```
backend/apps/node/
├── cmd/
│   └── node/
│       ├── main.go              # 服务启动入口
│       ├── wire.go              # 依赖注入配置
│       └── wire_gen.go          # Wire 生成的代码
├── internal/
│   ├── biz/                     # 业务逻辑层
│   │   ├── node.go              # 节点核心业务逻辑
│   │   ├── registry.go          # 节点注册表
│   │   ├── executor.go           # 执行器接口定义
│   │   ├── factory.go           # 执行器工厂
│   │   └── plugin.go            # 插件管理
│   ├── data/                    # 数据访问层
│   │   ├── node.go              # 节点数据访问
│   │   └── provider.go          # 数据提供者
│   ├── service/                  # 服务层（gRPC）
│   │   └── node.go              # gRPC 服务实现
│   ├── server/                   # 服务器配置
│   │   ├── grpc.go              # gRPC 服务器
│   │   └── http.go              # HTTP 服务器（可选）
│   ├── conf/                    # 配置
│   │   └── config.yaml          # 服务配置
│   └── nodes/                   # 节点执行器实现
│       ├── trigger/             # 触发节点
│       │   ├── webhook.go       # Webhook 触发
│       │   ├── timer.go         # 定时触发
│       │   ├── manual.go        # 手动触发
│       │   └── event.go         # 事件触发
│       ├── ai/                  # AI 节点
│       │   ├── text_generation.go    # 文本生成
│       │   ├── image_generation.go   # 图像生成
│       │   ├── chat.go               # 对话
│       │   └── embedding.go         # 向量化
│       ├── data/                # 数据处理节点
│       │   ├── transform.go     # 数据转换
│       │   ├── filter.go        # 数据过滤
│       │   ├── aggregate.go     # 数据聚合
│       │   ├── map.go           # 数据映射
│       │   └── sort.go          # 数据排序
│       ├── integration/         # 集成节点
│       │   ├── http.go          # HTTP 请求
│       │   ├── database.go      # 数据库操作
│       │   ├── message_queue.go # 消息队列
│       │   ├── email.go         # 邮件发送
│       │   └── webhook.go       # Webhook 发送
│       ├── control/              # 控制节点
│       │   ├── condition.go     # 条件判断
│       │   ├── loop.go          # 循环
│       │   ├── parallel.go      # 并行执行
│       │   ├── switch.go        # 分支选择
│       │   └── delay.go         # 延迟
│       └── tool/                # 工具节点
│           ├── script.go        # 脚本执行
│           ├── file_operation.go # 文件操作
│           └── system_command.go # 系统命令
└── api/
    └── node/
        └── v1/
            └── node.proto        # Protobuf API 定义
```

---

## 三、节点分类与实现

### 3.1 触发节点（Trigger Nodes）

**节点类型：**

#### 3.1.1 Webhook 节点
- **功能**：接收 HTTP 请求，触发工作流执行
- **参数**：
  - `path`：Webhook 路径
  - `method`：HTTP 方法（GET、POST 等）
  - `auth`：认证方式（无、API Key、Bearer Token）
  - `headers`：请求头验证规则
- **实现位置**：`internal/nodes/trigger/webhook.go`
- **特殊处理**：需要与 Gateway 或 Workflow Service 协调，注册 Webhook 路由

#### 3.1.2 定时器节点
- **功能**：按 Cron 表达式定时触发
- **参数**：
  - `cron`：Cron 表达式
  - `timezone`：时区
  - `enabled`：是否启用
- **实现位置**：`internal/nodes/trigger/timer.go`
- **特殊处理**：实际调度由 Scheduler Service 负责，此节点主要用于工作流内部定时

#### 3.1.3 手动触发节点
- **功能**：等待用户手动触发
- **参数**：
  - `trigger_key`：触发键（用于 API 调用）
  - `timeout`：超时时间
- **实现位置**：`internal/nodes/trigger/manual.go`
- **特殊处理**：需要与 Workflow Service 协调，管理等待状态

#### 3.1.4 事件触发节点
- **功能**：监听系统事件，触发工作流
- **参数**：
  - `event_type`：事件类型
  - `event_filter`：事件过滤条件
- **实现位置**：`internal/nodes/trigger/event.go`
- **特殊处理**：需要事件总线支持

---

### 3.2 AI 节点（AI Nodes）

**节点类型：**

#### 3.2.1 文本生成节点
- **功能**：调用 AI 模型生成文本
- **参数**：
  - `model`：模型名称（通过 AI Service 配置）
  - `prompt`：提示词
  - `max_tokens`：最大 token 数
  - `temperature`：温度参数
- **实现位置**：`internal/nodes/ai/text_generation.go`
- **调用方式**：通过 gRPC 调用 AI Service
- **返回数据**：生成的文本

#### 3.2.2 图像生成节点
- **功能**：调用 AI 模型生成图像
- **参数**：
  - `model`：模型名称
  - `prompt`：提示词
  - `size`：图像尺寸
  - `quality`：图像质量
- **实现位置**：`internal/nodes/ai/image_generation.go`
- **调用方式**：通过 gRPC 调用 AI Service
- **返回数据**：图像 URL 或 Base64

#### 3.2.3 对话节点
- **功能**：多轮对话
- **参数**：
  - `model`：模型名称
  - `messages`：对话历史
  - `system_prompt`：系统提示
- **实现位置**：`internal/nodes/ai/chat.go`
- **调用方式**：通过 gRPC 调用 AI Service
- **返回数据**：对话回复

#### 3.2.4 向量化节点
- **功能**：将文本转换为向量
- **参数**：
  - `model`：向量化模型
  - `text`：输入文本
- **实现位置**：`internal/nodes/ai/embedding.go`
- **调用方式**：通过 gRPC 调用 AI Service
- **返回数据**：向量数组

---

### 3.3 数据处理节点（Data Nodes）

**节点类型：**

#### 3.3.1 数据转换节点
- **功能**：转换数据格式和结构
- **参数**：
  - `mapping`：字段映射规则
  - `transform_rules`：转换规则（JSONPath、表达式等）
- **实现位置**：`internal/nodes/data/transform.go`
- **执行方式**：本地执行，使用 Go 的 JSON 处理库

#### 3.3.2 数据过滤节点
- **功能**：根据条件过滤数据
- **参数**：
  - `condition`：过滤条件（表达式）
  - `filter_type`：过滤类型（保留/删除匹配项）
- **实现位置**：`internal/nodes/data/filter.go`
- **执行方式**：本地执行

#### 3.3.3 数据聚合节点
- **功能**：对数据进行聚合计算
- **参数**：
  - `group_by`：分组字段
  - `aggregations`：聚合函数（sum、avg、count、max、min）
- **实现位置**：`internal/nodes/data/aggregate.go`
- **执行方式**：本地执行

#### 3.3.4 数据映射节点
- **功能**：对数组中的每个元素进行映射
- **参数**：
  - `mapping_expression`：映射表达式
- **实现位置**：`internal/nodes/data/map.go`
- **执行方式**：本地执行

#### 3.3.5 数据排序节点
- **功能**：对数据进行排序
- **参数**：
  - `sort_by`：排序字段
  - `order`：排序顺序（asc/desc）
- **实现位置**：`internal/nodes/data/sort.go`
- **执行方式**：本地执行

---

**第一部分结束，继续第二部分...**

