/**
 * 基础组件统一导出
 * 所有基础组件都应该在这里导出
 */

// 基础组件
export { Button } from './Button'
export { Input } from './Input'
export { Select } from './Select'
export { Dialog } from './Dialog'
export { Table } from './Table'
export { Form } from './Form'
export { Card } from './Card'
export { Loading } from './Loading'
export { Empty } from './Empty'
export { Checkbox } from './Checkbox'
export { Link } from './Link'
export { Icon } from './Icon'
export * as Message from './Message'

// 表单组件（重新导出以便从 base 导入）
export { FormItem } from '../form'

// 类型导出
export type { ButtonProps, ButtonEmits } from './Button/types'
export type { InputProps, InputEmits } from './Input/types'
export type { SelectProps, SelectEmits, SelectOption } from './Select/types'
export type { DialogProps, DialogEmits } from './Dialog/types'
export type { TableProps, TableEmits, TableColumn } from './Table/types'
export type { FormProps, FormEmits, FormRules } from './Form/types'
export type { CardProps } from './Card/types'
export type { LoadingProps } from './Loading/types'
export type { EmptyProps } from './Empty/types'
export type { CheckboxProps, CheckboxEmits } from './Checkbox/types'
export type { LinkProps, LinkEmits } from './Link/types'
export type { IconProps } from './Icon/types'
export type { MessageType, MessageOptions } from './Message/types'

// 待实现的组件（占位导出）
// export { default as Badge } from './Badge'
// export { default as Tag } from './Tag'
// export { default as Tooltip } from './Tooltip'
// export { default as Popover } from './Popover'
// export { default as Dropdown } from './Dropdown'
// export { default as Menu } from './Menu'
// export { default as Tabs } from './Tabs'
// export { default as Pagination } from './Pagination'
// export { default as Notification } from './Notification'
