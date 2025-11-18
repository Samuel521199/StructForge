# 前端目录结构说明

本文档说明前端项目的目录结构，帮助开发者快速了解项目组织方式。

## 📁 目录结构

```
frontend/
├── public/                    # 静态资源（不经过构建）
│   ├── favicon.ico
│   ├── logo.svg
│   └── robots.txt
│
├── src/
│   ├── api/                  # API接口层
│   │   ├── client/           # HTTP客户端配置
│   │   ├── services/         # API服务
│   │   └── types/            # API类型定义
│   │
│   ├── assets/               # 资源文件（经过构建）
│   │   ├── images/           # 图片资源
│   │   ├── fonts/            # 字体文件
│   │   └── styles/           # 全局样式
│   │
│   ├── components/           # 组件库
│   │   ├── common/           # 公共组件（核心）
│   │   │   ├── base/         # 基础组件
│   │   │   ├── data-display/ # 数据展示组件
│   │   │   ├── feedback/    # 反馈组件
│   │   │   ├── form/         # 表单组件
│   │   │   ├── navigation/   # 导航组件
│   │   │   ├── layout/       # 布局组件
│   │   │   └── business/     # 业务通用组件
│   │   ├── workflow/         # 工作流专用组件
│   │   ├── layout/           # 布局组件
│   │   └── business/         # 业务组件
│   │
│   ├── composables/          # Composition API组合函数
│   │   ├── core/             # 核心组合函数
│   │   ├── workflow/         # 工作流组合函数
│   │   ├── ui/               # UI相关组合函数
│   │   └── utils/            # 工具组合函数
│   │
│   ├── stores/               # Pinia状态管理
│   │   └── modules/         # 状态模块
│   │
│   ├── router/              # 路由配置
│   │   ├── routes/          # 路由定义
│   │   └── guards/          # 路由守卫
│   │
│   ├── views/               # 页面视图
│   │   ├── auth/           # 认证页面
│   │   ├── dashboard/      # 仪表盘
│   │   ├── workflow/       # 工作流页面
│   │   ├── user/           # 用户页面
│   │   ├── ai/             # AI页面
│   │   └── system/         # 系统页面
│   │
│   ├── utils/               # 工具函数
│   │   ├── validation/     # 验证工具
│   │   ├── format/         # 格式化工具
│   │   ├── storage/        # 存储工具
│   │   ├── http/           # HTTP工具
│   │   └── workflow/       # 工作流工具
│   │
│   ├── types/              # 全局类型定义
│   ├── constants/         # 常量定义
│   ├── directives/        # 自定义指令
│   ├── plugins/           # 插件
│   ├── App.vue            # 根组件
│   └── main.ts            # 入口文件
│
├── .env                    # 环境变量
├── .env.development
├── .env.production
├── index.html             # HTML模板
├── package.json
├── tsconfig.json          # TypeScript配置
├── vite.config.ts         # Vite配置
└── vitest.config.ts       # Vitest配置
```

## 📝 目录说明

### api/ - API接口层

负责与后端API通信，包括：
- `client/`: HTTP客户端配置（Axios实例、拦截器）
- `services/`: API服务（按模块划分）
- `types/`: API类型定义

### components/ - 组件库

#### common/ - 公共组件（核心）

这是整个项目的核心组件库，所有业务组件都基于这些组件构建。

**组件分类**：
- `base/`: 基础组件（Button, Input, Select等）
- `data-display/`: 数据展示组件（Table, List, Card等）
- `feedback/`: 反馈组件（Alert, Toast, Modal等）
- `form/`: 表单组件（FormField, DatePicker等）
- `navigation/`: 导航组件（Breadcrumb, Steps等）
- `layout/`: 布局组件（Container, Grid等）
- `business/`: 业务通用组件（SearchBar, FilterPanel等）

#### workflow/ - 工作流专用组件

专门用于工作流编辑器的组件：
- `editor/`: 编辑器组件（WorkflowEditor, Canvas等）
- `nodes/`: 节点组件（各种类型的节点）
- `execution/`: 执行相关组件（ExecutionMonitor等）
- `utils/`: 工作流工具组件

#### layout/ - 布局组件

应用的整体布局组件：
- `AppLayout/`: 应用主布局
- `PageLayout/`: 页面布局
- `SectionLayout/`: 区块布局
- `GridLayout/`: 网格布局

#### business/ - 业务组件

特定业务场景的组件：
- `user/`: 用户相关组件
- `ai/`: AI相关组件
- `system/`: 系统相关组件

### composables/ - 组合函数

Composition API组合函数，提供可复用的逻辑：

- `core/`: 核心组合函数（useAuth, useRequest等）
- `workflow/`: 工作流组合函数（useWorkflow, useNode等）
- `ui/`: UI相关组合函数（useModal, useToast等）
- `utils/`: 工具组合函数（useClipboard, useLocalStorage等）

### stores/ - 状态管理

Pinia状态管理，按模块划分：

- `modules/user.store.ts`: 用户状态
- `modules/auth.store.ts`: 认证状态
- `modules/workflow.store.ts`: 工作流状态
- `modules/execution.store.ts`: 执行状态
- `modules/ai.store.ts`: AI模型状态
- `modules/node.store.ts`: 节点状态
- `modules/ui.store.ts`: UI状态
- `modules/app.store.ts`: 应用状态

### router/ - 路由配置

- `routes/`: 路由定义（按模块划分）
- `guards/`: 路由守卫（认证守卫、权限守卫）

### views/ - 页面视图

所有页面组件，按功能模块划分：

- `auth/`: 认证页面（Login, Register）
- `dashboard/`: 仪表盘
- `workflow/`: 工作流页面（List, Editor, Detail等）
- `user/`: 用户页面（Profile, Settings）
- `ai/`: AI页面（ModelList, ModelConfig）
- `system/`: 系统页面（Settings, Logs）

### utils/ - 工具函数

工具函数库，按功能分类：

- `validation/`: 验证工具
- `format/`: 格式化工具
- `storage/`: 存储工具
- `http/`: HTTP工具
- `workflow/`: 工作流工具

## 🔄 组件导入规范

### 公共组件导入

```typescript
// 方式1：从分类导入
import { Button, Input, Dialog } from '@/components/common/base'
import { DataTable, Statistic } from '@/components/common/data-display'

// 方式2：从统一入口导入
import { Button, Input, DataTable } from '@/components/common'

// 方式3：从根入口导入
import { Button, Input } from '@/components'
```

### 工作流组件导入

```typescript
import { WorkflowEditor, Canvas, NodePalette } from '@/components/workflow'
```

### 布局组件导入

```typescript
import { AppLayout, PageLayout } from '@/components/layout'
```

## 📋 文件命名规范

- **组件文件**: PascalCase，如 `Button.vue`, `WorkflowEditor.vue`
- **工具文件**: camelCase，如 `date.ts`, `validators.ts`
- **类型文件**: camelCase + `.types.ts`，如 `user.types.ts`
- **常量文件**: camelCase + `.constants.ts`，如 `api.constants.ts`
- **导出文件**: `index.ts`

## 🎯 开发建议

1. **组件开发**: 优先使用公共组件，避免重复开发
2. **状态管理**: 使用Pinia Store管理状态，避免props drilling
3. **组合函数**: 将可复用逻辑提取到composables
4. **类型定义**: 所有API和组件都要有TypeScript类型定义
5. **代码组织**: 按功能模块组织代码，保持目录结构清晰

---

**最后更新**: 2024年

