/**
 * 认证路由守卫
 */

import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/modules/auth.store'

export function authGuard(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  try {
    const authStore = useAuthStore()
    const requiresAuth = to.meta.requiresAuth === true
    
    // 检查认证状态
    // 如果有 token，即使暂时没有用户信息，也认为已认证（用户信息会在页面加载时恢复）
    const hasToken = !!authStore.token
    const isAuthenticated = authStore.isAuthenticated || hasToken

    console.log('认证检查:', {
      path: to.path,
      name: to.name,
      requiresAuth,
      isAuthenticated,
      hasToken,
      authStoreIsAuthenticated: authStore.isAuthenticated,
    })

    // 如果访问的是登录页，直接允许
    if (to.name === 'Login' || to.path === '/auth/login') {
      // 如果已登录，重定向到仪表盘
      if (isAuthenticated) {
        console.log('已登录，重定向到仪表盘')
        next({ path: '/dashboard', replace: true })
      } else {
        next()
      }
      return
    }

    // 检查是否需要认证（包括父路由的 meta）
    const parentRequiresAuth = to.matched.some(record => record.meta.requiresAuth === true)
    const needsAuth = requiresAuth || parentRequiresAuth

    // 如果访问根路径且未登录，直接重定向到登录页
    if (to.path === '/' && !isAuthenticated) {
      console.log('访问根路径但未登录，重定向到登录页')
      next({ name: 'Login', replace: true })
      return
    }

    // 检查是否需要认证
    if (needsAuth && !isAuthenticated) {
      console.log('需要认证但未登录，重定向到登录页', { path: to.path, requiresAuth, parentRequiresAuth })
      next({ name: 'Login', query: { redirect: to.fullPath } })
    } else {
      next()
    }
  } catch (error) {
    console.error('认证守卫错误:', error)
    // 如果出错，重定向到登录页（避免白屏）
    if (to.path !== '/auth/login') {
      next({ name: 'Login' })
    } else {
      next()
    }
  }
}

