# TimeAgo 相对时间组件

## 功能说明

TimeAgo 是一个时间显示组件，可以显示相对时间（如"5分钟前"）或完整时间，支持自动更新。

## Props

| 参数 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| time | 时间（时间戳、日期字符串或Date对象） | `number \| string \| Date` | - | 是 |
| format | 完整时间格式 | `string` | `'YYYY-MM-DD HH:mm:ss'` | 否 |
| showFullTime | 是否显示完整时间 | `boolean` | `false` | 否 |
| updateInterval | 更新间隔（毫秒），0表示不自动更新 | `number` | `60000` | 否 |

## 使用示例

### 基础使用（相对时间）

```vue
<template>
  <TimeAgo :time="createTime" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { TimeAgo } from '@/components/common/business'

const createTime = ref(Date.now() - 3600000) // 1小时前
</script>
```

### 显示完整时间

```vue
<template>
  <TimeAgo
    :time="createTime"
    show-full-time
    format="YYYY-MM-DD HH:mm:ss"
  />
</template>
```

### 自定义格式

```vue
<template>
  <TimeAgo
    :time="createTime"
    show-full-time
    format="YYYY年MM月DD日 HH:mm"
  />
</template>
```

### 禁用自动更新

```vue
<template>
  <TimeAgo
    :time="createTime"
    :update-interval="0"
  />
</template>
```

### 自定义更新间隔

```vue
<template>
  <!-- 每30秒更新一次 -->
  <TimeAgo
    :time="createTime"
    :update-interval="30000"
  />
</template>
```

### 不同时间格式

```vue
<template>
  <!-- 时间戳 -->
  <TimeAgo :time="1640995200000" />
  
  <!-- 日期字符串 -->
  <TimeAgo time="2024-01-01 12:00:00" />
  
  <!-- Date对象 -->
  <TimeAgo :time="new Date()" />
</template>
```

## 相对时间规则

- **刚刚**: 小于60秒
- **X分钟前**: 小于60分钟
- **X小时前**: 小于24小时
- **X天前**: 小于30天
- **X个月前**: 小于12个月
- **X年前**: 大于等于12个月
- **未来**: 时间在未来

## 设计说明

- 鼠标悬停显示完整时间（通过 `title` 属性）
- 支持自动更新，默认每1分钟更新一次
- 支持自定义更新间隔，设置为0可禁用自动更新
- 支持多种时间格式输入：时间戳、日期字符串、Date对象
- 相对时间显示更友好，适合列表、评论等场景

