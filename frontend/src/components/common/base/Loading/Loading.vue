<template>
  <el-loading
    v-loading="loading"
    :text="text"
    :background="background"
    :spinner="spinner"
    :element="element"
    :target="target"
    :body="body"
    :fullscreen="fullscreen"
    :lock="lock"
  >
    <slot />
  </el-loading>
</template>

<script setup lang="ts">
import type { LoadingProps } from './types'

const props = withDefaults(defineProps<LoadingProps>(), {
  loading: false,
  text: '加载中...',
  background: 'rgba(0, 0, 0, 0.7)',
  fullscreen: false,
  lock: false,
})
</script>

