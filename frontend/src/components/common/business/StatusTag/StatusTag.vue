<template>
  <el-tag
    :type="statusConfig?.type || 'info'"
    :color="statusConfig?.color"
    :effect="effect"
    :size="size"
    :round="round"
  >
    <slot>
      {{ statusConfig?.label || status }}
    </slot>
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { StatusTagProps, StatusConfig } from './types'

const props = withDefaults(defineProps<StatusTagProps & {
  effect?: 'dark' | 'light' | 'plain'
  size?: 'large' | 'default' | 'small'
  round?: boolean
}>(), {
  effect: 'light',
  size: 'default',
  round: false,
})

// 默认状态映射（工作流相关）
const defaultStatusMap: Record<string, StatusConfig> = {
  // 工作流状态
  running: { label: '运行中', type: 'success' },
  stopped: { label: '已停止', type: 'info' },
  paused: { label: '已暂停', type: 'warning' },
  error: { label: '错误', type: 'danger' },
  pending: { label: '等待中', type: 'warning' },
  
  // 执行状态
  success: { label: '成功', type: 'success' },
  failed: { label: '失败', type: 'danger' },
  cancelled: { label: '已取消', type: 'info' },
  
  // 通用状态
  active: { label: '激活', type: 'success' },
  inactive: { label: '未激活', type: 'info' },
  enabled: { label: '已启用', type: 'success' },
  disabled: { label: '已禁用', type: 'info' },
  
  // 审核状态
  draft: { label: '草稿', type: 'info' },
  reviewing: { label: '审核中', type: 'warning' },
  approved: { label: '已通过', type: 'success' },
  rejected: { label: '已拒绝', type: 'danger' },
}

const statusConfig = computed(() => {
  const map = props.statusMap || defaultStatusMap
  return map[props.status] || { label: props.status, type: 'info' as const }
})
</script>

<style scoped lang="scss">
// 状态标签样式
</style>
