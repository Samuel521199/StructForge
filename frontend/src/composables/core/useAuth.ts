/**
 * 认证相关组合函数
 */

import { computed } from 'vue'
import { useAuthStore } from '@/stores/modules/auth.store'
import { useUserStore } from '@/stores/modules/user.store'

export function useAuth() {
  const authStore = useAuthStore()
  const userStore = useUserStore()

  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const user = computed(() => userStore.user)
  const token = computed(() => authStore.token)

  const login = async (username: string, password: string) => {
    return await authStore.login(username, password)
  }

  const logout = () => {
    authStore.logout()
  }

  return {
    isAuthenticated,
    user,
    token,
    login,
    logout,
  }
}

