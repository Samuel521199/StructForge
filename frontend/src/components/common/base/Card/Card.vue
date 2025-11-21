<template>
  <el-card
    :shadow="shadow"
    :body-style="bodyStyle"
    v-bind="$attrs"
  >
    <template v-if="header || $slots.header" #header>
      <slot name="header">
        <span v-if="header">{{ header }}</span>
      </slot>
    </template>
    <slot />
  </el-card>
</template>

<script setup lang="ts">
import type { CardProps } from './types'

withDefaults(defineProps<CardProps>(), {
  shadow: 'always',
  bodyStyle: () => ({}),
})
</script>

<style scoped lang="scss">
@use '@/assets/styles/glassmorphism' as *;

  // 覆盖 Element Plus Card 样式，实现玻璃拟物效果
  :deep(.el-card) {
    @include glassmorphism(25px, 180%, $glass-bg-dark);
    @include neon-border($neon-cyan, 0.4);
    @include soft-3d-shadow($neon-cyan, 0.15);
    border-radius: 16px;
    position: relative;
    overflow: hidden;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    
    // 内部光晕层
    &::before {
      content: '';
      position: absolute;
      inset: 0;
      background: radial-gradient(ellipse at top, rgba(0, 212, 255, 0.15) 0%, transparent 60%);
      pointer-events: none;
      z-index: 0;
    }
    
    // 赛博朋克装饰角（左上角）
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      width: 50px;
      height: 50px;
      border-top: 3px solid $neon-cyan;
      border-left: 3px solid $neon-cyan;
      border-radius: 16px 0 0 0;
      box-shadow: 
        -3px -3px 15px rgba(0, 212, 255, 0.6),
        -3px -3px 30px rgba(0, 212, 255, 0.4);
      pointer-events: none;
      z-index: 1;
      opacity: 0.8;
      animation: cornerGlow 2s ease-in-out infinite;
    }
    
    // 悬浮效果
    &:hover {
      transform: translateY(-6px);
      box-shadow: 
        0 15px 50px rgba(0, 0, 0, 0.6),
        0 0 50px rgba(0, 212, 255, 0.3) inset,
        0 0 100px rgba(0, 212, 255, 0.15);
      border-color: rgba(0, 212, 255, 0.7);
      
      &::after {
        opacity: 1;
        box-shadow: 
          -3px -3px 20px rgba(0, 212, 255, 0.8),
          -3px -3px 40px rgba(0, 212, 255, 0.6);
      }
    }
  
  // Card Header 样式
  .el-card__header {
    background: rgba(0, 212, 255, 0.05);
    border-bottom: 1px solid rgba(0, 212, 255, 0.2);
    color: rgba(0, 212, 255, 0.9);
    font-weight: 600;
    padding: 16px 20px;
    position: relative;
    z-index: 1;
  }
  
  // Card Body 样式
  .el-card__body {
    position: relative;
    z-index: 1;
    color: rgba(255, 255, 255, 0.85);
  }
}

@keyframes cornerGlow {
  0%, 100% {
    opacity: 0.6;
    box-shadow: 
      -3px -3px 15px rgba(0, 212, 255, 0.6),
      -3px -3px 30px rgba(0, 212, 255, 0.4);
  }
  50% {
    opacity: 1;
    box-shadow: 
      -3px -3px 25px rgba(0, 212, 255, 0.9),
      -3px -3px 50px rgba(0, 212, 255, 0.6);
  }
}
</style>

