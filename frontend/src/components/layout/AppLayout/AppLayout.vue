<template>
  <div class="app-layout">
    <AppHeader v-if="showHeader" />
    <div class="app-layout-body">
      <AppSidebar v-if="showSidebar" />
      <AppContent>
        <slot />
      </AppContent>
    </div>
    <AppFooter v-if="showFooter" />
  </div>
</template>

<script setup lang="ts">
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'
import AppContent from './AppContent.vue'
import AppFooter from './AppFooter.vue'

interface AppLayoutProps {
  showHeader?: boolean
  showSidebar?: boolean
  showFooter?: boolean
}

withDefaults(defineProps<AppLayoutProps>(), {
  showHeader: true,
  showSidebar: true,
  showFooter: false
})
</script>

<style scoped lang="scss">
@use '@/assets/styles/glassmorphism' as *;

.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  position: relative;
  z-index: 1;
  
  // 整体玻璃拟物容器
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background: rgba(20, 20, 30, 0.2);
    backdrop-filter: blur(10px) saturate(150%);
    -webkit-backdrop-filter: blur(10px) saturate(150%);
    z-index: -1;
    pointer-events: none;
  }
  
  .app-layout-body {
    display: flex;
    flex: 1;
    overflow: hidden;
    position: relative;
  }
}
</style>

