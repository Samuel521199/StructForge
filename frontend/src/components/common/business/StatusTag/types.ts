export interface StatusConfig {
  label: string
  type: 'success' | 'warning' | 'danger' | 'info'
  color?: string
}

export interface StatusTagProps {
  status: string
  statusMap?: Record<string, StatusConfig>
}

