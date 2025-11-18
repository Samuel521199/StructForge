/**
 * 用户状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/api/types/user.types'

export const useUserStore = defineStore('user', () => {
  // 状态
  const user = ref<User | null>(null)

  // Getters
  const isLoggedIn = computed(() => user.value !== null)
  const userId = computed(() => user.value?.id)
  const username = computed(() => user.value?.username)
  const roles = computed(() => user.value?.roles || [])
  const permissions = computed(() => user.value?.permissions || [])

  // Actions
  function setUser(userData: User | null) {
    user.value = userData
  }

  function updateUser(userData: Partial<User>) {
    if (user.value) {
      user.value = { ...user.value, ...userData }
    }
  }

  function clearUser() {
    user.value = null
  }

  return {
    // State
    user,
    // Getters
    isLoggedIn,
    userId,
    username,
    roles,
    permissions,
    // Actions
    setUser,
    updateUser,
    clearUser,
  }
})

