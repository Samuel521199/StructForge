/**
 * 执行相关组合函数
 */

import { computed } from 'vue'
import { useExecutionStore } from '@/stores/modules/execution.store'
import { workflowService } from '@/api/services/workflow.service'

export function useExecution() {
  const executionStore = useExecutionStore()

  const currentExecution = computed(() => executionStore.currentExecution)
  const executionLogs = computed(() => executionStore.executionLogs)
  const executionStatus = computed(() => executionStore.executionStatus)

  const executeWorkflow = async (workflowId: string, params?: Record<string, any>) => {
    const response = await workflowService.executeWorkflow(workflowId, params)
    if (response.data) {
      // 这里可以设置执行状态
      return response.data.executionId
    }
  }

  return {
    currentExecution,
    executionLogs,
    executionStatus,
    executeWorkflow,
  }
}

