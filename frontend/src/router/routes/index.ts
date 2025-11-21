/**
 * 路由定义统一导出
 * 整合所有模块的路由配置
 */

import type { RouteRecordRaw } from 'vue-router'
import workflowRoutes from './workflow.routes'
import userRoutes from './user.routes'
import systemRoutes from './system.routes'

// 认证路由
const authRoutes: RouteRecordRaw[] = [
  {
    path: '/auth',
    name: 'Auth',
    redirect: '/auth/login',
    meta: {
      title: '认证',
      hideInMenu: true
    },
    children: [
      {
        path: 'login',
        name: 'Login',
        component: () => import('@/views/auth/Login.vue'),
        meta: {
          title: '登录',
          hideInMenu: true
        }
      },
      {
        path: 'register',
        name: 'Register',
        component: () => import('@/views/auth/Register.vue'),
        meta: {
          title: '注册',
          hideInMenu: true
        }
      },
      {
        path: 'forgot-password',
        name: 'ForgotPassword',
        component: () => import('@/views/auth/ForgotPassword.vue'),
        meta: {
          title: '忘记密码',
          hideInMenu: true
        }
      },
      {
        path: 'reset-password',
        name: 'ResetPassword',
        component: () => import('@/views/auth/ResetPassword.vue'),
        meta: {
          title: '重置密码',
          hideInMenu: true
        }
      },
      {
        path: 'verify-email',
        name: 'VerifyEmail',
        component: () => import('@/views/auth/VerifyEmail.vue'),
        meta: {
          title: '邮箱验证',
          hideInMenu: true
        }
      }
    ]
  }
]

// 仪表盘路由
const dashboardRoutes: RouteRecordRaw[] = [
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/Dashboard.vue'),
    meta: {
      title: '仪表盘',
      icon: 'dashboard',
      requiresAuth: true
    }
  }
]

// AI模块路由
const aiRoutes: RouteRecordRaw[] = [
  {
    path: '/ai',
    name: 'AI',
    redirect: '/ai/models',
    meta: {
      title: 'AI模型',
      icon: 'ai',
      requiresAuth: true
    },
    children: [
      {
        path: 'models',
        name: 'ModelList',
        component: () => import('@/views/ai/ModelList/ModelList.vue'),
        meta: {
          title: '模型列表',
          requiresAuth: true
        }
      },
      {
        path: 'config',
        name: 'ModelConfig',
        component: () => import('@/views/ai/ModelConfig/ModelConfig.vue'),
        meta: {
          title: '模型配置',
          requiresAuth: true
        }
      }
    ]
  }
]

// 整合所有路由
export const routes: RouteRecordRaw[] = [
  // 认证路由（不需要布局）
  ...authRoutes,
  // 主布局路由（需要认证的页面）
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true, // 整个布局都需要认证
    },
    children: [
      ...dashboardRoutes,
      ...workflowRoutes,
      ...userRoutes,
      ...aiRoutes,
      ...systemRoutes,
    ]
  },
  // 404 路由
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: {
      title: '404',
      hideInMenu: true
    }
  }
]

