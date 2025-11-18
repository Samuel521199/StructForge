# Message 消息提示组件

消息提示组件，通过方法调用显示全局消息提示。

## 功能特性

- ✅ 多种消息类型（success, warning, error, info）
- ✅ 自动关闭
- ✅ 可手动关闭
- ✅ 支持HTML内容
- ✅ 消息合并

## 基础用法

### 方法调用

```vue
<script setup>
import { success, error, warning, info } from '@/components/common/base/Message'

// 成功消息
success('操作成功')

// 错误消息
error('操作失败')

// 警告消息
warning('请注意')

// 信息消息
info('提示信息')
</script>
```

### 使用组合函数

```vue
<script setup>
import { useMessage } from '@/components/common/base/Message'

const message = useMessage()

const handleSuccess = () => {
  message.success('操作成功')
}

const handleError = () => {
  message.error('操作失败')
}
</script>
```

## API

### 方法

| 方法名 | 说明 | 参数 |
|--------|------|------|
| success | 显示成功消息 | `(message: string, options?: MessageOptions)` |
| error | 显示错误消息 | `(message: string, options?: MessageOptions)` |
| warning | 显示警告消息 | `(message: string, options?: MessageOptions)` |
| info | 显示信息消息 | `(message: string, options?: MessageOptions)` |
| closeAll | 关闭所有消息 | `() => void` |

### MessageOptions

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| duration | 显示时间（毫秒），0则不自动关闭 | `number` | `3000` |
| showClose | 是否显示关闭按钮 | `boolean` | `false` |
| center | 文字是否居中 | `boolean` | `false` |
| grouping | 是否将消息合并为一条 | `boolean` | `false` |
| dangerouslyUseHTMLString | 是否将message作为HTML处理 | `boolean` | `false` |
| customClass | 自定义类名 | `string` | - |
| offset | 距离窗口顶部的偏移量 | `number` | - |
| onClose | 关闭时的回调 | `() => void` | - |

## 使用示例

### 基础用法

```typescript
import { success, error } from '@/components/common/base/Message'

// 成功消息
success('保存成功')

// 错误消息
error('保存失败，请重试')
```

### 自定义选项

```typescript
import { success } from '@/components/common/base/Message'

success('操作成功', {
  duration: 5000,
  showClose: true,
  center: true
})
```

### 不自动关闭

```typescript
import { warning } from '@/components/common/base/Message'

warning('重要提示', {
  duration: 0, // 不自动关闭
  showClose: true
})
```

### HTML内容

```typescript
import { info } from '@/components/common/base/Message'

info('<strong>重要</strong>：请检查配置', {
  dangerouslyUseHTMLString: true
})
```

### 关闭所有消息

```typescript
import { closeAll } from '@/components/common/base/Message'

closeAll()
```

### 在组合函数中使用

```vue
<script setup>
import { useMessage } from '@/components/common/base/Message'

const message = useMessage()

const handleSubmit = async () => {
  try {
    await submitForm()
    message.success('提交成功')
  } catch (error) {
    message.error('提交失败：' + error.message)
  }
}
</script>
```

## 设计说明

- 基于Element Plus Message封装
- 提供统一的API接口
- 支持组合函数方式使用
- 所有消息类型都有对应的方法

## 使用建议

- 成功操作使用 `success`
- 错误信息使用 `error`
- 警告信息使用 `warning`
- 一般提示使用 `info`
- 重要消息设置 `duration: 0` 并显示关闭按钮

