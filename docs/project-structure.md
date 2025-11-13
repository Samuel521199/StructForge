# StructForge 项目结构设计

## 1. 项目根目录结构

```
StructForge/
├── frontend/                 # 前端项目
├── backend/                  # 后端项目
├── docs/                     # 项目文档
├── docs-site/               # 独立文档站（可选）
│   ├── .vitepress/
│   ├── guide/
│   ├── api/
│   └── nodes/
├── scripts/                  # 脚本文件
├── docker/                   # Docker相关配置
├── .gitignore
├── README.md
└── docker-compose.yml        # 本地开发环境配置
```

## 2. 前端项目结构 (frontend/)

```
frontend/
├── public/                   # 静态资源
├── src/
│   ├── api/                  # API接口定义
│   │   ├── user.ts
│   │   ├── workflow.ts
│   │   ├── ai.ts
│   │   └── index.ts
│   ├── assets/               # 资源文件
│   │   ├── images/
│   │   ├── icons/
│   │   └── styles/
│   ├── components/           # 公共组件
│   │   ├── common/           # 通用组件
│   │   ├── workflow/         # 工作流相关组件
│   │   │   ├── WorkflowEditor.vue
│   │   │   ├── NodePalette.vue
│   │   │   ├── NodeProperties.vue
│   │   │   └── ExecutionMonitor.vue
│   │   └── layout/           # 布局组件
│   ├── composables/          # Composition API组合函数
│   │   ├── useAuth.ts
│   │   ├── useWorkflow.ts
│   │   └── useWebSocket.ts
│   ├── router/               # 路由配置
│   │   └── index.ts
│   ├── stores/               # Pinia状态管理
│   │   ├── user.ts
│   │   ├── workflow.ts
│   │   ├── ai.ts
│   │   └── index.ts
│   ├── utils/                # 工具函数
│   │   ├── request.ts        # HTTP请求封装
│   │   ├── websocket.ts      # WebSocket封装
│   │   └── common.ts
│   ├── views/                # 页面视图
│   │   ├── Login.vue
│   │   ├── Dashboard.vue
│   │   ├── WorkflowList.vue
│   │   ├── WorkflowEditor.vue
│   │   ├── WorkflowExecution.vue
│   │   ├── AIModels.vue
│   │   ├── Settings.vue
│   │   └── Documentation.vue  # 文档页面
│   │       ├── ApiDocs.vue    # API文档
│   │       ├── NodeDocs.vue   # 节点文档
│   │       └── WorkflowDocs.vue  # 工作流文档
│   ├── App.vue
│   └── main.ts
├── .env                      # 环境变量
├── .env.development
├── .env.production
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## 3. 后端项目结构 (backend/)

```
backend/
├── api/                      # API Gateway (Kratos)
│   ├── gateway/              # HTTP Gateway配置
│   │   ├── docs/             # API文档
│   │   │   ├── swagger.yaml  # OpenAPI规范
│   │   │   └── swagger.json
│   │   ├── handlers/         # 处理器
│   │   │   ├── docs.go       # 文档路由
│   │   │   └── ...
│   │   ├── middleware/       # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logger.go
│   │   └── routes.go
│   ├── proto/                # Protobuf定义
│   │   ├── user/
│   │   ├── workflow/
│   │   ├── ai/
│   │   └── common/
│   └── cmd/
│       └── gateway/          # Gateway启动入口
│           └── main.go
│
├── services/                 # 微服务目录
│   ├── user-service/         # 用户服务
│   │   ├── internal/        # 内部实现
│   │   │   ├── biz/         # 业务逻辑层
│   │   │   │   ├── user.go
│   │   │   │   └── auth.go
│   │   │   ├── data/        # 数据访问层
│   │   │   │   ├── user.go
│   │   │   │   └── db.go
│   │   │   ├── service/     # 服务层 (gRPC)
│   │   │   │   └── user.go
│   │   │   └── conf/        # 配置
│   │   │       └── config.yaml
│   │   ├── api/             # API定义
│   │   │   └── user/
│   │   │       └── v1/
│   │   │           └── user.proto
│   │   └── cmd/             # 启动入口
│   │       └── main.go
│   │
│   ├── workflow-service/     # 工作流服务
│   │   ├── internal/
│   │   │   ├── biz/
│   │   │   │   ├── workflow.go
│   │   │   │   ├── execution.go
│   │   │   │   ├── engine.go    # 执行引擎
│   │   │   │   └── documentation.go  # 工作流文档业务逻辑
│   │   │   ├── data/
│   │   │   │   ├── workflow.go
│   │   │   │   └── execution.go
│   │   │   └── service/
│   │   │       └── workflow.go
│   │   └── docs/             # 工作流文档（Markdown）
│   │       └── workflows/
│   │   ├── api/
│   │   │   └── workflow/
│   │   │       └── v1/
│   │   │           └── workflow.proto
│   │   └── cmd/
│   │       └── main.go
│   │
│   ├── ai-service/           # AI模型服务
│   │   ├── internal/
│   │   │   ├── biz/
│   │   │   │   ├── model.go
│   │   │   │   └── adapter.go
│   │   │   ├── adapters/    # 模型适配器
│   │   │   │   ├── openai/
│   │   │   │   │   └── adapter.go
│   │   │   │   ├── gemini/
│   │   │   │   │   └── adapter.go
│   │   │   │   ├── ollama/
│   │   │   │   │   └── adapter.go
│   │   │   │   └── custom/
│   │   │   │       └── adapter.go
│   │   │   ├── data/
│   │   │   │   └── model.go
│   │   │   └── service/
│   │   │       └── ai.go
│   │   ├── api/
│   │   │   └── ai/
│   │   │       └── v1/
│   │   │           └── ai.proto
│   │   └── cmd/
│   │       └── main.go
│   │
│   ├── node-service/         # 节点服务
│   │   ├── internal/
│   │   │   ├── biz/
│   │   │   │   ├── node.go
│   │   │   │   ├── executor.go
│   │   │   │   └── documentation.go  # 节点文档业务逻辑
│   │   │   ├── nodes/        # 节点实现
│   │   │   │   ├── trigger/  # 触发节点
│   │   │   │   │   ├── webhook.go
│   │   │   │   │   ├── timer.go
│   │   │   │   │   └── manual.go
│   │   │   │   ├── ai/       # AI节点
│   │   │   │   │   ├── text_generation.go
│   │   │   │   │   ├── image_generation.go
│   │   │   │   │   └── code_generation.go
│   │   │   │   ├── data/     # 数据处理节点
│   │   │   │   │   ├── transform.go
│   │   │   │   │   └── filter.go
│   │   │   │   ├── integration/ # 集成节点
│   │   │   │   │   ├── http.go
│   │   │   │   │   └── database.go
│   │   │   │   ├── control/  # 控制节点
│   │   │   │   │   ├── condition.go
│   │   │   │   │   └── loop.go
│   │   │   │   └── tool/     # 工具节点
│   │   │   │       └── script.go
│   │   │   ├── data/
│   │   │   │   ├── node.go
│   │   │   │   └── documentation.go  # 节点文档数据访问
│   │   │   └── service/
│   │   │       └── node.go
│   │   ├── api/
│   │   │   └── node/
│   │   │       └── v1/
│   │   │           └── node.proto
│   │   └── cmd/
│   │       └── main.go
│   │
│   ├── tool-service/         # 工具服务
│   │   ├── internal/
│   │   │   ├── biz/
│   │   │   │   └── tool.go
│   │   │   ├── tools/        # 工具实现
│   │   │   │   ├── code_executor.go
│   │   │   │   ├── file_handler.go
│   │   │   │   └── network.go
│   │   │   └── service/
│   │   │       └── tool.go
│   │   ├── api/
│   │   │   └── tool/
│   │   │       └── v1/
│   │   │           └── tool.proto
│   │   └── cmd/
│   │       └── main.go
│   │
│   ├── scheduler-service/    # 调度服务
│   │   ├── internal/
│   │   │   ├── biz/
│   │   │   │   ├── scheduler.go
│   │   │   │   └── job.go
│   │   │   └── service/
│   │   │       └── scheduler.go
│   │   ├── api/
│   │   │   └── scheduler/
│   │   │       └── v1/
│   │   │           └── scheduler.proto
│   │   └── cmd/
│   │       └── main.go
│   │
│   └── log-service/          # 日志服务
│       ├── internal/
│       │   ├── biz/
│       │   │   └── log.go
│       │   ├── data/
│       │   │   └── log.go
│       │   └── service/
│       │       └── log.go
│       ├── api/
│       │   └── log/
│       │       └── v1/
│       │           └── log.proto
│       └── cmd/
│           └── main.go
│
├── pkg/                      # 公共包
│   ├── errors/               # 错误定义
│   ├── logger/               # 日志工具
│   ├── database/             # 数据库工具
│   ├── cache/                # 缓存工具
│   ├── validator/            # 验证器
│   └── utils/                # 通用工具
│
├── configs/                  # 配置文件
│   ├── config.yaml           # 主配置
│   └── services/             # 各服务配置
│
├── deployments/              # 部署配置
│   ├── docker/               # Dockerfile
│   └── k8s/                  # Kubernetes配置（可选）
│
├── scripts/                  # 脚本
│   ├── init-db.sh            # 数据库初始化
│   └── generate-proto.sh     # 生成proto代码
│
├── go.mod
├── go.sum
└── Makefile
```

## 4. 数据库设计

### 4.1 数据库初始化脚本位置

```
backend/
└── migrations/               # 数据库迁移文件
    ├── 001_init_users.sql
    ├── 002_init_workflows.sql
    ├── 003_init_ai_models.sql
    └── ...
```

## 5. Docker配置结构

```
docker/
├── postgres/
│   └── init.sql              # 数据库初始化脚本
├── redis/
│   └── redis.conf
└── nginx/                    # 前端Nginx配置（生产环境）
    └── nginx.conf
```

## 6. 文档结构

```
docs/
├── architecture.md           # 架构设计文档
├── project-structure.md      # 项目结构文档（本文件）
├── api/                      # API文档
│   ├── user-api.md
│   ├── workflow-api.md
│   └── ai-api.md
├── development/              # 开发文档
│   ├── setup.md              # 环境搭建
│   ├── coding-standards.md   # 编码规范
│   └── testing.md            # 测试指南
└── deployment/               # 部署文档
    ├── local.md              # 本地部署
    └── production.md         # 生产部署
```

## 7. 关键目录说明

### 7.1 前端关键目录

- **src/components/workflow/**: 工作流编辑器核心组件
- **src/stores/**: 全局状态管理
- **src/api/**: 所有API接口定义

### 7.2 后端关键目录

- **services/*/internal/biz/**: 业务逻辑层，核心业务代码
- **services/*/internal/data/**: 数据访问层，数据库操作
- **services/*/internal/service/**: gRPC服务实现
- **services/ai-service/internal/adapters/**: AI模型适配器
- **services/node-service/internal/nodes/**: 各种节点实现
- **pkg/**: 跨服务共享的公共代码

## 8. 配置文件说明

### 8.1 前端配置

- `.env`: 基础环境变量
- `.env.development`: 开发环境配置
- `.env.production`: 生产环境配置

### 8.2 后端配置

- `configs/config.yaml`: 全局配置
- `services/*/internal/conf/config.yaml`: 各服务独立配置

## 9. 开发规范

### 9.1 命名规范

- **服务名**: kebab-case (user-service, workflow-service)
- **包名**: 小写字母，无下划线
- **文件名**: snake_case (Go) / camelCase (TypeScript)
- **结构体/类名**: PascalCase
- **函数/方法名**: camelCase

### 9.2 代码组织

- 遵循Kratos框架的目录结构
- 业务逻辑放在biz层
- 数据访问放在data层
- 服务接口放在service层

## 10. 文档系统结构

### 10.1 API文档
- **位置**: `backend/api/gateway/docs/`
- **格式**: OpenAPI 3.0 (YAML/JSON)
- **访问**: `/api/docs` (Swagger UI)

### 10.2 节点文档
- **存储**: PostgreSQL数据库
- **API**: 节点服务提供 (`/api/v1/nodes/docs`)
- **展示**: 前端文档页面

### 10.3 工作流文档
- **存储**: Markdown文件 (`backend/services/workflow-service/docs/`)
- **API**: 工作流服务提供 (`/api/v1/workflows/docs`)
- **展示**: 前端文档页面

### 10.4 独立文档站（可选）
- **位置**: `docs-site/`
- **技术**: VitePress
- **用途**: 对外产品文档

详细设计请参考 [文档系统设计文档](../docs/documentation-system.md)

## 11. 下一步工作

1. 初始化项目结构
2. 配置开发环境
3. 搭建基础框架
4. 实现核心功能模块
5. 集成文档系统

