/**
 * 工作流常量
 */

export const WORKFLOW_STATUS = {
  DRAFT: 'draft',
  ACTIVE: 'active',
  ARCHIVED: 'archived',
} as const

export const EXECUTION_STATUS = {
  PENDING: 'pending',
  RUNNING: 'running',
  SUCCESS: 'success',
  FAILED: 'failed',
  CANCELLED: 'cancelled',
} as const

export const NODE_CATEGORIES = {
  TRIGGER: 'trigger',
  AI: 'ai',
  DATA: 'data',
  INTEGRATION: 'integration',
  CONTROL: 'control',
  TOOL: 'tool',
} as const

