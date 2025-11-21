<template>
  <div class="sf-progress">
    <div class="progress-info" v-if="showInfo">
      <span class="progress-text">{{ text || `${percentage}%` }}</span>
      <span v-if="status" class="progress-status" :class="`status-${status}`">
        {{ statusText }}
      </span>
    </div>
    <div class="progress-bar" :class="`bar-${status || 'default'}`">
      <div
        class="progress-inner"
        :style="{
          width: `${percentage}%`,
          backgroundColor: color || getStatusColor(status),
        }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ProgressProps } from './types'

const props = withDefaults(defineProps<ProgressProps>(), {
  percentage: 0,
  showInfo: true,
  status: undefined,
})

const statusText = computed(() => {
  const statusMap: Record<string, string> = {
    success: '成功',
    exception: '异常',
    warning: '警告',
  }
  return statusMap[props.status || ''] || ''
})

const getStatusColor = (status?: string): string => {
  const colorMap: Record<string, string> = {
    success: '#67c23a',
    exception: '#f56c6c',
    warning: '#e6a23c',
  }
  return colorMap[status || ''] || '#409eff'
}
</script>

<style scoped lang="scss">
.sf-progress {
  .progress-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    font-size: 14px;

    .progress-text {
      color: var(--el-text-color-primary);
    }

    .progress-status {
      font-size: 12px;

      &.status-success {
        color: #67c23a;
      }

      &.status-exception {
        color: #f56c6c;
      }

      &.status-warning {
        color: #e6a23c;
      }
    }
  }

  .progress-bar {
    width: 100%;
    height: 8px;
    background-color: var(--el-fill-color-light);
    border-radius: 4px;
    overflow: hidden;

    .progress-inner {
      height: 100%;
      border-radius: 4px;
      transition: width 0.3s ease;
    }

    &.bar-success .progress-inner {
      background-color: #67c23a;
    }

    &.bar-exception .progress-inner {
      background-color: #f56c6c;
    }

    &.bar-warning .progress-inner {
      background-color: #e6a23c;
    }
  }
}
</style>

