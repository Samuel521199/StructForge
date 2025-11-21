/**
 * 路由配置
 */

import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { routes } from './routes'
import { authGuard } from './guards/auth.guard'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes as RouteRecordRaw[],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  console.log('路由导航:', { from: from.path, to: to.path, name: to.name })
  try {
    // 认证守卫
    authGuard(to, from, next)
  } catch (error) {
    console.error('路由守卫错误:', error)
    // 如果守卫出错，重定向到登录页
    if (to.path !== '/auth/login') {
      next({ name: 'Login', query: { redirect: to.fullPath } })
    } else {
      next()
    }
  }
})

// 路由错误处理
router.onError((error) => {
  console.error('路由错误:', error)
})

export default router

