/**
 * 格式化组合函数
 */

import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

export function useFormat() {
  const formatDate = (date: string | Date, format = 'YYYY-MM-DD HH:mm:ss') => {
    return dayjs(date).format(format)
  }

  const formatRelativeTime = (date: string | Date) => {
    return dayjs(date).fromNow()
  }

  const formatNumber = (num: number, decimals = 2) => {
    return num.toFixed(decimals)
  }

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  return {
    formatDate,
    formatRelativeTime,
    formatNumber,
    formatFileSize,
  }
}

