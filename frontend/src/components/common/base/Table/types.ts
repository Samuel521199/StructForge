export interface TableColumn {
  prop?: string
  label: string
  width?: string | number
  minWidth?: string | number
  fixed?: boolean | 'left' | 'right'
  sortable?: boolean
  formatter?: (row: unknown, column: TableColumn, cellValue: unknown) => unknown
  align?: 'left' | 'center' | 'right'
  slot?: string // 自定义插槽名称
}

export interface TableProps<T = unknown> {
  data: T[]
  columns: TableColumn[]
  stripe?: boolean
  border?: boolean
  size?: 'large' | 'default' | 'small'
  showHeader?: boolean
  highlightCurrentRow?: boolean
  emptyText?: string
  loading?: boolean
  height?: string | number
  maxHeight?: string | number
}

export interface TableEmits<T = unknown> {
  (e: 'selection-change', selection: T[]): void
  (e: 'row-click', row: T, column: TableColumn, event: Event): void
  (e: 'sort-change', sortInfo: { column: TableColumn; prop: string; order: string }): void
}

