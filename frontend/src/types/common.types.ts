/**
 * 通用类型定义
 */

export interface Pagination {
  page: number
  pageSize: number
  total: number
}

export interface Sort {
  field: string
  order: 'asc' | 'desc'
}

export interface Filter {
  [key: string]: any
}

