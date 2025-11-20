# 通用组件库使用指南

## 📋 目录

- [概述](#概述)
- [组件分类](#组件分类)
- [使用规范](#使用规范)
- [组件列表](#组件列表)
- [最佳实践](#最佳实践)
- [编辑器化支持](#编辑器化支持)

## 概述

通用组件库（`@/components/common/base`）是对 Element Plus 组件的统一封装，提供：

- ✅ **统一的 API**：所有组件使用一致的接口设计
- ✅ **易于维护**：集中管理组件逻辑和样式
- ✅ **便于扩展**：可以轻松添加自定义功能
- ✅ **编辑器化支持**：为未来的可视化编辑器做准备

## 组件分类

### 基础组件 (`base/`)

封装 Element Plus 的基础 UI 组件，提供统一的接口和样式。

- **表单组件**：Form, FormItem, Input, Select, Checkbox, Button
- **展示组件**：Card, Table, Empty, Loading, Icon
- **反馈组件**：Dialog, Message
- **导航组件**：Link

### 业务组件 (`business/`)

基于基础组件构建的业务相关组件。

- ActionBar, SearchBar, StatusTag, TimeAgo 等

### 效果组件 (`effects/`)

特殊视觉效果组件。

- CodeRain（代码雨效果）

## 使用规范

### 1. 统一导入

```typescript
// ✅ 推荐：从统一入口导入
import { 
  Button, 
  Input, 
  Form, 
  FormItem, 
  Checkbox, 
  Link, 
  Icon 
} from '@/components/common/base'

// ❌ 禁止：直接使用 Element Plus
import { ElButton, ElInput } from 'element-plus'
```

### 2. 组件命名

- 使用 PascalCase：`<Button>`, `<Input>`, `<Form>`
- 不使用 `el-` 前缀：`<Button>` 而不是 `<el-button>`

### 3. 图标使用

```vue
<template>
  <!-- ✅ 推荐：使用 Icon 组件 -->
  <Icon :icon="User" :size="20" color="#00FF00" />
  
  <!-- ❌ 禁止：直接使用 el-icon -->
  <el-icon><User /></el-icon>
</template>

<script setup lang="ts">
import { Icon } from '@/components/common/base'
import { User } from '@element-plus/icons-vue'
</script>
```

### 4. 消息提示

```typescript
// ✅ 推荐：使用封装的 Message 方法
import { success, error, warning, info } from '@/components/common/base/Message'

success('操作成功')
error('操作失败')
warning('警告信息')
info('提示信息')

// ❌ 禁止：直接使用 ElMessage
import { ElMessage } from 'element-plus'
ElMessage.success('操作成功')
```

## 组件列表

### 已实现组件

| 组件 | 说明 | 文档 |
|------|------|------|
| Button | 按钮 | [Button/README.md](./Button/README.md) |
| Input | 输入框 | [Input/README.md](./Input/README.md) |
| Form | 表单 | [Form/README.md](./Form/README.md) |
| FormItem | 表单项 | [FormItem/README.md](../form/FormItem/README.md) |
| Select | 选择器 | [Select/README.md](./Select/README.md) |
| Checkbox | 复选框 | [Checkbox/README.md](./Checkbox/README.md) |
| Link | 链接 | [Link/README.md](./Link/README.md) |
| Icon | 图标 | [Icon/README.md](./Icon/README.md) |
| Card | 卡片 | [Card/README.md](./Card/README.md) |
| Table | 表格 | [Table/README.md](./Table/README.md) |
| Dialog | 对话框 | [Dialog/README.md](./Dialog/README.md) |
| Loading | 加载中 | [Loading/README.md](./Loading/README.md) |
| Empty | 空状态 | [Empty/README.md](./Empty/README.md) |
| Message | 消息提示 | [Message/README.md](./Message/README.md) |

### 待实现组件

以下组件目录已创建，但尚未实现：

- Badge, Tag, Tooltip, Popover
- Dropdown, Menu, Tabs
- Pagination, Notification

## 最佳实践

### 1. 表单页面

```vue
<template>
  <Form ref="formRef" :model="form" :rules="rules">
    <FormItem label="用户名" prop="username">
      <Input v-model="form.username" placeholder="请输入用户名" />
    </FormItem>
    
    <FormItem>
      <Button type="primary" @click="handleSubmit">提交</Button>
    </FormItem>
  </Form>
</template>

<script setup lang="ts">
import { Form, FormItem, Input, Button } from '@/components/common/base'
</script>
```

### 2. 列表页面

```vue
<template>
  <Card>
    <Loading :loading="loading">
      <Table :data="tableData" :columns="columns">
        <template #empty>
          <Empty description="暂无数据" />
        </template>
      </Table>
    </Loading>
  </Card>
</template>

<script setup lang="ts">
import { Card, Loading, Table, Empty } from '@/components/common/base'
</script>
```

### 3. 图标使用

```vue
<template>
  <Input :prefix-icon="User" />
  <Icon :icon="Loading" :is-loading="true" :size="24" />
</template>

<script setup lang="ts">
import { Input, Icon } from '@/components/common/base'
import { User, Loading } from '@element-plus/icons-vue'
</script>
```

## 编辑器化支持

为了支持未来的可视化编辑器，所有组件都遵循以下规范：

### 1. Props 类型定义

每个组件都有完整的 TypeScript 类型定义：

```typescript
// Button/types.ts
export interface ButtonProps {
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'text'
  size?: 'large' | 'default' | 'small'
  disabled?: boolean
  loading?: boolean
  // ...
}
```

### 2. 组件元数据

每个组件目录包含：

- `README.md`：组件文档和使用示例
- `types.ts`：TypeScript 类型定义
- `index.ts`：统一导出

### 3. 统一的组件结构

```typescript
// 所有组件都遵循相同的结构
export interface ComponentProps {
  // Props 定义
}

export interface ComponentEmits {
  // Events 定义
}

// 组件实现
export default defineComponent<ComponentProps, ComponentEmits>({
  // ...
})
```

### 4. 编辑器配置支持

未来可以基于这些类型定义自动生成编辑器配置：

```json
{
  "component": "Button",
  "props": {
    "type": {
      "type": "select",
      "options": ["primary", "success", "warning", "danger", "info", "text"],
      "default": "default"
    },
    "size": {
      "type": "select",
      "options": ["large", "default", "small"],
      "default": "default"
    }
  }
}
```

## 迁移指南

### 从 Element Plus 迁移到通用组件

#### 1. 更新导入

```typescript
// 之前
import { ElButton, ElInput } from 'element-plus'

// 之后
import { Button, Input } from '@/components/common/base'
```

#### 2. 更新模板

```vue
<!-- 之前 -->
<el-button type="primary">提交</el-button>
<el-input v-model="value" />

<!-- 之后 -->
<Button type="primary">提交</Button>
<Input v-model="value" />
```

#### 3. 更新图标

```vue
<!-- 之前 -->
<el-icon><User /></el-icon>

<!-- 之后 -->
<Icon :icon="User" />
```

#### 4. 更新消息提示

```typescript
// 之前
import { ElMessage } from 'element-plus'
ElMessage.success('操作成功')

// 之后
import { success } from '@/components/common/base/Message'
success('操作成功')
```

## 开发规范

### 1. 新增组件

创建新组件时，请遵循以下结构：

```
ComponentName/
├── ComponentName.vue    # 组件实现
├── types.ts             # 类型定义
├── index.ts             # 导出文件
└── README.md            # 组件文档
```

### 2. 组件实现规范

- 使用 `defineProps` 和 `defineEmits` 定义 Props 和 Events
- 提供完整的 TypeScript 类型定义
- 支持 `v-bind="$attrs"` 传递原生属性
- 提供合理的默认值

### 3. 文档要求

每个组件必须包含：

- Props 说明
- Events 说明
- 使用示例
- 最佳实践

## 相关文档

- [使用示例](./USAGE_EXAMPLES.md)
- [前端架构设计](../../FRONTEND_ARCHITECTURE.md)
- [组件库设计文档](./COMPONENT_LIBRARY_DESIGN.md)

---

**最后更新**: 2024年

