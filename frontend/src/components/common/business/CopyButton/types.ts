/**
 * CopyButton 组件类型定义
 */

export interface CopyButtonProps {
  /** 要复制的内容 */
  value: string | object
  /** 按钮文本 */
  text?: string
  /** 复制成功后的文本 */
  successText?: string
  /** 按钮类型 */
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'text' | 'default'
  /** 按钮大小 */
  size?: 'large' | 'default' | 'small'
  /** 是否禁用 */
  disabled?: boolean
  /** 是否圆角 */
  round?: boolean
  /** 是否朴素按钮 */
  plain?: boolean
  /** 是否显示消息提示 */
  showMessage?: boolean
}

