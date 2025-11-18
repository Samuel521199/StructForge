# Loading 加载组件

加载组件，用于显示加载状态。

## 功能特性

- ✅ 局部加载
- ✅ 全屏加载
- ✅ 自定义加载文字
- ✅ 自定义背景色
- ✅ 锁定滚动

## 基础用法

```vue
<template>
  <Loading :loading="isLoading">
    <div>内容区域</div>
  </Loading>
</template>

<script setup>
import { ref } from 'vue'

const isLoading = ref(false)
</script>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| loading | 是否显示加载状态 | `boolean` | `false` |
| text | 加载文字 | `string` | `'加载中...'` |
| background | 背景色 | `string` | `'rgba(0, 0, 0, 0.7)'` |
| spinner | 自定义加载图标类名 | `string` | - |
| element | Loading 覆盖的 DOM 节点 | `string \| HTMLElement` | - |
| target | Loading 需要覆盖的 DOM 节点 | `string \| HTMLElement` | - |
| body | 是否将遮罩层插入到 body 中 | `boolean` | - |
| fullscreen | 是否全屏显示 | `boolean` | `false` |
| lock | 是否锁定屏幕滚动 | `boolean` | `false` |

### Slots

| 插槽名 | 说明 |
|--------|------|
| default | 需要显示加载状态的内容 |

## 使用示例

### 局部加载

```vue
<template>
  <Loading :loading="isLoading" text="加载中...">
    <Table :data="tableData" :columns="columns" />
  </Loading>
</template>
```

### 全屏加载

```vue
<Loading :loading="isLoading" fullscreen text="正在加载..." />
```

### 自定义加载文字和背景

```vue
<Loading
  :loading="isLoading"
  text="正在处理..."
  background="rgba(255, 255, 255, 0.8)"
>
  <div>内容</div>
</Loading>
```

### 在表格中使用

```vue
<template>
  <Table :data="tableData" :columns="columns" :loading="isLoading" />
</template>
```

## 设计说明

- 基于Element Plus Loading封装
- 支持局部和全屏加载
- 可以通过v-loading指令使用（推荐）

## 使用建议

- 对于表格、列表等组件，建议使用组件自带的loading属性
- 对于自定义区域，使用Loading组件包裹
- 对于全屏加载，使用fullscreen属性

