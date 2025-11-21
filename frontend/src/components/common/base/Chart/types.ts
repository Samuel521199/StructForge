/**
 * Chart 组件类型定义
 */

import type { EChartsOption } from 'echarts'

export interface ChartProps {
  /** 图表配置 */
  option: EChartsOption | Record<string, unknown>
  /** 宽度 */
  width?: string
  /** 高度 */
  height?: string
  /** 主题 */
  theme?: 'light' | 'dark' | string
}

