<template>
  <el-icon
    :size="size"
    :color="color"
    :class="[isLoading && 'is-loading', customClass]"
    v-bind="$attrs"
  >
    <component :is="icon" v-if="icon" />
  </el-icon>
</template>

<script setup lang="ts">
import type { IconProps } from './types'
import { computed } from 'vue'

const props = withDefaults(defineProps<IconProps>(), {
  size: 16,
  isLoading: false,
})

const customClass = computed(() => {
  return props.class || ''
})
</script>

