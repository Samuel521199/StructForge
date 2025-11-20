# Checkbox 复选框组件

基于 Element Plus `el-checkbox` 的封装组件。

## 使用方式

```vue
<script setup>
import { Checkbox } from '@/components/common/base'
import { ref } from 'vue'

const checked = ref(false)
</script>

<template>
  <Checkbox v-model="checked">同意协议</Checkbox>
</template>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `boolean \| string \| number` | - | 绑定值 |
| label | `string \| number \| boolean` | - | 选中状态的值（当 modelValue 为数组时使用） |
| disabled | `boolean` | `false` | 是否禁用 |
| readonly | `boolean` | `false` | 是否只读 |
| border | `boolean` | `false` | 是否显示边框 |
| size | `'large' \| 'default' \| 'small'` | `'default'` | 复选框的尺寸 |
| name | `string` | - | 原生 name 属性 |
| id | `string` | - | 原生 id 属性 |
| indeterminate | `boolean` | `false` | 是否不确定状态 |
| checked | `boolean` | - | 是否选中 |
| tabindex | `string \| number` | - | 原生 tabindex 属性 |
| validateEvent | `boolean` | `true` | 是否触发表单验证 |

## Events

| 事件名 | 说明 | 回调参数 |
|--------|------|----------|
| update:modelValue | 更新绑定值时触发 | `(value: boolean \| string \| number)` |
| change | 值改变时触发 | `(value: boolean \| string \| number)` |

## 示例

### 基础用法

```vue
<Checkbox v-model="checked">选项</Checkbox>
```

### 禁用状态

```vue
<Checkbox v-model="checked" disabled>禁用选项</Checkbox>
```

### 带边框

```vue
<Checkbox v-model="checked" border>带边框选项</Checkbox>
```

