/**
 * ActionBar 组件类型定义
 */

import type { Component } from 'vue'

export interface ActionItem {
  /** 操作标签 */
  label: string
  /** 按钮类型 */
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'text'
  /** 按钮大小 */
  size?: 'large' | 'default' | 'small'
  /** 是否禁用 */
  disabled?: boolean
  /** 是否加载中 */
  loading?: boolean
  /** 图标 */
  icon?: Component | string
  /** 是否圆角 */
  round?: boolean
  /** 是否朴素按钮 */
  plain?: boolean
  /** 是否隐藏 */
  hidden?: boolean
  /** 点击回调 */
  onClick?: () => void
}

export interface ActionBarProps {
  /** 操作项列表 */
  actions: ActionItem[]
  /** 对齐方式 */
  align?: 'left' | 'center' | 'right'
  /** 按钮大小 */
  size?: 'large' | 'default' | 'small'
}

export interface ActionBarEmits {
  /** 操作点击事件 */
  (e: 'action', action: ActionItem, index: number): void
}

