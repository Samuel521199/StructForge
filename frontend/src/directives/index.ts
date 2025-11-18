/**
 * 自定义指令统一导出
 */

import type { App } from 'vue'
import { setupLoadingDirective } from './v-loading'
import { setupPermissionDirective } from './v-permission'

/**
 * 注册所有自定义指令
 */
export function setupDirectives(app: App) {
  setupLoadingDirective(app)
  setupPermissionDirective(app)
}
