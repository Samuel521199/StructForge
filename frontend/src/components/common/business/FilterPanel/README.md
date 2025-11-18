# FilterPanel 筛选面板组件

## 功能说明

FilterPanel 是一个独立的筛选面板组件，支持多种筛选类型，可折叠显示，适用于复杂的筛选场景。

## Props

| 参数 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| filters | 筛选配置列表 | `FilterOption[]` | - | 是 |
| modelValue | 筛选值（v-model） | `Record<string, any>` | - | 是 |
| collapsible | 是否可折叠 | `boolean` | `true` | 否 |

### FilterOption 类型

```typescript
interface FilterOption {
  label: string                                    // 筛选标签
  prop: string                                     // 筛选字段名
  type?: 'select' | 'input' | 'date' | 'daterange' // 筛选类型
  options?: Array<{ label: string; value: any }>  // 选项数据（select类型）
  placeholder?: string                             // 占位符
  required?: boolean                                // 是否必填
  multiple?: boolean                               // 是否多选（select类型）
}
```

## Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 筛选值更新 | `(value: Record<string, any>) => void` |
| change | 筛选值变化 | `(value: Record<string, any>) => void` |
| reset | 重置筛选 | `() => void` |
| apply | 应用筛选 | `(value: Record<string, any>) => void` |

## 使用示例

### 基础使用

```vue
<template>
  <FilterPanel
    :filters="filters"
    v-model="filterValues"
    @apply="handleApply"
    @reset="handleReset"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { FilterPanel, type FilterOption } from '@/components/common/business'

const filterValues = ref<Record<string, any>>({})

const filters: FilterOption[] = [
  {
    label: '状态',
    prop: 'status',
    type: 'select',
    options: [
      { label: '运行中', value: 'running' },
      { label: '已停止', value: 'stopped' },
    ],
  },
  {
    label: '名称',
    prop: 'name',
    type: 'input',
    placeholder: '请输入名称',
  },
  {
    label: '创建时间',
    prop: 'createTime',
    type: 'daterange',
  },
]

const handleApply = (values: Record<string, any>) => {
  console.log('应用筛选:', values)
  // 执行筛选逻辑
}

const handleReset = () => {
  console.log('重置筛选')
  // 重置筛选逻辑
}
</script>
```

### 多选筛选

```vue
<template>
  <FilterPanel
    :filters="filters"
    v-model="filterValues"
    @apply="handleApply"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { FilterPanel, type FilterOption } from '@/components/common/business'

const filterValues = ref<Record<string, any>>({})

const filters: FilterOption[] = [
  {
    label: '标签',
    prop: 'tags',
    type: 'select',
    multiple: true, // 多选
    options: [
      { label: '重要', value: 'important' },
      { label: '紧急', value: 'urgent' },
      { label: '待办', value: 'todo' },
    ],
  },
]
</script>
```

### 不可折叠

```vue
<template>
  <FilterPanel
    :filters="filters"
    v-model="filterValues"
    :collapsible="false"
  />
</template>
```

### 必填筛选

```vue
<template>
  <FilterPanel
    :filters="filters"
    v-model="filterValues"
    @apply="handleApply"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { FilterPanel, type FilterOption } from '@/components/common/business'

const filterValues = ref<Record<string, any>>({})

const filters: FilterOption[] = [
  {
    label: '状态',
    prop: 'status',
    type: 'select',
    required: true, // 必填
    options: [
      { label: '运行中', value: 'running' },
      { label: '已停止', value: 'stopped' },
    ],
  },
]

const handleApply = (values: Record<string, any>) => {
  // 验证必填项
  if (!values.status) {
    ElMessage.warning('请选择状态')
    return
  }
  // 执行筛选逻辑
}
</script>
```

## 设计说明

- 支持多种筛选类型：选择器、输入框、日期选择器、日期范围选择器
- 可折叠设计，节省页面空间
- 筛选值变化会触发 `change` 事件，点击"应用"触发 `apply` 事件
- 点击"重置"会清空所有筛选值并触发 `reset` 事件
- 使用网格布局，自动适应不同屏幕尺寸

