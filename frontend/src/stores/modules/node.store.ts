/**
 * 节点状态管理
 */

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { NodeType } from '@/api/services/node.service'

export const useNodeStore = defineStore('node', () => {
  // 状态
  const nodeTypes = ref<NodeType[]>([])

  // Actions
  function setNodeTypes(types: NodeType[]) {
    nodeTypes.value = types
  }

  function getNodeTypeById(id: string): NodeType | undefined {
    return nodeTypes.value.find(type => type.id === id)
  }

  function getNodeTypesByCategory(category: string): NodeType[] {
    return nodeTypes.value.filter(type => type.category === category)
  }

  return {
    // State
    nodeTypes,
    // Actions
    setNodeTypes,
    getNodeTypeById,
    getNodeTypesByCategory,
  }
})

