<template>
  <div
    class="node-port"
    :class="{
      'is-input': isInput,
      'is-output': !isInput,
      [`port-type-${portType}`]: portType,
    }"
    @mousedown.stop="handleMouseDown"
    @mouseup.stop="handleMouseUp"
  >
    <div class="port-label" v-if="portName">
      <span v-if="isInput" class="label-text">{{ portName }}</span>
    </div>
    <div class="port-dot" :class="{ 'is-connected': isConnected }" />
    <div class="port-label" v-if="portName">
      <span v-if="!isInput" class="label-text">{{ portName }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface NodePortProps {
  /** 端口ID */
  portId: string
  /** 端口名称 */
  portName?: string
  /** 端口类型 */
  portType?: string
  /** 是否为输入端口 */
  isInput?: boolean
}

const props = withDefaults(defineProps<NodePortProps>(), {
  isInput: true,
  portType: 'default',
})

const emit = defineEmits<{
  connect: [portId: string]
}>()

const isConnected = ref(false)

const handleMouseDown = (event: MouseEvent) => {
  if (!props.isInput) {
    // 输出端口：开始连接
    emit('connect', props.portId)
  }
  event.stopPropagation()
}

const handleMouseUp = (event: MouseEvent) => {
  if (props.isInput) {
    // 输入端口：完成连接
    emit('connect', props.portId)
    isConnected.value = true
  }
  event.stopPropagation()
}
</script>

<style scoped lang="scss">
.node-port {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  position: relative;
  cursor: crosshair;

  &.is-input {
    justify-content: flex-start;

    .port-label {
      order: 1;
    }

    .port-dot {
      order: 2;
    }
  }

  &.is-output {
    justify-content: flex-end;

    .port-label {
      order: 2;
    }

    .port-dot {
      order: 1;
    }
  }

  .port-label {
    .label-text {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }

  .port-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background-color: var(--el-border-color);
    border: 2px solid #fff;
    transition: all 0.2s;
    position: relative;
    z-index: 10;

    &:hover {
      background-color: var(--el-color-primary);
      transform: scale(1.2);
    }

    &.is-connected {
      background-color: var(--el-color-primary);
    }
  }

  // 不同端口类型的颜色
  &.port-type-string .port-dot {
    border-color: #67c23a;
  }

  &.port-type-number .port-dot {
    border-color: #409eff;
  }

  &.port-type-boolean .port-dot {
    border-color: #e6a23c;
  }

  &.port-type-object .port-dot {
    border-color: #f56c6c;
  }
}
</style>
