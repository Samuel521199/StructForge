<template>
  <el-form
    ref="formRef"
    :model="model"
    :rules="rules"
    :label-width="labelWidth"
    :label-position="labelPosition"
    :size="size"
    :disabled="disabled"
    v-bind="$attrs"
    @validate="handleValidate"
    @submit="handleSubmit"
  >
    <slot />
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance } from 'element-plus'
import type { FormProps, FormEmits } from './types'

withDefaults(defineProps<FormProps>(), {
  labelWidth: '100px',
  labelPosition: 'right',
  size: 'default',
  disabled: false,
})

const emit = defineEmits<FormEmits>()

const formRef = ref<FormInstance>()

const handleValidate = (prop: string, isValid: boolean, message: string) => {
  emit('validate', prop, isValid, message)
}

const handleSubmit = (event: Event) => {
  event.preventDefault()
  emit('submit', event)
}

// 暴露表单方法
defineExpose({
  validate: (callback?: (valid: boolean, fields?: any) => void) => {
    return formRef.value?.validate(callback)
  },
  validateField: (props: string | string[], callback?: (valid: boolean, fields?: any) => void) => {
    return formRef.value?.validateField(props, callback)
  },
  resetFields: () => formRef.value?.resetFields(),
  clearValidate: (props?: string | string[]) => formRef.value?.clearValidate(props),
  scrollToField: (prop: string) => formRef.value?.scrollToField(prop),
})
</script>

