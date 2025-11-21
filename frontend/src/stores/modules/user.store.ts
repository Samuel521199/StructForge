/**
 * 用户状态管理
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/api/types/user.types'

// 从 localStorage 恢复用户信息
function loadUserFromStorage(): User | null {
  try {
    const userStr = localStorage.getItem('user')
    if (userStr) {
      return JSON.parse(userStr) as User
    }
  } catch (error) {
    console.error('Failed to load user from localStorage:', error)
    localStorage.removeItem('user')
  }
  return null
}

// 保存用户信息到 localStorage
function saveUserToStorage(user: User | null) {
  if (user) {
    try {
      localStorage.setItem('user', JSON.stringify(user))
    } catch (error) {
      console.error('Failed to save user to localStorage:', error)
    }
  } else {
    localStorage.removeItem('user')
  }
}

export const useUserStore = defineStore('user', () => {
  // 状态 - 从 localStorage 恢复
  const user = ref<User | null>(loadUserFromStorage())

  // Getters
  const isLoggedIn = computed(() => user.value !== null)
  const userId = computed(() => user.value?.id)
  const username = computed(() => user.value?.username)
  const roles = computed(() => user.value?.roles || [])
  const permissions = computed(() => user.value?.permissions || [])

  // Actions
  function setUser(userData: User | null) {
    user.value = userData
    saveUserToStorage(userData)
  }

  function updateUser(userData: Partial<User>) {
    if (user.value) {
      user.value = { ...user.value, ...userData }
      saveUserToStorage(user.value)
    }
  }

  function clearUser() {
    user.value = null
    saveUserToStorage(null)
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

