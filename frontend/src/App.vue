<template>
  <router-view v-if="isRouterReady" />
  <div v-else class="loading-container">
    <p>正在加载应用...</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isRouterReady = ref(false)

// 立即尝试设置就绪状态
router.isReady().then(() => {
  console.log('路由已就绪，当前路径:', router.currentRoute.value.path)
  isRouterReady.value = true
}).catch((error) => {
  console.error('路由初始化失败:', error)
  // 即使失败也显示，避免白屏
  isRouterReady.value = true
})

onMounted(() => {
  // 再次确认路由就绪
  router.isReady().then(() => {
    console.log('路由已就绪（onMounted），当前路径:', router.currentRoute.value.path)
    isRouterReady.value = true
  }).catch((error) => {
    console.error('路由初始化失败（onMounted）:', error)
    isRouterReady.value = true
  })
  
  // 设置超时，避免无限等待
  setTimeout(() => {
    if (!isRouterReady.value) {
      console.warn('路由就绪超时，强制显示')
      isRouterReady.value = true
    }
  }, 1000) // 缩短超时时间到 1 秒
})

// 监听路由变化，确保路由视图更新
watch(() => router.currentRoute.value.path, (newPath) => {
  console.log('路由路径变化:', newPath)
  if (!isRouterReady.value) {
    isRouterReady.value = true
  }
}, { immediate: true })
</script>

<style lang="scss">
@use '@/assets/styles/glassmorphism' as *;

// 全局背景样式
body, #app {
  @include global-background;
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  color: rgba(255, 255, 255, 0.9);
}

// 加载容器样式
.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  font-family: Arial, sans-serif;
  color: rgba(255, 255, 255, 0.7);
  position: relative;
  z-index: 1;
}
</style>

