/**
 * TimeAgo 组件类型定义
 */

export interface TimeAgoProps {
  /** 时间（时间戳、日期字符串或Date对象） */
  time: number | string | Date
  /** 完整时间格式 */
  format?: string
  /** 是否显示完整时间 */
  showFullTime?: boolean
  /** 更新间隔（毫秒），0表示不自动更新 */
  updateInterval?: number
}

