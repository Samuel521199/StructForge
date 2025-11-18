<template>
  <aside class="app-sidebar">
    <!-- TODO: 实现应用侧边栏 -->
    <slot />
  </aside>
</template>

<script setup lang="ts">
// TODO: 实现侧边栏逻辑
</script>

<style scoped lang="scss">
.app-sidebar {
  width: 240px;
  background-color: #fff;
  border-right: 1px solid #e4e7ed;
  overflow-y: auto;
}
</style>

