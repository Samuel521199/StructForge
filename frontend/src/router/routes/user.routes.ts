import type { RouteRecordRaw } from 'vue-router'

/**
 * 用户模块路由
 */
const userRoutes: RouteRecordRaw[] = [
  {
    path: '/user',
    name: 'User',
    redirect: '/user/profile',
    meta: {
      title: '用户',
      icon: 'user',
      requiresAuth: true
    },
    children: [
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/user/UserProfile/UserProfile.vue'),
        meta: {
          title: '个人资料',
          requiresAuth: true
        }
      },
      {
        path: 'settings',
        name: 'UserSettings',
        component: () => import('@/views/user/UserSettings/UserSettings.vue'),
        meta: {
          title: '设置',
          requiresAuth: true
        }
      }
    ]
  }
]

export default userRoutes

