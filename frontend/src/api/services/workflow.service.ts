/**
 * 工作流服务 API
 */

import { apiClient } from '../client'
import type { ApiResponse } from '../client/types'
import type {
  Workflow,
  Node,
  Edge,
  WorkflowListParams,
  WorkflowListResponse,
} from '../types/workflow.types'

export interface CreateWorkflowRequest {
  name: string
  description?: string
  nodes?: Node[]
  edges?: Edge[]
}

export interface UpdateWorkflowRequest {
  name?: string
  description?: string
  nodes?: Node[]
  edges?: Edge[]
  status?: Workflow['status']
}

export interface ExecuteWorkflowRequest {
  params?: Record<string, unknown>
}

export interface ExecuteWorkflowResponse {
  executionId: string
}

export const workflowService = {
  /**
   * 获取工作流列表
   */
  getWorkflows(params?: WorkflowListParams): Promise<ApiResponse<WorkflowListResponse>> {
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
  executeWorkflow(id: string, data?: ExecuteWorkflowRequest): Promise<ApiResponse<ExecuteWorkflowResponse>> {
    return apiClient.post(`/v1/workflows/${id}/execute`, data)
  },
}

