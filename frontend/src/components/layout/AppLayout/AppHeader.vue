<template>
  <header class="app-header">
    <!-- TODO: 实现应用头部 -->
    <div class="header-left">
      <slot name="left" />
    </div>
    <div class="header-center">
      <slot name="center" />
    </div>
    <div class="header-right">
      <slot name="right" />
    </div>
  </header>
</template>

<script setup lang="ts">
// TODO: 实现头部逻辑
</script>

<style scoped lang="scss">
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e4e7ed;
}
</style>

