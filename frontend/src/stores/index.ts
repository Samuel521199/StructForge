/**
 * Pinia Store 统一导出
 */

import { createPinia } from 'pinia'

export const pinia = createPinia()

// 导出所有 store
export * from './modules/user.store'
export * from './modules/auth.store'
export * from './modules/workflow.store'
export * from './modules/execution.store'
export * from './modules/ai.store'
export * from './modules/node.store'
export * from './modules/ui.store'
export * from './modules/app.store'

