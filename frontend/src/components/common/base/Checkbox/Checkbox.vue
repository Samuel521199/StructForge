<template>
  <el-checkbox
    :model-value="modelValue"
    :label="label"
    :disabled="disabled"
    :readonly="readonly"
    :border="border"
    :size="size"
    :name="name"
    :id="id"
    :indeterminate="indeterminate"
    :checked="checked"
    :tabindex="tabindex"
    :validate-event="validateEvent"
    v-bind="$attrs"
    @update:model-value="handleUpdate"
    @change="handleChange"
  >
    <slot />
  </el-checkbox>
</template>

<script setup lang="ts">
import type { CheckboxProps, CheckboxEmits } from './types'

const props = withDefaults(defineProps<CheckboxProps>(), {
  disabled: false,
  readonly: false,
  border: false,
  size: 'default',
  indeterminate: false,
  validateEvent: true,
})

const emit = defineEmits<CheckboxEmits>()

const handleUpdate = (value: boolean | string | number) => {
  emit('update:modelValue', value)
}

const handleChange = (value: boolean | string | number) => {
  emit('change', value)
}
</script>

