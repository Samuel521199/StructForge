<template>
  <el-table
    :data="data"
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
    />
    <template v-if="$slots.empty" #empty>
      <slot name="empty" />
    </template>
  </el-table>
</template>

<script setup lang="ts">
import type { TableProps, TableEmits, TableColumn } from './types'

const props = withDefaults(defineProps<TableProps>(), {
  stripe: false,
  border: false,
  size: 'default',
  showHeader: true,
  highlightCurrentRow: false,
  emptyText: '暂无数据',
  loading: false,
})

const emit = defineEmits<TableEmits>()

const handleSelectionChange = (selection: any[]) => {
  emit('selection-change', selection)
}

const handleRowClick = (row: any, column: TableColumn, event: Event) => {
  emit('row-click', row, column, event)
}

const handleSortChange = (sortInfo: { column: TableColumn; prop: string; order: string }) => {
  emit('sort-change', sortInfo)
}
</script>

