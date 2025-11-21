/**
 * 用户相关类型定义
 */

export interface UserProfile {
  id?: string
  userId?: string
  nickname?: string
  avatarUrl?: string
  bio?: string
  phone?: string
  gender?: string
  birthday?: string
  location?: string
  website?: string
  createdAt?: string
  updatedAt?: string
}

export interface User {
  id: string
  username: string
  email: string
  avatar?: string
  profile?: UserProfile
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

