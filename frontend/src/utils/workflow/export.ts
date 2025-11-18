/**
 * 工作流导出工具
 */

export const workflowExport = {
  /**
   * 导出为JSON
   */
  toJSON(workflow: any): string {
    return JSON.stringify(workflow, null, 2)
  },

  /**
   * 导出为文件
   */
  toFile(workflow: any, filename: string) {
    const json = this.toJSON(workflow)
    const blob = new Blob([json], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    link.click()
    URL.revokeObjectURL(url)
  },
}

