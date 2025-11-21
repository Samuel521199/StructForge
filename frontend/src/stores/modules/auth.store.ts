/**
 * 认证状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userService } from '@/api/services/user.service'
import { useUserStore } from './user.store'
import type { User } from '@/api/types/user.types'

export const useAuthStore = defineStore('auth', () => {
  const userStore = useUserStore()

  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const isAuthenticated = computed(() => {
    // 检查 token 和用户信息都存在
    const hasToken = !!token.value
    const hasUser = userStore.isLoggedIn
    return hasToken && hasUser
  })

  // Actions
  async function login(username: string, password: string) {
    try {
      const response = await userService.login({ username, password })
      console.log('Login response (raw):', JSON.stringify(response, null, 2))
      
      // apiClient 返回格式: ApiResponse<LoginResponse>
      // 格式: { code: 200, message: "success", data: LoginResponse }
      // LoginResponse 格式: { token: string, user: UserInfo }
      // 但后端实际返回: { success: true, message: "...", token: "...", user: {...} }
      // 经过拦截器包装后: { code: 200, message: "success", data: { success: true, token: "...", user: {...} } }
      
      const responseData = response as any
      
      // 获取实际的登录数据（可能在 data 字段中）
      let loginData: any = null
      if (responseData?.data) {
        // 如果 data 字段存在，检查是否是嵌套结构
        if (responseData.data.data && typeof responseData.data.data === 'object') {
          // 双重嵌套
          loginData = responseData.data.data
        } else if (responseData.data.token || responseData.data.success !== undefined) {
          // 单层嵌套，data 就是登录数据
          loginData = responseData.data
        } else {
          loginData = responseData.data
        }
      } else if (responseData?.token || responseData?.success !== undefined) {
        // 直接是登录数据（未包装）
        loginData = responseData
      }
      
      console.log('Parsed login data:', loginData)
      
      if (!loginData) {
        console.error('无法解析登录响应:', responseData)
        return false
      }
      
      // 检查登录是否成功（后端返回 success 字段）
      const isSuccess = loginData.success === true || loginData.Success === true
      if (!isSuccess) {
        const errorMsg = loginData.message || loginData.Message || '登录失败'
        console.warn('Login failed:', errorMsg)
        return false
      }
      
      // 获取 token
      const tokenValue = loginData.token || loginData.Token
      if (!tokenValue) {
        console.error('Login response missing token:', loginData)
        return false
      }
      
      // 保存 token
      token.value = tokenValue
      localStorage.setItem('token', tokenValue)
      console.log('Token saved')
      
      // 转换用户数据
      const userData = loginData.user || loginData.User
      if (userData) {
        const user = userData
        
        // 处理时间戳
        const formatTimestamp = (ts: any): string => {
          if (!ts) return new Date().toISOString()
          if (typeof ts === 'string') return ts
          if (typeof ts === 'object' && ts.seconds !== undefined) {
            // Protobuf Timestamp: { seconds: number, nanos: number }
            return new Date(ts.seconds * 1000 + (ts.nanos || 0) / 1000000).toISOString()
          }
          return new Date(ts).toISOString()
        }
        
        const frontendUser: User = {
          id: String(user.id || user.Id || ''),
          username: user.username || user.Username || '',
          email: user.email || user.Email || '',
          avatar: (user.profile || user.Profile)?.avatarUrl || (user.profile || user.Profile)?.AvatarUrl || '',
          profile: (user.profile || user.Profile) ? {
            id: String((user.profile || user.Profile)?.id || (user.profile || user.Profile)?.Id || ''),
            userId: String((user.profile || user.Profile)?.userId || (user.profile || user.Profile)?.UserId || ''),
            nickname: (user.profile || user.Profile)?.nickname || (user.profile || user.Profile)?.Nickname || '',
            avatarUrl: (user.profile || user.Profile)?.avatarUrl || (user.profile || user.Profile)?.AvatarUrl || '',
            bio: (user.profile || user.Profile)?.bio || (user.profile || user.Profile)?.Bio || '',
            phone: (user.profile || user.Profile)?.phone || (user.profile || user.Profile)?.Phone || '',
            gender: (user.profile || user.Profile)?.gender || (user.profile || user.Profile)?.Gender || '',
            birthday: (user.profile || user.Profile)?.birthday || (user.profile || user.Profile)?.Birthday || '',
            location: (user.profile || user.Profile)?.location || (user.profile || user.Profile)?.Location || '',
            website: (user.profile || user.Profile)?.website || (user.profile || user.Profile)?.Website || '',
            createdAt: formatTimestamp((user.profile || user.Profile)?.createdAt || (user.profile || user.Profile)?.CreatedAt),
            updatedAt: formatTimestamp((user.profile || user.Profile)?.updatedAt || (user.profile || user.Profile)?.UpdatedAt),
          } : undefined,
          roles: [],
          permissions: [],
          createdAt: formatTimestamp(user.createdAt || user.CreatedAt),
          updatedAt: formatTimestamp(user.updatedAt || user.UpdatedAt),
        }
        
        userStore.setUser(frontendUser)
        console.log('User set in store:', { id: frontendUser.id, username: frontendUser.username })
      } else {
        console.warn('Login response missing user data, but token saved')
      }
      
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  async function logout() {
    token.value = null
    userStore.clearUser()
    localStorage.removeItem('token')
  }

  function checkAuth() {
    if (token.value) {
      // 可以在这里验证 token 是否有效
      // 如果token有效，可以获取用户信息
    }
  }

  return {
    // State
    token,
    isAuthenticated,
    // Actions
    login,
    logout,
    checkAuth,
  }
})

