/**
 * Checkbox 组件类型定义
 */

export interface CheckboxProps {
  /** 绑定值 */
  modelValue?: boolean | string | number
  /** 选中状态的值（当 modelValue 为数组时使用） */
  label?: string | number | boolean
  /** 是否禁用 */
  disabled?: boolean
  /** 是否只读 */
  readonly?: boolean
  /** 是否显示边框 */
  border?: boolean
  /** 复选框的尺寸 */
  size?: 'large' | 'default' | 'small'
  /** 原生 name 属性 */
  name?: string
  /** 原生 id 属性 */
  id?: string
  /** 是否不确定状态 */
  indeterminate?: boolean
  /** 是否选中 */
  checked?: boolean
  /** 原生 tabindex 属性 */
  tabindex?: string | number
  /** 是否触发表单验证 */
  validateEvent?: boolean
}

export interface CheckboxEmits {
  /** 更新绑定值 */
  (e: 'update:modelValue', value: boolean | string | number): void
  /** 值改变时触发 */
  (e: 'change', value: boolean | string | number): void
}

