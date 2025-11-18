import type { RouteRecordRaw } from 'vue-router'

/**
 * 系统模块路由
 */
const systemRoutes: RouteRecordRaw[] = [
  {
    path: '/system',
    name: 'System',
    redirect: '/system/settings',
    meta: {
      title: '系统',
      icon: 'system',
      requiresAuth: true,
      requiresAdmin: true // 需要管理员权限
    },
    children: [
      {
        path: 'settings',
        name: 'SystemSettings',
        component: () => import('@/views/system/SystemSettings/SystemSettings.vue'),
        meta: {
          title: '系统设置',
          requiresAuth: true,
          requiresAdmin: true
        }
      },
      {
        path: 'logs',
        name: 'SystemLogs',
        component: () => import('@/views/system/SystemLogs/SystemLogs.vue'),
        meta: {
          title: '系统日志',
          requiresAuth: true,
          requiresAdmin: true
        }
      }
    ]
  }
]

export default systemRoutes

