/**
 * 工作流执行相关类型定义
 */

export interface Execution {
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
  /** 执行结果 */
  result?: Record<string, unknown>
  /** 错误信息 */
  error?: string
  /** 执行数据 */
  executionData?: Record<string, unknown>
  /** 节点执行状态 */
  nodeExecutions?: NodeExecution[]
  /** 创建时间 */
  createdAt: string
  /** 更新时间 */
  updatedAt: string
}

export interface NodeExecution {
  /** 节点ID */
  nodeId: string
  /** 节点名称 */
  nodeName: string
  /** 节点类型 */
  nodeType: string
  /** 状态 */
  status: 'pending' | 'running' | 'success' | 'failed' | 'skipped'
  /** 开始时间 */
  startTime?: string
  /** 结束时间 */
  endTime?: string
  /** 执行时长（秒） */
  duration?: number
  /** 输入数据 */
  inputs?: Record<string, unknown>
  /** 输出数据 */
  outputs?: Record<string, unknown>
  /** 错误信息 */
  error?: string
}

export interface ExecutionLog {
  /** 日志ID */
  id: string
  /** 执行ID */
  executionId: string
  /** 节点ID */
  nodeId?: string
  /** 日志级别 */
  level: 'debug' | 'info' | 'warn' | 'error'
  /** 日志消息 */
  message: string
  /** 时间戳 */
  timestamp: string
  /** 额外数据 */
  data?: Record<string, unknown>
}

export interface ExecutionListParams {
  page?: number
  pageSize?: number
  status?: Execution['status']
  workflowId?: string
  startTime?: string
  endTime?: string
}

export interface ExecutionListResponse {
  list: Execution[]
  total: number
  page: number
  pageSize: number
}

