/**
 * HTTP 客户端类型定义
 */

import type { AxiosRequestConfig, AxiosResponse } from 'axios'

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface ApiError {
  code: number
  message: string
  details?: any
}

export type RequestConfig = AxiosRequestConfig
export type Response<T = any> = AxiosResponse<ApiResponse<T>>

