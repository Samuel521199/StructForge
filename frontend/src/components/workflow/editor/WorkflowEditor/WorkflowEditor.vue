<template>
  <div
    class="workflow-editor"
    ref="editorRef"
    @drop="handleDrop"
    @dragover="handleDragOver"
  >
    <VueFlow
      v-model="localNodes"
      v-model:edges="localEdges"
      :default-viewport="{ zoom: 1 }"
      :min-zoom="0.2"
      :max-zoom="4"
      :fit-view-on-init="true"
      @nodes-change="handleNodesChange"
      @edges-change="handleEdgesChange"
      @node-click="handleNodeClick"
      @connect="handleConnect"
      @pane-click="handlePaneClick"
    >
      <!-- 节点模板 -->
      <template #node-default="{ data }">
        <div class="workflow-node">
          <div class="workflow-node__header">
            <el-icon class="workflow-node__icon">
              <Circle />
            </el-icon>
            <span class="workflow-node__label">{{ data.label || data.type }}</span>
          </div>
          <div class="workflow-node__body">
            <Handle
              v-for="(handle, index) in data.inputs"
              :key="`input-${index}`"
              :id="handle.id"
              type="target"
              :position="Position.Left"
              class="workflow-node__handle"
            />
            <Handle
              v-for="(handle, index) in data.outputs"
              :key="`output-${index}`"
              :id="handle.id"
              type="source"
              :position="Position.Right"
              class="workflow-node__handle"
            />
          </div>
        </div>
      </template>

      <!-- 背景 -->
      <Background />
      <!-- 小地图 -->
      <MiniMap />
      <!-- 控制面板 -->
      <Controls />
    </VueFlow>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Circle } from '@element-plus/icons-vue'
import { VueFlow, Background, MiniMap, Controls, Handle, Position, useVueFlow } from '@vue-flow/core'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'

interface Node {
  id: string
  type: string
  label?: string
  position: { x: number; y: number }
  data: any
  inputs?: Array<{ id: string; label: string }>
  outputs?: Array<{ id: string; label: string }>
}

interface Edge {
  id: string
  source: string
  target: string
  sourceHandle?: string
  targetHandle?: string
}

interface Props {
  nodes?: Node[]
  edges?: Edge[]
}

interface Emits {
  (e: 'update:nodes', nodes: Node[]): void
  (e: 'update:edges', edges: Edge[]): void
  (e: 'node-click', node: Node): void
  (e: 'node-add', node: Node): void
  (e: 'node-delete', nodeId: string): void
}

const props = withDefaults(defineProps<Props>(), {
  nodes: () => [],
  edges: () => [],
})

const emit = defineEmits<Emits>()

const editorRef = ref<HTMLElement>()

const localNodes = computed({
  get: () => props.nodes,
  set: (value) => emit('update:nodes', value),
})

const localEdges = computed({
  get: () => props.edges,
  set: (value) => emit('update:edges', value),
})

// 获取节点图标
const getNodeIcon = (type: string) => {
  const iconMap: Record<string, string> = {
    webhook: 'Link',
    timer: 'Clock',
    manual: 'Pointer',
    'text-generation': 'Document',
    'image-generation': 'Picture',
    'code-generation': 'Code',
    transform: 'Refresh',
    filter: 'Filter',
    aggregate: 'DataAnalysis',
    condition: 'Switch',
    loop: 'RefreshRight',
    parallel: 'Grid',
    http: 'Connection',
    database: 'Database',
    file: 'Folder',
    script: 'Document',
    'code-executor': 'Cpu',
  }
  return iconMap[type] || 'Circle'
}

const handleNodesChange = (changes: any[]) => {
  // 处理节点变化
}

const handleEdgesChange = (changes: any[]) => {
  // 处理边变化
}

const handleNodeClick = (event: any) => {
  emit('node-click', event.node)
}


const handleConnect = (connection: any) => {
  const newEdge: Edge = {
    id: `${connection.source}-${connection.target}`,
    source: connection.source,
    target: connection.target,
    sourceHandle: connection.sourceHandle,
    targetHandle: connection.targetHandle,
  }
  const edges = [...localEdges.value, newEdge]
  emit('update:edges', edges)
}

const handlePaneClick = () => {
  // 点击画布空白处
}

// 处理拖放
const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'copy'
  }
}

const { screenToFlowCoordinate } = useVueFlow()

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  
  if (!event.dataTransfer) return

  try {
    const nodeData = JSON.parse(event.dataTransfer.getData('application/json'))
    
    // 获取画布位置
    const position = screenToFlowCoordinate({
      x: event.clientX,
      y: event.clientY,
    })

    // 创建新节点
    const newNode: Node = {
      id: `${nodeData.type}-${Date.now()}`,
      type: nodeData.type,
      label: nodeData.label,
      position,
      data: {
        type: nodeData.type,
        label: nodeData.label,
        inputs: getDefaultInputs(nodeData.type),
        outputs: getDefaultOutputs(nodeData.type),
      },
    }

    const nodes = [...localNodes.value, newNode]
    emit('update:nodes', nodes)
    emit('node-add', newNode)
  } catch (error) {
    console.error('Drop node error:', error)
  }
}

// 获取默认输入端口
const getDefaultInputs = (type: string): Array<{ id: string; label: string }> => {
  // 触发节点没有输入
  if (['webhook', 'timer', 'manual'].includes(type)) {
    return []
  }
  // 其他节点默认有一个输入
  return [{ id: 'input-1', label: '输入' }]
}

// 获取默认输出端口
const getDefaultOutputs = (type: string): Array<{ id: string; label: string }> => {
  // 所有节点默认有一个输出
  return [{ id: 'output-1', label: '输出' }]
}
</script>

<style scoped lang="scss">
.workflow-editor {
  width: 100%;
  height: 100%;
  position: relative;
}

.workflow-node {
  background: var(--el-bg-color);
  border: 2px solid var(--el-border-color);
  border-radius: 8px;
  padding: 8px 12px;
  min-width: 120px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;

  &:hover {
    border-color: var(--el-color-primary);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  &__header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }

  &__icon {
    font-size: 16px;
    color: var(--el-color-primary);
  }

  &__label {
    font-size: 14px;
    font-weight: 500;
    color: var(--el-text-color-primary);
  }

  &__body {
    position: relative;
  }

  &__handle {
    width: 8px;
    height: 8px;
    background: var(--el-color-primary);
    border: 2px solid var(--el-bg-color);
    border-radius: 50%;
  }
}
</style>
