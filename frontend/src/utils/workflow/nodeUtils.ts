/**
 * 节点工具函数
 */

export const nodeUtils = {
  /**
   * 生成节点ID
   */
  generateNodeId(): string {
    return `node_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  },

  /**
   * 生成边ID
   */
  generateEdgeId(): string {
    return `edge_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  },
}

