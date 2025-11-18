<template>
  <div class="workflow-list-page">
    <Card>
      <template #header>
        <div class="list-header">
          <h2>工作流列表</h2>
          <Button type="primary" :icon="Plus" @click="handleCreate">
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
        <el-pagination
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
import { Card, Button, Table, SearchBar, StatusTag } from '@/components/common'
import { useWorkflowStore } from '@/stores/modules/workflow.store'
import { success, error } from '@/components/common/base/Message'
import type { TableColumn } from '@/components/common/base/Table/types'

const router = useRouter()
const workflowStore = useWorkflowStore()

const loading = ref(false)
const searchKeyword = ref('')
const workflowList = ref<any[]>([])

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
    // TODO: 调用API获取工作流列表
    // const response = await workflowService.getWorkflowList({
    //   page: pagination.page,
    //   pageSize: pagination.pageSize,
    //   keyword: searchKeyword.value,
    // })
    // workflowList.value = response.data.list
    // pagination.total = response.data.total

    // 模拟数据
    workflowList.value = [
      {
        id: '1',
        name: '示例工作流1',
        description: '这是一个示例工作流',
        status: 'running',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    ]
    pagination.total = 1
  } catch (err) {
    error('加载工作流列表失败')
    console.error('Load workflows error:', err)
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

const handleEdit = (row: any) => {
  router.push(`/workflow/editor/${row.id}`)
}

const handleDelete = async (row: any) => {
  // TODO: 实现删除确认对话框
  try {
    // await workflowService.deleteWorkflow(row.id)
    success('删除成功')
    loadWorkflows()
  } catch (err) {
    error('删除失败')
  }
}

const handleRowClick = (row: any) => {
  router.push(`/workflow/detail/${row.id}`)
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
