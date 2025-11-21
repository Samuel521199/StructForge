/**
 * Statistic 组件类型定义
 */

export interface StatisticProps {
  /** 标题 */
  title?: string
  /** 数值 */
  value: number | string
  /** 后缀 */
  suffix?: string
  /** 描述 */
  description?: string
  /** 图标 */
  icon?: unknown
  /** 数值颜色 */
  valueColor?: string
  /** 小数精度 */
  precision?: number
  /** 趋势 */
  trend?: {
    direction: 'up' | 'down'
    value: number
  }
}

