/**
 * 工作流验证工具
 */

export const workflowValidation = {
  /**
   * 验证工作流
   */
  validate(workflow: any): { valid: boolean; errors: string[] } {
    const errors: string[] = []

    if (!workflow.name || workflow.name.trim() === '') {
      errors.push('工作流名称不能为空')
    }

    if (!workflow.nodes || workflow.nodes.length === 0) {
      errors.push('工作流至少需要一个节点')
    }

    return {
      valid: errors.length === 0,
      errors,
    }
  },
}

