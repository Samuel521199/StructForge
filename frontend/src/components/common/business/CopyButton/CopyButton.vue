<template>
  <Button
    :type="type"
    :size="size"
    :icon="copied ? SuccessFilled : CopyDocument"
    :loading="loading"
    :disabled="disabled"
    :round="round"
    :plain="plain"
    @click="handleCopy"
  >
    <slot>
      {{ copied ? successText : text }}
    </slot>
  </Button>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { CopyDocument, SuccessFilled } from '@element-plus/icons-vue'
import { Button } from '@/components/common/base'
import { success } from '@/components/common/base/Message'
import type { CopyButtonProps } from './types'

const props = withDefaults(defineProps<CopyButtonProps>(), {
  text: '复制',
  successText: '已复制',
  type: 'default',
  size: 'default',
  round: false,
  plain: false,
  showMessage: true,
})

const loading = ref(false)
const copied = ref(false)

const handleCopy = async () => {
  if (loading.value || !props.value) {
    return
  }

  try {
    loading.value = true
    copied.value = false

    // 复制到剪贴板
    const text = typeof props.value === 'string' 
      ? props.value 
      : JSON.stringify(props.value, null, 2)

    await navigator.clipboard.writeText(text)

    copied.value = true

    if (props.showMessage) {
      success(props.successText || '已复制')
    }

    // 2秒后恢复状态
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (error) {
    console.error('复制失败:', error)
    if (props.showMessage) {
      // 使用error消息提示
      const { error: showError } = await import('@/components/common/base/Message')
      showError('复制失败，请手动复制')
    }
    throw error
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
// CopyButton样式
</style>

