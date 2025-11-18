/**
 * 节点服务 API
 */

import { apiClient } from '../client'
import type { ApiResponse } from '../client/types'

export interface NodeType {
  id: string
  name: string
  category: string
  icon?: string
  description?: string
  configSchema: any
}

export const nodeService = {
  /**
   * 获取节点类型列表
   */
  getNodeTypes(): Promise<ApiResponse<NodeType[]>> {
    return apiClient.get('/v1/nodes/types')
  },

  /**
   * 获取节点类型详情
   */
  getNodeType(id: string): Promise<ApiResponse<NodeType>> {
    return apiClient.get(`/v1/nodes/types/${id}`)
  },
}

