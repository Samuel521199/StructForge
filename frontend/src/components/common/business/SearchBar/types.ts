export interface FilterOption {
  label: string
  prop: string
  type?: 'select' | 'date' | 'daterange' | 'input'
  options?: Array<{ label: string; value: any }>
  placeholder?: string
}

export interface SearchBarProps {
  modelValue: string
  placeholder?: string
  searchable?: boolean
  filterable?: boolean
  filters?: FilterOption[]
}

export interface SearchBarEmits {
  (e: 'update:modelValue', value: string): void
  (e: 'search', keyword: string): void
  (e: 'filter', filters: Record<string, any>): void
}

