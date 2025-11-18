/**
 * UI 状态管理
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUIStore = defineStore('ui', () => {
  // 状态
  const sidebarCollapsed = ref(false)
  const theme = ref<'light' | 'dark'>('light')
  const loading = ref(false)

  // Actions
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setSidebarCollapsed(collapsed: boolean) {
    sidebarCollapsed.value = collapsed
  }

  function setTheme(newTheme: 'light' | 'dark') {
    theme.value = newTheme
    document.documentElement.setAttribute('data-theme', newTheme)
  }

  function setLoading(value: boolean) {
    loading.value = value
  }

  return {
    // State
    sidebarCollapsed,
    theme,
    loading,
    // Actions
    toggleSidebar,
    setSidebarCollapsed,
    setTheme,
    setLoading,
  }
})

