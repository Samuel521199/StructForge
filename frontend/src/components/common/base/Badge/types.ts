/**
 * Badge 组件类型定义
 */

export interface BadgeProps {
  /** 显示值 */
  value?: number | string
  /** 最大值，超过显示 max+ */
  max?: number
  /** 是否显示为点 */
  dot?: boolean
  /** 类型 */
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  /** 自定义颜色 */
  color?: string
  /** 是否固定在右上角 */
  fixed?: boolean
}

