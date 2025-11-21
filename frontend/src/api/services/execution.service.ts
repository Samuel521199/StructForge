/**
 * 工作流执行服务 API
 */

import { apiClient } from '../client'
import type { ApiResponse } from '../client/types'
import type {
  Execution,
  ExecutionLog,
  ExecutionListParams,
  ExecutionListResponse,
} from '../types/execution.types'

export const executionService = {
  /**
   * 获取执行列表
   */
  getExecutions(params?: ExecutionListParams): Promise<ApiResponse<ExecutionListResponse>> {
    return apiClient.get('/v1/executions', { params })
  },

  /**
   * 获取执行详情
   */
  getExecution(id: string): Promise<ApiResponse<Execution>> {
    return apiClient.get(`/v1/executions/${id}`)
  },

  /**
   * 获取执行日志
   */
  getExecutionLogs(executionId: string, nodeId?: string): Promise<ApiResponse<ExecutionLog[]>> {
    return apiClient.get(`/v1/executions/${executionId}/logs`, { params: { nodeId } })
  },

  /**
   * 取消执行
   */
  cancelExecution(id: string): Promise<ApiResponse<void>> {
    return apiClient.post(`/v1/executions/${id}/cancel`)
  },

  /**
   * 重试执行
   */
  retryExecution(id: string): Promise<ApiResponse<Execution>> {
    return apiClient.post(`/v1/executions/${id}/retry`)
  },
}

