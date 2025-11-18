<template>
  <div class="workflow-editor-page">
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <div class="toolbar-left">
        <Button :icon="ArrowLeft" @click="handleBack">返回</Button>
        <span class="workflow-name">{{ workflowName || '未命名工作流' }}</span>
      </div>
      <div class="toolbar-right">
        <Button @click="handleSave" :loading="saving">保存</Button>
        <Button type="primary" @click="handleRun" :loading="running">运行</Button>
      </div>
    </div>

    <!-- 编辑器主体 -->
    <div class="editor-main">
      <!-- 节点面板 -->
      <NodePalette />

      <!-- 画布区域 -->
      <div class="editor-canvas">
        <WorkflowEditor
          v-model:nodes="nodes"
          v-model:edges="edges"
          @node-click="handleNodeClick"
          @node-add="handleNodeAdd"
          @node-delete="handleNodeDelete"
        />
      </div>

      <!-- 属性面板 -->
      <div class="editor-properties">
        <Card v-if="selectedNode" header="节点属性">
          <div class="properties-content">
            <Form :model="nodeProperties" label-width="80px">
              <FormItem label="节点ID">
                <Input :model-value="selectedNode.id" disabled />
              </FormItem>
              <FormItem label="节点类型">
                <Input :model-value="selectedNode.type" disabled />
              </FormItem>
              <FormItem label="节点名称">
                <Input
                  :model-value="selectedNode.data?.label || selectedNode.label"
                  @update:model-value="updateNodeLabel"
                />
              </FormItem>
              <!-- 根据节点类型显示不同的属性 -->
              <template v-if="selectedNode.type === 'http'">
                <FormItem label="URL">
                  <Input
                    :model-value="selectedNode.data?.url || ''"
                    placeholder="请输入URL"
                    @update:model-value="updateNodeData('url', $event)"
                  />
                </FormItem>
                <FormItem label="方法">
                  <Select
                    :model-value="selectedNode.data?.method || 'GET'"
                    :options="[
                      { label: 'GET', value: 'GET' },
                      { label: 'POST', value: 'POST' },
                      { label: 'PUT', value: 'PUT' },
                      { label: 'DELETE', value: 'DELETE' },
                    ]"
                    @update:model-value="updateNodeData('method', $event)"
                  />
                </FormItem>
              </template>
            </Form>
          </div>
        </Card>
        <Empty v-else description="请选择一个节点" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { Button, Card, Form, FormItem, Input, Select, Empty } from '@/components/common'
import { NodePalette } from '@/components/workflow/editor/NodePalette'
import { WorkflowEditor } from '@/components/workflow/editor/WorkflowEditor'
import { useWorkflowStore } from '@/stores/modules/workflow.store'
import { success, error } from '@/components/common/base/Message'

const route = useRoute()
const router = useRouter()
const workflowStore = useWorkflowStore()

const workflowName = ref('')
const saving = ref(false)
const running = ref(false)
const selectedNode = ref<any>(null)

const nodes = ref<any[]>([])
const edges = ref<any[]>([])

// 更新节点标签
const updateNodeLabel = (label: string) => {
  if (selectedNode.value) {
    if (!selectedNode.value.data) {
      selectedNode.value.data = {}
    }
    selectedNode.value.data.label = label
    selectedNode.value.label = label
    // 更新nodes数组中的节点
    const index = nodes.value.findIndex(n => n.id === selectedNode.value!.id)
    if (index !== -1) {
      nodes.value[index] = { ...selectedNode.value }
    }
  }
}

// 更新节点数据
const updateNodeData = (key: string, value: any) => {
  if (selectedNode.value) {
    if (!selectedNode.value.data) {
      selectedNode.value.data = {}
    }
    selectedNode.value.data[key] = value
    // 更新nodes数组中的节点
    const index = nodes.value.findIndex(n => n.id === selectedNode.value!.id)
    if (index !== -1) {
      nodes.value[index] = { ...selectedNode.value }
    }
  }
}

// 加载工作流
const loadWorkflow = async () => {
  const workflowId = route.params.id as string
  if (!workflowId) {
    // 新建工作流
    nodes.value = []
    edges.value = []
    return
  }

  try {
    // TODO: 调用API加载工作流
    // const response = await workflowService.getWorkflow(workflowId)
    // workflowName.value = response.data.name
    // nodes.value = response.data.nodes || []
    // edges.value = response.data.edges || []
  } catch (err) {
    error('加载工作流失败')
    console.error('Load workflow error:', err)
  }
}

const handleBack = () => {
  router.push('/workflow/list')
}

const handleSave = async () => {
  saving.value = true
  try {
    // TODO: 调用API保存工作流
    // await workflowService.saveWorkflow({
    //   id: route.params.id as string,
    //   name: workflowName.value,
    //   nodes: nodes.value,
    //   edges: edges.value,
    // })
    success('保存成功')
  } catch (err) {
    error('保存失败')
    console.error('Save workflow error:', err)
  } finally {
    saving.value = false
  }
}

const handleRun = async () => {
  running.value = true
  try {
    // TODO: 调用API运行工作流
    // const response = await workflowService.runWorkflow({
    //   id: route.params.id as string,
    // })
    success('工作流已启动')
  } catch (err) {
    error('运行失败')
    console.error('Run workflow error:', err)
  } finally {
    running.value = false
  }
}

const handleNodeClick = (node: any) => {
  selectedNode.value = node
}

const handleNodeAdd = (node: any) => {
  nodes.value.push(node)
}

const handleNodeDelete = (nodeId: string) => {
  nodes.value = nodes.value.filter(n => n.id !== nodeId)
  edges.value = edges.value.filter(
    e => e.source !== nodeId && e.target !== nodeId
  )
  if (selectedNode.value?.id === nodeId) {
    selectedNode.value = null
  }
}

onMounted(() => {
  loadWorkflow()
})
</script>

<style scoped lang="scss">
.workflow-editor-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--el-bg-color-page);
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 50px;
  padding: 0 16px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);

  .toolbar-left {
    display: flex;
    align-items: center;
    gap: 16px;

    .workflow-name {
      font-size: 16px;
      font-weight: 500;
      color: var(--el-text-color-primary);
    }
  }

  .toolbar-right {
    display: flex;
    gap: 8px;
  }
}

.editor-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-canvas {
  flex: 1;
  position: relative;
}

.editor-properties {
  width: 300px;
  background: var(--el-bg-color);
  border-left: 1px solid var(--el-border-color-light);
  padding: 16px;
  overflow-y: auto;

  .properties-content {
    padding: 16px 0;
  }
}
</style>
