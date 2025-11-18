/**
 * 验证器函数
 */

export const validators = {
  required: (value: any): boolean => {
    if (typeof value === 'string') {
      return value.trim() !== ''
    }
    return value !== null && value !== undefined
  },

  email: (value: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(value)
  },

  minLength: (value: string, min: number): boolean => {
    return value.length >= min
  },

  maxLength: (value: string, max: number): boolean => {
    return value.length <= max
  },

  pattern: (value: string, pattern: RegExp): boolean => {
    return pattern.test(value)
  },
}

