/**
 * Loading 组件类型定义
 */

export interface LoadingProps {
  /** 是否显示加载状态 */
  loading?: boolean
  /** 加载文字 */
  text?: string
  /** 背景色 */
  background?: string
  /** 自定义加载图标类名 */
  spinner?: string
  /** Loading 覆盖的 DOM 节点 */
  element?: string | HTMLElement
  /** Loading 需要覆盖的 DOM 节点 */
  target?: string | HTMLElement
  /** 是否将遮罩层插入到 body 中 */
  body?: boolean
  /** 是否全屏显示 */
  fullscreen?: boolean
  /** 是否锁定屏幕滚动 */
  lock?: boolean
}

