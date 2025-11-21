/**
 * HTTP 客户端类型定义
 */

import type { AxiosRequestConfig, AxiosResponse } from 'axios'

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  success?: boolean
  timestamp?: string
  traceId?: string
}

// 错误类型分类
export enum ErrorType {
  NETWORK = 'NETWORK',        // 网络错误
  TIMEOUT = 'TIMEOUT',        // 超时错误
  BUSINESS = 'BUSINESS',      // 业务错误
  VALIDATION = 'VALIDATION',  // 验证错误
  AUTH = 'AUTH',              // 认证错误
  PERMISSION = 'PERMISSION',  // 权限错误
  SERVER = 'SERVER',          // 服务器错误
  UNKNOWN = 'UNKNOWN',        // 未知错误
}

export interface ApiError {
  code: number
  message: string
  type?: ErrorType
  details?: Record<string, unknown>
  timestamp?: string
  traceId?: string
}

export type RequestConfig = AxiosRequestConfig
export type Response<T = unknown> = AxiosResponse<ApiResponse<T>>

