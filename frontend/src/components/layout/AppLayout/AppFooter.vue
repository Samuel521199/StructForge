<template>
  <footer class="app-footer">
    <div class="footer-content">
      <slot>
        <div class="footer-text">
          <span>© 2025 StructForge. All rights reserved.</span>
        </div>
      </slot>
    </div>
  </footer>
</template>

<script setup lang="ts">
// 应用底部组件
</script>

<style scoped lang="scss">
.app-footer {
  height: 50px;
  background-color: #fff;
  border-top: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  justify-content: center;

  .footer-content {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    padding: 0 24px;

    .footer-text {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }
}
</style>

