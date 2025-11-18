# CopyButton 复制按钮组件

## 功能说明

CopyButton 是一个复制按钮组件，点击后可以将内容复制到剪贴板，并显示复制成功状态。

## Props

| 参数 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| value | 要复制的内容 | `string \| object` | - | 是 |
| text | 按钮文本 | `string` | `'复制'` | 否 |
| successText | 复制成功后的文本 | `string` | `'已复制'` | 否 |
| type | 按钮类型 | `'primary' \| 'success' \| 'warning' \| 'danger' \| 'info' \| 'text' \| 'default'` | `'default'` | 否 |
| size | 按钮大小 | `'large' \| 'default' \| 'small'` | `'default'` | 否 |
| disabled | 是否禁用 | `boolean` | `false` | 否 |
| round | 是否圆角 | `boolean` | `false` | 否 |
| plain | 是否朴素按钮 | `boolean` | `false` | 否 |
| showMessage | 是否显示消息提示 | `boolean` | `true` | 否 |

## Events

无（通过消息提示反馈）

## Slots

| 插槽名 | 说明 |
|--------|------|
| default | 自定义按钮内容 |

## 使用示例

### 基础使用

```vue
<template>
  <CopyButton :value="text" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CopyButton } from '@/components/common/business'

const text = ref('要复制的内容')
</script>
```

### 复制对象

```vue
<template>
  <CopyButton :value="data" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CopyButton } from '@/components/common/business'

const data = ref({
  name: 'StructForge',
  version: '1.0.0',
})
</script>
```

### 自定义按钮样式

```vue
<template>
  <CopyButton
    :value="text"
    type="primary"
    size="small"
    round
  />
</template>
```

### 自定义文本

```vue
<template>
  <CopyButton
    :value="text"
    text="复制代码"
    success-text="代码已复制"
  />
</template>
```

### 禁用消息提示

```vue
<template>
  <CopyButton
    :value="text"
    :show-message="false"
  />
</template>
```

### 自定义按钮内容

```vue
<template>
  <CopyButton :value="text">
    <el-icon><DocumentCopy /></el-icon>
    复制到剪贴板
  </CopyButton>
</template>
```

### 复制代码块

```vue
<template>
  <div>
    <pre>{{ code }}</pre>
    <CopyButton :value="code" text="复制代码" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CopyButton } from '@/components/common/business'

const code = ref(`function hello() {
  console.log('Hello, World!')
}`)
</script>
```

## 设计说明

- 使用 `navigator.clipboard.writeText` API 进行复制
- 复制成功后按钮图标会变为成功图标，2秒后恢复
- 复制失败会显示错误提示
- 对象类型会自动转换为格式化的 JSON 字符串
- 支持通过插槽自定义按钮内容
- 基于 Element Plus 的 Button 组件实现

## 注意事项

- 复制功能需要 HTTPS 环境或 localhost
- 某些浏览器可能需要用户交互才能复制
- 复制失败时会显示错误提示，但不会抛出异常

