/**
 * Avatar 组件类型定义
 */

export interface AvatarProps {
  /** 头像尺寸 */
  size?: number | 'large' | 'default' | 'small' | string
  /** 头像形状 */
  shape?: 'circle' | 'square'
  /** 头像图片地址 */
  src?: string
  /** 图标类名 */
  icon?: string | object
  /** 图片如何适应容器 */
  fit?: 'fill' | 'contain' | 'cover' | 'none' | 'scale-down'
  /** 图片无法显示时的替代文本 */
  alt?: string
  /** 显示文本（当没有 src 时） */
  text?: string
}

