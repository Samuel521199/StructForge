# Button 按钮组件

基础按钮组件，支持多种类型、尺寸和状态。

## 功能特性

- ✅ 多种按钮类型（primary, success, warning, danger, info, text）
- ✅ 多种尺寸（large, default, small）
- ✅ 加载状态
- ✅ 禁用状态
- ✅ 图标支持
- ✅ 圆角/圆形按钮
- ✅ 朴素按钮样式

## 基础用法

```vue
<template>
  <Button type="primary" @click="handleClick">主要按钮</Button>
  <Button type="success">成功按钮</Button>
  <Button type="warning">警告按钮</Button>
  <Button type="danger">危险按钮</Button>
</template>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| type | 按钮类型 | `'primary' \| 'success' \| 'warning' \| 'danger' \| 'info' \| 'text'` | `'default'` |
| size | 按钮尺寸 | `'large' \| 'default' \| 'small'` | `'default'` |
| disabled | 是否禁用 | `boolean` | `false` |
| loading | 是否加载中 | `boolean` | `false` |
| icon | 图标类名 | `string` | - |
| round | 是否圆角按钮 | `boolean` | `false` |
| circle | 是否圆形按钮 | `boolean` | `false` |
| plain | 是否朴素按钮 | `boolean` | `false` |
| nativeType | 原生 type 属性 | `'button' \| 'submit' \| 'reset'` | `'button'` |

### Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| click | 点击事件 | `(event: MouseEvent)` |

### Slots

| 插槽名 | 说明 |
|--------|------|
| default | 按钮内容 |

## 使用示例

### 不同尺寸

```vue
<Button size="large">大按钮</Button>
<Button size="default">默认按钮</Button>
<Button size="small">小按钮</Button>
```

### 加载状态

```vue
<Button type="primary" :loading="isLoading" @click="handleSubmit">
  提交
</Button>
```

### 图标按钮

```vue
<Button type="primary" icon="Search">搜索</Button>
<Button type="primary" icon="Edit" circle />
```

### 禁用状态

```vue
<Button disabled>禁用按钮</Button>
<Button type="primary" disabled>禁用主要按钮</Button>
```

## 设计说明

- 按钮类型遵循Element Plus设计规范
- 支持所有Element Plus Button的原生功能
- 保持API一致性，便于使用和维护
