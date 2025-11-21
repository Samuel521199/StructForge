/**
 * MenuItem 组件类型定义
 */

export interface MenuItemProps {
  /** 标签文字 */
  label: string
  /** 图标组件 */
  icon?: unknown
  /** 图标图片 URL（如果使用图片而不是图标组件） */
  iconImage?: string
  /** 图标大小 */
  iconSize?: number
  /** 图标类型（用于自动应用颜色） */
  iconType?: 'default' | 'dashboard' | 'workflow' | 'user' | 'ai' | 'system'
  /** 自定义图标颜色 */
  iconColor?: string
  /** 是否激活 */
  active?: boolean
  /** 是否禁用 */
  disabled?: boolean
  /** 徽章值 */
  badge?: number | string
  /** 是否显示右侧箭头 */
  showArrow?: boolean
}

export interface MenuItemEmits {
  (e: 'click', event: MouseEvent): void
}

