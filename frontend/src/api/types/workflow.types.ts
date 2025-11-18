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
  data: Record<string, any>
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

export interface Variable {
  name: string
  type: string
  value: any
  description?: string
}

export interface Execution {
  id: string
  workflowId: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  startTime?: string
  endTime?: string
  result?: any
  error?: string
}

