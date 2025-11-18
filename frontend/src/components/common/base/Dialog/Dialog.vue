<template>
  <el-dialog
    :model-value="modelValue"
    :title="title"
    :width="width"
    :fullscreen="fullscreen"
    :top="top"
    :modal="modal"
    :close-on-click-modal="closeOnClickModal"
    :close-on-press-escape="closeOnPressEscape"
    :show-close="showClose"
    :append-to-body="appendToBody"
    :lock-scroll="lockScroll"
    v-bind="$attrs"
    @update:model-value="handleUpdate"
    @open="handleOpen"
    @close="handleClose"
    @opened="handleOpened"
    @closed="handleClosed"
  >
    <template v-if="$slots.header || title" #header>
      <slot name="header">
        <span v-if="title">{{ title }}</span>
      </slot>
    </template>
    <slot />
    <template v-if="$slots.footer" #footer>
      <slot name="footer" />
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import type { DialogProps, DialogEmits } from './types'

const props = withDefaults(defineProps<DialogProps>(), {
  width: '50%',
  fullscreen: false,
  top: '15vh',
  modal: true,
  closeOnClickModal: true,
  closeOnPressEscape: true,
  showClose: true,
  appendToBody: false,
  lockScroll: true,
})

const emit = defineEmits<DialogEmits>()

const handleUpdate = (value: boolean) => {
  emit('update:modelValue', value)
}

const handleOpen = () => {
  emit('open')
}

const handleClose = () => {
  emit('close')
}

const handleOpened = () => {
  emit('opened')
}

const handleClosed = () => {
  emit('closed')
}
</script>

