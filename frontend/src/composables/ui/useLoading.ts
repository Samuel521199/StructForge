/**
 * 加载状态组合函数
 */

import { ref } from 'vue'
import { useUIStore } from '@/stores/modules/ui.store'

export function useLoading() {
  const uiStore = useUIStore()
  const loading = ref(false)

  const show = () => {
    loading.value = true
    uiStore.setLoading(true)
  }

  const hide = () => {
    loading.value = false
    uiStore.setLoading(false)
  }

  return {
    loading,
    show,
    hide,
  }
}

