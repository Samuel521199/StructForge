import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/modules/auth.store'
import { usePermissionStore } from '@/stores/modules/user.store'

/**
 * 权限路由守卫
 * 检查用户是否有权限访问特定路由
 */
export function setupPermissionGuard(router: Router) {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()
    const permissionStore = usePermissionStore()

    // 检查路由是否需要认证
    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
      next({
        name: 'Login',
        query: { redirect: to.fullPath }
      })
      return
    }

    // 检查路由是否需要管理员权限
    if (to.meta.requiresAdmin && !permissionStore.isAdmin) {
      next({
        name: 'Dashboard',
        replace: true
      })
      return
    }

    // 检查用户是否有特定权限
    if (to.meta.permission) {
      const hasPermission = permissionStore.hasPermission(to.meta.permission as string)
      if (!hasPermission) {
        next({
          name: 'Dashboard',
          replace: true
        })
        return
      }
    }

    next()
  })
}

