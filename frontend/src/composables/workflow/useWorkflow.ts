/**
 * 工作流操作组合函数
 */

import { computed } from 'vue'
import { useWorkflowStore } from '@/stores/modules/workflow.store'
import { workflowService } from '@/api/services/workflow.service'

export function useWorkflow() {
  const workflowStore = useWorkflowStore()

  const workflows = computed(() => workflowStore.workflows)
  const currentWorkflow = computed(() => workflowStore.currentWorkflow)

  const loadWorkflows = async () => {
    const response = await workflowService.getWorkflows()
    if (response.data) {
      workflowStore.setWorkflows(response.data.list)
    }
  }

  const createWorkflow = async (data: any) => {
    const response = await workflowService.createWorkflow(data)
    if (response.data) {
      workflowStore.addWorkflow(response.data)
      return response.data
    }
  }

  const updateWorkflow = async (id: string, data: any) => {
    const response = await workflowService.updateWorkflow(id, data)
    if (response.data) {
      workflowStore.updateWorkflow(id, response.data)
    }
  }

  const deleteWorkflow = async (id: string) => {
    await workflowService.deleteWorkflow(id)
    workflowStore.removeWorkflow(id)
  }

  return {
    workflows,
    currentWorkflow,
    loadWorkflows,
    createWorkflow,
    updateWorkflow,
    deleteWorkflow,
  }
}

