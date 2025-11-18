# AvatarGroup 头像组组件

## 功能说明

AvatarGroup 是一个头像组组件，用于显示多个头像，支持堆叠显示、最大数量限制等功能。

## Props

| 参数 | 说明 | 类型 | 默认值 | 必填 |
|------|------|------|--------|------|
| avatars | 头像列表 | `AvatarItem[]` | - | 是 |
| size | 头像大小 | `'large' \| 'default' \| 'small' \| number` | `'default'` | 否 |
| shape | 头像形状 | `'circle' \| 'square'` | `'circle'` | 否 |
| stacked | 是否堆叠显示 | `boolean` | `false` | 否 |
| max | 最大显示数量（超出显示+数字） | `number` | `0`（不限制） | 否 |

### AvatarItem 类型

```typescript
interface AvatarItem {
  src?: string                    // 头像图片地址
  text?: string                   // 头像文本（无图片时显示）
  icon?: Component | string       // 头像图标
  shape?: 'circle' | 'square'     // 头像形状
  label?: string                  // 标签文本
  data?: any                      // 额外数据
}
```

## Events

| 事件名 | 说明 | 参数 |
|--------|------|------|
| avatarClick | 头像点击事件 | `(avatar: AvatarItem, index: number) => void` |
| moreClick | 更多点击事件 | `() => void` |

## 使用示例

### 基础使用

```vue
<template>
  <AvatarGroup :avatars="avatars" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { AvatarGroup, type AvatarItem } from '@/components/common/business'

const avatars: AvatarItem[] = [
  { src: 'https://example.com/avatar1.jpg', text: '用户1' },
  { src: 'https://example.com/avatar2.jpg', text: '用户2' },
  { text: '用户3' }, // 无图片，显示文本
]
</script>
```

### 堆叠显示

```vue
<template>
  <AvatarGroup :avatars="avatars" stacked />
</template>
```

### 限制显示数量

```vue
<template>
  <AvatarGroup
    :avatars="avatars"
    :max="3"
    @more-click="handleMoreClick"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { AvatarGroup, type AvatarItem } from '@/components/common/business'

const avatars: AvatarItem[] = [
  { text: '用户1' },
  { text: '用户2' },
  { text: '用户3' },
  { text: '用户4' },
  { text: '用户5' },
]

const handleMoreClick = () => {
  console.log('查看更多')
}
</script>
```

### 不同大小

```vue
<template>
  <AvatarGroup :avatars="avatars" size="large" />
  <AvatarGroup :avatars="avatars" size="default" />
  <AvatarGroup :avatars="avatars" size="small" />
  <AvatarGroup :avatars="avatars" :size="50" />
</template>
```

### 方形头像

```vue
<template>
  <AvatarGroup :avatars="avatars" shape="square" />
</template>
```

### 带标签的头像

```vue
<template>
  <AvatarGroup :avatars="avatars" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { AvatarGroup, type AvatarItem } from '@/components/common/business'

const avatars: AvatarItem[] = [
  { text: '张三', label: '管理员' },
  { text: '李四', label: '成员' },
  { text: '王五', label: '访客' },
]
</script>
```

### 点击事件

```vue
<template>
  <AvatarGroup
    :avatars="avatars"
    @avatar-click="handleAvatarClick"
    @more-click="handleMoreClick"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { AvatarGroup, type AvatarItem } from '@/components/common/business'

const avatars: AvatarItem[] = [
  { text: '用户1', data: { id: 1 } },
  { text: '用户2', data: { id: 2 } },
]

const handleAvatarClick = (avatar: AvatarItem, index: number) => {
  console.log('点击头像:', avatar, index)
  // 可以访问 avatar.data 获取额外数据
}

const handleMoreClick = () => {
  console.log('点击更多')
}
</script>
```

## 设计说明

- 基于 Element Plus 的 `el-avatar` 组件实现
- 堆叠显示时，头像会有重叠效果，z-index 自动调整
- 超出最大数量时，显示 "+N" 的头像
- 支持图片、文本、图标三种头像显示方式
- 鼠标悬停时头像会有轻微上移效果
- 支持通过 `data` 属性传递额外数据，便于点击事件处理

