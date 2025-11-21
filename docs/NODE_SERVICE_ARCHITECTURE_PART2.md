# Node Service 架构设计文档 - 第二部分

## 四、节点分类与实现（续）

### 4.1 集成节点（Integration Nodes）

**节点类型：**

#### 4.1.1 HTTP 请求节点
- **功能**：发送 HTTP 请求
- **参数**：
  - `url`：请求 URL
  - `method`：HTTP 方法
  - `headers`：请求头
  - `body`：请求体
  - `timeout`：超时时间
  - `retry`：重试配置
- **实现位置**：`internal/nodes/integration/http.go`
- **执行方式**：本地执行，使用 Go 的 `net/http` 包
- **特殊处理**：支持认证（Basic Auth、Bearer Token、API Key）

#### 4.1.2 数据库操作节点
- **功能**：执行数据库操作
- **参数**：
  - `database_type`：数据库类型（PostgreSQL、MySQL、MongoDB 等）
  - `connection_string`：连接字符串（加密存储）
  - `query_type`：操作类型（SELECT、INSERT、UPDATE、DELETE）
  - `query`：SQL 查询或操作
  - `parameters`：查询参数
- **实现位置**：`internal/nodes/integration/database.go`
- **执行方式**：本地执行，使用相应的数据库驱动
- **安全考虑**：连接字符串加密存储，执行时解密

#### 4.1.3 消息队列节点
- **功能**：发送/接收消息队列消息
- **参数**：
  - `queue_type`：队列类型（RabbitMQ、Redis Streams、Kafka 等）
  - `operation`：操作类型（send/receive）
  - `queue_name`：队列名称
  - `message`：消息内容
- **实现位置**：`internal/nodes/integration/message_queue.go`
- **执行方式**：本地执行，使用相应的消息队列客户端

#### 4.1.4 邮件发送节点
- **功能**：发送邮件
- **参数**：
  - `smtp_server`：SMTP 服务器
  - `smtp_port`：SMTP 端口
  - `username`：用户名
  - `password`：密码（加密）
  - `to`：收件人
  - `subject`：主题
  - `body`：正文
  - `attachments`：附件
- **实现位置**：`internal/nodes/integration/email.go`
- **执行方式**：本地执行，使用 Go 的 `net/smtp` 包

#### 4.1.5 Webhook 发送节点
- **功能**：向外部 URL 发送 Webhook
- **参数**：
  - `url`：目标 URL
  - `method`：HTTP 方法
  - `headers`：请求头
  - `body`：请求体
  - `retry`：重试配置
- **实现位置**：`internal/nodes/integration/webhook.go`
- **执行方式**：本地执行，类似 HTTP 请求节点

---

### 4.2 控制节点（Control Nodes）

**节点类型：**

#### 4.2.1 条件判断节点
- **功能**：根据条件分支执行
- **参数**：
  - `conditions`：条件列表
    - `condition`：条件表达式
    - `output_port`：满足条件时的输出端口
  - `default_port`：默认输出端口
- **实现位置**：`internal/nodes/control/condition.go`
- **执行方式**：本地执行，解析条件表达式
- **特殊处理**：支持多个分支，类似 if-else if-else

#### 4.2.2 循环节点
- **功能**：循环执行子节点
- **参数**：
  - `loop_type`：循环类型（for、while、foreach）
  - `condition`：循环条件
  - `max_iterations`：最大迭代次数
  - `break_condition`：跳出条件
- **实现位置**：`internal/nodes/control/loop.go`
- **执行方式**：本地执行，管理循环状态
- **特殊处理**：需要与 Workflow Service 协调，支持嵌套循环

#### 4.2.3 并行执行节点
- **功能**：并行执行多个分支
- **参数**：
  - `branches`：分支列表
  - `wait_all`：是否等待所有分支完成
  - `timeout`：超时时间
- **实现位置**：`internal/nodes/control/parallel.go`
- **执行方式**：本地执行，使用 Goroutine 管理并行
- **特殊处理**：需要协调多个分支的执行结果

#### 4.2.4 分支选择节点
- **功能**：根据值选择分支（类似 switch）
- **参数**：
  - `switch_value`：用于判断的值
  - `cases`：分支列表
    - `value`：匹配值
    - `output_port`：输出端口
  - `default_port`：默认输出端口
- **实现位置**：`internal/nodes/control/switch.go`
- **执行方式**：本地执行

#### 4.2.5 延迟节点
- **功能**：延迟指定时间后继续执行
- **参数**：
  - `delay_type`：延迟类型（固定时间、动态时间）
  - `delay`：延迟时间（秒或毫秒）
- **实现位置**：`internal/nodes/control/delay.go`
- **执行方式**：本地执行，使用 `time.Sleep` 或定时器

---

### 4.3 工具节点（Tool Nodes）

**节点类型：**

#### 4.3.1 脚本执行节点
- **功能**：执行脚本代码
- **参数**：
  - `script_type`：脚本类型（JavaScript、Python、Lua 等）
  - `script_code`：脚本代码
  - `timeout`：超时时间
  - `sandbox`：是否沙箱执行
- **实现位置**：`internal/nodes/tool/script.go`
- **执行方式**：可以本地执行（Go 脚本引擎）或调用 Tool Service
- **安全考虑**：沙箱执行，限制资源使用

#### 4.3.2 文件操作节点
- **功能**：文件读写操作
- **参数**：
  - `operation`：操作类型（read、write、delete、list）
  - `file_path`：文件路径
  - `content`：文件内容（写入时）
  - `encoding`：编码方式
- **实现位置**：`internal/nodes/tool/file_operation.go`
- **执行方式**：本地执行，使用 Go 的文件操作 API
- **安全考虑**：路径验证，防止路径遍历攻击

#### 4.3.3 系统命令节点
- **功能**：执行系统命令
- **参数**：
  - `command`：命令
  - `args`：命令参数
  - `timeout`：超时时间
  - `working_dir`：工作目录
- **实现位置**：`internal/nodes/tool/system_command.go`
- **执行方式**：可以本地执行或调用 Tool Service
- **安全考虑**：命令白名单，限制可执行命令

---

## 五、核心接口设计

### 5.1 执行器接口

```go
// NodeExecutor 节点执行器接口
type NodeExecutor interface {
    // Execute 执行节点
    Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error)
    
    // Validate 验证节点配置
    Validate(ctx context.Context, config map[string]interface{}) error
    
    // GetMetadata 获取节点元数据
    GetMetadata() *NodeMetadata
    
    // GetInputSchema 获取输入数据模式
    GetInputSchema() *DataSchema
    
    // GetOutputSchema 获取输出数据模式
    GetOutputSchema() *DataSchema
}
```

### 5.2 执行请求结构

```go
// ExecuteRequest 执行请求
type ExecuteRequest struct {
    // 节点信息
    NodeID      string                 // 节点ID
    NodeType    string                 // 节点类型
    NodeName    string                 // 节点名称
    
    // 节点配置
    Config      map[string]interface{} // 节点配置参数
    
    // 输入数据
    InputData   map[string]interface{} // 输入数据（来自上游节点）
    
    // 执行上下文
    ExecutionID string                 // 执行ID
    WorkflowID  string                 // 工作流ID
    Context     *ExecutionContext      // 执行上下文
    
    // 元数据
    Metadata    map[string]interface{} // 元数据
}
```

### 5.3 执行响应结构

```go
// ExecuteResponse 执行响应
type ExecuteResponse struct {
    // 执行结果
    Success     bool                   // 是否成功
    OutputData  map[string]interface{} // 输出数据
    Error       *NodeError             // 错误信息（如果失败）
    
    // 执行信息
    Duration    int64                  // 执行时长（毫秒）
    Logs        []LogEntry              // 执行日志
    Metrics     *ExecutionMetrics      // 执行指标
    
    // 状态信息
    Status      string                 // 执行状态
    NextNodes   []string               // 下一个节点ID（用于控制节点）
}
```

### 5.4 节点元数据结构

```go
// NodeMetadata 节点元数据
type NodeMetadata struct {
    Type        string                 // 节点类型
    Name        string                 // 节点名称
    Description string                 // 节点描述
    Category    string                 // 节点分类
    Icon        string                 // 节点图标
    Version     string                 // 节点版本
    
    // 输入输出
    InputPorts  []Port                 // 输入端口
    OutputPorts []Port                 // 输出端口
    
    // 配置参数
    ConfigSchema *ConfigSchema         // 配置参数模式
    
    // 文档
    Documentation string               // 节点文档
    Examples     []Example            // 使用示例
}
```

---

## 六、节点注册表设计

### 6.1 注册表结构

```go
// NodeRegistry 节点注册表
type NodeRegistry struct {
    executors map[string]NodeExecutor  // 节点类型 -> 执行器映射
    metadata  map[string]*NodeMetadata  // 节点类型 -> 元数据映射
    mutex     sync.RWMutex             // 读写锁
}
```

### 6.2 注册流程

```go
// RegisterNode 注册节点
func (r *NodeRegistry) RegisterNode(nodeType string, executor NodeExecutor) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    // 验证节点类型
    if err := r.validateNodeType(nodeType); err != nil {
        return err
    }
    
    // 获取元数据
    metadata := executor.GetMetadata()
    
    // 注册执行器和元数据
    r.executors[nodeType] = executor
    r.metadata[nodeType] = metadata
    
    return nil
}
```

### 6.3 节点发现

```go
// GetExecutor 获取执行器
func (r *NodeRegistry) GetExecutor(nodeType string) (NodeExecutor, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    executor, exists := r.executors[nodeType]
    if !exists {
        return nil, fmt.Errorf("node type %s not found", nodeType)
    }
    
    return executor, nil
}

// GetMetadata 获取元数据
func (r *NodeRegistry) GetMetadata(nodeType string) (*NodeMetadata, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    metadata, exists := r.metadata[nodeType]
    if !exists {
        return nil, fmt.Errorf("node type %s not found", nodeType)
    }
    
    return metadata, nil
}

// ListNodes 列出所有节点
func (r *NodeRegistry) ListNodes() []*NodeMetadata {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    nodes := make([]*NodeMetadata, 0, len(r.metadata))
    for _, metadata := range r.metadata {
        nodes = append(nodes, metadata)
    }
    
    return nodes
}
```

---

## 七、执行器工厂设计

### 7.1 工厂接口

```go
// ExecutorFactory 执行器工厂
type ExecutorFactory struct {
    registry *NodeRegistry
    pool     map[string]*ExecutorPool // 执行器池（可选，用于性能优化）
}

// CreateExecutor 创建执行器
func (f *ExecutorFactory) CreateExecutor(nodeType string) (NodeExecutor, error) {
    // 从注册表获取执行器
    executor, err := f.registry.GetExecutor(nodeType)
    if err != nil {
        return nil, err
    }
    
    // 如果使用池化，从池中获取
    if pool, exists := f.pool[nodeType]; exists {
        return pool.Get(), nil
    }
    
    return executor, nil
}

// ReleaseExecutor 释放执行器（如果使用池化）
func (f *ExecutorFactory) ReleaseExecutor(nodeType string, executor NodeExecutor) {
    if pool, exists := f.pool[nodeType]; exists {
        pool.Put(executor)
    }
}
```

### 7.2 执行器池（可选优化）

```go
// ExecutorPool 执行器池
type ExecutorPool struct {
    pool sync.Pool
}

// Get 从池中获取执行器
func (p *ExecutorPool) Get() NodeExecutor {
    executor := p.pool.Get()
    if executor == nil {
        return nil
    }
    return executor.(NodeExecutor)
}

// Put 将执行器放回池中
func (p *ExecutorPool) Put(executor NodeExecutor) {
    p.pool.Put(executor)
}
```

---

**第二部分结束，继续第三部分...**

