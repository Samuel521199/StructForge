<template>
  <el-select
    :model-value="modelValue"
    :placeholder="placeholder"
    :multiple="multiple"
    :disabled="disabled"
    :clearable="clearable"
    :filterable="filterable"
    :allow-create="allowCreate"
    :size="size"
    :loading="loading"
    v-bind="$attrs"
    @update:model-value="handleUpdate"
    @change="handleChange"
    @visible-change="handleVisibleChange"
    @remove-tag="handleRemoveTag"
    @clear="handleClear"
  >
    <template v-for="option in options" :key="option.value">
      <el-option-group v-if="option.children" :label="option.label">
        <el-option
          v-for="child in option.children"
          :key="child.value"
          :label="child.label"
          :value="child.value"
          :disabled="child.disabled"
        />
      </el-option-group>
      <el-option
        v-else
        :label="option.label"
        :value="option.value"
        :disabled="option.disabled"
      />
    </template>
  </el-select>
</template>

<script setup lang="ts">
import type { SelectProps, SelectEmits } from './types'

const props = withDefaults(defineProps<SelectProps>(), {
  placeholder: '请选择',
  multiple: false,
  disabled: false,
  clearable: false,
  filterable: false,
  allowCreate: false,
  size: 'default',
  loading: false,
})

const emit = defineEmits<SelectEmits>()

const handleUpdate = (value: string | number | Array<string | number>) => {
  emit('update:modelValue', value)
}

const handleChange = (value: string | number | Array<string | number>) => {
  emit('change', value)
}

const handleVisibleChange = (visible: boolean) => {
  emit('visible-change', visible)
}

const handleRemoveTag = (tag: string | number) => {
  emit('remove-tag', tag)
}

const handleClear = () => {
  emit('clear')
}
</script>

