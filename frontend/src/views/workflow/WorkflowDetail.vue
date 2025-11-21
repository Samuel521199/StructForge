<template>
  <div class="workflow-detail-page">
    <!-- 页面头部 -->
    <div class="detail-header">
      <div class="header-left">
        <Button :icon="ArrowLeft" @click="handleBack">返回</Button>
        <div class="workflow-info">
          <h1 class="workflow-name">{{ workflow?.name || '加载中...' }}</h1>
          <div class="workflow-meta">
            <StatusTag :status="workflow?.status || 'draft'" />
            <span class="meta-item">
              <Icon :icon="Clock" :size="14" />
              创建于 {{ formatTime(workflow?.createdAt) }}
            </span>
            <span class="meta-item">
              <Icon :icon="Refresh" :size="14" />
              更新于 {{ formatTime(workflow?.updatedAt) }}
            </span>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <Button @click="handleEdit">编辑</Button>
        <Button type="primary" @click="handleRun" :loading="running">运行</Button>
      </div>
    </div>

    <!-- 工作流内容 -->
    <div v-if="loading" class="loading-container">
      <Loading />
    </div>

    <div v-else-if="workflow" class="detail-content">
      <!-- 基本信息 -->
      <Card class="info-card">
        <template #header>
          <span>基本信息</span>
        </template>
        <div class="info-grid">
          <div class="info-item">
            <label>工作流ID</label>
            <div class="info-value">
              <span>{{ workflow.id }}</span>
              <CopyButton :value="workflow.id" />
            </div>
          </div>
          <div class="info-item">
            <label>名称</label>
            <div class="info-value">{{ workflow.name }}</div>
          </div>
          <div class="info-item" v-if="workflow.description">
            <label>描述</label>
            <div class="info-value">{{ workflow.description }}</div>
          </div>
          <div class="info-item">
            <label>状态</label>
            <div class="info-value">
              <StatusTag :status="workflow.status" />
            </div>
          </div>
          <div class="info-item">
            <label>节点数量</label>
            <div class="info-value">{{ workflow.nodes?.length || 0 }}</div>
          </div>
          <div class="info-item">
            <label>连接数量</label>
            <div class="info-value">{{ workflow.edges?.length || 0 }}</div>
          </div>
        </div>
      </Card>

      <!-- 节点列表 -->
      <Card class="nodes-card">
        <template #header>
          <span>节点列表</span>
        </template>
        <Table
          :data="workflow.nodes || []"
          :columns="nodeColumns"
          :empty-text="'暂无节点'"
        >
          <template #type="{ row }">
            <Tag :type="getNodeTypeColor(row.type)">{{ row.type }}</Tag>
          </template>
        </Table>
      </Card>

      <!-- 执行历史 -->
      <Card class="executions-card">
        <template #header>
          <div class="card-header">
            <span>执行历史</span>
            <Link type="primary" @click="goToExecutions">查看全部</Link>
          </div>
        </template>
        <Table
          :data="executionList"
          :columns="executionColumns"
          :loading="loadingExecutions"
          stripe
          @row-click="handleExecutionClick"
        >
          <template #status="{ row }">
            <StatusTag :status="row.status" />
          </template>
          <template #duration="{ row }">
            <span v-if="row.duration">{{ formatDuration(row.duration) }}</span>
            <span v-else>-</span>
          </template>
        </Table>
        <div class="table-pagination">
          <Pagination
            v-model:current-page="executionPagination.page"
            v-model:page-size="executionPagination.pageSize"
            :total="executionPagination.total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @size-change="handleExecutionSizeChange"
            @current-change="handleExecutionPageChange"
          />
        </div>
      </Card>
    </div>

    <Empty v-else description="工作流不存在" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Clock, Refresh } from '@element-plus/icons-vue'
import {
  Card,
  Button,
  Icon,
  Table,
  Pagination,
  Loading,
  Empty,
  Link,
  Tag,
} from '@/components/common/base'
import { StatusTag, CopyButton, TimeAgo } from '@/components/common/business'
import { workflowService } from '@/api/services/workflow.service'
import { executionService } from '@/api/services/execution.service'
import type { Workflow } from '@/api/types/workflow.types'
import type { Execution } from '@/api/types/execution.types'
import type { TableColumn } from '@/components/common/base'
import { error } from '@/components/common/base/Message'

const route = useRoute()
const router = useRouter()

const workflowId = route.params.id as string

const loading = ref(false)
const running = ref(false)
const loadingExecutions = ref(false)

const workflow = ref<Workflow | null>(null)
const executionList = ref<Execution[]>([])
const executionPagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 表格列定义
const nodeColumns: TableColumn[] = [
  { prop: 'name', label: '节点名称', minWidth: 150 },
  { prop: 'type', label: '类型', width: 120, slot: 'type' },
  { prop: 'id', label: '节点ID', width: 200 },
]

const executionColumns: TableColumn[] = [
  { prop: 'startTime', label: '开始时间', width: 180 },
  { prop: 'status', label: '状态', width: 100, slot: 'status' },
  { prop: 'duration', label: '执行时长', width: 120, slot: 'duration' },
]

// 加载工作流详情
const loadWorkflow = async () => {
  loading.value = true
  try {
    const response = await workflowService.getWorkflow(workflowId)
    if (response.data) {
      workflow.value = response.data
    }
  } catch (err) {
    error('加载工作流失败')
    console.error('Load workflow error:', err)
  } finally {
    loading.value = false
  }
}

// 加载执行历史
const loadExecutions = async () => {
  loadingExecutions.value = true
  try {
    const response = await executionService.getExecutions({
      page: executionPagination.page,
      pageSize: executionPagination.pageSize,
      workflowId,
    })
    if (response.data) {
      executionList.value = response.data.list
      executionPagination.total = response.data.total
    }
  } catch (err) {
    console.error('Load executions error:', err)
  } finally {
    loadingExecutions.value = false
  }
}

// 格式化时间
const formatTime = (time?: string): string => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

// 格式化时长
const formatDuration = (seconds: number): string => {
  if (seconds < 1) {
    return `${(seconds * 1000).toFixed(0)}ms`
  }
  if (seconds < 60) {
    return `${seconds.toFixed(2)}s`
  }
  const minutes = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${minutes}m ${secs.toFixed(0)}s`
}

// 获取节点类型颜色
const getNodeTypeColor = (type: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  const colorMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    trigger: 'success',
    ai: 'primary',
    data: 'warning',
    control: 'danger',
    integration: 'info',
    tool: 'primary',
  }
  return colorMap[type] || 'info'
}

// 返回
const handleBack = () => {
  router.push('/workflow/list')
}

// 编辑
const handleEdit = () => {
  router.push(`/workflow/editor/${workflowId}`)
}

// 运行
const handleRun = async () => {
  running.value = true
  try {
    const response = await workflowService.executeWorkflow(workflowId)
    if (response.data) {
      router.push(`/workflow/execution/${response.data.executionId}`)
    }
  } catch (err) {
    error('执行工作流失败')
    console.error('Execute workflow error:', err)
  } finally {
    running.value = false
  }
}

// 执行点击
const handleExecutionClick = (row: unknown) => {
  const execution = row as Execution
  router.push(`/workflow/execution/${execution.id}`)
}

// 跳转到执行列表
const goToExecutions = () => {
  router.push(`/workflow/executions?workflowId=${workflowId}`)
}

// 执行列表分页
const handleExecutionSizeChange = (size: number) => {
  executionPagination.pageSize = size
  executionPagination.page = 1
  loadExecutions()
}

const handleExecutionPageChange = (page: number) => {
  executionPagination.page = page
  loadExecutions()
}

// 初始化
onMounted(() => {
  loadWorkflow()
  loadExecutions()
})
</script>

<style scoped lang="scss">
.workflow-detail-page {
  padding: 24px;
  background-color: var(--el-bg-color-page);
  min-height: 100%;

  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 24px;
    padding-bottom: 24px;
    border-bottom: 1px solid var(--el-border-color-light);

    .header-left {
      display: flex;
      align-items: flex-start;
      gap: 16px;
      flex: 1;

      .workflow-info {
        flex: 1;

        .workflow-name {
          margin: 0 0 12px 0;
          font-size: 24px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }

        .workflow-meta {
          display: flex;
          align-items: center;
          gap: 16px;
          font-size: 14px;
          color: var(--el-text-color-secondary);

          .meta-item {
            display: flex;
            align-items: center;
            gap: 4px;
          }
        }
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .loading-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 400px;
  }

  .detail-content {
    display: flex;
    flex-direction: column;
    gap: 24px;

    .info-card {
      .info-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 24px;

        .info-item {
          display: flex;
          flex-direction: column;
          gap: 8px;

          label {
            font-size: 14px;
            color: var(--el-text-color-secondary);
            font-weight: 500;
          }

          .info-value {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 14px;
            color: var(--el-text-color-primary);
          }
        }
      }
    }

    .table-pagination {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style>
