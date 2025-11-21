/**
 * 环境类型声明
 */

/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, unknown>
  export default component
}

interface ImportMetaEnv {
  /** 后端 API 基础地址 */
  readonly VITE_API_BASE_URL?: string
  /** WebSocket 服务器地址 */
  readonly VITE_WS_URL?: string
  /** 应用标题 */
  readonly VITE_APP_TITLE?: string
  /** 应用描述 */
  readonly VITE_APP_DESCRIPTION?: string
  /** 环境标识 (development | production | test) */
  readonly VITE_APP_ENV?: 'development' | 'production' | 'test'
  /** 是否启用调试模式 */
  readonly VITE_APP_DEBUG?: string
  /** 是否启用工作流编辑器 */
  readonly VITE_FEATURE_WORKFLOW_EDITOR?: string
  /** 是否启用 AI 模型管理 */
  readonly VITE_FEATURE_AI_MODELS?: string
  /** 是否启用用户管理 */
  readonly VITE_FEATURE_USER_MANAGEMENT?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

