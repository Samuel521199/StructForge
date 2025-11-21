<template>
  <div
    class="base-node"
    :class="{
      'is-selected': selected,
      [`node-type-${nodeType}`]: nodeType,
    }"
  >
    <!-- 节点头部 -->
    <div class="node-header">
      <div class="node-icon">
        <Icon v-if="icon" :icon="icon" :size="16" />
      </div>
      <div class="node-title">
        <span class="node-label">{{ label }}</span>
        <span v-if="nodeType" class="node-type">{{ nodeType }}</span>
      </div>
    </div>

    <!-- 输入端口 -->
    <div v-if="hasInputs" class="node-inputs">
      <NodePort
        v-for="port in inputs"
        :key="port.id"
        :port-id="port.id"
        :port-name="port.name"
        :port-type="port.type"
        is-input
        @connect="handleInputConnect"
      />
    </div>

    <!-- 节点内容 -->
    <div class="node-content">
      <slot />
    </div>

    <!-- 输出端口 -->
    <div v-if="hasOutputs" class="node-outputs">
      <NodePort
        v-for="port in outputs"
        :key="port.id"
        :port-id="port.id"
        :port-name="port.name"
        :port-type="port.type"
        :is-input="false"
        @connect="handleOutputConnect"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@/components/common/base'
import NodePort from './NodePort.vue'
import type { NodePort as NodePortType } from '@/api/types/workflow.types'

interface BaseNodeProps {
  /** 节点ID */
  id: string
  /** 节点类型 */
  type: string
  /** 节点标签 */
  label: string
  /** 节点图标 */
  icon?: unknown
  /** 是否选中 */
  selected?: boolean
  /** 输入端口列表 */
  inputs?: NodePortType[]
  /** 输出端口列表 */
  outputs?: NodePortType[]
  /** 节点数据 */
  data?: Record<string, unknown>
}

const props = withDefaults(defineProps<BaseNodeProps>(), {
  selected: false,
  inputs: () => [],
  outputs: () => [],
  data: () => ({}),
})

const emit = defineEmits<{
  connect: [portId: string, isInput: boolean]
  click: []
  delete: []
}>()

const nodeType = computed(() => props.type)
const hasInputs = computed(() => props.inputs && props.inputs.length > 0)
const hasOutputs = computed(() => props.outputs && props.outputs.length > 0)

const handleInputConnect = (portId: string) => {
  emit('connect', portId, true)
}

const handleOutputConnect = (portId: string) => {
  emit('connect', portId, false)
}
</script>

<style scoped lang="scss">
.base-node {
  min-width: 180px;
  background: #fff;
  border: 2px solid var(--el-border-color);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
  cursor: pointer;

  &:hover {
    border-color: var(--el-color-primary);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }

  &.is-selected {
    border-color: var(--el-color-primary);
    box-shadow: 0 0 0 2px var(--el-color-primary-light-8);
  }

  .node-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px;
    background: var(--el-fill-color-light);
    border-bottom: 1px solid var(--el-border-color-light);
    border-radius: 6px 6px 0 0;

    .node-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 24px;
      height: 24px;
      color: var(--el-color-primary);
    }

    .node-title {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 2px;

      .node-label {
        font-size: 14px;
        font-weight: 500;
        color: var(--el-text-color-primary);
      }

      .node-type {
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }
  }

  .node-inputs {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 8px 0;
  }

  .node-content {
    padding: 12px;
    min-height: 40px;
  }

  .node-outputs {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 8px 0;
  }

  // 不同节点类型的样式
  &.node-type-trigger {
    border-left: 4px solid #67c23a;
  }

  &.node-type-ai {
    border-left: 4px solid #409eff;
  }

  &.node-type-data {
    border-left: 4px solid #e6a23c;
  }

  &.node-type-control {
    border-left: 4px solid #f56c6c;
  }

  &.node-type-integration {
    border-left: 4px solid #909399;
  }
}
</style>
