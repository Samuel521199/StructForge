export interface TableColumn {
  prop?: string
  label: string
  width?: string | number
  minWidth?: string | number
  fixed?: boolean | 'left' | 'right'
  sortable?: boolean
  formatter?: (row: any, column: TableColumn, cellValue: any) => any
  align?: 'left' | 'center' | 'right'
}

export interface TableProps {
  data: any[]
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

export interface TableEmits {
  (e: 'selection-change', selection: any[]): void
  (e: 'row-click', row: any, column: TableColumn, event: Event): void
  (e: 'sort-change', sortInfo: { column: TableColumn; prop: string; order: string }): void
}

