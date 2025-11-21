<template>
  <el-avatar
    :size="size"
    :shape="shape"
    :src="src"
    :icon="icon"
    :fit="fit"
    :alt="alt"
  >
    <slot>{{ text }}</slot>
  </el-avatar>
</template>

<script setup lang="ts">
import type { AvatarProps } from './types'

const props = withDefaults(defineProps<AvatarProps>(), {
  size: 'default',
  shape: 'circle',
  fit: 'cover',
})

defineOptions({
  name: 'Avatar',
})
</script>

