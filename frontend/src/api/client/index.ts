/**
 * HTTP 客户端配置
 * 创建 Axios 实例，配置请求/响应拦截器
 */

import axios, { type AxiosInstance, type InternalAxiosRequestConfig, type AxiosError } from 'axios'
import { useAuthStore } from '@/stores/modules/auth.store'
import { handleError } from './errorHandler'
import type { ApiResponse } from './types'
import { ErrorType } from './types'

// 创建 Axios 实例
const createAxiosInstance = (): AxiosInstance => {
  const instance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api',
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  // 请求拦截器
  instance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const authStore = useAuthStore()
      
      // 添加认证 token
      if (authStore.token && config.headers) {
        config.headers.Authorization = `Bearer ${authStore.token}`
      }

      // 添加请求 ID（用于追踪）
      if (config.headers) {
        config.headers['X-Request-ID'] = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
      }

      return config
    },
    (error) => {
      return Promise.reject(error)
    }
  )

  // 响应拦截器
  instance.interceptors.response.use(
    (response) => {
      // Kratos HTTP Gateway 返回的格式可能是：
      // 1. 直接是业务数据: { success: true, token: "...", user: {...} }
      // 2. 包装在标准响应中: { code: 200, message: "...", data: {...} }
      
      const data = response.data
      
      // 如果已经是标准 ApiResponse 格式（有 code 和 data 字段）
      if (data && typeof data === 'object' && 'code' in data && 'data' in data) {
        // 返回完整的 AxiosResponse，但替换 data 字段
        response.data = data as ApiResponse
        return response
      }
      
      // 如果是业务数据（如登录响应: { success: true, token: "...", user: {...} }）
      // 包装成标准格式
      response.data = { 
        code: response.status || 200, 
        message: (data as any)?.message || 'success', 
        data: data 
      } as ApiResponse
      
      return response
    },
    (error: AxiosError) => {
      const apiError = handleError(error, false) // 先不显示消息，由调用方决定

      // 处理认证错误
      if (apiError.type === ErrorType.AUTH) {
        const authStore = useAuthStore()
        authStore.logout()
        // 跳转到登录页（使用 router 而不是 window.location 避免页面刷新）
        // 注意：这里不能直接使用 router，需要通过动态导入
        if (window.location.pathname !== '/auth/login') {
          // 使用 replace 避免在历史记录中留下记录
          window.location.replace('/auth/login')
        }
      }
      
      // 处理网络连接错误
      if (apiError.type === ErrorType.NETWORK) {
        const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api'
        console.error('网络连接失败:', {
          baseURL,
          error: error.message,
          code: error.code,
        })
        console.error('请检查:')
        console.error('1. 后端 Gateway 服务是否运行在 http://localhost:8000')
        console.error('2. 网络连接是否正常')
        console.error('3. 防火墙是否阻止了连接')
        
        // 在开发环境下显示更友好的错误提示
        if (import.meta.env.DEV) {
          // 延迟显示，避免在页面加载时立即弹出
          setTimeout(async () => {
            try {
              const { error: showErrorMsg } = await import('@/components/common/base/Message')
              showErrorMsg('无法连接到服务器，请确保后端服务正在运行')
            } catch (importError) {
              console.error('无法导入 Message 组件:', importError)
            }
          }, 1000)
        }
      }

      return Promise.reject(apiError)
    }
  )

  return instance
}

export const apiClient = createAxiosInstance()

