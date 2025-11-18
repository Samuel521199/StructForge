/**
 * Message 组件类型定义
 */

export type MessageType = 'success' | 'warning' | 'error' | 'info'

export interface MessageOptions {
  /** 显示时间，单位为毫秒。设为 0 则不会自动关闭 */
  duration?: number
  /** 是否显示关闭按钮 */
  showClose?: boolean
  /** 文字是否居中 */
  center?: boolean
  /** 是否将 message 合并为一条 */
  grouping?: boolean
  /** 是否将 message 属性作为 HTML 片段处理 */
  dangerouslyUseHTMLString?: boolean
  /** 自定义类名 */
  customClass?: string
  /** 距离窗口顶部的偏移量 */
  offset?: number
  /** 是否在关闭时销毁 */
  onClose?: () => void
}

