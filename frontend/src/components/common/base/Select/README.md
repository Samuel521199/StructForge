# Select 选择器组件

下拉选择组件，支持单选、多选、分组等功能。

## 功能特性

- ✅ 单选/多选
- ✅ 可搜索
- ✅ 可清空
- ✅ 可创建新选项
- ✅ 分组选项
- ✅ 加载状态
- ✅ 多种尺寸

## 基础用法

```vue
<template>
  <Select v-model="value" :options="options" placeholder="请选择" />
</template>

<script setup>
import { ref } from 'vue'

const value = ref('')
const options = [
  { label: '选项1', value: '1' },
  { label: '选项2', value: '2' },
  { label: '选项3', value: '3' },
]
</script>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 绑定值 | `string \| number \| Array<string \| number>` | - |
| options | 选项数据 | `SelectOption[]` | - |
| placeholder | 输入框占位文本 | `string` | `'请选择'` |
| multiple | 是否多选 | `boolean` | `false` |
| disabled | 是否禁用 | `boolean` | `false` |
| clearable | 是否可以清空选项 | `boolean` | `false` |
| filterable | 是否可搜索 | `boolean` | `false` |
| allowCreate | 是否允许用户创建新条目 | `boolean` | `false` |
| size | 输入框尺寸 | `'large' \| 'default' \| 'small'` | `'default'` |
| loading | 是否加载中 | `boolean` | `false` |

### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | v-model更新 | `(value: string \| number \| Array<string \| number>)` |
| change | 选中值发生变化时触发 | `(value: string \| number \| Array<string \| number>)` |
| visible-change | 下拉框出现/隐藏时触发 | `(visible: boolean)` |
| remove-tag | 多选模式下移除tag时触发 | `(tag: string \| number)` |
| clear | 可清空的单选模式下用户点击清空按钮时触发 | - |

### Types

```typescript
interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
  children?: SelectOption[] // 分组选项
}
```

## 使用示例

### 多选

```vue
<Select
  v-model="value"
  :options="options"
  multiple
  placeholder="请选择"
/>
```

### 可搜索

```vue
<Select
  v-model="value"
  :options="options"
  filterable
  placeholder="请输入关键词搜索"
/>
```

### 分组选项

```vue
<template>
  <Select v-model="value" :options="groupOptions" />
</template>

<script setup>
const groupOptions = [
  {
    label: '热门城市',
    children: [
      { label: '北京', value: 'beijing' },
      { label: '上海', value: 'shanghai' },
    ]
  },
  {
    label: '其他城市',
    children: [
      { label: '广州', value: 'guangzhou' },
      { label: '深圳', value: 'shenzhen' },
    ]
  }
]
</script>
```

### 可创建新选项

```vue
<Select
  v-model="value"
  :options="options"
  filterable
  allow-create
  placeholder="请选择或输入"
/>
```

## 设计说明

- 支持Element Plus Select的所有功能
- 简化了选项配置，使用options数组
- 支持分组选项嵌套结构

