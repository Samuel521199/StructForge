/**
 * AI 服务 API
 */

import { apiClient } from '../client'
import type { ApiResponse } from '../client/types'

export interface AIModel {
  id: string
  name: string
  type: string
  provider: string
  status: 'active' | 'inactive'
  config: Record<string, any>
}

export const aiService = {
  /**
   * 获取 AI 模型列表
   */
  getModels(): Promise<ApiResponse<AIModel[]>> {
    return apiClient.get('/v1/ai/models')
  },

  /**
   * 获取 AI 模型详情
   */
  getModel(id: string): Promise<ApiResponse<AIModel>> {
    return apiClient.get(`/v1/ai/models/${id}`)
  },
}

