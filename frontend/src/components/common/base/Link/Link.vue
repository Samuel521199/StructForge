<template>
  <el-link
    :type="type"
    :underline="underline"
    :disabled="disabled"
    :href="href"
    :target="target"
    :icon="icon"
    v-bind="$attrs"
    @click="handleClick"
  >
    <slot />
  </el-link>
</template>

<script setup lang="ts">
import type { LinkProps, LinkEmits } from './types'

const props = withDefaults(defineProps<LinkProps>(), {
  type: 'default',
  underline: true,
  disabled: false,
})

const emit = defineEmits<LinkEmits>()

const handleClick = (event: MouseEvent) => {
  if (!props.disabled) {
    emit('click', event)
  }
}
</script>

