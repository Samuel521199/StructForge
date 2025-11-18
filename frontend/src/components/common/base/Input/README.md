# Input 输入框组件

文本输入组件，支持多种输入类型和功能。

## 功能特性

- ✅ 多种输入类型（text, password, number, email等）
- ✅ 清空功能
- ✅ 密码显示/隐藏
- ✅ 字数限制
- ✅ 前缀/后缀图标
- ✅ 前置/后置内容
- ✅ 多种尺寸

## 基础用法

```vue
<template>
  <Input v-model="value" placeholder="请输入内容" />
</template>

<script setup>
import { ref } from 'vue'

const value = ref('')
</script>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 绑定值 | `string \| number` | - |
| type | 输入框类型 | `'text' \| 'password' \| 'number' \| 'email' \| 'url' \| 'tel'` | `'text'` |
| placeholder | 输入框占位文本 | `string` | - |
| disabled | 是否禁用 | `boolean` | `false` |
| readonly | 是否只读 | `boolean` | `false` |
| clearable | 是否可清空 | `boolean` | `false` |
| showPassword | 是否显示切换密码图标 | `boolean` | `false` |
| prefixIcon | 自定义前缀图标 | `string` | - |
| suffixIcon | 自定义后缀图标 | `string` | - |
| maxlength | 最大输入长度 | `number` | - |
| minlength | 最小输入长度 | `number` | - |
| showWordLimit | 是否显示统计字数 | `boolean` | `false` |
| validateEvent | 是否触发表单验证 | `boolean` | `true` |
| size | 输入框尺寸 | `'large' \| 'default' \| 'small'` | `'default'` |

### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | v-model更新 | `(value: string \| number)` |
| focus | 获得焦点 | `(event: FocusEvent)` |
| blur | 失去焦点 | `(event: FocusEvent)` |
| clear | 清空 | - |
| input | 输入事件 | `(value: string \| number)` |

### Slots

| 插槽名 | 说明 |
|--------|------|
| prefix | 输入框头部内容 |
| suffix | 输入框尾部内容 |
| prepend | 输入框前置内容 |
| append | 输入框后置内容 |

## 使用示例

### 清空功能

```vue
<Input v-model="value" placeholder="请输入内容" clearable />
```

### 密码输入

```vue
<Input
  v-model="password"
  type="password"
  placeholder="请输入密码"
  show-password
/>
```

### 字数限制

```vue
<Input
  v-model="value"
  placeholder="最多输入20个字符"
  :maxlength="20"
  show-word-limit
/>
```

### 带图标

```vue
<Input
  v-model="value"
  placeholder="请输入内容"
  prefix-icon="Search"
/>
```

### 前置/后置内容

```vue
<template>
  <Input v-model="value" placeholder="请输入内容">
    <template #prepend>Http://</template>
    <template #append>.com</template>
  </Input>
</template>
```

## 设计说明

- 完全兼容Element Plus Input的所有功能
- 支持v-model双向绑定
- 提供丰富的插槽支持自定义内容
