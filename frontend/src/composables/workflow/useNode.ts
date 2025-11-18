/**
 * 节点操作组合函数
 */

import { computed } from 'vue'
import { useNodeStore } from '@/stores/modules/node.store'
import { nodeService } from '@/api/services/node.service'

export function useNode() {
  const nodeStore = useNodeStore()

  const nodeTypes = computed(() => nodeStore.nodeTypes)

  const loadNodeTypes = async () => {
    const response = await nodeService.getNodeTypes()
    if (response.data) {
      nodeStore.setNodeTypes(response.data)
    }
  }

  const getNodeTypeById = (id: string) => {
    return nodeStore.getNodeTypeById(id)
  }

  const getNodeTypesByCategory = (category: string) => {
    return nodeStore.getNodeTypesByCategory(category)
  }

  return {
    nodeTypes,
    loadNodeTypes,
    getNodeTypeById,
    getNodeTypesByCategory,
  }
}

