/**
 * Icon 组件类型定义
 */

import type { Component } from 'vue'

export interface IconProps {
  /** 图标组件（来自 @element-plus/icons-vue） */
  icon?: Component | string
  /** 图标尺寸 */
  size?: number | string
  /** 图标颜色 */
  color?: string
  /** 是否加载中（显示旋转动画） */
  isLoading?: boolean
  /** 自定义类名 */
  class?: string
}

