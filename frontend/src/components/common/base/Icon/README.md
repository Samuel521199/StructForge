# Icon 组件

封装 Element Plus Icons，提供统一的图标使用方式。

## 使用方式

```vue
<template>
  <!-- 使用 Element Plus Icons -->
  <Icon :icon="User" :size="20" color="#00FF00" />
  
  <!-- 加载中状态 -->
  <Icon :icon="Loading" :size="24" :is-loading="true" />
</template>

<script setup lang="ts">
import { Icon } from '@/components/common/base'
import { User, Loading } from '@element-plus/icons-vue'
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| icon | Component \| string | - | 图标组件（来自 @element-plus/icons-vue） |
| size | number \| string | 16 | 图标尺寸 |
| color | string | - | 图标颜色 |
| isLoading | boolean | false | 是否加载中（显示旋转动画） |
| class | string | - | 自定义类名 |

## 示例

```vue
<template>
  <div>
    <!-- 基础使用 -->
    <Icon :icon="User" />
    
    <!-- 自定义尺寸和颜色 -->
    <Icon :icon="Lock" :size="24" color="#00FF00" />
    
    <!-- 加载动画 -->
    <Icon :icon="Loading" :is-loading="true" :size="40" />
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@/components/common/base'
import { User, Lock, Loading } from '@element-plus/icons-vue'
</script>
```

