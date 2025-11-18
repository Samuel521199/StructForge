/**
 * 文本格式化工具
 */

export const truncate = (text: string, length: number) => {
  if (text.length <= length) return text
  return text.slice(0, length) + '...'
}

export const capitalize = (text: string) => {
  return text.charAt(0).toUpperCase() + text.slice(1)
}

