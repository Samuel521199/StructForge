/**
 * HTTP 错误处理
 */

import type { AxiosError } from 'axios'
import { useToast } from '@/composables/ui/useToast'

export function handleHttpError(error: AxiosError) {
  const toast = useToast()

  if (error.response) {
    const status = error.response.status
    const message = (error.response.data as any)?.message || '请求失败'

    switch (status) {
      case 400:
        toast.error(`请求参数错误: ${message}`)
        break
      case 401:
        toast.error('未授权，请重新登录')
        break
      case 403:
        toast.error('没有权限访问')
        break
      case 404:
        toast.error('资源不存在')
        break
      case 500:
        toast.error(`服务器错误: ${message}`)
        break
      default:
        toast.error(`请求失败: ${message}`)
    }
  } else if (error.request) {
    toast.error('网络错误，请检查网络连接')
  } else {
    toast.error(`请求错误: ${error.message}`)
  }
}

