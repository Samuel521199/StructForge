<template>
  <div class="model-list-page">
    <Card>
      <template #header>
        <div class="page-header">
          <h2>AI模型列表</h2>
        </div>
      </template>

      <div class="page-content">
        <p>模型列表功能开发中...</p>
        <!-- TODO: 实现模型列表功能 -->
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { Card } from '@/components/common/base'

// TODO: 实现模型列表功能
</script>

<style scoped lang="scss">
.model-list-page {
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

