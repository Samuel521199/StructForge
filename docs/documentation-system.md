# StructForge 文档系统设计

## 1. 文档系统概述

文档系统需要提供：
1. **API文档**: RESTful API的完整文档，支持在线测试
2. **节点使用说明**: 每个节点的详细使用说明、参数说明、示例
3. **工作流使用说明**: 工作流设计指南、最佳实践、示例

## 2. 架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      文档访问层                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Swagger UI  │  │  前端文档页   │  │  独立文档站   │      │
│  │  (API文档)   │  │  (集成文档)   │  │  (可选)      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway                               │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │ /api/docs    │  │ /api/nodes   │                        │
│  │ (OpenAPI)    │  │ (节点说明)   │                        │
│  └──────────────┘  └──────────────┘                        │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      文档服务层                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │文档服务      │  │节点服务      │  │工作流服务    │      │
│  │(可选独立)    │  │(节点说明)    │  │(工作流说明)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                               │
│  ┌──────────────┐  ┌──────────────┐                        │
│  │PostgreSQL    │  │文件系统/MD   │                        │
│  │(节点元数据)  │  │(文档内容)    │                        │
│  └──────────────┘  └──────────────┘                        │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 推荐方案：混合架构

**核心原则**:
- **API文档**: 使用OpenAPI 3.0，通过后端API Gateway提供Swagger UI
- **节点说明**: 通过API动态提供，存储在数据库，支持版本管理
- **工作流说明**: 通过API提供，前端集成展示
- **独立文档站**: 可选，用于完整的产品文档

## 3. API文档设计

### 3.1 OpenAPI集成方案

#### 方案A：后端集成（推荐）

**实现方式**:
- 在API Gateway中集成Swagger UI
- 使用OpenAPI 3.0规范
- 自动从代码生成或手动维护

**优点**:
- ✅ 与API代码同步，易于维护
- ✅ 支持在线测试API
- ✅ 独立部署，不影响前端
- ✅ 可以集成认证，保护文档

**技术选型**:
- **Go**: 使用 `swaggo/swag` 或 `getkin/kin-openapi`
- **Swagger UI**: 通过静态文件或CDN提供
- **OpenAPI生成**: 从protobuf或代码注释生成

**访问路径**:
- 开发环境: `http://localhost:8000/api/docs`
- 生产环境: `https://api.structforge.com/docs`

#### 方案B：前端集成

**实现方式**:
- 前端集成Swagger UI组件
- 从后端获取OpenAPI JSON
- 在管理后台中展示

**优点**:
- ✅ 统一的前端体验
- ✅ 可以自定义样式

**缺点**:
- ❌ 增加前端复杂度
- ❌ 需要额外的路由管理

### 3.2 OpenAPI规范设计

```yaml
# openapi.yaml 示例结构
openapi: 3.0.3
info:
  title: StructForge API
  version: 1.0.0
  description: StructForge工作流平台API文档

servers:
  - url: http://localhost:8000/api
    description: 本地开发环境
  - url: https://api.structforge.com
    description: 生产环境

paths:
  /v1/users/login:
    post:
      summary: 用户登录
      tags:
        - 用户管理
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: 登录成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/LoginResponse'

components:
  schemas:
    LoginRequest:
      type: object
      required:
        - username
        - password
      properties:
        username:
          type: string
          example: admin
        password:
          type: string
          format: password
          example: password123

  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

### 3.3 实现位置

**推荐**: 集成到API Gateway中

```
backend/
├── api/
│   ├── gateway/
│   │   ├── docs/              # Swagger文档
│   │   │   ├── swagger.yaml   # OpenAPI规范
│   │   │   └── swagger.json
│   │   ├── handlers/
│   │   │   └── docs.go        # 文档路由处理
│   │   └── middleware/
│   │       └── swagger.go     # Swagger中间件
```

## 4. 节点使用说明设计

### 4.1 节点说明数据结构

```go
// 节点说明结构
type NodeDocumentation struct {
    NodeType      string                 // 节点类型标识
    Name          string                 // 节点名称
    Category      string                 // 节点分类
    Description   string                 // 节点描述
    Icon          string                 // 图标URL
    Version       string                 // 版本号
    
    // 输入说明
    Inputs        []PortDocumentation   // 输入端口说明
    // 输出说明
    Outputs       []PortDocumentation   // 输出端口说明
    // 配置参数说明
    Parameters    []ParameterDocumentation  // 参数说明
    
    // 使用说明
    Usage         string                 // 使用说明（Markdown）
    Examples      []Example              // 使用示例
    BestPractices []string               // 最佳实践
    
    // 元数据
    Tags          []string               // 标签
    Author        string                 // 作者
    UpdatedAt     time.Time              // 更新时间
}

type PortDocumentation struct {
    Name        string      // 端口名称
    Type        string      // 数据类型
    Required    bool        // 是否必需
    Description string      // 说明
    Default     interface{} // 默认值
}

type ParameterDocumentation struct {
    Name        string      // 参数名
    Type        string      // 参数类型
    Required    bool        // 是否必需
    Description string      // 说明
    Default     interface{} // 默认值
    Options     []string    // 可选值（如果是枚举）
    Validation  string      // 验证规则
    Example     interface{} // 示例值
}

type Example struct {
    Title       string                 // 示例标题
    Description string                 // 示例描述
    Config      map[string]interface{} // 配置示例
    Input       interface{}            // 输入示例
    Output      interface{}            // 输出示例
    Code        string                 // 代码示例（如果有）
}
```

### 4.2 存储方案

#### 方案A：数据库存储（推荐）

**优点**:
- ✅ 动态更新，无需重新部署
- ✅ 支持版本管理
- ✅ 支持多语言
- ✅ 可以关联节点代码

**实现**:
```sql
-- 节点文档表
CREATE TABLE node_documentations (
    id SERIAL PRIMARY KEY,
    node_type VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    icon VARCHAR(500),
    version VARCHAR(20) NOT NULL,
    
    -- JSON字段存储详细说明
    inputs JSONB,
    outputs JSONB,
    parameters JSONB,
    usage TEXT,  -- Markdown格式
    examples JSONB,
    best_practices JSONB,
    tags TEXT[],
    
    author VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 节点文档版本表
CREATE TABLE node_documentation_versions (
    id SERIAL PRIMARY KEY,
    node_type VARCHAR(100) NOT NULL,
    version VARCHAR(20) NOT NULL,
    content JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(node_type, version)
);
```

#### 方案B：文件存储

**优点**:
- ✅ 易于版本控制（Git）
- ✅ 易于编写和编辑

**缺点**:
- ❌ 需要重新部署才能更新
- ❌ 难以动态管理

**实现**:
```
backend/
└── docs/
    └── nodes/
        ├── trigger/
        │   ├── webhook.md
        │   └── timer.md
        ├── ai/
        │   ├── text_generation.md
        │   └── image_generation.md
        └── ...
```

### 4.3 API设计

```go
// 获取所有节点说明
GET /api/v1/nodes/docs
Response: {
    "nodes": [
        {
            "node_type": "trigger.webhook",
            "name": "Webhook触发器",
            "category": "trigger",
            "description": "通过HTTP请求触发工作流",
            ...
        }
    ]
}

// 获取特定节点说明
GET /api/v1/nodes/docs/{node_type}
Response: {
    "node_type": "trigger.webhook",
    "name": "Webhook触发器",
    "description": "...",
    "inputs": [...],
    "outputs": [...],
    "parameters": [...],
    "usage": "# 使用说明\n...",
    "examples": [...]
}

// 获取节点说明（带版本）
GET /api/v1/nodes/docs/{node_type}?version=1.0.0
```

### 4.4 实现位置

**推荐**: 集成到节点服务中

```
backend/
└── services/
    └── node-service/
        ├── internal/
        │   ├── biz/
        │   │   ├── node.go
        │   │   └── documentation.go  # 节点文档业务逻辑
        │   └── data/
        │       └── documentation.go  # 节点文档数据访问
        └── api/
            └── node/
                └── v1/
                    └── node.proto     # 包含文档相关接口
```

## 5. 工作流使用说明设计

### 5.1 工作流说明内容

- **快速开始**: 如何创建第一个工作流
- **核心概念**: 节点、边、变量、执行等
- **设计指南**: 工作流设计最佳实践
- **示例工作流**: 常用场景的示例
- **故障排查**: 常见问题和解决方案

### 5.2 存储方案

**推荐**: 混合方案
- **静态内容**: Markdown文件（Git管理）
- **动态内容**: 数据库（工作流模板、示例）

```
docs/
└── workflows/
    ├── getting-started.md      # 快速开始
    ├── concepts.md             # 核心概念
    ├── design-guide.md         # 设计指南
    ├── examples/               # 示例
    │   ├── simple-workflow.md
    │   └── ai-workflow.md
    └── troubleshooting.md      # 故障排查
```

### 5.3 API设计

```go
// 获取工作流文档列表
GET /api/v1/workflows/docs
Response: {
    "categories": [
        {
            "name": "快速开始",
            "articles": [
                {
                    "id": "getting-started",
                    "title": "快速开始",
                    "summary": "..."
                }
            ]
        }
    ]
}

// 获取特定文档
GET /api/v1/workflows/docs/{article_id}
Response: {
    "id": "getting-started",
    "title": "快速开始",
    "content": "# 快速开始\n...",  // Markdown内容
    "updated_at": "2024-01-01T00:00:00Z"
}

// 获取工作流模板
GET /api/v1/workflows/templates
Response: {
    "templates": [
        {
            "id": "simple-ai-workflow",
            "name": "简单AI工作流",
            "description": "...",
            "workflow": {...},  // 工作流定义
            "documentation": "..."  // 说明文档
        }
    ]
}
```

### 5.4 实现位置

**推荐**: 集成到工作流服务中

```
backend/
└── services/
    └── workflow-service/
        ├── internal/
        │   └── biz/
        │       └── documentation.go  # 工作流文档业务逻辑
        └── docs/                     # 静态文档文件
            └── workflows/
```

## 6. 前端集成方案

### 6.1 文档页面设计

**位置**: 前端管理后台中

```
frontend/
└── src/
    └── views/
        ├── Documentation.vue        # 文档首页
        ├── ApiDocs.vue             # API文档页（集成Swagger UI）
        ├── NodeDocs.vue            # 节点文档页
        │   ├── NodeList.vue        # 节点列表
        │   └── NodeDetail.vue     # 节点详情
        └── WorkflowDocs.vue        # 工作流文档页
```

### 6.2 文档组件设计

```vue
<!-- NodeDocumentation.vue -->
<template>
  <div class="node-docs">
    <NodeList :nodes="nodes" @select="selectNode" />
    <NodeDetail v-if="selectedNode" :node="selectedNode" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useNodeDocs } from '@/composables/useNodeDocs'

const { nodes, loadNodeDocs } = useNodeDocs()
const selectedNode = ref(null)

const selectNode = (node) => {
  selectedNode.value = node
}

onMounted(() => {
  loadNodeDocs()
})
</script>
```

### 6.3 文档展示组件

- **Markdown渲染**: 使用 `marked` 或 `markdown-it`
- **代码高亮**: 使用 `highlight.js` 或 `prismjs`
- **示例运行**: 集成代码执行器（可选）

## 7. 独立文档站（可选）

### 7.1 使用场景

- 对外产品文档
- 开发者文档
- 完整的用户手册

### 7.2 技术选型

**方案A：VitePress（推荐）**
- 基于Vite + Vue3
- 支持Markdown
- 自动生成导航
- 支持搜索

**方案B：Docusaurus**
- React生态
- 功能完善
- 社区活跃

**方案C：MkDocs**
- Python生态
- 简单易用

### 7.3 项目结构

```
docs-site/                    # 独立文档站
├── .vitepress/
│   └── config.ts            # VitePress配置
├── guide/
│   ├── getting-started.md
│   └── ...
├── api/
│   └── reference.md         # API参考（链接到Swagger）
├── nodes/
│   └── index.md             # 节点文档索引
└── package.json
```

## 8. 推荐方案总结

### 8.1 API文档

**方案**: 后端集成Swagger UI
- **位置**: API Gateway (`/api/docs`)
- **技术**: OpenAPI 3.0 + Swagger UI
- **优势**: 独立、可测试、易维护

### 8.2 节点说明

**方案**: 数据库存储 + API提供 + 前端展示
- **存储**: PostgreSQL（支持版本管理）
- **API**: 节点服务提供
- **展示**: 前端文档页面
- **优势**: 动态更新、版本管理、多语言支持

### 8.3 工作流说明

**方案**: 文件存储 + API提供 + 前端展示
- **存储**: Markdown文件（Git管理）
- **API**: 工作流服务提供
- **展示**: 前端文档页面
- **优势**: 易于编写、版本控制

### 8.4 独立文档站

**方案**: 可选，使用VitePress
- **用途**: 对外产品文档
- **内容**: 从API和数据库同步
- **部署**: 独立部署或CDN

## 9. 实施建议

### 9.1 第一阶段：基础文档

1. **API文档**
   - 集成Swagger UI到API Gateway
   - 编写OpenAPI规范
   - 提供在线测试

2. **节点说明**
   - 设计数据库表结构
   - 实现API接口
   - 前端展示页面

### 9.2 第二阶段：完善文档

1. **工作流说明**
   - 编写Markdown文档
   - 实现API接口
   - 前端展示

2. **文档管理**
   - 文档编辑界面
   - 版本管理
   - 多语言支持

### 9.3 第三阶段：独立文档站

1. **搭建文档站**
   - 使用VitePress
   - 同步API和节点文档
   - 部署到CDN

## 10. 技术实现细节

### 10.1 OpenAPI生成（Go）

```go
// 使用swaggo/swag
//go:generate swag init -g cmd/gateway/main.go -o api/gateway/docs

// main.go
// @title StructForge API
// @version 1.0
// @description StructForge工作流平台API
// @host localhost:8000
// @BasePath /api/v1
func main() {
    // ...
}
```

### 10.2 节点文档API实现

```go
// services/node-service/internal/service/node.go
func (s *NodeService) GetNodeDocumentation(ctx context.Context, req *pb.GetNodeDocRequest) (*pb.NodeDocumentation, error) {
    doc, err := s.biz.GetNodeDocumentation(ctx, req.NodeType, req.Version)
    if err != nil {
        return nil, err
    }
    return convertToProto(doc), nil
}
```

### 10.3 前端文档组件

```typescript
// composables/useNodeDocs.ts
export function useNodeDocs() {
  const nodes = ref<NodeDoc[]>([])
  
  const loadNodeDocs = async () => {
    const response = await api.get('/api/v1/nodes/docs')
    nodes.value = response.data.nodes
  }
  
  return { nodes, loadNodeDocs }
}
```

## 11. 总结

**推荐架构**:
- ✅ **API文档**: 后端Swagger UI（独立、专业）
- ✅ **节点说明**: 数据库 + API + 前端展示（动态、灵活）
- ✅ **工作流说明**: Markdown + API + 前端展示（易维护）
- ⚪ **独立文档站**: 可选，用于对外文档

**优势**:
1. 文档与代码同步
2. 支持动态更新
3. 易于维护和扩展
4. 用户体验良好

