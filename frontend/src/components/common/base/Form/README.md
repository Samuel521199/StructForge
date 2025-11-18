# Form 表单组件

表单组件，用于收集、校验和提交数据。

## 功能特性

- ✅ 表单验证
- ✅ 多种布局方式
- ✅ 多种尺寸
- ✅ 禁用状态
- ✅ 表单方法（validate, resetFields等）

## 基础用法

```vue
<template>
  <Form ref="formRef" :model="form" :rules="rules" label-width="100px">
    <FormItem label="用户名" prop="username">
      <Input v-model="form.username" />
    </FormItem>
    <FormItem label="密码" prop="password">
      <Input v-model="form.password" type="password" />
    </FormItem>
    <FormItem>
      <Button type="primary" @click="handleSubmit">提交</Button>
      <Button @click="handleReset">重置</Button>
    </FormItem>
  </Form>
</template>

<script setup>
import { ref } from 'vue'

const formRef = ref()
const form = ref({
  username: '',
  password: '',
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  await formRef.value.validate()
  // 提交表单
}

const handleReset = () => {
  formRef.value.resetFields()
}
</script>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| model | 表单数据对象 | `Record<string, any>` | - |
| rules | 表单验证规则 | `FormRules` | - |
| labelWidth | 表单域标签的宽度 | `string` | `'100px'` |
| labelPosition | 表单域标签的位置 | `'left' \| 'right' \| 'top'` | `'right'` |
| size | 用于控制该表单内组件的尺寸 | `'large' \| 'default' \| 'small'` | `'default'` |
| disabled | 是否禁用该表单内的所有组件 | `boolean` | `false` |

### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| validate | 任一表单项被校验后触发 | `(prop: string, isValid: boolean, message: string)` |

### Methods

通过ref调用表单方法：

| 方法名 | 说明 | 参数 |
|--------|------|------|
| validate | 对整个表单进行校验 | `() => Promise<void>` |
| validateField | 对部分表单字段进行校验 | `(props: string \| string[]) => Promise<void>` |
| resetFields | 重置表单 | `() => void` |
| clearValidate | 移除表单项的校验结果 | `(props?: string \| string[]) => void` |
| scrollToField | 滚动到指定字段 | `(prop: string) => void` |

### Slots

| 插槽名 | 说明 |
|--------|------|
| default | 表单内容 |

### Types

```typescript
interface FormRules {
  [key: string]: Array<{
    required?: boolean
    message?: string
    trigger?: 'blur' | 'change'
    validator?: (rule: any, value: any, callback: (error?: Error) => void) => void
  }>
}
```

## 使用示例

### 表单验证

```vue
<template>
  <Form ref="formRef" :model="form" :rules="rules">
    <FormItem label="邮箱" prop="email">
      <Input v-model="form.email" />
    </FormItem>
  </Form>
</template>

<script setup>
const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}
</script>
```

### 自定义验证器

```vue
<script setup>
const rules = {
  password: [
    {
      validator: (rule, value, callback) => {
        if (value.length < 6) {
          callback(new Error('密码长度不能少于6位'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}
</script>
```

### 标签位置

```vue
<!-- 标签在左侧 -->
<Form :model="form" label-position="left">
  <!-- ... -->
</Form>

<!-- 标签在顶部 -->
<Form :model="form" label-position="top">
  <!-- ... -->
</Form>
```

## 设计说明

- 完全兼容Element Plus Form的所有功能
- 提供完整的表单验证能力
- 支持异步验证
- 暴露所有表单方法供外部调用

