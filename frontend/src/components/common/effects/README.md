# 通用效果组件库

本目录用于存放所有特殊视觉效果组件，如动画、粒子效果、背景特效等。

## 📁 目录结构

```
effects/
├── CodeRain/          # 代码雨效果（黑客帝国风格）
│   ├── CodeRain.vue   # 组件实现
│   ├── types.ts       # 类型定义
│   ├── index.ts       # 导出文件
│   └── README.md      # 组件说明（可选）
├── index.ts           # 统一导出
└── README.md          # 本文件
```

## 🎯 使用规范

### 1. 组件命名

- **目录名**: PascalCase，描述性名称
  - 示例: `CodeRain`, `ParticleEffect`, `GlitchEffect`
- **组件文件名**: 与目录名相同，使用 `.vue` 扩展名
  - 示例: `CodeRain.vue`

### 2. 文件结构

每个效果组件应该包含以下文件：

```
EffectName/
├── EffectName.vue     # 组件实现（必需）
├── types.ts           # TypeScript 类型定义（必需）
├── index.ts           # 导出文件（必需）
└── README.md          # 组件文档（推荐）
```

### 3. 导出规范

**types.ts** 应包含：
- `EffectNameProps`: 组件属性接口
- `EffectNameEmits`: 组件事件接口（如果有）

**index.ts** 应包含：
```typescript
export { default as EffectName } from './EffectName.vue'
export type { EffectNameProps, EffectNameEmits } from './types'
```

**effects/index.ts** 应包含：
```typescript
export { EffectName } from './EffectName'
export type { EffectNameProps, EffectNameEmits } from './EffectName/types'
```

### 4. 组件设计原则

- **性能优先**: 使用 Canvas 或 WebGL 实现高性能动画
- **可配置**: 提供丰富的配置选项，支持自定义
- **响应式**: 自动适应容器大小变化
- **资源管理**: 正确清理定时器、事件监听器等资源
- **无障碍**: 考虑可访问性，提供禁用选项

### 5. 使用示例

```vue
<template>
  <div class="container">
    <CodeRain
      :color="'#00ff41'"
      :speed="2"
      :density="0.015"
    />
    <!-- 其他内容 -->
  </div>
</template>

<script setup lang="ts">
import { CodeRain } from '@/components/common/effects'
</script>

<style scoped>
.container {
  position: relative;
  width: 100%;
  height: 100vh;
}
</style>
```

## 📝 添加新效果组件

1. 在 `effects/` 目录下创建新的子目录
2. 按照文件结构创建必要文件
3. 在 `effects/index.ts` 中添加导出
4. 编写组件文档（README.md）

## 🎨 现有效果组件

### CodeRain

代码雨效果，类似《黑客帝国》中的绿色代码雨。

**特性**:
- 可自定义颜色、速度、密度
- 支持自定义字符集
- 自动适应容器大小
- 高性能 Canvas 渲染

**使用**:
```vue
<CodeRain
  :color="'#00ff41'"
  :speed="2"
  :density="0.015"
  :opacity="0.8"
/>
```

## ⚠️ 注意事项

1. **性能考虑**: 效果组件通常涉及动画，要注意性能影响
2. **资源清理**: 组件卸载时必须清理所有资源（定时器、事件监听器等）
3. **响应式**: 确保组件能正确响应容器大小变化
4. **可访问性**: 提供禁用选项，避免影响用户体验

