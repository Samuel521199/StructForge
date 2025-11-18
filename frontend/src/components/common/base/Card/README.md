# Card 卡片组件

卡片组件，用于展示内容区块。

## 功能特性

- ✅ 自定义头部
- ✅ 自定义内容
- ✅ 多种阴影效果
- ✅ 自定义样式

## 基础用法

```vue
<template>
  <Card header="卡片标题">
    <p>卡片内容</p>
  </Card>
</template>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| header | 卡片标题 | `string` | - |
| shadow | 卡片阴影效果 | `'always' \| 'hover' \| 'never'` | `'always'` |
| bodyStyle | 卡片body的样式 | `Record<string, any>` | `{}` |

### Slots

| 插槽名 | 说明 |
|--------|------|
| default | 卡片内容 |
| header | 自定义头部内容 |

## 使用示例

### 基础卡片

```vue
<Card header="工作流列表">
  <p>这是卡片内容</p>
</Card>
```

### 自定义头部

```vue
<Card>
  <template #header>
    <div style="display: flex; justify-content: space-between;">
      <span>标题</span>
      <Button size="small">操作</Button>
    </div>
  </template>
  <p>卡片内容</p>
</Card>
```

### 无头部卡片

```vue
<Card>
  <p>没有头部的卡片</p>
</Card>
```

### 不同阴影效果

```vue
<!-- 始终显示阴影 -->
<Card shadow="always">内容</Card>

<!-- 悬停时显示阴影 -->
<Card shadow="hover">内容</Card>

<!-- 无阴影 -->
<Card shadow="never">内容</Card>
```

### 自定义样式

```vue
<Card :body-style="{ padding: '20px' }">
  <p>自定义内边距的卡片</p>
</Card>
```

## 设计说明

- 基于Element Plus Card封装
- 保持简洁的API
- 支持灵活的插槽自定义

