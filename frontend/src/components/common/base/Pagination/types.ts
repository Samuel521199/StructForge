/**
 * Pagination 组件类型定义
 */

export interface PaginationProps {
  /** 当前页数 */
  currentPage?: number
  /** 每页显示条目个数 */
  pageSize?: number
  /** 总条目数 */
  total?: number
  /** 每页显示个数选择器的选项设置 */
  pageSizes?: number[]
  /** 组件布局，子组件名用逗号分隔 */
  layout?: string
  /** 是否为分页按钮添加背景色 */
  background?: boolean
  /** 是否禁用 */
  disabled?: boolean
  /** 只有一页时是否隐藏 */
  hideOnSinglePage?: boolean
}

export interface PaginationEmits {
  /** 每页条数改变时触发 */
  (e: 'size-change', size: number): void
  /** 当前页改变时触发 */
  (e: 'current-change', page: number): void
  /** 点击上一页按钮时触发 */
  (e: 'prev-click', page: number): void
  /** 点击下一页按钮时触发 */
  (e: 'next-click', page: number): void
}

