# StructForge - AI工作流编排平台

## 项目简介

StructForge 是一个全功能的AI工作流编排平台，支持多种AI大模型（开源、自建、公共API）和各种工作节点，提供可视化的流程设计界面。类似n8n，但专注于AI工作流场景。

## 核心特性

- 🤖 **多AI模型支持**: 支持OpenAI、Gemini、Ollama等各类AI模型
- 🔧 **丰富的工作节点**: 触发、AI、数据处理、集成、控制、工具等节点
- 🎨 **可视化编辑器**: 拖拽式工作流设计，直观易用
- ⚡ **高性能执行**: 支持同步/异步执行，适合各种场景
- 🔐 **安全可靠**: 完善的认证授权、数据加密、多租户隔离
- 🚀 **易于扩展**: 插件化架构，支持自定义节点和模型

## 技术架构

### 前端
- **框架**: Vue 3.x + TypeScript
- **构建工具**: Vite
- **状态管理**: Pinia
- **UI组件**: Element Plus
- **工作流编辑器**: Vue Flow

### 后端
- **框架**: Go + Kratos微服务框架
- **协议**: gRPC + HTTP Gateway
- **数据库**: PostgreSQL + Redis
- **消息队列**: RabbitMQ（可选）

### 部署
- **容器化**: Docker + Docker Compose
- **平台支持**: Windows、Linux

## 项目结构

```text
StructForge/
├── frontend/              # 前端项目（Vue3）
│   ├── src/
│   │   ├── api/          # API接口定义
│   │   ├── components/   # 组件
│   │   ├── views/        # 页面视图
│   │   ├── stores/       # Pinia状态管理
│   │   └── utils/        # 工具函数
│   └── package.json
│
├── backend/               # 后端项目（Go + Kratos）
│   ├── api/              # API定义（Protobuf）
│   │   ├── user/
│   │   │   └── v1/
│   │   │       └── user.proto
│   │   ├── workflow/
│   │   │   └── v1/
│   │   │       └── workflow.proto
│   │   ├── ai/
│   │   │   └── v1/
│   │   │       └── ai.proto
│   │   ├── node/
│   │   │   └── v1/
│   │   │       └── node.proto
│   │   ├── gateway/
│   │   │   └── v1/
│   │   │       └── gateway.proto
│   │   └── common/
│   │       └── v1/
│   │           └── common.proto
│   │
│   ├── apps/             # 微服务应用
│   │   ├── gateway/      # API Gateway服务
│   │   │   ├── cmd/
│   │   │   │   └── gateway/
│   │   │   │       ├── main.go
│   │   │   │       └── wire.go
│   │   │   └── internal/
│   │   │       ├── conf/     # 配置
│   │   │       ├── handler/ # HTTP处理器
│   │   │       ├── data/     # 数据访问（gRPC客户端）
│   │   │       └── server/   # 服务器配置
│   │   │
│   │   ├── user/         # 用户服务
│   │   │   ├── cmd/
│   │   │   │   └── user/
│   │   │   │       ├── main.go
│   │   │   │       └── wire.go
│   │   │   └── internal/
│   │   │       ├── biz/      # 业务逻辑层
│   │   │       ├── data/     # 数据访问层
│   │   │       ├── service/  # gRPC服务层
│   │   │       ├── server/   # 服务器配置
│   │   │       └── conf/     # 配置
│   │   │
│   │   ├── workflow/     # 工作流服务
│   │   │   ├── cmd/
│   │   │   │   └── workflow/
│   │   │   │       ├── main.go
│   │   │   │       └── wire.go
│   │   │   └── internal/
│   │   │       ├── biz/      # 业务逻辑层（执行引擎）
│   │   │       ├── data/     # 数据访问层
│   │   │       ├── service/  # gRPC服务层
│   │   │       ├── server/   # 服务器配置
│   │   │       └── conf/     # 配置
│   │   │
│   │   ├── ai/           # AI模型服务
│   │   │   ├── cmd/
│   │   │   │   └── ai/
│   │   │   │       ├── main.go
│   │   │   │       └── wire.go
│   │   │   └── internal/
│   │   │       ├── biz/      # 业务逻辑层
│   │   │       ├── adapters/ # AI模型适配器
│   │   │       │   ├── openai/
│   │   │       │   ├── gemini/
│   │   │       │   └── ollama/
│   │   │       ├── data/     # 数据访问层
│   │   │       ├── service/  # gRPC服务层
│   │   │       └── conf/     # 配置
│   │   │
│   │   ├── node/         # 节点服务
│   │   │   ├── cmd/
│   │   │   │   └── node/
│   │   │   │       ├── main.go
│   │   │   │       └── wire.go
│   │   │   └── internal/
│   │   │       ├── biz/      # 业务逻辑层
│   │   │       ├── nodes/    # 节点实现
│   │   │       │   ├── trigger/
│   │   │       │   ├── ai/
│   │   │       │   ├── data/
│   │   │       │   └── control/
│   │   │       ├── data/     # 数据访问层
│   │   │       ├── service/  # gRPC服务层
│   │   │       └── conf/     # 配置
│   │   │
│   │   ├── tool/         # 工具服务
│   │   ├── scheduler/    # 调度服务
│   │   └── log/          # 日志服务
│   │
│   ├── common/           # 公共代码
│   │   ├── middleware/  # 中间件
│   │   │   ├── auth/     # 认证中间件
│   │   │   ├── cors/     # CORS中间件
│   │   │   ├── logger/    # 日志中间件
│   │   │   └── ratelimit/ # 限流中间件
│   │   ├── data/         # 数据访问公共代码
│   │   │   ├── database/ # 数据库连接
│   │   │   └── redis/    # Redis连接
│   │   ├── utils/        # 工具函数
│   │   │   ├── crypto.go  # 加密工具
│   │   │   ├── response.go # 响应工具
│   │   │   └── validator.go # 验证工具
│   │   └── log/          # 日志工具
│   │
│   ├── configs/          # 配置文件
│   │   ├── local/        # 本地环境配置
│   │   │   ├── gateway.yaml
│   │   │   ├── user.yaml
│   │   │   ├── workflow.yaml
│   │   │   └── ...
│   │   └── test/         # 测试环境配置
│   │
│   ├── deploy/           # 部署相关
│   │   ├── configs/      # 部署配置
│   │   └── scripts/      # 部署脚本
│   │
│   ├── script/           # 开发脚本
│   │   ├── proto.sh      # 生成proto代码
│   │   ├── wire.sh       # 生成wire代码
│   │   └── build.sh      # 构建脚本
│   │
│   ├── third_party/      # 第三方依赖
│   │   ├── google/       # Google API
│   │   └── validate/    # 验证规则
│   │
│   ├── go.mod            # Go模块定义
│   └── Makefile          # Make命令
│
├── docs/                 # 项目文档
├── docker/               # Docker配置
├── docker-compose.yml    # Docker Compose配置
└── README.md             # 项目说明
```

## 文档导航

- [架构设计文档](docs/architecture.md) - 整体架构设计和技术选型
- [项目结构文档](docs/project-structure.md) - 详细的目录结构说明
- [设计分析文档](docs/design-analysis.md) - 核心模块设计和实现分析
- [文档系统设计](docs/documentation-system.md) - API文档和节点/工作流使用说明系统设计

## 快速开始

### 环境要求

- Docker Desktop
- Node.js 18+
- Go 1.21+

### 本地开发

```bash
# 1. 启动基础设施
docker-compose up -d postgres redis

# 2. 生成proto代码和wire代码
cd backend
make proto
make wire

# 3. 启动后端服务
# 方式1: 使用Makefile启动所有服务
make run

# 方式2: 单独启动某个服务
cd apps/gateway && go run cmd/gateway/main.go
cd apps/user && go run cmd/user/main.go
# ... 其他服务类似

# 4. 启动前端
cd frontend
npm install
npm run dev
```

### 服务说明

- **gateway**: API Gateway，统一入口，默认端口 `8000`
- **user**: 用户服务，gRPC端口 `9001`
- **workflow**: 工作流服务，gRPC端口 `9002`
- **ai**: AI模型服务，gRPC端口 `9003`
- **node**: 节点服务，gRPC端口 `9004`
- **tool**: 工具服务，gRPC端口 `9005`
- **scheduler**: 调度服务，gRPC端口 `9006`
- **log**: 日志服务，gRPC端口 `9007`

## 核心服务说明

### Gateway服务 (apps/gateway)
API Gateway，提供统一的HTTP入口，负责路由、认证、限流等功能。

### 用户服务 (apps/user)
负责用户认证、授权、用户信息管理等功能。

### 工作流服务 (apps/workflow)
工作流定义、版本管理、执行引擎核心服务。

### AI模型服务 (apps/ai)
统一接入各种AI模型，提供统一的调用接口。

### 节点服务 (apps/node)
实现各种工作流节点，包括触发、AI、数据处理等节点。

### 工具服务 (apps/tool)
提供代码执行、文件处理等工具能力。

### 调度服务 (apps/scheduler)
定时任务调度，支持Cron表达式。

### 日志服务 (apps/log)
工作流执行日志收集、查询和分析。

## 文档系统

### API文档
- **技术**: OpenAPI 3.0 + Swagger UI
- **位置**: 集成到API Gateway
- **访问**: `http://localhost:8000/api/docs`

### 节点使用说明
- **存储**: PostgreSQL数据库
- **API**: 节点服务提供 (`/api/v1/nodes/docs`)
- **展示**: 前端文档页面
- **特点**: 动态更新、版本管理、支持多语言

### 工作流使用说明
- **存储**: Markdown文件
- **API**: 工作流服务提供 (`/api/v1/workflows/docs`)
- **展示**: 前端文档页面
- **内容**: 快速开始、设计指南、示例、故障排查

详细设计请参考 [文档系统设计文档](docs/documentation-system.md)

## 开发计划

- [ ] 项目初始化
- [ ] 基础框架搭建
- [ ] 用户服务实现
- [ ] 工作流服务实现
- [ ] AI模型服务实现
- [ ] 节点服务实现
- [ ] 前端编辑器实现
- [ ] 测试和优化

## 许可证

待定

## 贡献

欢迎提交Issue和Pull Request！

