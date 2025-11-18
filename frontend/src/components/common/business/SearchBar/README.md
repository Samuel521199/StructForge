# SearchBar 搜索栏组件

## 功能说明

SearchBar 是一个集成了搜索和筛选功能的复合组件，适用于列表页面的搜索场景。

## Props

| 参数 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| modelValue | 搜索关键词（v-model） | `string` | - | 是 |
| placeholder | 搜索框占位符 | `string` | `'请输入搜索关键词'` | 否 |
| searchable | 是否显示搜索框 | `boolean` | `true` | 否 |
| filterable | 是否显示筛选功能 | `boolean` | `false` | 否 |
| filters | 筛选配置列表 | `FilterOption[]` | - | 否 |

### FilterOption 类型

```typescript
interface FilterOption {
  label: string           // 筛选标签
  prop: string            // 筛选字段名
  type?: 'select' | 'date' | 'daterange' | 'input'  // 筛选类型
  options?: Array<{ label: string; value: any }>    // 选项数据（select类型）
  placeholder?: string    // 占位符
}
```

## Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | 搜索关键词更新 | `(value: string) => void` |
| search | 搜索事件 | `(keyword: string) => void` |
| filter | 筛选事件 | `(filters: Record<string, any>) => void` |

## 使用示例

### 基础搜索

```vue
<template>
  <SearchBar
    v-model="keyword"
    placeholder="搜索工作流名称"
    @search="handleSearch"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { SearchBar } from '@/components/common/business'

const keyword = ref('')

const handleSearch = (keyword: string) => {
  console.log('搜索:', keyword)
  // 执行搜索逻辑
}
</script>
```

### 带筛选的搜索

```vue
<template>
  <SearchBar
    v-model="keyword"
    :filters="filters"
    filterable
    @search="handleSearch"
    @filter="handleFilter"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { SearchBar, type FilterOption } from '@/components/common/business'

const keyword = ref('')

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
    label: '创建时间',
    prop: 'createTime',
    type: 'daterange',
  },
]

const handleSearch = (keyword: string) => {
  console.log('搜索:', keyword)
}

const handleFilter = (filters: Record<string, any>) => {
  console.log('筛选:', filters)
}
</script>
```

## 设计说明

- 搜索框支持回车搜索和清除功能
- 筛选面板可展开/收起，显示激活的筛选条件数量
- 筛选支持选择器、输入框、日期选择器等多种类型
- 筛选值变化后需要点击"应用"按钮才会触发筛选事件

