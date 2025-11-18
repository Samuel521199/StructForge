<template>
  <el-form-item
    :label="label"
    :prop="prop"
    :required="required"
    :rules="rules"
    :error="error"
    :show-message="showMessage"
    :inline-message="inlineMessage"
    :size="size"
    :label-width="labelWidth"
    :label-position="labelPosition"
    v-bind="$attrs"
  >
    <slot />
  </el-form-item>
</template>

<script setup lang="ts">
import type { FormItemProps } from './types'

const props = withDefaults(defineProps<FormItemProps>(), {
  required: false,
  showMessage: true,
  inlineMessage: false,
  size: 'default',
})
</script>

<style scoped lang="scss">
/* FormItem 样式 */
</style>

