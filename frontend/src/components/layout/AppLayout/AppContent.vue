<template>
  <main class="app-content">
    <div class="content-wrapper">
      <slot />
    </div>
  </main>
</template>

<script setup lang="ts">
// 内容区组件，提供统一的页面容器
</script>

<style scoped lang="scss">
@use '@/assets/styles/glassmorphism' as *;

.app-content {
  flex: 1;
  overflow-y: auto;
  @include glassmorphism(20px, 150%, rgba(15, 15, 25, 0.2));
  position: relative;
  
  // 内容区域光晕和网格
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background-image: 
      // 细网格
      repeating-linear-gradient(
        0deg,
        transparent 0px,
        transparent 19px,
        rgba(0, 212, 255, 0.02) 20px,
        rgba(0, 212, 255, 0.02) 21px
      ),
      repeating-linear-gradient(
        90deg,
        transparent 0px,
        transparent 19px,
        rgba(0, 212, 255, 0.02) 20px,
        rgba(0, 212, 255, 0.02) 21px
      ),
      // 光晕
      radial-gradient(ellipse at center, rgba(0, 212, 255, 0.05) 0%, transparent 70%);
    pointer-events: none;
    z-index: 0;
  }
  
  .content-wrapper {
    min-height: 100%;
    padding: 32px;
    position: relative;
    z-index: 1;
  }
  
  // 霓虹滚动条
  &::-webkit-scrollbar {
    width: 8px;
  }
  
  &::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.3);
    border-radius: 4px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: linear-gradient(180deg, 
      rgba(0, 212, 255, 0.6) 0%, 
      rgba(0, 255, 136, 0.6) 50%,
      rgba(183, 148, 246, 0.6) 100%
    );
    border-radius: 4px;
    box-shadow: 0 0 10px rgba(0, 212, 255, 0.4);
    
    &:hover {
      box-shadow: 0 0 15px rgba(0, 212, 255, 0.6);
    }
  }
}
</style>

