# ActionBar 操作栏组件

## 功能说明

ActionBar 是一个操作按钮组组件，用于统一管理页面操作按钮，支持灵活的对齐方式和按钮配置。

## Props

| 参数 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| actions | 操作项列表 | `ActionItem[]` | - | 是 |
| align | 对齐方式 | `'left' \| 'center' \| 'right'` | `'right'` | 否 |
| size | 按钮大小 | `'large' \| 'default' \| 'small'` | `'default'` | 否 |

### ActionItem 类型

```typescript
interface ActionItem {
  label: string                                    // 按钮文本
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'text'
  size?: 'large' | 'default' | 'small'
  disabled?: boolean                               // 是否禁用
  loading?: boolean                                // 是否加载中
  icon?: Component | string                        // 图标
  round?: boolean                                 // 是否圆角
  plain?: boolean                                 // 是否朴素按钮
  hidden?: boolean                                // 是否隐藏
  onClick?: () => void                            // 点击回调
}
```

## Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| action | 操作点击事件 | `(action: ActionItem, index: number) => void` |

## 使用示例

### 基础使用

```vue
<template>
  <ActionBar :actions="actions" @action="handleAction" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ActionBar, type ActionItem } from '@/components/common/business'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'

const actions: ActionItem[] = [
  {
    label: '新增',
    type: 'primary',
    icon: Plus,
    onClick: () => {
      console.log('新增')
    },
  },
  {
    label: '编辑',
    type: 'default',
    icon: Edit,
  },
  {
    label: '删除',
    type: 'danger',
    icon: Delete,
    disabled: true,
  },
]

const handleAction = (action: ActionItem, index: number) => {
  console.log('操作:', action.label, index)
}
</script>
```

### 不同对齐方式

```vue
<template>
  <!-- 左对齐 -->
  <ActionBar :actions="actions" align="left" />
  
  <!-- 居中 -->
  <ActionBar :actions="actions" align="center" />
  
  <!-- 右对齐（默认） -->
  <ActionBar :actions="actions" align="right" />
</template>
```

### 动态控制按钮状态

```vue
<template>
  <ActionBar :actions="actions" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ActionBar, type ActionItem } from '@/components/common/business'

const selectedCount = ref(0)
const isSubmitting = ref(false)

const actions = computed<ActionItem[]>(() => [
  {
    label: '提交',
    type: 'primary',
    loading: isSubmitting.value,
    disabled: selectedCount.value === 0,
    onClick: async () => {
      isSubmitting.value = true
      // 执行提交逻辑
      await new Promise(resolve => setTimeout(resolve, 1000))
      isSubmitting.value = false
    },
  },
  {
    label: '取消',
    type: 'default',
    disabled: isSubmitting.value,
  },
])
</script>
```

### 条件显示按钮

```vue
<template>
  <ActionBar :actions="actions" />
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ActionBar, type ActionItem } from '@/components/common/business'

const userRole = ref('admin')

const actions = computed<ActionItem[]>(() => [
  {
    label: '查看',
    type: 'default',
  },
  {
    label: '编辑',
    type: 'primary',
    hidden: userRole.value !== 'admin', // 非管理员隐藏
  },
  {
    label: '删除',
    type: 'danger',
    hidden: userRole.value !== 'admin', // 非管理员隐藏
  },
])
</script>
```

## 设计说明

- 支持通过 `onClick` 属性或 `action` 事件处理点击
- 支持动态控制按钮的显示、禁用、加载状态
- 支持通过 `hidden` 属性条件显示按钮
- 按钮样式和行为完全基于 Element Plus 的 Button 组件

