/**
 * 抽屉组合函数
 */

import { ref } from 'vue'

export function useDrawer() {
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

