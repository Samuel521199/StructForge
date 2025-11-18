/**
 * 工作流验证组合函数
 */

export function useWorkflowValidation() {
  const validateWorkflow = (workflow: any): { valid: boolean; errors: string[] } => {
    const errors: string[] = []

    if (!workflow.name || workflow.name.trim() === '') {
      errors.push('工作流名称不能为空')
    }

    if (!workflow.nodes || workflow.nodes.length === 0) {
      errors.push('工作流至少需要一个节点')
    }

    // 检查是否有触发节点
    const hasTriggerNode = workflow.nodes.some((node: any) => 
      node.type.startsWith('trigger')
    )
    if (!hasTriggerNode) {
      errors.push('工作流必须至少包含一个触发节点')
    }

    return {
      valid: errors.length === 0,
      errors,
    }
  }

  return {
    validateWorkflow,
  }
}

