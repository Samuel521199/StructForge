/**
 * 插件统一导出和注册
 */

import type { App } from 'vue'
import { setupElementPlus } from './element-plus'
import { setupVueFlow } from './vue-flow'

/**
 * 注册所有插件
 */
export function setupPlugins(app: App) {
  setupElementPlus(app)
  setupVueFlow(app)
}
