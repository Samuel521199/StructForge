<template>
  <div class="sf-action-bar" :class="[`is-${align}`]">
    <div class="sf-action-bar__content">
      <template v-for="(action, index) in actions" :key="index">
        <Button
          v-if="!action.hidden"
          :type="action.type || 'default'"
          :size="action.size || size"
          :disabled="action.disabled"
          :loading="action.loading"
          :icon="action.icon"
          :round="action.round"
          :plain="action.plain"
          @click="handleAction(action)"
        >
          {{ action.label }}
        </Button>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Button } from '@/components/common/base'
import type { ActionBarProps, ActionItem } from './types'

const props = withDefaults(defineProps<ActionBarProps>(), {
  align: 'right',
  size: 'default',
})

const emit = defineEmits<{
  action: [action: ActionItem, index: number]
}>()

const handleAction = (action: ActionItem) => {
  if (action.onClick) {
    action.onClick()
  }
  const index = props.actions.findIndex(a => a === action)
  emit('action', action, index)
}
</script>

<style scoped lang="scss">
.sf-action-bar {
  display: flex;
  align-items: center;

  &__content {
    display: flex;
    gap: 8px;
  }

  &.is-left {
    justify-content: flex-start;
  }

  &.is-center {
    justify-content: center;
  }

  &.is-right {
    justify-content: flex-end;
  }
}
</style>

