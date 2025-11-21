/**
 * 仪表盘服务 API
 */

import { apiClient } from '../client'
import type { ApiResponse } from '../client/types'
import type {
  DashboardStats,
  ExecutionListItem,
  SuccessRateItem,
  ErrorTrendData,
  ExecutionDurationItem,
  WebhookActivity,
} from '../types/dashboard.types'

export interface ExecutionListParams {
  page?: number
  pageSize?: number
  status?: ExecutionListItem['status']
  workflowId?: string
}

export interface SuccessRateParams {
  period?: '24h' | '7d' | '30d'
}

export interface ErrorTrendParams {
  period?: '24h' | '7d' | '30d'
}

export const dashboardService = {
  /**
   * 获取仪表盘统计数据
   */
  getStats(): Promise<ApiResponse<DashboardStats>> {
    return apiClient.get('/v1/dashboard/stats')
  },

  /**
   * 获取执行历史列表
   */
  getExecutions(params?: ExecutionListParams): Promise<ApiResponse<{ list: ExecutionListItem[]; total: number }>> {
    return apiClient.get('/v1/dashboard/executions', { params })
  },

  /**
   * 获取执行成功率
   */
  getSuccessRate(params?: SuccessRateParams): Promise<ApiResponse<SuccessRateItem[]>> {
    return apiClient.get('/v1/dashboard/success-rate', { params })
  },

  /**
   * 获取错误趋势
   */
  getErrorTrend(params?: ErrorTrendParams): Promise<ApiResponse<ErrorTrendData[]>> {
    return apiClient.get('/v1/dashboard/error-trend', { params })
  },

  /**
   * 获取平均执行时长
   */
  getExecutionDuration(): Promise<ApiResponse<ExecutionDurationItem[]>> {
    return apiClient.get('/v1/dashboard/execution-duration')
  },

  /**
   * 获取 Webhook 活动
   */
  getWebhookActivity(limit?: number): Promise<ApiResponse<WebhookActivity[]>> {
    return apiClient.get('/v1/dashboard/webhook-activity', { params: { limit } })
  },
}

