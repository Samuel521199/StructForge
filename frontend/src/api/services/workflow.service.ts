/**
 * 工作流服务 API
 */

import { apiClient } from '../client'
import type { ApiResponse } from '../client/types'

export interface Workflow {
  id: string
  name: string
  description?: string
  nodes: any[]
  edges: any[]
  status: 'draft' | 'active' | 'archived'
  createdAt: string
  updatedAt: string
}

export interface CreateWorkflowRequest {
  name: string
  description?: string
  nodes?: any[]
  edges?: any[]
}

export interface UpdateWorkflowRequest {
  name?: string
  description?: string
  nodes?: any[]
  edges?: any[]
  status?: string
}

export const workflowService = {
  /**
   * 获取工作流列表
   */
  getWorkflows(params?: { page?: number; pageSize?: number; search?: string }): Promise<ApiResponse<{ list: Workflow[]; total: number }>> {
    return apiClient.get('/v1/workflows', { params })
  },

  /**
   * 获取工作流详情
   */
  getWorkflow(id: string): Promise<ApiResponse<Workflow>> {
    return apiClient.get(`/v1/workflows/${id}`)
  },

  /**
   * 创建工作流
   */
  createWorkflow(data: CreateWorkflowRequest): Promise<ApiResponse<Workflow>> {
    return apiClient.post('/v1/workflows', data)
  },

  /**
   * 更新工作流
   */
  updateWorkflow(id: string, data: UpdateWorkflowRequest): Promise<ApiResponse<Workflow>> {
    return apiClient.put(`/v1/workflows/${id}`, data)
  },

  /**
   * 删除工作流
   */
  deleteWorkflow(id: string): Promise<ApiResponse<void>> {
    return apiClient.delete(`/v1/workflows/${id}`)
  },

  /**
   * 执行工作流
   */
  executeWorkflow(id: string, params?: Record<string, any>): Promise<ApiResponse<{ executionId: string }>> {
    return apiClient.post(`/v1/workflows/${id}/execute`, params)
  },
}

