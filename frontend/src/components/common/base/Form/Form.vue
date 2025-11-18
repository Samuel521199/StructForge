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
  >
    <slot />
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance } from 'element-plus'
import type { FormProps, FormEmits } from './types'

const props = withDefaults(defineProps<FormProps>(), {
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

// 暴露表单方法
defineExpose({
  validate: () => formRef.value?.validate(),
  validateField: (props: string | string[]) => formRef.value?.validateField(props),
  resetFields: () => formRef.value?.resetFields(),
  clearValidate: (props?: string | string[]) => formRef.value?.clearValidate(props),
  scrollToField: (prop: string) => formRef.value?.scrollToField(prop),
})
</script>

