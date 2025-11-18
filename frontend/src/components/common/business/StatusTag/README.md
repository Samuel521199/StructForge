# StatusTag 状态标签组件

## 功能说明

StatusTag 是一个状态标签组件，用于显示各种状态信息，支持自定义状态映射和样式。

## Props

| 参数 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| status | 状态值 | `string` | - | 是 |
| statusMap | 自定义状态映射 | `Record<string, StatusConfig>` | - | 否 |
| effect | 标签效果 | `'dark' \| 'light' \| 'plain'` | `'light'` | 否 |
| size | 标签大小 | `'large' \| 'default' \| 'small'` | `'default'` | 否 |
| round | 是否圆角 | `boolean` | `false` | 否 |

### StatusConfig 类型

```typescript
interface StatusConfig {
  label: string                                    // 显示文本
  type: 'success' | 'warning' | 'danger' | 'info' // 标签类型
  color?: string                                   // 自定义颜色
}
```

## 默认状态映射

组件内置了以下状态映射：

### 工作流状态
- `running`: 运行中（success）
- `stopped`: 已停止（info）
- `paused`: 已暂停（warning）
- `error`: 错误（danger）
- `pending`: 等待中（warning）

### 执行状态
- `success`: 成功（success）
- `failed`: 失败（danger）
- `cancelled`: 已取消（info）

### 通用状态
- `active`: 激活（success）
- `inactive`: 未激活（info）
- `enabled`: 已启用（success）
- `disabled`: 已禁用（info）

### 审核状态
- `draft`: 草稿（info）
- `reviewing`: 审核中（warning）
- `approved`: 已通过（success）
- `rejected`: 已拒绝（danger）

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 自定义标签内容 |

## 使用示例

### 基础使用

```vue
<template>
  <StatusTag status="running" />
  <StatusTag status="error" />
  <StatusTag status="pending" />
</template>

<script setup lang="ts">
import { StatusTag } from '@/components/common/business'
</script>
```

### 自定义状态映射

```vue
<template>
  <StatusTag
    :status="workflowStatus"
    :status-map="customStatusMap"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { StatusTag, type StatusConfig } from '@/components/common/business'

const workflowStatus = ref('processing')

const customStatusMap: Record<string, StatusConfig> = {
  processing: { label: '处理中', type: 'warning' },
  completed: { label: '已完成', type: 'success' },
  failed: { label: '失败', type: 'danger' },
}
</script>
```

### 自定义样式

```vue
<template>
  <StatusTag
    status="running"
    effect="dark"
    size="large"
    round
  />
</template>
```

### 自定义内容

```vue
<template>
  <StatusTag status="running">
    <el-icon><Check /></el-icon>
    运行中
  </StatusTag>
</template>
```

## 设计说明

- 使用默认状态映射可以快速显示常见状态
- 支持通过 `statusMap` 自定义状态映射
- 支持通过插槽自定义标签内容
- 基于 Element Plus 的 `el-tag` 组件实现

