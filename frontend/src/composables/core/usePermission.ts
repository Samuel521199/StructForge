/**
 * 权限相关组合函数
 */

import { computed } from 'vue'
import { useUserStore } from '@/stores/modules/user.store'

export function usePermission() {
  const userStore = useUserStore()

  const permissions = computed(() => userStore.permissions)
  const roles = computed(() => userStore.roles)

  const hasPermission = (permission: string): boolean => {
    return permissions.value.includes(permission)
  }

  const hasRole = (role: string): boolean => {
    return roles.value.includes(role)
  }

  const hasAnyPermission = (perms: string[]): boolean => {
    return perms.some(perm => permissions.value.includes(perm))
  }

  const hasAllPermissions = (perms: string[]): boolean => {
    return perms.every(perm => permissions.value.includes(perm))
  }

  return {
    permissions,
    roles,
    hasPermission,
    hasRole,
    hasAnyPermission,
    hasAllPermissions,
  }
}

