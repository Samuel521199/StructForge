# Link 链接组件

基于 Element Plus `el-link` 的封装组件。

## 使用方式

```vue
<script setup>
import { Link } from '@/components/common/base'
</script>

<template>
  <Link type="primary" @click="handleClick">点击链接</Link>
</template>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| type | `'primary' \| 'success' \| 'warning' \| 'danger' \| 'info' \| 'default'` | `'default'` | 类型 |
| underline | `boolean` | `true` | 是否下划线 |
| disabled | `boolean` | `false` | 是否禁用 |
| href | `string` | - | 原生 href 属性 |
| target | `string` | - | 原生 target 属性 |
| icon | `string \| object` | - | 图标类名或组件 |

## Events

| 事件名 | 说明 | 回调参数 |
|--------|------|----------|
| click | 点击事件 | `(event: MouseEvent)` |

## 示例

### 基础用法

```vue
<Link type="primary">主要链接</Link>
<Link type="success">成功链接</Link>
<Link type="warning">警告链接</Link>
<Link type="danger">危险链接</Link>
```

### 禁用下划线

```vue
<Link type="primary" :underline="false">无下划线链接</Link>
```

### 禁用状态

```vue
<Link type="primary" disabled>禁用链接</Link>
```

