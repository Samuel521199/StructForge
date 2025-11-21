# Node Service 架构设计文档 - 第三部分

## 八、插件系统设计

### 8.1 插件架构

**插件类型：**
- **Go 插件**：编译为 `.so` 文件的 Go 代码
- **脚本插件**：JavaScript、Python 等脚本
- **远程插件**：通过 gRPC 调用的远程服务

### 8.2 插件接口

```go
// Plugin 插件接口
type Plugin interface {
    // GetInfo 获取插件信息
    GetInfo() *PluginInfo
    
    // CreateExecutor 创建执行器
    CreateExecutor(config map[string]interface{}) (NodeExecutor, error)
    
    // Validate 验证插件配置
    Validate(config map[string]interface{}) error
}

// PluginInfo 插件信息
type PluginInfo struct {
    ID          string   // 插件ID
    Name        string   // 插件名称
    Version     string   // 插件版本
    NodeTypes   []string // 支持的节点类型
    Description string   // 插件描述
}
```

### 8.3 插件加载

```go
// PluginManager 插件管理器
type PluginManager struct {
    plugins map[string]Plugin
    mutex   sync.RWMutex
}

// LoadPlugin 加载插件
func (m *PluginManager) LoadPlugin(pluginPath string) error {
    // 加载插件文件
    plugin, err := plugin.Open(pluginPath)
    if err != nil {
        return err
    }
    
    // 查找插件符号
    symbol, err := plugin.Lookup("Plugin")
    if err != nil {
        return err
    }
    
    // 创建插件实例
    p := symbol.(Plugin)
    
    // 注册插件
    info := p.GetInfo()
    m.plugins[info.ID] = p
    
    // 注册插件支持的节点类型
    for _, nodeType := range info.NodeTypes {
        executor, _ := p.CreateExecutor(nil)
        m.registry.RegisterNode(nodeType, executor)
    }
    
    return nil
}
```

---

## 九、与外部服务的交互

### 9.1 与 AI Service 的交互

**交互方式：gRPC**

```go
// AI 节点执行器示例
type TextGenerationExecutor struct {
    aiClient ai.AIClient // AI Service 的 gRPC 客户端
}

func (e *TextGenerationExecutor) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    // 提取参数
    model := req.Config["model"].(string)
    prompt := req.Config["prompt"].(string)
    
    // 调用 AI Service
    aiReq := &ai.TextGenerationRequest{
        Model:  model,
        Prompt: prompt,
    }
    
    aiResp, err := e.aiClient.GenerateText(ctx, aiReq)
    if err != nil {
        return &ExecuteResponse{
            Success: false,
            Error:   &NodeError{Message: err.Error()},
        }, nil
    }
    
    // 返回结果
    return &ExecuteResponse{
        Success: true,
        OutputData: map[string]interface{}{
            "text": aiResp.Text,
        },
    }, nil
}
```

### 9.2 与 Tool Service 的交互

**交互方式：gRPC**

```go
// 工具节点执行器示例
type ScriptExecutor struct {
    toolClient tool.ToolClient // Tool Service 的 gRPC 客户端
}

func (e *ScriptExecutor) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    // 提取参数
    scriptType := req.Config["script_type"].(string)
    scriptCode := req.Config["script_code"].(string)
    
    // 调用 Tool Service
    toolReq := &tool.ExecuteScriptRequest{
        ScriptType: scriptType,
        ScriptCode: scriptCode,
    }
    
    toolResp, err := e.toolClient.ExecuteScript(ctx, toolReq)
    if err != nil {
        return &ExecuteResponse{
            Success: false,
            Error:   &NodeError{Message: err.Error()},
        }, nil
    }
    
    // 返回结果
    return &ExecuteResponse{
        Success: true,
        OutputData: map[string]interface{}{
            "result": toolResp.Result,
        },
    }, nil
}
```

### 9.3 与 Workflow Service 的交互

**交互方式：gRPC（Node Service 作为服务端）**

```go
// gRPC 服务实现
type NodeService struct {
    factory *ExecutorFactory
}

// ExecuteNode 执行节点
func (s *NodeService) ExecuteNode(ctx context.Context, req *pb.ExecuteNodeRequest) (*pb.ExecuteNodeResponse, error) {
    // 创建执行器
    executor, err := s.factory.CreateExecutor(req.NodeType)
    if err != nil {
        return nil, err
    }
    
    // 构建执行请求
    executeReq := &ExecuteRequest{
        NodeID:    req.NodeId,
        NodeType:  req.NodeType,
        Config:    req.Config,
        InputData: req.InputData,
        Context:   convertContext(req.Context),
    }
    
    // 执行节点
    resp, err := executor.Execute(ctx, executeReq)
    if err != nil {
        return nil, err
    }
    
    // 转换为 gRPC 响应
    return &pb.ExecuteNodeResponse{
        Success:    resp.Success,
        OutputData:  resp.OutputData,
        Error:       convertError(resp.Error),
        Duration:    resp.Duration,
        Logs:        resp.Logs,
    }, nil
}
```

---

## 十、API 设计（Protobuf）

### 10.1 API 定义

```protobuf
syntax = "proto3";

package node.v1;

import "google/api/annotations.proto";

// Node Service API
service NodeService {
  // 执行节点
  rpc ExecuteNode(ExecuteNodeRequest) returns (ExecuteNodeResponse) {
    option (google.api.http) = {
      post: "/api/v1/nodes/execute"
      body: "*"
    };
  }
  
  // 验证节点配置
  rpc ValidateNode(ValidateNodeRequest) returns (ValidateNodeResponse) {
    option (google.api.http) = {
      post: "/api/v1/nodes/validate"
      body: "*"
    };
  }
  
  // 获取节点元数据
  rpc GetNodeMetadata(GetNodeMetadataRequest) returns (GetNodeMetadataResponse) {
    option (google.api.http) = {
      get: "/api/v1/nodes/{node_type}/metadata"
    };
  }
  
  // 列出所有节点
  rpc ListNodes(ListNodesRequest) returns (ListNodesResponse) {
    option (google.api.http) = {
      get: "/api/v1/nodes"
    };
  }
}

// 执行节点请求
message ExecuteNodeRequest {
  string node_id = 1;           // 节点ID
  string node_type = 2;         // 节点类型
  string node_name = 3;         // 节点名称
  map<string, Value> config = 4; // 节点配置
  map<string, Value> input_data = 5; // 输入数据
  ExecutionContext context = 6; // 执行上下文
}

// 执行节点响应
message ExecuteNodeResponse {
  bool success = 1;             // 是否成功
  map<string, Value> output_data = 2; // 输出数据
  NodeError error = 3;          // 错误信息
  int64 duration = 4;           // 执行时长（毫秒）
  repeated LogEntry logs = 5;    // 执行日志
  repeated string next_nodes = 6; // 下一个节点ID（控制节点）
}

// 执行上下文
message ExecutionContext {
  string execution_id = 1;       // 执行ID
  string workflow_id = 2;        // 工作流ID
  map<string, Value> variables = 3; // 全局变量
  map<string, Value> node_outputs = 4; // 节点输出
}

// 节点元数据
message NodeMetadata {
  string type = 1;              // 节点类型
  string name = 2;               // 节点名称
  string description = 3;       // 节点描述
  string category = 4;           // 节点分类
  string icon = 5;              // 节点图标
  string version = 6;            // 节点版本
  repeated Port input_ports = 7; // 输入端口
  repeated Port output_ports = 8; // 输出端口
  ConfigSchema config_schema = 9; // 配置参数模式
}

// 端口定义
message Port {
  string name = 1;               // 端口名称
  string type = 2;               // 数据类型
  bool required = 3;              // 是否必需
  string description = 4;         // 端口描述
}
```

---

## 十一、错误处理设计

### 11.1 错误类型

```go
// NodeError 节点错误
type NodeError struct {
    Code      string                 // 错误代码
    Message   string                 // 错误消息
    Details   map[string]interface{} // 错误详情
    Retryable bool                   // 是否可重试
    Cause     error                  // 原始错误
}

// 错误代码定义
const (
    ErrCodeInvalidConfig    = "INVALID_CONFIG"    // 配置无效
    ErrCodeExecutionFailed  = "EXECUTION_FAILED"  // 执行失败
    ErrCodeTimeout          = "TIMEOUT"           // 超时
    ErrCodeResourceExhausted = "RESOURCE_EXHAUSTED" // 资源耗尽
    ErrCodeExternalError    = "EXTERNAL_ERROR"    // 外部服务错误
)
```

### 11.2 错误处理策略

**节点级别错误处理：**
- 每个节点可以配置错误处理策略
- 策略选项：
  - `continue`：继续执行下一个节点
  - `retry`：重试当前节点
  - `fail`：标记执行失败
  - `skip`：跳过当前节点

**重试机制：**
- 固定间隔重试
- 指数退避重试
- 最大重试次数限制

---

## 十二、性能优化

### 12.1 执行器池化

**适用场景：**
- 执行器创建成本高
- 执行器可以复用
- 高并发场景

**实现方式：**
- 使用 `sync.Pool` 实现对象池
- 执行器实现 `Reset()` 方法，重置状态
- 从池中获取，使用后放回

### 12.2 并发控制

**节点级别并发：**
- 每个节点可以设置最大并发数
- 使用信号量控制并发

**服务级别并发：**
- 全局最大并发执行数
- 防止服务过载

### 12.3 缓存策略

**配置缓存：**
- 节点配置缓存（避免重复解析）
- 元数据缓存（避免重复查询）

**结果缓存：**
- 相同输入的节点执行结果缓存
- 适用于纯函数节点（无副作用）

---

## 十三、安全考虑

### 13.1 输入验证

**参数验证：**
- 类型验证
- 范围验证
- 格式验证（正则表达式）
- 业务规则验证

**数据验证：**
- 输入数据大小限制
- 输入数据结构验证
- 恶意数据检测

### 13.2 资源限制

**执行限制：**
- 执行时间限制
- 内存使用限制
- CPU 使用限制
- 网络请求限制

**文件操作限制：**
- 文件大小限制
- 路径白名单
- 防止路径遍历攻击

### 13.3 敏感信息保护

**配置加密：**
- API Key 加密存储
- 数据库密码加密
- 连接字符串加密

**执行隔离：**
- 每个执行在独立上下文中
- 防止数据泄露
- 支持多租户隔离

---

**第三部分结束，继续第四部分...**

