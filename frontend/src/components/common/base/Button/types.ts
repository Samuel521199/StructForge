/**
 * Button 组件类型定义
 */

export interface ButtonProps {
  /** 按钮类型 */
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'text'
  /** 按钮尺寸 */
  size?: 'large' | 'default' | 'small'
  /** 是否禁用 */
  disabled?: boolean
  /** 是否加载中 */
  loading?: boolean
  /** 图标类名 */
  icon?: string
  /** 是否圆角按钮 */
  round?: boolean
  /** 是否圆形按钮 */
  circle?: boolean
  /** 是否朴素按钮 */
  plain?: boolean
  /** 原生 type 属性 */
  nativeType?: 'button' | 'submit' | 'reset'
}

export interface ButtonEmits {
  /** 点击事件 */
  (e: 'click', event: MouseEvent): void
}

