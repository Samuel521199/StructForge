/**
 * 用户相关类型定义
 */

export interface User {
  id: string
  username: string
  email: string
  avatar?: string
  roles: string[]
  permissions: string[]
  createdAt: string
  updatedAt: string
}

export interface Role {
  id: string
  name: string
  permissions: string[]
}

export interface Permission {
  id: string
  name: string
  resource: string
  action: string
}

