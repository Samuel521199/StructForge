/**
 * API 错误处理工具
 */

import type { AxiosError } from 'axios'
import { ErrorType, type ApiError } from './types'
import { error as showError } from '@/components/common/base/Message'

/**
 * 判断是否为网络错误
 */
function isNetworkError(error: AxiosError): boolean {
  return !error.response && (error.code === 'ERR_NETWORK' || error.message.includes('Network Error'))
}

/**
 * 判断是否为超时错误
 */
function isTimeoutError(error: AxiosError): boolean {
  return error.code === 'ECONNABORTED' || error.message.includes('timeout')
}

/**
 * 判断是否为业务错误（4xx）
 */
function isBusinessError(status: number): boolean {
  return status >= 400 && status < 500
}

/**
 * 判断是否为服务器错误（5xx）
 */
function isServerError(status: number): boolean {
  return status >= 500 && status < 600
}

/**
 * 判断是否为认证错误
 */
function isAuthError(status: number): boolean {
  return status === 401
}

/**
 * 判断是否为权限错误
 */
function isPermissionError(status: number): boolean {
  return status === 403
}

/**
 * 判断是否为验证错误
 */
function isValidationError(status: number): boolean {
  return status === 422
}

/**
 * 分类错误类型
 */
export function classifyError(error: AxiosError): ErrorType {
  if (isNetworkError(error)) {
    return ErrorType.NETWORK
  }

  if (isTimeoutError(error)) {
    return ErrorType.TIMEOUT
  }

  const status = error.response?.status

  if (!status) {
    return ErrorType.UNKNOWN
  }

  if (isAuthError(status)) {
    return ErrorType.AUTH
  }

  if (isPermissionError(status)) {
    return ErrorType.PERMISSION
  }

  if (isValidationError(status)) {
    return ErrorType.VALIDATION
  }

  if (isBusinessError(status)) {
    return ErrorType.BUSINESS
  }

  if (isServerError(status)) {
    return ErrorType.SERVER
  }

  return ErrorType.UNKNOWN
}

/**
 * 格式化错误信息
 */
export function formatError(error: AxiosError): ApiError {
  const errorType = classifyError(error)
  const status = error.response?.status || 0
  const responseData = error.response?.data as { message?: string; code?: number; traceId?: string } | undefined

  const apiError: ApiError = {
    code: responseData?.code || status || -1,
    message: responseData?.message || error.message || '未知错误',
    type: errorType,
    traceId: responseData?.traceId,
  }

  // 添加详细信息
  if (error.response?.data && typeof error.response.data === 'object') {
    apiError.details = error.response.data as Record<string, unknown>
  }

  return apiError
}

/**
 * 获取用户友好的错误消息
 */
export function getUserFriendlyMessage(error: ApiError): string {
  switch (error.type) {
    case ErrorType.NETWORK:
      return '无法连接到服务器，请确保后端服务正在运行（http://localhost:8000）'
    case ErrorType.TIMEOUT:
      return '请求超时，请稍后重试'
    case ErrorType.AUTH:
      return '登录已过期，请重新登录'
    case ErrorType.PERMISSION:
      return '您没有权限执行此操作'
    case ErrorType.VALIDATION:
      return error.message || '数据验证失败，请检查输入'
    case ErrorType.SERVER:
      return '服务器错误，请稍后重试'
    case ErrorType.BUSINESS:
      return error.message || '操作失败'
    default:
      return error.message || '未知错误，请稍后重试'
  }
}

/**
 * 处理错误并显示提示
 */
export function handleError(error: AxiosError, showMessage = true): ApiError {
  const apiError = formatError(error)

  if (showMessage) {
    const message = getUserFriendlyMessage(apiError)
    showError(message)

    // 开发环境下显示详细错误信息
    if (import.meta.env.DEV) {
      console.error('API Error:', {
        type: apiError.type,
        code: apiError.code,
        message: apiError.message,
        details: apiError.details,
        traceId: apiError.traceId,
        originalError: error,
      })
    }
  }

  return apiError
}

