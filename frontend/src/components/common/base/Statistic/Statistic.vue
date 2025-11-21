<template>
  <div class="sf-statistic">
    <div class="statistic-title" v-if="title">
      <Icon v-if="icon" :icon="icon" :size="16" class="title-icon" />
      <span>{{ title }}</span>
    </div>
    <div class="statistic-content">
      <div class="statistic-value" :style="{ color: valueColor }">
        <span class="value-number">{{ formattedValue }}</span>
        <span v-if="suffix" class="value-suffix">{{ suffix }}</span>
      </div>
      <div v-if="description" class="statistic-description">
        {{ description }}
      </div>
      <div v-if="trend" class="statistic-trend" :class="`trend-${trend.direction}`">
        <Icon :icon="trend.direction === 'up' ? ArrowUp : ArrowDown" :size="12" />
        <span>{{ trend.value }}%</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@/components/common/base'
import { ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import type { StatisticProps } from './types'

const props = withDefaults(defineProps<StatisticProps>(), {
  precision: 0,
  valueColor: '#00FF00',
})

const formattedValue = computed(() => {
  if (typeof props.value === 'number') {
    return props.value.toFixed(props.precision)
  }
  return props.value
})
</script>

<style scoped lang="scss">
.sf-statistic {
  .statistic-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: var(--el-text-color-secondary);
    margin-bottom: 12px;

    .title-icon {
      color: var(--el-color-primary);
    }
  }

  .statistic-content {
    .statistic-value {
      display: flex;
      align-items: baseline;
      gap: 4px;
      font-size: 32px;
      font-weight: 600;
      line-height: 1.2;

      .value-number {
        font-variant-numeric: tabular-nums;
      }

      .value-suffix {
        font-size: 16px;
        font-weight: 400;
        margin-left: 4px;
      }
    }

    .statistic-description {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-top: 8px;
    }

    .statistic-trend {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;
      margin-top: 8px;

      &.trend-up {
        color: #67c23a;
      }

      &.trend-down {
        color: #f56c6c;
      }
    }
  }
}
</style>

