<template>
  <el-tag
    :type="type"
    :size="size"
    :effect="effect"
    :round="round"
    :closable="closable"
    :disable-transitions="disableTransitions"
    :hit="hit"
    :color="color"
    @close="handleClose"
  >
    <slot />
  </el-tag>
</template>

<script setup lang="ts">
import type { TagProps, TagEmits } from './types'

const props = withDefaults(defineProps<TagProps>(), {
  type: 'default',
  size: 'default',
  effect: 'light',
  round: false,
  closable: false,
  disableTransitions: false,
  hit: false,
})

const emit = defineEmits<TagEmits>()

const handleClose = (event: MouseEvent) => {
  emit('close', event)
}
</script>

<style scoped lang="scss">
// Tag 样式由 Element Plus 提供
</style>

