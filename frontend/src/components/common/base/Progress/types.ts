/**
 * Progress 组件类型定义
 */

export interface ProgressProps {
  /** 百分比（0-100） */
  percentage: number
  /** 是否显示文字 */
  showInfo?: boolean
  /** 文字内容 */
  text?: string
  /** 状态 */
  status?: 'success' | 'exception' | 'warning'
  /** 自定义颜色 */
  color?: string
}

