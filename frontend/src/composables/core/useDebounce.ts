/**
 * 防抖组合函数
 */

import { ref, watch } from 'vue'

export function useDebounce<T>(value: T, delay = 300) {
  const debouncedValue = ref<T>(value)

  let timer: ReturnType<typeof setTimeout> | null = null

  watch(
    () => value,
    (newValue) => {
      if (timer) {
        clearTimeout(timer)
      }
      timer = setTimeout(() => {
        debouncedValue.value = newValue
      }, delay)
    },
    { immediate: true }
  )

  return debouncedValue
}

