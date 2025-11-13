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

```
StructForge/
├── frontend/          # 前端项目（Vue3）
├── backend/           # 后端项目（Go + Kratos）
│   ├── api/           # API Gateway
│   └── services/      # 微服务
│       ├── user-service/      # 用户服务
│       ├── workflow-service/  # 工作流服务
│       ├── ai-service/        # AI模型服务
│       ├── node-service/      # 节点服务
│       ├── tool-service/      # 工具服务
│       ├── scheduler-service/ # 调度服务
│       └── log-service/       # 日志服务
├── docs/              # 项目文档
├── docker/            # Docker配置
└── scripts/           # 脚本文件
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

# 2. 启动后端服务
cd backend
make run-services

# 3. 启动前端
cd frontend
npm install
npm run dev
```

## 核心服务说明

### 用户服务 (user-service)
负责用户认证、授权、用户信息管理等功能。

### 工作流服务 (workflow-service)
工作流定义、版本管理、执行引擎核心服务。

### AI模型服务 (ai-service)
统一接入各种AI模型，提供统一的调用接口。

### 节点服务 (node-service)
实现各种工作流节点，包括触发、AI、数据处理等节点。

### 工具服务 (tool-service)
提供代码执行、文件处理等工具能力。

### 调度服务 (scheduler-service)
定时任务调度，支持Cron表达式。

### 日志服务 (log-service)
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

