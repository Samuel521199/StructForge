/**
 * 模态框组合函数
 */

import { ref } from 'vue'

export function useModal() {
  const visible = ref(false)

  const show = () => {
    visible.value = true
  }

  const hide = () => {
    visible.value = false
  }

  const toggle = () => {
    visible.value = !visible.value
  }

  return {
    visible,
    show,
    hide,
    toggle,
  }
}

