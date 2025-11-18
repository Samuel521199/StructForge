/**
 * 执行状态管理
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Execution {
  id: string
  workflowId: string
  status: 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
  startTime?: string
  endTime?: string
  result?: any
  error?: string
}

export interface LogEntry {
  id: string
  executionId: string
  level: 'info' | 'warn' | 'error'
  message: string
  timestamp: string
  data?: any
}

export const useExecutionStore = defineStore('execution', () => {
  // 状态
  const currentExecution = ref<Execution | null>(null)
  const executionLogs = ref<LogEntry[]>([])
  const executionStatus = ref<'idle' | 'running' | 'completed' | 'failed'>('idle')

  // Actions
  function setCurrentExecution(execution: Execution | null) {
    currentExecution.value = execution
    if (execution) {
      executionStatus.value = execution.status === 'running' ? 'running' : 
                             execution.status === 'success' ? 'completed' :
                             execution.status === 'failed' ? 'failed' : 'idle'
    } else {
      executionStatus.value = 'idle'
    }
  }

  function setExecutionLogs(logs: LogEntry[]) {
    executionLogs.value = logs
  }

  function addExecutionLog(log: LogEntry) {
    executionLogs.value.push(log)
  }

  function clearExecutionLogs() {
    executionLogs.value = []
  }

  return {
    // State
    currentExecution,
    executionLogs,
    executionStatus,
    // Actions
    setCurrentExecution,
    setExecutionLogs,
    addExecutionLog,
    clearExecutionLogs,
  }
})

