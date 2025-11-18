# Dialog 对话框组件

模态对话框组件，用于显示重要信息或收集用户输入。

## 功能特性

- ✅ 模态对话框
- ✅ 自定义宽度和位置
- ✅ 全屏模式
- ✅ 可自定义头部和底部
- ✅ 多种关闭方式
- ✅ 锁定滚动

## 基础用法

```vue
<template>
  <Dialog v-model="visible" title="提示">
    <p>这是一段内容</p>
    <template #footer>
      <Button @click="visible = false">取消</Button>
      <Button type="primary" @click="handleConfirm">确认</Button>
    </template>
  </Dialog>
</template>

<script setup>
import { ref } from 'vue'

const visible = ref(false)

const handleConfirm = () => {
  // 处理确认逻辑
  visible.value = false
}
</script>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| modelValue | 是否显示对话框 | `boolean` | - |
| title | 对话框标题 | `string` | - |
| width | 对话框宽度 | `string \| number` | `'50%'` |
| fullscreen | 是否全屏 | `boolean` | `false` |
| top | 对话框距离顶部的位置 | `string` | `'15vh'` |
| modal | 是否显示遮罩层 | `boolean` | `true` |
| closeOnClickModal | 是否可以通过点击遮罩层关闭 | `boolean` | `true` |
| closeOnPressEscape | 是否可以通过按下 ESC 关闭 | `boolean` | `true` |
| showClose | 是否显示关闭按钮 | `boolean` | `true` |
| appendToBody | 是否将对话框插入到 body 元素上 | `boolean` | `false` |
| lockScroll | 是否在对话框出现时将 body 滚动锁定 | `boolean` | `true` |

### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| update:modelValue | v-model更新 | `(value: boolean)` |
| open | 对话框打开时触发 | - |
| close | 对话框关闭时触发 | - |
| opened | 对话框打开动画结束时触发 | - |
| closed | 对话框关闭动画结束时触发 | - |

### Slots

| 插槽名 | 说明 |
|--------|------|
| default | 对话框内容 |
| header | 自定义头部内容 |
| footer | 自定义底部内容 |

## 使用示例

### 自定义宽度

```vue
<Dialog v-model="visible" title="提示" width="500px">
  <p>内容</p>
</Dialog>
```

### 全屏对话框

```vue
<Dialog v-model="visible" title="全屏对话框" fullscreen>
  <p>全屏内容</p>
</Dialog>
```

### 自定义头部和底部

```vue
<Dialog v-model="visible">
  <template #header>
    <h2>自定义标题</h2>
  </template>
  <p>内容</p>
  <template #footer>
    <Button @click="visible = false">取消</Button>
    <Button type="primary" @click="handleConfirm">确认</Button>
  </template>
</Dialog>
```

### 确认对话框

```vue
<Dialog v-model="deleteVisible" title="确认删除" width="400px">
  <p>确定要删除这个工作流吗？此操作不可恢复。</p>
  <template #footer>
    <Button @click="deleteVisible = false">取消</Button>
    <Button type="danger" @click="handleDelete">确认删除</Button>
  </template>
</Dialog>
```

## 设计说明

- 完全兼容Element Plus Dialog的所有功能
- 支持v-model双向绑定
- 提供灵活的插槽支持自定义内容

