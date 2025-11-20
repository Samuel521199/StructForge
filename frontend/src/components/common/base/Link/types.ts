/**
 * Link 组件类型定义
 */

export interface LinkProps {
  /** 类型 */
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'default'
  /** 是否下划线 */
  underline?: boolean
  /** 是否禁用 */
  disabled?: boolean
  /** 原生 href 属性 */
  href?: string
  /** 原生 target 属性 */
  target?: string
  /** 图标类名 */
  icon?: string | object
}

export interface LinkEmits {
  /** 点击事件 */
  (e: 'click', event: MouseEvent): void
}

