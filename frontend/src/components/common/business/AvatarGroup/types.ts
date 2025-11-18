/**
 * AvatarGroup 组件类型定义
 */

import type { Component } from 'vue'

export interface AvatarItem {
  /** 头像图片地址 */
  src?: string
  /** 头像文本（无图片时显示） */
  text?: string
  /** 头像图标 */
  icon?: Component | string
  /** 头像形状 */
  shape?: 'circle' | 'square'
  /** 标签文本 */
  label?: string
  /** 额外数据 */
  data?: any
}

export interface AvatarGroupProps {
  /** 头像列表 */
  avatars: AvatarItem[]
  /** 头像大小 */
  size?: 'large' | 'default' | 'small' | number
  /** 头像形状 */
  shape?: 'circle' | 'square'
  /** 是否堆叠显示 */
  stacked?: boolean
  /** 最大显示数量（超出显示+数字） */
  max?: number
}

export interface AvatarGroupEmits {
  /** 头像点击事件 */
  (e: 'avatarClick', avatar: AvatarItem, index: number): void
  /** 更多点击事件 */
  (e: 'moreClick'): void
}

