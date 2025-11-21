<template>
  <div class="workflow-list-page">
    <Card>
      <template #header>
        <div class="list-header">
          <h2>工作流列表</h2>
          <Button type="primary" @click="handleCreate">
            <template #icon>
              <Icon :icon="Plus" />
            </template>
            创建工作流
          </Button>
        </div>
      </template>

      <!-- 搜索和筛选 -->
      <div class="list-filters">
        <SearchBar
          v-model="searchKeyword"
          :filters="filters"
          filterable
          @search="handleSearch"
          @filter="handleFilter"
        />
      </div>

      <!-- 工作流列表 -->
      <Table
        :data="workflowList"
        :columns="columns"
        :loading="loading"
        stripe
        @row-click="handleRowClick"
      >
        <template #status="{ row }">
          <StatusTag :status="row.status" />
        </template>

        <template #actions="{ row }">
          <Button
            type="primary"
            link
            size="small"
            @click.stop="handleEdit(row)"
          >
            编辑
          </Button>
          <Button
            type="danger"
            link
            size="small"
            @click.stop="handleDelete(row)"
          >
            删除
          </Button>
        </template>
      </Table>

      <!-- 分页 -->
      <div class="list-pagination">
        <Pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { Card, Button, Table, Pagination, SearchBar, StatusTag } from '@/components/common'
import { success } from '@/components/common/base/Message'
import { workflowService } from '@/api/services/workflow.service'
import type { TableColumn } from '@/components/common/base/Table/types'
import type { WorkflowListItem } from '@/api/types/workflow.types'

const router = useRouter()

const loading = ref(false)
const searchKeyword = ref('')
const workflowList = ref<WorkflowListItem[]>([])

const filters = [
  {
    label: '状态',
    prop: 'status',
    type: 'select' as const,
    options: [
      { label: '运行中', value: 'running' },
      { label: '已停止', value: 'stopped' },
      { label: '已暂停', value: 'paused' },
      { label: '错误', value: 'error' },
    ],
  },
]

const columns: TableColumn[] = [
  { prop: 'name', label: '名称', width: 200 },
  { prop: 'description', label: '描述', minWidth: 200 },
  { prop: 'status', label: '状态', width: 120, slot: 'status' },
  { prop: 'createdAt', label: '创建时间', width: 180 },
  { prop: 'updatedAt', label: '更新时间', width: 180 },
  { label: '操作', width: 150, fixed: 'right', slot: 'actions' },
]

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 加载工作流列表
const loadWorkflows = async () => {
  loading.value = true
  try {
    const response = await workflowService.getWorkflows({
      page: pagination.page,
      pageSize: pagination.pageSize,
      search: searchKeyword.value,
    })
    
    if (response.data) {
      workflowList.value = response.data.list
      pagination.total = response.data.total
    } else {
      // 如果没有数据，使用空数组
      workflowList.value = []
      pagination.total = 0
    }
  } catch (err) {
    // 错误已在 errorHandler 中处理，这里只记录日志
    console.error('Load workflows error:', err)
    // 如果请求失败，使用空数组
    workflowList.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

const handleSearch = (keyword: string) => {
  searchKeyword.value = keyword
  pagination.page = 1
  loadWorkflows()
}

const handleFilter = (filters: Record<string, any>) => {
  console.log('Filter:', filters)
  pagination.page = 1
  loadWorkflows()
}

const handleCreate = () => {
  router.push('/workflow/editor')
}

const handleEdit = (row: WorkflowListItem) => {
  router.push(`/workflow/editor/${row.id}`)
}

const handleDelete = async (row: WorkflowListItem) => {
  // TODO: 实现删除确认对话框
  try {
    await workflowService.deleteWorkflow(row.id)
    success('删除成功')
    loadWorkflows()
  } catch (err) {
    // 错误已在 errorHandler 中处理，这里只记录日志
    console.error('Delete workflow error:', err)
  }
}

const handleRowClick = (row: unknown) => {
  const workflow = row as WorkflowListItem
  router.push(`/workflow/detail/${workflow.id}`)
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadWorkflows()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadWorkflows()
}

onMounted(() => {
  loadWorkflows()
})
</script>

<style scoped lang="scss">
.workflow-list-page {
  padding: 24px;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.list-filters {
  margin-bottom: 16px;
}

.list-pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
