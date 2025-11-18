/**
 * 认证状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userService } from '@/api/services/user.service'
import { useUserStore } from './user.store'

export const useAuthStore = defineStore('auth', () => {
  const userStore = useUserStore()

  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const isAuthenticated = computed(() => {
    return !!token.value && userStore.isLoggedIn
  })

  // Actions
  async function login(username: string, password: string) {
    try {
      const response = await userService.login({ username, password })
      if (response.data) {
        token.value = response.data.token
        userStore.setUser(response.data.user)
        localStorage.setItem('token', response.data.token)
        return true
      }
      return false
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

