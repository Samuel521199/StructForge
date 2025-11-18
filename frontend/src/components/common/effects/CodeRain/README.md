# CodeRain 矩阵代码雨特效组件

影视级矩阵代码雨特效，支持深度分层、光晕尾迹、日文片假名和动态干扰效果。

## 特性

### 🎬 影视级效果

1. **三层深度系统（Depth & Layering）**
   - **远景层（Far Layer）**：较小、较淡、较慢，模拟远距离效果
   - **中景层（Mid Layer）**：正常大小和密度，标准霓虹绿
   - **近景层（Near Layer）**：较大、较亮、较快，突出显示

2. **光晕与尾迹（Glow & Trail）**
   - **发光头（Leader）**：每条代码串的第一个字符使用纯白色，带有强烈光晕
   - **渐变尾迹**：后续字符透明度指数衰减，形成流畅的拖尾效果
   - **Canvas 光晕**：使用 `shadowBlur` 和 `shadowColor` 实现霓虹灯发光效果

3. **字符与随机性（Characters & Randomness）**
   - **日文片假名**：默认使用日文片假名字符集（アイウエオ...），还原电影原版效果
   - **动态字符替换**：每帧有 10% 的概率随机替换部分字符，模拟数据流变化

4. **动态干扰与失真（Glitch & Distortion）**
   - **屏幕抖动**：周期性（每 5-10 秒）触发微小的随机位移（±1-2 像素），模拟数据流不稳定
   - **颜色脉冲**：主绿色调的亮度缓慢变化，模拟电流流淌效果

## 使用方法

### 基础用法

```vue
<template>
  <CodeRain />
</template>

<script setup lang="ts">
import { CodeRain } from '@/components/common/effects'
</script>
```

### 完整配置

```vue
<template>
  <CodeRain
    :color="'#00ff41'"
    :backgroundColor="'rgba(0, 0, 0, 0.08)'"
    :fontSize="14"
    :fontWeight="'bold'"
    :speed="2.5"
    :speedVariation="0.6"
    :density="0.006"
    :opacity="0.9"
    :fadeSpeed="0.04"
    :minLength="20"
    :maxLength="40"
    :enableLayers="true"
    :enableGlow="true"
    :enableGlitch="true"
    :glowIntensity="0.8"
  />
</template>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `fontSize` | `number` | `14` | 字体大小（像素） |
| `fontFamily` | `string` | `'monospace'` | 字体族 |
| `fontWeight` | `string \| number` | `'bold'` | 字体粗细 |
| `color` | `string` | `'#00ff41'` | 代码颜色（十六进制） |
| `backgroundColor` | `string` | `'#000000'` | 背景颜色（十六进制或 rgba） |
| `speed` | `number` | `2.5` | 下落速度（像素/帧） |
| `speedVariation` | `number` | `0.6` | 速度变化范围（0-1），值越大速度差异越大 |
| `density` | `number` | `0.008` | 雨滴密度（0-1），值越大密度越高 |
| `opacity` | `number` | `0.9` | 透明度（0-1） |
| `fadeSpeed` | `number` | `0.04` | 拖尾淡出速度（0-1），值越大淡出越快 |
| `characters` | `string` | 日文片假名 | 字符集，默认使用日文片假名 |
| `minLength` | `number` | `20` | 代码串最小长度 |
| `maxLength` | `number` | `40` | 代码串最大长度 |
| `enableLayers` | `boolean` | `true` | 是否启用三层深度效果 |
| `enableGlow` | `boolean` | `true` | 是否启用光晕效果 |
| `enableGlitch` | `boolean` | `true` | 是否启用动态干扰（抖动、脉冲） |
| `glowIntensity` | `number` | `0.8` | 光晕强度（0-1） |

## 性能优化

- 使用 Canvas API 实现高性能渲染
- 响应式设计，自动适应屏幕尺寸变化
- 合理的密度控制，避免过度渲染
- 使用 `requestAnimationFrame` 实现流畅动画

## 技术实现

### 核心技术栈

- **Vue 3** + **TypeScript**：组件框架
- **HTML Canvas API**：高性能 2D 渲染
- **requestAnimationFrame**：流畅动画循环

### 关键算法

1. **分层渲染**：按层级顺序绘制（远景 → 中景 → 近景），实现深度感
2. **光晕效果**：使用 Canvas `shadowBlur` 和 `shadowColor` 属性
3. **渐变尾迹**：指数衰减函数 `Math.pow(1 - position, 1.5)`
4. **动态干扰**：周期性触发随机位移和颜色脉冲

## 注意事项

1. **性能考虑**：密度（`density`）不宜设置过高，建议在 0.003-0.01 之间
2. **字符集**：默认使用日文片假名，如需自定义可通过 `characters` 属性设置
3. **响应式**：组件会自动适应容器大小变化，无需手动处理
4. **层级效果**：启用 `enableLayers` 会增加渲染负担，但视觉效果更佳

## 示例场景

### 登录页面背景

```vue
<template>
  <div class="login-page">
    <CodeRain
      :enableLayers="true"
      :enableGlow="true"
      :enableGlitch="true"
      :density="0.006"
    />
    <div class="login-container">
      <!-- 登录表单 -->
    </div>
  </div>
</template>

<style scoped>
.login-page {
  position: relative;
  min-height: 100vh;
  background: #000000;
}

.login-container {
  position: relative;
  z-index: 1;
}
</style>
```

## 更新日志

### v2.0.0 - 影视级特效版本

- ✨ 新增三层深度系统
- ✨ 新增光晕与尾迹效果
- ✨ 新增日文片假名字符集
- ✨ 新增动态干扰效果（抖动、脉冲）
- 🎨 优化视觉效果，更接近电影原版
- ⚡ 性能优化，支持更高密度渲染
