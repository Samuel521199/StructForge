/**
 * Tag 组件类型定义
 */

export interface TagProps {
  /** 类型 */
  type?: 'success' | 'info' | 'warning' | 'danger' | ''
  /** 尺寸 */
  size?: 'large' | 'default' | 'small'
  /** 主题 */
  effect?: 'dark' | 'light' | 'plain'
  /** 是否可关闭 */
  closable?: boolean
  /** 是否禁用渐变动画 */
  disableTransitions?: boolean
  /** 是否有边框描边 */
  hit?: boolean
  /** 是否圆形 */
  round?: boolean
  /** 背景色 */
  color?: string
}

export interface TagEmits {
  (e: 'close', event: MouseEvent): void
}

