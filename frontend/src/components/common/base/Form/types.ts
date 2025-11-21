import type { FormRules as ElementFormRules } from 'element-plus'

// 使用 Element Plus 的 FormRules 类型，但允许更灵活的使用
export type FormRules = ElementFormRules

export interface FormProps {
  model: Record<string, any>
  rules?: FormRules | Record<string, any>
  labelWidth?: string
  labelPosition?: 'left' | 'right' | 'top'
  size?: 'large' | 'default' | 'small'
  disabled?: boolean
}

export interface FormEmits {
  (e: 'validate', prop: string, isValid: boolean, message: string): void
  (e: 'submit', event: Event): void
}

