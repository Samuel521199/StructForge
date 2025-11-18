# Table 表格组件

数据表格组件，用于展示结构化数据。

## 功能特性

- ✅ 基础表格展示
- ✅ 斑马纹
- ✅ 边框
- ✅ 固定列
- ✅ 排序
- ✅ 选择
- ✅ 加载状态
- ✅ 空状态
- ✅ 自定义列渲染

## 基础用法

```vue
<template>
  <Table :data="tableData" :columns="columns" />
</template>

<script setup>
import { ref } from 'vue'

const tableData = ref([
  { id: 1, name: '工作流1', status: 'running', createTime: '2024-01-01' },
  { id: 2, name: '工作流2', status: 'stopped', createTime: '2024-01-02' },
])

const columns = [
  { prop: 'name', label: '名称', width: 200 },
  { prop: 'status', label: '状态', width: 100 },
  { prop: 'createTime', label: '创建时间', width: 180 },
]
</script>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| data | 表格数据 | `any[]` | - |
| columns | 表格列配置 | `TableColumn[]` | - |
| stripe | 是否为斑马纹表格 | `boolean` | `false` |
| border | 是否带有纵向边框 | `boolean` | `false` |
| size | 表格尺寸 | `'large' \| 'default' \| 'small'` | `'default'` |
| showHeader | 是否显示表头 | `boolean` | `true` |
| highlightCurrentRow | 是否要高亮当前行 | `boolean` | `false` |
| emptyText | 空数据时显示的文本 | `string` | `'暂无数据'` |
| loading | 是否显示加载状态 | `boolean` | `false` |
| height | 表格高度 | `string \| number` | - |
| maxHeight | 表格最大高度 | `string \| number` | - |

### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| selection-change | 当选择项发生变化时会触发该事件 | `(selection: any[])` |
| row-click | 当某一行被点击时会触发该事件 | `(row: any, column: TableColumn, event: Event)` |
| sort-change | 当表格的排序条件发生变化的时候会触发该事件 | `(sortInfo: { column: TableColumn, prop: string, order: string })` |

### Slots

| 插槽名 | 说明 |
|--------|------|
| empty | 空数据时显示的内容 |

### Types

```typescript
interface TableColumn {
  prop?: string // 对应列内容的字段名
  label: string // 显示的标题
  width?: string | number // 对应列的宽度
  minWidth?: string | number // 对应列的最小宽度
  fixed?: boolean | 'left' | 'right' // 列是否固定
  sortable?: boolean // 对应列是否可以排序
  formatter?: (row: any, column: TableColumn, cellValue: any) => any // 用来格式化内容
  align?: 'left' | 'center' | 'right' // 对齐方式
}
```

## 使用示例

### 斑马纹表格

```vue
<Table :data="tableData" :columns="columns" stripe />
```

### 带边框表格

```vue
<Table :data="tableData" :columns="columns" border />
```

### 固定列

```vue
<template>
  <Table :data="tableData" :columns="columns" />
</template>

<script setup>
const columns = [
  { prop: 'id', label: 'ID', width: 80, fixed: 'left' },
  { prop: 'name', label: '名称', width: 200 },
  { prop: 'action', label: '操作', width: 150, fixed: 'right' },
]
</script>
```

### 可排序

```vue
<template>
  <Table :data="tableData" :columns="columns" @sort-change="handleSortChange" />
</template>

<script setup>
const columns = [
  { prop: 'name', label: '名称', sortable: true },
  { prop: 'createTime', label: '创建时间', sortable: true },
]
</script>
```

### 加载状态

```vue
<Table :data="tableData" :columns="columns" :loading="loading" />
```

### 空状态

```vue
<template>
  <Table :data="tableData" :columns="columns">
    <template #empty>
      <Empty description="暂无数据" />
    </template>
  </Table>
</template>
```

## 设计说明

- 基于Element Plus Table封装
- 简化列配置，使用columns数组
- 支持所有Element Plus Table的原生功能
- 可以通过插槽自定义列内容（需要在columns中配置）

