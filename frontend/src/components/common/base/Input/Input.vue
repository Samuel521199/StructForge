<template>
  <el-input
    :model-value="modelValue"
    :type="type"
    :placeholder="placeholder"
    :disabled="disabled"
    :readonly="readonly"
    :clearable="clearable"
    :show-password="showPassword"
    :prefix-icon="prefixIcon"
    :suffix-icon="suffixIcon"
    :maxlength="maxlength"
    :minlength="minlength"
    :show-word-limit="showWordLimit"
    :validate-event="validateEvent"
    :size="size"
    v-bind="$attrs"
    @update:model-value="handleUpdate"
    @focus="handleFocus"
    @blur="handleBlur"
    @clear="handleClear"
    @input="handleInput"
  >
    <template v-if="$slots.prefix" #prefix>
      <slot name="prefix" />
    </template>
    <template v-if="$slots.suffix" #suffix>
      <slot name="suffix" />
    </template>
    <template v-if="$slots.prepend" #prepend>
      <slot name="prepend" />
    </template>
    <template v-if="$slots.append" #append>
      <slot name="append" />
    </template>
  </el-input>
</template>

<script setup lang="ts">
import type { InputProps, InputEmits } from './types'

const props = withDefaults(defineProps<InputProps>(), {
  type: 'text',
  disabled: false,
  readonly: false,
  clearable: false,
  showPassword: false,
  showWordLimit: false,
  validateEvent: true,
  size: 'default',
})

const emit = defineEmits<InputEmits>()

const handleUpdate = (value: string | number) => {
  emit('update:modelValue', value)
}

const handleFocus = (event: FocusEvent) => {
  emit('focus', event)
}

const handleBlur = (event: FocusEvent) => {
  emit('blur', event)
}

const handleClear = () => {
  emit('clear')
}

const handleInput = (value: string | number) => {
  emit('input', value)
}
</script>

<style scoped lang="scss">
// 注意：这里不添加样式，样式由使用该组件的页面控制
// 如果需要组件级别的样式，可以在使用页面通过 :deep() 覆盖
</style>

