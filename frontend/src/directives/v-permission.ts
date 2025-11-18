import type { App, Directive } from 'vue'
import { usePermissionStore } from '@/stores/modules/user.store'

/**
 * 权限指令
 * 使用方式: v-permission="'workflow:create'"
 */
const permissionDirective: Directive = {
  mounted(el, binding) {
    const permissionStore = usePermissionStore()
    const permission = binding.value

    if (permission && !permissionStore.hasPermission(permission)) {
      el.style.display = 'none'
      // 或者直接移除元素
      // el.parentNode?.removeChild(el)
    }
  },
  updated(el, binding) {
    const permissionStore = usePermissionStore()
    const permission = binding.value

    if (permission && !permissionStore.hasPermission(permission)) {
      el.style.display = 'none'
    } else {
      el.style.display = ''
    }
  }
}

export function setupPermissionDirective(app: App) {
  app.directive('permission', permissionDirective)
}

