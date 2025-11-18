<template>
  <div class="node-port" :class="{ 'is-input': isInput, 'is-output': !isInput }">
    <!-- TODO: 实现节点端口组件 -->
    <div class="port-dot" />
  </div>
</template>

<script setup lang="ts">
interface NodePortProps {
  isInput?: boolean
  portId: string
}

defineProps<NodePortProps>()
</script>

<style scoped lang="scss">
.node-port {
  .port-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: #409eff;
  }
}
</style>

