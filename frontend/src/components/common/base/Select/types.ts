export interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
  children?: SelectOption[]
}

export interface SelectProps {
  modelValue: string | number | Array<string | number>
  options: SelectOption[]
  placeholder?: string
  multiple?: boolean
  disabled?: boolean
  clearable?: boolean
  filterable?: boolean
  allowCreate?: boolean
  size?: 'large' | 'default' | 'small'
  loading?: boolean
}

export interface SelectEmits {
  (e: 'update:modelValue', value: string | number | Array<string | number>): void
  (e: 'change', value: string | number | Array<string | number>): void
  (e: 'visible-change', visible: boolean): void
  (e: 'remove-tag', tag: string | number): void
  (e: 'clear'): void
}

