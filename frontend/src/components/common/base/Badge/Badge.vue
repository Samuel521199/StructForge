<template>
  <div class="sf-badge" :class="badgeClass">
    <slot />
    <span
      v-if="value !== undefined && value !== null"
      class="badge-content"
      :class="{
        'is-dot': dot,
        [`badge-${type}`]: type,
      }"
      :style="badgeStyle"
    >
      <template v-if="!dot">{{ formattedValue }}</template>
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { BadgeProps } from './types'

const props = withDefaults(defineProps<BadgeProps>(), {
  dot: false,
  type: 'danger',
  max: 99,
})

const formattedValue = computed(() => {
  if (props.dot) return ''
  if (typeof props.value === 'number' && props.value > props.max) {
    return `${props.max}+`
  }
  return props.value
})

const badgeClass = computed(() => {
  return {
    'is-fixed': props.fixed,
  }
})

const badgeStyle = computed(() => {
  if (props.color) {
    return {
      backgroundColor: props.color,
    }
  }
  return {}
})
</script>

<style scoped lang="scss">
.sf-badge {
  position: relative;
  display: inline-block;

  &.is-fixed {
    .badge-content {
      position: absolute;
      top: 0;
      right: 0;
      transform: translate(50%, -50%);
    }
  }

  .badge-content {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
    height: 18px;
    padding: 0 6px;
    font-size: 12px;
    font-weight: 500;
    color: #fff;
    background-color: var(--el-color-danger);
    border-radius: 9px;
    white-space: nowrap;

    &.is-dot {
      width: 8px;
      height: 8px;
      min-width: 8px;
      padding: 0;
      border-radius: 50%;
    }

    &.badge-primary {
      background-color: var(--el-color-primary);
    }

    &.badge-success {
      background-color: var(--el-color-success);
    }

    &.badge-warning {
      background-color: var(--el-color-warning);
    }

    &.badge-danger {
      background-color: var(--el-color-danger);
    }

    &.badge-info {
      background-color: var(--el-color-info);
    }
  }
}
</style>

