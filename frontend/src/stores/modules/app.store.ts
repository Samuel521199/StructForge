/**
 * 应用全局状态管理
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 状态
  const appName = ref('StructForge')
  const appVersion = ref('1.0.0')
  const language = ref<'zh-CN' | 'en-US'>('zh-CN')

  // Actions
  function setLanguage(lang: 'zh-CN' | 'en-US') {
    language.value = lang
  }

  return {
    // State
    appName,
    appVersion,
    language,
    // Actions
    setLanguage,
  }
})

