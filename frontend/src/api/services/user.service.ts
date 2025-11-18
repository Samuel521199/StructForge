/**
 * 用户服务 API
 */

import { apiClient } from '../client'
import type { ApiResponse } from '../client/types'
import type { User } from '../types/user.types'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

export interface UserInfo extends User {}

export const userService = {
  /**
   * 用户登录
   */
  login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return apiClient.post('/v1/users/login', data)
  },

  /**
   * 用户注册
   */
  register(data: RegisterRequest): Promise<ApiResponse<UserInfo>> {
    return apiClient.post('/v1/users/register', data)
  },

  /**
   * 获取用户信息
   */
  getUserInfo(): Promise<ApiResponse<UserInfo>> {
    return apiClient.get('/v1/users/me')
  },

  /**
   * 更新用户信息
   */
  updateUserInfo(data: Partial<UserInfo>): Promise<ApiResponse<UserInfo>> {
    return apiClient.put('/v1/users/me', data)
  },
}

interface RegisterRequest {
  username: string
  email: string
  password: string
}

