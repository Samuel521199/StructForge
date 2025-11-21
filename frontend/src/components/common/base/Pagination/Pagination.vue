<template>
  <el-pagination
    v-model:current-page="currentPage"
    v-model:page-size="pageSize"
    :total="total"
    :page-sizes="pageSizes"
    :layout="layout"
    :background="background"
    :disabled="disabled"
    :hide-on-single-page="hideOnSinglePage"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
    @prev-click="handlePrevClick"
    @next-click="handleNextClick"
  />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { PaginationProps, PaginationEmits } from './types'

const props = withDefaults(defineProps<PaginationProps>(), {
  currentPage: 1,
  pageSize: 10,
  total: 0,
  pageSizes: () => [10, 20, 50, 100],
  layout: 'total, sizes, prev, pager, next, jumper',
  background: false,
  disabled: false,
  hideOnSinglePage: false,
})

const emit = defineEmits<PaginationEmits>()

const currentPage = ref(props.currentPage)
const pageSize = ref(props.pageSize)

watch(() => props.currentPage, (val) => {
  currentPage.value = val
})

watch(() => props.pageSize, (val) => {
  pageSize.value = val
})

const handleSizeChange = (size: number) => {
  pageSize.value = size
  emit('size-change', size)
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  emit('current-change', page)
}

const handlePrevClick = (page: number) => {
  emit('prev-click', page)
}

const handleNextClick = (page: number) => {
  emit('next-click', page)
}

defineOptions({
  name: 'Pagination',
})
</script>

