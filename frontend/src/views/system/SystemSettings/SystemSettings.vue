<template>
  <div class="system-settings-page">
    <Card>
      <template #header>
        <div class="page-header">
          <h2>系统设置</h2>
        </div>
      </template>

      <div class="page-content">
        <p>系统设置功能开发中...</p>
        <!-- TODO: 实现系统设置功能 -->
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { Card } from '@/components/common/base'

// TODO: 实现系统设置功能
</script>

<style scoped lang="scss">
.system-settings-page {
  padding: 24px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-content {
  padding: 24px;
}
</style>

