/**
 * FormItem 组件类型定义
 */

export interface FormItemProps {
  /** 标签文本 */
  label?: string
  /** 表单域 model 字段，在使用 validate、resetFields 方法的情况下，该属性是必填的 */
  prop?: string
  /** 是否必填，如不设置，则会根据校验规则自动生成 */
  required?: boolean
  /** 表单验证规则 */
  rules?: any
  /** 表单域验证错误信息, 设置该值会使表单验证状态变为error，并显示该错误信息 */
  error?: string
  /** 是否显示校验错误信息 */
  showMessage?: boolean
  /** 以行内形式展示校验信息 */
  inlineMessage?: boolean
  /** 用于控制该表单域下组件的尺寸 */
  size?: 'large' | 'default' | 'small'
  /** 标签宽度，例如 '50px'。作为 Form 直接子元素的 form-item 会继承该值。支持 auto。 */
  labelWidth?: string
  /** 标签的位置 */
  labelPosition?: 'left' | 'right' | 'top'
}

export interface FormItemEmits {
  (e: 'blur'): void
  (e: 'change', value: any): void
}

