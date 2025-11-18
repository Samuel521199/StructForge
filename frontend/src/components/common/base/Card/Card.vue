<template>
  <el-card
    :shadow="shadow"
    :body-style="bodyStyle"
    v-bind="$attrs"
  >
    <template v-if="header || $slots.header" #header>
      <slot name="header">
        <span v-if="header">{{ header }}</span>
      </slot>
    </template>
    <slot />
  </el-card>
</template>

<script setup lang="ts">
import type { CardProps } from './types'

const props = withDefaults(defineProps<CardProps>(), {
  shadow: 'always',
  bodyStyle: () => ({}),
})
</script>

