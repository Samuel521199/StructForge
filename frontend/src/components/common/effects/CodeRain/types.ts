/**
 * CodeRain 组件类型定义
 */

export interface CodeRainProps {
  /** 字体大小 */
  fontSize?: number
  /** 字体族 */
  fontFamily?: string
  /** 字体粗细 */
  fontWeight?: string | number
  /** 代码颜色（十六进制） */
  color?: string
  /** 背景颜色（十六进制） */
  backgroundColor?: string
  /** 下落速度 */
  speed?: number
  /** 速度变化范围（0-1），值越大速度差异越大 */
  speedVariation?: number
  /** 雨滴密度（0-1） */
  density?: number
  /** 透明度（0-1） */
  opacity?: number
  /** 拖尾淡出速度（0-1） */
  fadeSpeed?: number
  /** 字符集 */
  characters?: string
  /** 最小长度 */
  minLength?: number
  /** 最大长度 */
  maxLength?: number
  /** 是否启用三层深度效果 */
  enableLayers?: boolean
  /** 是否启用光晕效果 */
  enableGlow?: boolean
  /** 是否启用动态干扰 */
  enableGlitch?: boolean
  /** 光晕强度（0-1） */
  glowIntensity?: number
}

export interface CodeRainEmits {
  // 可以添加事件定义
}

