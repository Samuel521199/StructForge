<template>
  <el-button
    :type="type"
    :size="size"
    :disabled="disabled"
    :loading="loading"
    :icon="typeof icon === 'string' ? icon : undefined"
    :round="round"
    :circle="circle"
    :plain="plain"
    :native-type="nativeType"
    @click="handleClick"
  >
    <template v-if="icon && typeof icon === 'object'" #icon>
      <component :is="icon" />
    </template>
    <slot />
  </el-button>
</template>

<script setup lang="ts">
import type { ButtonProps, ButtonEmits } from './types'

const props = withDefaults(defineProps<ButtonProps>(), {
  type: 'default',
  size: 'default',
  disabled: false,
  loading: false,
  round: false,
  circle: false,
  plain: false,
  nativeType: 'button',
})

const emit = defineEmits<ButtonEmits>()

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

