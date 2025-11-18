# Empty 空状态组件

空状态组件，用于显示空数据状态。

## 功能特性

- ✅ 自定义图片
- ✅ 自定义描述文字
- ✅ 自定义内容
- ✅ 多种尺寸

## 基础用法

```vue
<template>
  <Empty description="暂无数据" />
</template>
```

## API

### Props

| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| image | 图片地址 | `string` | - |
| imageSize | 图片大小（宽度） | `number` | `200` |
| description | 描述文字 | `string` | `'暂无数据'` |

### Slots

| 插槽名 | 说明 |
|--------|------|
| default | 自定义内容 |
| image | 自定义图片 |
| description | 自定义描述 |

## 使用示例

### 基础用法

```vue
<Empty description="暂无工作流" />
```

### 自定义图片

```vue
<Empty
  image="/images/empty-workflow.png"
  description="还没有创建工作流"
/>
```

### 自定义内容

```vue
<Empty description="暂无数据">
  <Button type="primary" @click="handleCreate">创建第一个</Button>
</Empty>
```

### 在表格中使用

```vue
<template>
  <Table :data="tableData" :columns="columns">
    <template #empty>
      <Empty description="暂无数据" />
    </template>
  </Table>
</template>
```

### 自定义图片和描述

```vue
<Empty>
  <template #image>
    <img src="/images/custom-empty.png" alt="空状态" />
  </template>
  <template #description>
    <p>自定义描述内容</p>
  </template>
  <Button type="primary">立即创建</Button>
</Empty>
```

## 设计说明

- 基于Element Plus Empty封装
- 提供灵活的插槽支持
- 适用于各种空数据场景

## 使用场景

- 列表为空时
- 搜索结果为空时
- 数据加载失败时
- 权限不足时

