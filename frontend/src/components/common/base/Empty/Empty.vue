<template>
  <el-empty
    :image="image"
    :image-size="imageSize"
    :description="description"
  >
    <template v-if="$slots.image" #image>
      <slot name="image" />
    </template>
    <template v-if="$slots.description" #description>
      <slot name="description" />
    </template>
    <template v-if="$slots.default" #default>
      <slot />
    </template>
  </el-empty>
</template>

<script setup lang="ts">
import type { EmptyProps } from './types'

const props = withDefaults(defineProps<EmptyProps>(), {
  description: '暂无数据',
  imageSize: 200,
})
</script>

