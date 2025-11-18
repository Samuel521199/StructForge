/**
 * 主题组合函数
 */

import { computed } from 'vue'
import { useUIStore } from '@/stores/modules/ui.store'

export function useTheme() {
  const uiStore = useUIStore()

  const theme = computed(() => uiStore.theme)

  const setTheme = (newTheme: 'light' | 'dark') => {
    uiStore.setTheme(newTheme)
  }

  const toggleTheme = () => {
    setTheme(theme.value === 'light' ? 'dark' : 'light')
  }

  return {
    theme,
    setTheme,
    toggleTheme,
  }
}

