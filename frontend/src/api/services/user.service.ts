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

export interface UserInfo extends User {
  nickname?: string
  avatar_url?: string
  bio?: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface RegisterResponse {
  success: boolean
  message: string
  user?: UserInfo
}

export interface VerifyEmailRequest {
  token: string
}

export interface VerifyEmailResponse {
  success: boolean
  message: string
}

export interface ResendVerificationEmailRequest {
  email: string
}

export interface ResendVerificationEmailResponse {
  success: boolean
  message: string
}

export interface RequestPasswordResetRequest {
  email: string
}

export interface RequestPasswordResetResponse {
  success: boolean
  message: string
}

export interface ResetPasswordRequest {
  token: string
  newPassword: string
}

export interface ResetPasswordResponse {
  success: boolean
  message: string
}

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
  register(data: RegisterRequest): Promise<ApiResponse<RegisterResponse>> {
    return apiClient.post('/v1/users/register', data)
  },

  /**
   * 验证邮箱
   */
  verifyEmail(data: VerifyEmailRequest): Promise<ApiResponse<VerifyEmailResponse>> {
    return apiClient.post('/v1/users/verify-email', data)
  },

  /**
   * 重新发送验证邮件
   */
  resendVerificationEmail(data: ResendVerificationEmailRequest): Promise<ApiResponse<ResendVerificationEmailResponse>> {
    return apiClient.post('/v1/users/resend-verification', data)
  },

  /**
   * 请求重置密码
   */
  requestPasswordReset(data: RequestPasswordResetRequest): Promise<ApiResponse<RequestPasswordResetResponse>> {
    return apiClient.post('/v1/users/request-password-reset', data)
  },

  /**
   * 重置密码
   */
  resetPassword(data: ResetPasswordRequest): Promise<ApiResponse<ResetPasswordResponse>> {
    return apiClient.post('/v1/users/reset-password', data)
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

