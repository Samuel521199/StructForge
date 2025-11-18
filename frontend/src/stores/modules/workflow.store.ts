/**
 * 工作流状态管理
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Workflow } from '@/api/types/workflow.types'

export const useWorkflowStore = defineStore('workflow', () => {
  // 状态
  const workflows = ref<Workflow[]>([])
  const currentWorkflow = ref<Workflow | null>(null)
  const currentNode = ref<any>(null)

  // Actions
  function setWorkflows(list: Workflow[]) {
    workflows.value = list
  }

  function setCurrentWorkflow(workflow: Workflow | null) {
    currentWorkflow.value = workflow
  }

  function setCurrentNode(node: any) {
    currentNode.value = node
  }

  function addWorkflow(workflow: Workflow) {
    workflows.value.push(workflow)
  }

  function updateWorkflow(id: string, data: Partial<Workflow>) {
    const index = workflows.value.findIndex(w => w.id === id)
    if (index !== -1) {
      workflows.value[index] = { ...workflows.value[index], ...data }
    }
    if (currentWorkflow.value?.id === id) {
      currentWorkflow.value = { ...currentWorkflow.value, ...data }
    }
  }

  function removeWorkflow(id: string) {
    workflows.value = workflows.value.filter(w => w.id !== id)
    if (currentWorkflow.value?.id === id) {
      currentWorkflow.value = null
    }
  }

  return {
    // State
    workflows,
    currentWorkflow,
    currentNode,
    // Actions
    setWorkflows,
    setCurrentWorkflow,
    setCurrentNode,
    addWorkflow,
    updateWorkflow,
    removeWorkflow,
  }
})

