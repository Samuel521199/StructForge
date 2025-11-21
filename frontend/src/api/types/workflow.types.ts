/**
 * 工作流相关类型定义
 */

export interface Workflow {
  id: string
  name: string
  description?: string
  nodes: Node[]
  edges: Edge[]
  variables?: Variable[]
  status: 'draft' | 'active' | 'archived'
  createdAt: string
  updatedAt: string
}

export interface Node {
  id: string
  type: string
  name: string
  position: { x: number; y: number }
  data: Record<string, unknown>
  inputs?: NodePort[]
  outputs?: NodePort[]
}

export interface Edge {
  id: string
  source: string
  target: string
  sourceHandle?: string
  targetHandle?: string
}

export interface NodePort {
  id: string
  name: string
  type: string
}

export type VariableValue = string | number | boolean | null | undefined | Record<string, unknown> | unknown[]

export interface Variable {
  name: string
  type: string
  value: VariableValue
  description?: string
}

export interface Execution {
  id: string
  workflowId: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  startTime?: string
  endTime?: string
  result?: Record<string, unknown>
  error?: string
}

// 工作流列表项（简化版，用于列表展示）
export interface WorkflowListItem {
  id: string
  name: string
  description?: string
  status: 'draft' | 'active' | 'archived' | 'running' | 'stopped' | 'paused' | 'error'
  createdAt: string
  updatedAt: string
}

// 工作流列表查询参数
export interface WorkflowListParams {
  page?: number
  pageSize?: number
  search?: string
  status?: WorkflowListItem['status']
}

// 工作流列表响应
export interface WorkflowListResponse {
  list: WorkflowListItem[]
  total: number
  page: number
  pageSize: number
}

