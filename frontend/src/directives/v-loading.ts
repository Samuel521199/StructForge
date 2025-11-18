import type { App, Directive } from 'vue'
import { Loading } from 'element-plus'

/**
 * 加载指令
 * 使用方式: v-loading="isLoading"
 */
const loadingDirective: Directive = {
  mounted(el, binding) {
    const loadingInstance = Loading.service({
      target: el,
      text: binding.value?.text || '加载中...',
      background: binding.value?.background || 'rgba(0, 0, 0, 0.7)'
    })

    el._loadingInstance = loadingInstance
  },
  updated(el, binding) {
    if (binding.value) {
      if (!el._loadingInstance) {
        const loadingInstance = Loading.service({
          target: el,
          text: binding.value?.text || '加载中...',
          background: binding.value?.background || 'rgba(0, 0, 0, 0.7)'
        })
        el._loadingInstance = loadingInstance
      }
    } else {
      if (el._loadingInstance) {
        el._loadingInstance.close()
        el._loadingInstance = null
      }
    }
  },
  unmounted(el) {
    if (el._loadingInstance) {
      el._loadingInstance.close()
      el._loadingInstance = null
    }
  }
}

export function setupLoadingDirective(app: App) {
  app.directive('loading', loadingDirective)
}

