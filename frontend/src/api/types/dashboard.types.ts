/**
 * 仪表盘相关类型定义
 */

export interface DashboardStats {
  /** 正在运行的工作流数量 */
  running: number
  /** 排队等待的工作流数量 */
  queued: number
  /** 总执行次数（最近7天） */
  totalExecutions: number
  /** 失败执行次数（最近7天） */
  failedExecutions: number
  /** 失败率 */
  failureRate: number
  /** 平均执行时间（秒） */
  avgExecutionTime: number
  /** 总工作流数量 */
  totalWorkflows: number
  /** 活跃工作流数量 */
  activeWorkflows: number
}

export interface ExecutionListItem {
  /** 执行ID */
  id: string
  /** 工作流ID */
  workflowId: string
  /** 工作流名称 */
  workflowName: string
  /** 状态 */
  status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  /** 开始时间 */
  startTime: string
  /** 结束时间 */
  endTime?: string
  /** 执行时长（秒） */
  duration?: number
  /** 执行数据 */
  executionData?: Record<string, unknown>
}

export interface SuccessRateItem {
  /** 工作流名称 */
  workflowName: string
  /** 总执行次数 */
  totalExecutions: number
  /** 成功次数 */
  successful: number
  /** 失败次数 */
  failed: number
  /** 成功率 */
  successRate: number
}

export interface ErrorTrendData {
  /** 时间 */
  time: string
  /** 错误数量 */
  count: number
}

export interface ExecutionDurationItem {
  /** 工作流名称 */
  workflowName: string
  /** 平均执行时长（毫秒） */
  avgDuration: number
  /** 最小执行时长（毫秒） */
  minDuration: number
  /** 最大执行时长（毫秒） */
  maxDuration: number
}

export interface WebhookActivity {
  /** 工作流名称 */
  workflowName: string
  /** Webhook 路径 */
  webhookPath: string
  /** HTTP 方法 */
  httpMethod: string
  /** 执行时间 */
  executionTime: string
  /** Webhook ID */
  webhookId?: string
}

