<template>
  <el-table
    :data="tableData"
    :stripe="stripe"
    :border="border"
    :size="size"
    :show-header="showHeader"
    :highlight-current-row="highlightCurrentRow"
    :empty-text="emptyText"
    :loading="loading"
    :height="height"
    :max-height="maxHeight"
    v-bind="$attrs"
    @selection-change="handleSelectionChange"
    @row-click="handleRowClick"
    @sort-change="handleSortChange"
  >
    <el-table-column
      v-for="column in columns"
      :key="column.prop || column.label"
      :prop="column.prop"
      :label="column.label"
      :width="column.width"
      :min-width="column.minWidth"
      :fixed="column.fixed"
      :sortable="column.sortable"
      :formatter="column.formatter"
      :align="column.align"
    >
      <template v-if="column.slot && $slots[column.slot]" #[column.slot]="{ row, column: col, $index }">
        <slot :name="column.slot" :row="row" :column="col" :index="$index" />
      </template>
    </el-table-column>
    <template v-if="$slots.empty" #empty>
      <slot name="empty" />
    </template>
  </el-table>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TableProps, TableEmits, TableColumn } from './types'

const props = withDefaults(defineProps<TableProps>(), {
  stripe: false,
  border: false,
  size: 'default',
  showHeader: true,
  highlightCurrentRow: false,
  emptyText: '暂无数据',
  loading: false,
  data: () => [],
})

const emit = defineEmits<TableEmits>()

// 确保 data 始终是数组，避免 Element Plus 内部错误
const tableData = computed(() => {
  if (!props.data) {
    return []
  }
  if (Array.isArray(props.data)) {
    return props.data
  }
  // 如果不是数组，返回空数组并输出警告
  console.warn('Table component: data prop must be an array, got:', typeof props.data)
  return []
})

const handleSelectionChange = (selection: unknown[]) => {
  emit('selection-change', selection)
}

const handleRowClick = (row: unknown, column: TableColumn, event: Event) => {
  emit('row-click', row, column, event)
}

const handleSortChange = (sortInfo: { column: TableColumn; prop: string; order: string }) => {
  emit('sort-change', sortInfo)
}
</script>

