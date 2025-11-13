# StructForge 设计分析文档

## 1. 核心设计理念

### 1.1 可扩展性
- **插件化架构**: 节点、AI模型、工具都采用插件化设计，便于扩展
- **适配器模式**: AI模型通过适配器统一接口，新增模型只需实现适配器
- **策略模式**: 不同执行策略可插拔替换

### 1.2 高可用性
- **微服务架构**: 服务独立部署，故障隔离
- **异步执行**: 长时间任务异步处理，不阻塞
- **重试机制**: 失败任务自动重试
- **健康检查**: 服务健康状态监控

### 1.3 易用性
- **可视化编辑**: 拖拽式工作流设计，降低使用门槛
- **实时反馈**: WebSocket实时推送执行状态
- **模板系统**: 提供常用工作流模板

## 2. 技术选型分析

### 2.1 前端框架选择：Vue 3.x

**优势**:
- Composition API提供更好的逻辑复用
- 性能优化（Proxy响应式、Tree-shaking）
- 生态丰富，组件库完善
- 学习曲线平缓

**工作流编辑器选择**:
- **Vue Flow**: 基于Vue的流程图库，功能完善
- **LogicFlow**: 腾讯开源，支持自定义节点
- **自研**: 完全可控，但开发成本高

**推荐**: Vue Flow + 自定义扩展

### 2.2 后端框架选择：Kratos

**为什么选择Kratos**:
- **微服务原生**: 专为微服务设计
- **gRPC + HTTP**: 同时支持gRPC和HTTP Gateway
- **配置管理**: 内置配置系统
- **服务发现**: 支持多种服务发现机制
- **中间件**: 丰富的中间件生态
- **代码生成**: 基于protobuf自动生成代码

**Kratos分层架构**:
```
api/          # API定义（protobuf）
├── internal/
│   ├── biz/  # 业务逻辑层（核心）
│   ├── data/ # 数据访问层
│   └── service/ # 服务层（gRPC实现）
└── cmd/      # 启动入口
```

### 2.3 数据库选择

**PostgreSQL**:
- 支持JSON类型，适合存储工作流定义
- 事务支持完善
- 性能优秀
- 支持全文搜索

**Redis**:
- 缓存热点数据
- 任务队列
- 分布式锁
- Session存储

### 2.4 消息队列（可选）

**RabbitMQ**:
- 适合复杂路由场景
- 消息持久化
- 管理界面友好

**使用场景**:
- 异步工作流执行
- 事件通知
- 日志收集

## 3. 核心模块设计分析

### 3.1 工作流执行引擎

#### 3.1.1 执行模式

**同步执行**:
- 适用场景: 简单快速的工作流（< 30秒）
- 实现: 直接在当前请求中执行
- 优点: 实时返回结果
- 缺点: 占用连接时间长

**异步执行**:
- 适用场景: 复杂耗时的工作流
- 实现: 提交到任务队列，返回任务ID
- 优点: 不阻塞，支持长时间运行
- 缺点: 需要轮询或WebSocket获取结果

#### 3.1.2 节点执行流程

```
节点执行流程:
1. 参数验证
   ├── 必填参数检查
   ├── 参数类型验证
   └── 参数值范围检查

2. 前置处理
   ├── 输入数据转换
   ├── 上下文准备
   └── 依赖检查

3. 节点执行
   ├── 调用节点执行器
   ├── 超时控制
   └── 错误捕获

4. 后置处理
   ├── 结果验证
   ├── 数据转换
   └── 日志记录

5. 结果传递
   ├── 序列化输出
   ├── 传递给下游节点
   └── 更新执行上下文
```

#### 3.1.3 执行上下文设计

```go
type ExecutionContext struct {
    WorkflowID    string
    ExecutionID   string
    UserID        string
    Variables     map[string]interface{}  // 全局变量
    NodeOutputs   map[string]interface{}  // 节点输出
    Status        ExecutionStatus
    StartTime     time.Time
    EndTime       *time.Time
    Error         *ExecutionError
}
```

### 3.2 AI模型接入设计

#### 3.2.1 统一接口设计

```go
// 模型请求
type ModelRequest struct {
    Model     string                 // 模型标识
    Prompt    string                 // 提示词
    Params    map[string]interface{} // 模型参数
    Context   map[string]interface{} // 上下文
}

// 模型响应
type ModelResponse struct {
    Content   string                 // 响应内容
    Metadata  map[string]interface{} // 元数据
    Usage     *UsageInfo             // 使用统计
    Error     error                  // 错误信息
}
```

#### 3.2.2 适配器实现模式

```go
// 基础适配器接口
type ModelAdapter interface {
    Call(ctx context.Context, req *ModelRequest) (*ModelResponse, error)
    ValidateConfig(config *ModelConfig) error
    HealthCheck(ctx context.Context) error
}

// OpenAI适配器示例
type OpenAIAdapter struct {
    client *openai.Client
    config *OpenAIConfig
}

func (a *OpenAIAdapter) Call(ctx context.Context, req *ModelRequest) (*ModelResponse, error) {
    // OpenAI API调用逻辑
}
```

#### 3.2.3 模型配置管理

```yaml
# 模型配置示例
models:
  - id: openai-gpt4
    type: openai
    name: GPT-4
    endpoint: https://api.openai.com/v1
    api_key: ${OPENAI_API_KEY}
    max_tokens: 4096
    temperature: 0.7
    
  - id: gemini-pro
    type: gemini
    name: Gemini Pro
    endpoint: https://generativelanguage.googleapis.com
    api_key: ${GEMINI_API_KEY}
    
  - id: local-llama
    type: ollama
    name: Llama 2
    endpoint: http://localhost:11434
    model: llama2
```

### 3.3 节点系统设计

#### 3.3.1 节点接口定义

```go
// 节点执行器接口
type NodeExecutor interface {
    Execute(ctx context.Context, node *Node, input *NodeInput) (*NodeOutput, error)
    Validate(node *Node) error
    GetSchema() *NodeSchema
}

// 节点定义
type Node struct {
    ID          string
    Type        string
    Name        string
    Config      map[string]interface{}
    Inputs      []*NodePort
    Outputs     []*NodePort
    Position    *Position
}
```

#### 3.3.2 节点类型分类

**触发节点**:
- Webhook: 接收HTTP请求触发
- Timer: 定时触发（Cron表达式）
- Manual: 手动触发
- Event: 事件触发

**AI节点**:
- TextGeneration: 文本生成
- ImageGeneration: 图像生成
- CodeGeneration: 代码生成
- DataAnalysis: 数据分析

**数据处理节点**:
- Transform: 数据转换（JSON Path, 模板等）
- Filter: 数据过滤
- Aggregate: 数据聚合
- Sort: 数据排序

**集成节点**:
- HTTP: HTTP请求
- Database: 数据库操作
- File: 文件操作
- Email: 邮件发送

**控制节点**:
- Condition: 条件判断（if/else）
- Loop: 循环执行
- Parallel: 并行执行
- Delay: 延迟执行

**工具节点**:
- CodeExecutor: 代码执行（Python, JS等）
- Script: 脚本执行
- ExternalTool: 外部工具调用

### 3.4 前端工作流编辑器设计

#### 3.4.1 编辑器架构

```
WorkflowEditor (主组件)
├── Canvas (画布)
│   ├── NodeRenderer (节点渲染)
│   ├── EdgeRenderer (边渲染)
│   └── Background (背景网格)
├── NodePalette (节点面板)
│   └── NodeCategory (节点分类)
├── PropertiesPanel (属性面板)
│   └── NodeConfigForm (节点配置表单)
└── Toolbar (工具栏)
    ├── SaveButton
    ├── RunButton
    └── PreviewButton
```

#### 3.4.2 数据流设计

```typescript
// 工作流数据结构
interface Workflow {
  id: string;
  name: string;
  description: string;
  nodes: Node[];
  edges: Edge[];
  variables: Variable[];
  settings: WorkflowSettings;
}

// 节点数据
interface Node {
  id: string;
  type: string;
  position: { x: number; y: number };
  data: {
    label: string;
    config: Record<string, any>;
  };
}

// 边数据
interface Edge {
  id: string;
  source: string;  // 源节点ID
  target: string;  // 目标节点ID
  sourceHandle: string;  // 源端口
  targetHandle: string;  // 目标端口
}
```

#### 3.4.3 实时执行反馈

**WebSocket通信**:
```typescript
// 执行状态推送
interface ExecutionUpdate {
  executionId: string;
  status: 'running' | 'completed' | 'failed';
  currentNodeId?: string;
  progress?: number;
  logs?: LogEntry[];
  result?: any;
}
```

## 4. 数据模型详细设计

### 4.1 工作流定义存储

**方案1: JSON存储**
- 优点: 灵活，易于修改
- 缺点: 查询困难，版本控制复杂

**方案2: 关系型存储**
- 优点: 查询方便，支持复杂查询
- 缺点: 结构复杂

**推荐: 混合方案**
- 工作流定义（nodes, edges）存JSON
- 元数据（id, name, version等）存关系表
- 使用PostgreSQL的JSONB类型

### 4.2 执行日志设计

**日志级别**:
- DEBUG: 详细调试信息
- INFO: 一般信息
- WARN: 警告信息
- ERROR: 错误信息

**日志存储**:
- 实时日志: Redis（最近N条）
- 历史日志: PostgreSQL（持久化）
- 归档日志: 文件系统或对象存储

## 5. 性能优化考虑

### 5.1 前端优化

- **虚拟滚动**: 大量节点时使用虚拟列表
- **懒加载**: 节点库按需加载
- **防抖节流**: 频繁操作优化
- **WebSocket连接池**: 复用连接

### 5.2 后端优化

- **连接池**: 数据库、Redis连接池
- **缓存策略**: 热点数据缓存
- **异步处理**: 耗时操作异步化
- **批量操作**: 批量处理减少IO

### 5.3 数据库优化

- **索引设计**: 关键字段建立索引
- **分表策略**: 日志表按时间分表
- **读写分离**: 主从复制（生产环境）

## 6. 安全考虑

### 6.1 认证授权

- **JWT Token**: 无状态认证
- **Refresh Token**: 刷新机制
- **RBAC**: 基于角色的访问控制
- **API Key**: 服务间调用认证

### 6.2 数据安全

- **加密存储**: 敏感信息加密
- **传输加密**: HTTPS/WSS
- **输入验证**: 防止注入攻击
- **输出转义**: XSS防护

### 6.3 资源隔离

- **多租户**: 数据隔离
- **资源限制**: 防止资源滥用
- **速率限制**: API限流

## 7. 监控和运维

### 7.1 日志系统

- **结构化日志**: JSON格式
- **日志聚合**: ELK或Loki
- **日志查询**: 支持复杂查询

### 7.2 监控指标

- **服务指标**: CPU、内存、请求数
- **业务指标**: 工作流执行数、成功率
- **AI模型指标**: 调用次数、响应时间、成本

### 7.3 告警机制

- **服务异常**: 服务宕机告警
- **业务异常**: 执行失败率告警
- **资源告警**: 资源使用率告警

## 8. 扩展性设计

### 8.1 节点扩展机制

**方式1: 内置节点**
- 在node-service中实现
- 需要重新编译部署

**方式2: 插件节点**
- 动态加载插件
- 支持热更新

**方式3: 远程节点**
- 通过HTTP/gRPC调用
- 完全解耦

### 8.2 AI模型扩展

- **适配器注册**: 新适配器注册到系统
- **配置驱动**: 通过配置添加新模型
- **SDK支持**: 提供SDK便于开发

## 9. 开发流程

### 9.1 本地开发环境

1. **环境准备**
   - Docker Desktop
   - Node.js 18+
   - Go 1.21+
   - PostgreSQL客户端

2. **启动服务**
   ```bash
   # 启动基础设施
   docker-compose up -d postgres redis
   
   # 启动后端服务
   make run-services
   
   # 启动前端
   cd frontend && npm run dev
   ```

3. **开发调试**
   - 前端: Vite HMR
   - 后端: 热重载（air工具）

### 9.2 代码规范

- **Go**: 遵循Go官方规范，使用gofmt
- **TypeScript**: ESLint + Prettier
- **提交规范**: Conventional Commits

## 10. 总结

### 10.1 核心优势

1. **灵活扩展**: 插件化架构，易于扩展
2. **统一接口**: AI模型统一接入，使用简单
3. **可视化**: 拖拽式设计，降低使用门槛
4. **高性能**: 异步执行，支持大规模工作流
5. **易维护**: 微服务架构，模块清晰

### 10.2 技术难点

1. **工作流执行引擎**: 复杂流程控制
2. **实时状态同步**: WebSocket推送
3. **节点数据传递**: 类型转换和验证
4. **AI模型适配**: 不同模型接口统一

### 10.3 后续优化方向

1. **性能优化**: 执行引擎优化
2. **功能扩展**: 更多节点类型
3. **用户体验**: 编辑器功能增强
4. **企业功能**: 多租户、审计日志等

