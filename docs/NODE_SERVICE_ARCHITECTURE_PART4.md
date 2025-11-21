# Node Service 架构设计文档 - 第四部分

## 十四、节点实现示例

### 14.1 简单节点示例：数据转换节点

```go
// TransformExecutor 数据转换执行器
type TransformExecutor struct{}

func (e *TransformExecutor) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    startTime := time.Now()
    
    // 提取配置
    mapping, ok := req.Config["mapping"].(map[string]interface{})
    if !ok {
        return &ExecuteResponse{
            Success: false,
            Error: &NodeError{
                Code:    ErrCodeInvalidConfig,
                Message: "mapping config is required",
            },
        }, nil
    }
    
    // 转换数据
    outputData := make(map[string]interface{})
    for targetField, sourceField := range mapping {
        sourcePath := sourceField.(string)
        value := getValueByPath(req.InputData, sourcePath)
        setValueByPath(outputData, targetField.(string), value)
    }
    
    duration := time.Since(startTime).Milliseconds()
    
    return &ExecuteResponse{
        Success:    true,
        OutputData: outputData,
        Duration:   duration,
    }, nil
}

func (e *TransformExecutor) Validate(ctx context.Context, config map[string]interface{}) error {
    // 验证配置
    if mapping, ok := config["mapping"]; !ok {
        return fmt.Errorf("mapping is required")
    } else if _, ok := mapping.(map[string]interface{}); !ok {
        return fmt.Errorf("mapping must be a map")
    }
    return nil
}

func (e *TransformExecutor) GetMetadata() *NodeMetadata {
    return &NodeMetadata{
        Type:        "transform",
        Name:        "数据转换",
        Description: "转换数据格式和结构",
        Category:    "data",
        Icon:        "transform",
        Version:     "1.0.0",
        InputPorts: []Port{
            {Name: "input", Type: "any", Required: true, Description: "输入数据"},
        },
        OutputPorts: []Port{
            {Name: "output", Type: "any", Description: "转换后的数据"},
        },
    }
}
```

### 14.2 复杂节点示例：HTTP 请求节点

```go
// HTTPExecutor HTTP 请求执行器
type HTTPExecutor struct {
    httpClient *http.Client
}

func NewHTTPExecutor() *HTTPExecutor {
    return &HTTPExecutor{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (e *HTTPExecutor) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    startTime := time.Now()
    
    // 提取配置
    url := req.Config["url"].(string)
    method := req.Config["method"].(string)
    headers := req.Config["headers"].(map[string]interface{})
    body := req.Config["body"]
    timeout := req.Config["timeout"].(float64)
    
    // 创建请求
    httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
    if err != nil {
        return &ExecuteResponse{
            Success: false,
            Error: &NodeError{
                Code:    ErrCodeInvalidConfig,
                Message: fmt.Sprintf("invalid request: %v", err),
            },
        }, nil
    }
    
    // 设置请求头
    for key, value := range headers {
        httpReq.Header.Set(key, fmt.Sprintf("%v", value))
    }
    
    // 设置请求体
    if body != nil {
        bodyBytes, _ := json.Marshal(body)
        httpReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
        httpReq.ContentLength = int64(len(bodyBytes))
    }
    
    // 设置超时
    if timeout > 0 {
        ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
        defer cancel()
        httpReq = httpReq.WithContext(ctx)
    }
    
    // 发送请求
    resp, err := e.httpClient.Do(httpReq)
    if err != nil {
        return &ExecuteResponse{
            Success: false,
            Error: &NodeError{
                Code:      ErrCodeExternalError,
                Message:   fmt.Sprintf("HTTP request failed: %v", err),
                Retryable: true,
            },
        }, nil
    }
    defer resp.Body.Close()
    
    // 读取响应
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return &ExecuteResponse{
            Success: false,
            Error: &NodeError{
                Code:    ErrCodeExecutionFailed,
                Message: fmt.Sprintf("failed to read response: %v", err),
            },
        }, nil
    }
    
    // 解析响应
    var responseData interface{}
    if err := json.Unmarshal(respBody, &responseData); err != nil {
        // 如果不是 JSON，返回原始字符串
        responseData = string(respBody)
    }
    
    duration := time.Since(startTime).Milliseconds()
    
    return &ExecuteResponse{
        Success: true,
        OutputData: map[string]interface{}{
            "status_code": resp.StatusCode,
            "headers":     resp.Header,
            "body":        responseData,
        },
        Duration: duration,
    }, nil
}
```

### 14.3 调用外部服务节点示例：AI 文本生成节点

```go
// TextGenerationExecutor AI 文本生成执行器
type TextGenerationExecutor struct {
    aiClient ai.AIClient // AI Service 的 gRPC 客户端
}

func NewTextGenerationExecutor(aiClient ai.AIClient) *TextGenerationExecutor {
    return &TextGenerationExecutor{
        aiClient: aiClient,
    }
}

func (e *TextGenerationExecutor) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
    startTime := time.Now()
    
    // 提取配置
    model := req.Config["model"].(string)
    prompt := req.Config["prompt"].(string)
    maxTokens := int32(req.Config["max_tokens"].(float64))
    temperature := float32(req.Config["temperature"].(float64))
    
    // 构建提示词（支持变量替换）
    finalPrompt := replaceVariables(prompt, req.Context.Variables)
    
    // 调用 AI Service
    aiReq := &ai.TextGenerationRequest{
        Model:      model,
        Prompt:     finalPrompt,
        MaxTokens:  maxTokens,
        Temperature: temperature,
    }
    
    aiResp, err := e.aiClient.GenerateText(ctx, aiReq)
    if err != nil {
        return &ExecuteResponse{
            Success: false,
            Error: &NodeError{
                Code:      ErrCodeExternalError,
                Message:   fmt.Sprintf("AI service error: %v", err),
                Retryable: true,
            },
        }, nil
    }
    
    duration := time.Since(startTime).Milliseconds()
    
    return &ExecuteResponse{
        Success: true,
        OutputData: map[string]interface{}{
            "text":      aiResp.Text,
            "usage":     aiResp.Usage,
            "model":     aiResp.Model,
            "finish_reason": aiResp.FinishReason,
        },
        Duration: duration,
    }, nil
}
```

---

## 十五、节点初始化流程

### 15.1 服务启动时的节点注册

```go
// 在 main.go 或 wire.go 中初始化节点
func initNodes(registry *NodeRegistry, aiClient ai.AIClient, toolClient tool.ToolClient) error {
    // 注册触发节点
    registry.RegisterNode("webhook", NewWebhookExecutor())
    registry.RegisterNode("timer", NewTimerExecutor())
    registry.RegisterNode("manual", NewManualExecutor())
    
    // 注册 AI 节点
    registry.RegisterNode("text_generation", NewTextGenerationExecutor(aiClient))
    registry.RegisterNode("image_generation", NewImageGenerationExecutor(aiClient))
    registry.RegisterNode("chat", NewChatExecutor(aiClient))
    
    // 注册数据处理节点
    registry.RegisterNode("transform", NewTransformExecutor())
    registry.RegisterNode("filter", NewFilterExecutor())
    registry.RegisterNode("aggregate", NewAggregateExecutor())
    
    // 注册集成节点
    registry.RegisterNode("http", NewHTTPExecutor())
    registry.RegisterNode("database", NewDatabaseExecutor())
    registry.RegisterNode("email", NewEmailExecutor())
    
    // 注册控制节点
    registry.RegisterNode("condition", NewConditionExecutor())
    registry.RegisterNode("loop", NewLoopExecutor())
    registry.RegisterNode("parallel", NewParallelExecutor())
    
    // 注册工具节点
    registry.RegisterNode("script", NewScriptExecutor(toolClient))
    registry.RegisterNode("file_operation", NewFileOperationExecutor())
    
    return nil
}
```

### 15.2 动态节点注册（插件）

```go
// 加载插件节点
func loadPluginNodes(pluginManager *PluginManager, registry *NodeRegistry) error {
    // 扫描插件目录
    pluginDir := "./plugins"
    files, err := os.ReadDir(pluginDir)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if strings.HasSuffix(file.Name(), ".so") {
            // 加载插件
            if err := pluginManager.LoadPlugin(filepath.Join(pluginDir, file.Name())); err != nil {
                log.Warn("Failed to load plugin", "file", file.Name(), "error", err)
                continue
            }
        }
    }
    
    return nil
}
```

---

## 十六、节点配置模式（Schema）

### 16.1 配置模式定义

```go
// ConfigSchema 配置参数模式
type ConfigSchema struct {
    Properties map[string]*PropertySchema // 属性定义
    Required   []string                    // 必需属性
}

// PropertySchema 属性模式
type PropertySchema struct {
    Type        string                 // 类型：string, number, boolean, object, array
    Description string                 // 描述
    Default     interface{}            // 默认值
    Enum        []interface{}          // 枚举值
    Minimum     *float64               // 最小值（number 类型）
    Maximum     *float64               // 最大值（number 类型）
    Pattern     string                 // 正则表达式（string 类型）
    Items       *PropertySchema        // 数组元素模式（array 类型）
    Properties  map[string]*PropertySchema // 对象属性（object 类型）
}
```

### 16.2 配置模式示例

**HTTP 节点配置模式：**
```go
func (e *HTTPExecutor) GetConfigSchema() *ConfigSchema {
    return &ConfigSchema{
        Properties: map[string]*PropertySchema{
            "url": {
                Type:        "string",
                Description: "请求 URL",
                Required:    true,
                Pattern:     "^https?://",
            },
            "method": {
                Type:        "string",
                Description: "HTTP 方法",
                Default:     "GET",
                Enum:        []interface{}{"GET", "POST", "PUT", "DELETE", "PATCH"},
            },
            "headers": {
                Type:        "object",
                Description: "请求头",
            },
            "body": {
                Type:        "object",
                Description: "请求体",
            },
            "timeout": {
                Type:        "number",
                Description: "超时时间（秒）",
                Default:     30,
                Minimum:     floatPtr(1),
                Maximum:     floatPtr(300),
            },
        },
        Required: []string{"url"},
    }
}
```

---

## 十七、数据模式（Schema）

### 17.1 输入输出数据模式

```go
// DataSchema 数据模式
type DataSchema struct {
    Type       string                 // 类型：object, array, string, number, boolean
    Properties map[string]*DataSchema // 对象属性
    Items      *DataSchema            // 数组元素
    Required   []string               // 必需字段
}

// GetInputSchema 获取输入数据模式
func (e *TransformExecutor) GetInputSchema() *DataSchema {
    return &DataSchema{
        Type: "object",
        Properties: map[string]*DataSchema{
            "input": {
                Type: "any",
                Description: "输入数据",
            },
        },
        Required: []string{"input"},
    }
}

// GetOutputSchema 获取输出数据模式
func (e *TransformExecutor) GetOutputSchema() *DataSchema {
    return &DataSchema{
        Type: "object",
        Properties: map[string]*DataSchema{
            "output": {
                Type: "any",
                Description: "转换后的数据",
            },
        },
    }
}
```

---

## 十八、测试策略

### 18.1 单元测试

**执行器测试：**
- 测试正常执行流程
- 测试错误处理
- 测试参数验证
- 测试边界情况

**注册表测试：**
- 测试节点注册
- 测试节点查找
- 测试节点列表

### 18.2 集成测试

**服务集成测试：**
- 测试与 AI Service 的交互
- 测试与 Tool Service 的交互
- 测试与 Workflow Service 的交互

**端到端测试：**
- 测试完整的工作流执行
- 测试节点间数据传递
- 测试错误恢复

---

## 十九、监控和日志

### 19.1 监控指标

**节点执行指标：**
- 节点执行次数
- 节点执行成功率
- 节点执行平均时长
- 节点执行错误率

**服务指标：**
- 服务 QPS
- 服务响应时间
- 服务错误率
- 并发执行数

### 19.2 日志记录

**执行日志：**
- 节点执行开始/结束
- 节点执行参数
- 节点执行结果
- 节点执行错误

**日志级别：**
- DEBUG：详细调试信息
- INFO：一般信息
- WARN：警告信息
- ERROR：错误信息

---

## 二十、总结

### 20.1 核心设计要点

1. **统一服务入口**：所有节点通过 Node Service 统一管理
2. **插件化架构**：支持自定义节点扩展
3. **清晰的接口**：统一的执行器接口，便于开发
4. **灵活的扩展**：可以水平扩展，也可以插件化扩展
5. **完善的错误处理**：统一的错误处理和重试机制

### 20.2 实施步骤

**Phase 1：基础框架**
1. 实现 Node Registry
2. 实现 Executor Factory
3. 实现基础执行器接口
4. 实现 gRPC 服务

**Phase 2：核心节点**
1. 实现触发节点（webhook、timer、manual）
2. 实现数据处理节点（transform、filter）
3. 实现控制节点（condition、loop）
4. 实现集成节点（http、database）

**Phase 3：高级节点**
1. 实现 AI 节点（调用 AI Service）
2. 实现工具节点（调用 Tool Service）
3. 实现复杂控制节点（parallel、switch）

**Phase 4：优化和扩展**
1. 实现插件系统
2. 性能优化（池化、缓存）
3. 监控和日志完善

---

**文档完成**

