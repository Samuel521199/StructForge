# StructForge 前端项目

## 项目简介

StructForge 工作流平台前端项目，基于 Vue 3 + TypeScript + Vite 构建。

## 技术栈

- **框架**: Vue 3.x (Composition API)
- **语言**: TypeScript
- **构建工具**: Vite
- **状态管理**: Pinia
- **UI组件库**: Element Plus
- **工作流编辑器**: Vue Flow
- **路由**: Vue Router
- **HTTP客户端**: Axios

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API接口层
│   ├── assets/           # 资源文件
│   ├── components/       # 组件库
│   ├── composables/      # 组合函数
│   ├── stores/           # 状态管理
│   ├── router/           # 路由配置
│   ├── views/            # 页面视图
│   ├── utils/            # 工具函数
│   ├── types/            # 类型定义
│   └── constants/        # 常量定义
```

## 开发指南

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

### 构建生产版本

```bash
npm run build
```

### 代码检查

```bash
npm run lint
```

### 运行测试

```bash
npm run test
```

## 目录说明

详细的目录结构说明请参考 [FRONTEND_ARCHITECTURE.md](./FRONTEND_ARCHITECTURE.md)

## 组件库

公共组件库设计请参考 [COMPONENT_LIBRARY_DESIGN.md](./COMPONENT_LIBRARY_DESIGN.md)

---

**最后更新**: 2024年

