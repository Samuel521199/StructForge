/**
 * Message 消息提示组件
 * 通过方法调用显示消息提示
 */

import { ElMessage, ElMessageBox } from 'element-plus'
import type { MessageOptions, MessageType } from './types'

/**
 * 显示消息提示
 */
export function showMessage(
  message: string,
  type: MessageType = 'info',
  options?: MessageOptions
) {
  return ElMessage({
    message,
    type,
    duration: options?.duration ?? 3000,
    showClose: options?.showClose ?? false,
    center: options?.center ?? false,
    grouping: options?.grouping ?? false,
    dangerouslyUseHTMLString: options?.dangerouslyUseHTMLString ?? false,
    ...options,
  })
}

/**
 * 成功消息
 */
export function success(message: string, options?: MessageOptions) {
  return showMessage(message, 'success', options)
}

/**
 * 警告消息
 */
export function warning(message: string, options?: MessageOptions) {
  return showMessage(message, 'warning', options)
}

/**
 * 错误消息
 */
export function error(message: string, options?: MessageOptions) {
  return showMessage(message, 'error', options)
}

/**
 * 信息消息
 */
export function info(message: string, options?: MessageOptions) {
  return showMessage(message, 'info', options)
}

/**
 * 关闭所有消息
 */
export function closeAll() {
  ElMessage.closeAll()
}

/**
 * 使用组合函数方式
 */
export function useMessage() {
  return {
    success,
    warning,
    error,
    info,
    closeAll,
  }
}

// 默认导出
export default {
  success,
  warning,
  error,
  info,
  closeAll,
  useMessage,
}

