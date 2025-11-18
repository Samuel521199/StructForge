/**
 * API 常量
 */

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api'

export const API_ENDPOINTS = {
  // 用户相关
  USER_LOGIN: '/v1/users/login',
  USER_REGISTER: '/v1/users/register',
  USER_INFO: '/v1/users/me',

  // 工作流相关
  WORKFLOW_LIST: '/v1/workflows',
  WORKFLOW_DETAIL: (id: string) => `/v1/workflows/${id}`,
  WORKFLOW_CREATE: '/v1/workflows',
  WORKFLOW_UPDATE: (id: string) => `/v1/workflows/${id}`,
  WORKFLOW_DELETE: (id: string) => `/v1/workflows/${id}`,
  WORKFLOW_EXECUTE: (id: string) => `/v1/workflows/${id}/execute`,

  // AI相关
  AI_MODELS: '/v1/ai/models',
  AI_MODEL_DETAIL: (id: string) => `/v1/ai/models/${id}`,

  // 节点相关
  NODE_TYPES: '/v1/nodes/types',
  NODE_TYPE_DETAIL: (id: string) => `/v1/nodes/types/${id}`,
} as const

