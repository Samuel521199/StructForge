<template>
  <main class="app-content">
    <slot />
  </main>
</template>

<script setup lang="ts">
// TODO: 实现内容区逻辑
</script>

<style scoped lang="scss">
.app-content {
  flex: 1;
  overflow-y: auto;
  background-color: #f5f7fa;
  padding: 20px;
}
</style>

