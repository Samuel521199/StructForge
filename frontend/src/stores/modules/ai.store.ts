/**
 * AI 模型状态管理
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AIModel } from '@/api/services/ai.service'

export const useAIStore = defineStore('ai', () => {
  // 状态
  const models = ref<AIModel[]>([])
  const currentModel = ref<AIModel | null>(null)

  // Actions
  function setModels(list: AIModel[]) {
    models.value = list
  }

  function setCurrentModel(model: AIModel | null) {
    currentModel.value = model
  }

  return {
    // State
    models,
    currentModel,
    // Actions
    setModels,
    setCurrentModel,
  }
})

