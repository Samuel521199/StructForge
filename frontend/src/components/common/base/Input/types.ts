import type { Component } from 'vue'

export interface InputProps {
  modelValue: string | number
  type?: 'text' | 'password' | 'number' | 'email' | 'url' | 'tel'
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  clearable?: boolean
  showPassword?: boolean
  prefixIcon?: string | Component
  suffixIcon?: string | Component
  maxlength?: number
  minlength?: number
  showWordLimit?: boolean
  validateEvent?: boolean
  size?: 'large' | 'default' | 'small'
}

export interface InputEmits {
  (e: 'update:modelValue', value: string | number): void
  (e: 'focus', event: FocusEvent): void
  (e: 'blur', event: FocusEvent): void
  (e: 'clear'): void
  (e: 'input', value: string | number): void
}

