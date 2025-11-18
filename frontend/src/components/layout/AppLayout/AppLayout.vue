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
import { AppHeader } from './AppHeader.vue'
import { AppSidebar } from './AppSidebar.vue'
import { AppContent } from './AppContent.vue'
import { AppFooter } from './AppFooter.vue'

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
.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  
  .app-layout-body {
    display: flex;
    flex: 1;
    overflow: hidden;
  }
}
</style>

