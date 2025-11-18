/**
 * 路由类型定义
 */

export interface RouteMeta {
  title?: string
  requiresAuth?: boolean
  requiresPermission?: string[]
  requiresRole?: string[]
  keepAlive?: boolean
}

