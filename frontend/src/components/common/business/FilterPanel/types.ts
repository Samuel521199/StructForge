/**
 * FilterPanel 组件类型定义
 */

export interface FilterOption {
  /** 筛选标签 */
  label: string
  /** 筛选字段名 */
  prop: string
  /** 筛选类型 */
  type?: 'select' | 'input' | 'date' | 'daterange'
  /** 选项数据（select类型需要） */
  options?: Array<{ label: string; value: any }>
  /** 占位符 */
  placeholder?: string
  /** 是否必填 */
  required?: boolean
  /** 是否多选（select类型） */
  multiple?: boolean
}

export interface FilterPanelProps {
  /** 筛选配置 */
  filters: FilterOption[]
  /** 筛选值 */
  modelValue: Record<string, any>
  /** 是否可折叠 */
  collapsible?: boolean
}

export interface FilterPanelEmits {
  /** v-model更新 */
  (e: 'update:modelValue', value: Record<string, any>): void
  /** 筛选值变化 */
  (e: 'change', value: Record<string, any>): void
  /** 重置 */
  (e: 'reset'): void
  /** 应用 */
  (e: 'apply', value: Record<string, any>): void
}

