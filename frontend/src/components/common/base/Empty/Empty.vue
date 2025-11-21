<template>
  <el-empty
    :image="image"
    :image-size="imageSize"
    :description="description"
  >
    <template v-if="$slots.image" #image>
      <slot name="image" />
    </template>
    <template v-if="$slots.description" #description>
      <slot name="description" />
    </template>
    <template v-if="$slots.default" #default>
      <slot />
    </template>
  </el-empty>
</template>

<script setup lang="ts">
import type { EmptyProps } from './types'

const props = withDefaults(defineProps<EmptyProps>(), {
  description: '暂无数据',
  imageSize: 200,
})
</script>

<style scoped lang="scss">
@use '@/assets/styles/glassmorphism' as *;

// 增强 Empty 组件的赛博朋克风格
:deep(.el-empty) {
  padding: 40px 20px;
  position: relative;
  
  // 背景网格装饰
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background-image: 
      repeating-linear-gradient(
        0deg,
        transparent 0px,
        transparent 19px,
        rgba(0, 212, 255, 0.03) 20px,
        rgba(0, 212, 255, 0.03) 21px
      ),
      repeating-linear-gradient(
        90deg,
        transparent 0px,
        transparent 19px,
        rgba(0, 212, 255, 0.03) 20px,
        rgba(0, 212, 255, 0.03) 21px
      );
    pointer-events: none;
    z-index: 0;
  }
  
  // 光效装饰
  &::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 200px;
    height: 200px;
    background: radial-gradient(circle, rgba(0, 212, 255, 0.1) 0%, transparent 70%);
    pointer-events: none;
    z-index: 0;
    animation: emptyGlow 3s ease-in-out infinite;
  }
  
  .el-empty__image {
    position: relative;
    z-index: 1;
    filter: drop-shadow(0 0 20px rgba(0, 212, 255, 0.3));
  }
  
  .el-empty__description {
    position: relative;
    z-index: 1;
    color: rgba(255, 255, 255, 0.6);
    font-size: 14px;
    text-shadow: 0 0 10px rgba(0, 212, 255, 0.3);
  }
}

@keyframes emptyGlow {
  0%, 100% {
    opacity: 0.3;
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    opacity: 0.6;
    transform: translate(-50%, -50%) scale(1.2);
  }
}
</style>

