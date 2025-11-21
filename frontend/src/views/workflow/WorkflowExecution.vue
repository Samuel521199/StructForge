<template>
  <div class="workflow-execution-page">
    <!-- 页面头部 -->
    <div class="execution-header">
      <div class="header-left">
        <Button :icon="ArrowLeft" @click="handleBack">返回</Button>
        <div class="execution-info">
          <h1 class="execution-title">执行详情</h1>
          <div class="execution-meta">
            <StatusTag :status="execution?.status || 'pending'" />
            <span class="meta-item">
              <Icon :icon="Clock" :size="14" />
              开始时间: {{ formatTime(execution?.startTime) }}
            </span>
            <span v-if="execution?.endTime" class="meta-item">
              <Icon :icon="CircleCheck" :size="14" />
              结束时间: {{ formatTime(execution?.endTime) }}
            </span>
            <span v-if="execution?.duration" class="meta-item">
              <Icon :icon="Timer" :size="14" />
              执行时长: {{ formatDuration(execution.duration) }}
            </span>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <Button v-if="execution?.status === 'running'" type="danger" @click="handleCancel">
          取消执行
        </Button>
        <Button v-if="execution?.status === 'failed'" type="primary" @click="handleRetry" :loading="retrying">
          重试
        </Button>
        <Button @click="handleRefresh" :loading="refreshing">
          <template #icon>
            <Icon :icon="Refresh" />
          </template>
          刷新
        </Button>
      </div>
    </div>

    <!-- 执行内容 -->
    <div v-if="loading" class="loading-container">
      <Loading />
    </div>

    <div v-else-if="execution" class="execution-content">
      <!-- 基本信息 -->
      <Card class="info-card">
        <template #header>
          <span>基本信息</span>
        </template>
        <div class="info-grid">
          <div class="info-item">
            <label>执行ID</label>
            <div class="info-value">
              <span>{{ execution.id }}</span>
              <CopyButton :value="execution.id" />
            </div>
          </div>
          <div class="info-item">
            <label>工作流名称</label>
            <div class="info-value">
              <Link type="primary" @click="goToWorkflow">{{ execution.workflowName }}</Link>
            </div>
          </div>
          <div class="info-item">
            <label>状态</label>
            <div class="info-value">
              <StatusTag :status="execution.status" />
            </div>
          </div>
          <div class="info-item" v-if="execution.duration">
            <label>执行时长</label>
            <div class="info-value">{{ formatDuration(execution.duration) }}</div>
          </div>
          <div class="info-item" v-if="execution.error">
            <label>错误信息</label>
            <div class="info-value error-text">{{ execution.error }}</div>
          </div>
        </div>
      </Card>

      <!-- 节点执行状态 -->
      <Card v-if="execution.nodeExecutions && execution.nodeExecutions.length > 0" class="nodes-card">
        <template #header>
          <span>节点执行状态</span>
        </template>
        <Table
          :data="execution.nodeExecutions"
          :columns="nodeExecutionColumns"
        >
          <template #status="{ row }">
            <StatusTag :status="row.status" />
          </template>
          <template #duration="{ row }">
            <span v-if="row.duration">{{ formatDuration(row.duration) }}</span>
            <span v-else>-</span>
          </template>
        </Table>
      </Card>

      <!-- 执行日志 -->
      <Card class="logs-card">
        <template #header>
          <div class="card-header">
            <span>执行日志</span>
            <div class="log-filters">
              <Select
                v-model="logLevel"
                :options="logLevelOptions"
                placeholder="日志级别"
                style="width: 120px"
                clearable
                @change="loadLogs"
              />
            </div>
          </div>
        </template>
        <div class="logs-container">
          <div
            v-for="log in filteredLogs"
            :key="log.id"
            class="log-item"
            :class="`log-${log.level}`"
          >
            <span class="log-time">{{ formatTime(log.timestamp) }}</span>
            <span class="log-level">{{ log.level.toUpperCase() }}</span>
            <span v-if="log.nodeId" class="log-node">[{{ log.nodeId }}]</span>
            <span class="log-message">{{ log.message }}</span>
          </div>
          <Empty v-if="filteredLogs.length === 0" description="暂无日志" />
        </div>
      </Card>

      <!-- 执行结果 -->
      <Card v-if="execution.result" class="result-card">
        <template #header>
          <span>执行结果</span>
        </template>
        <div class="result-content">
          <pre class="result-json">{{ formatJSON(execution.result) }}</pre>
        </div>
      </Card>
    </div>

    <Empty v-else description="执行记录不存在" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Clock, Refresh, CircleCheck, Timer } from '@element-plus/icons-vue'
import {
  Card,
  Button,
  Icon,
  Table,
  Loading,
  Empty,
  Link,
  Select,
  SelectOption,
} from '@/components/common/base'
import { StatusTag, CopyButton } from '@/components/common/business'
import { executionService } from '@/api/services/execution.service'
import type { Execution, ExecutionLog } from '@/api/types/execution.types'
import type { TableColumn } from '@/components/common/base'
import { error, success } from '@/components/common/base/Message'

const route = useRoute()
const router = useRouter()

const executionId = route.params.id as string

const loading = ref(false)
const refreshing = ref(false)
const retrying = ref(false)
const logLevel = ref<string>('')

const execution = ref<Execution | null>(null)
const logs = ref<ExecutionLog[]>([])

// 过滤后的日志
const filteredLogs = computed(() => {
  if (!logLevel.value) {
    return logs.value
  }
  return logs.value.filter((log) => log.level === logLevel.value)
})

// 表格列定义
const nodeExecutionColumns: TableColumn[] = [
  { prop: 'nodeName', label: '节点名称', minWidth: 150 },
  { prop: 'nodeType', label: '节点类型', width: 120 },
  { prop: 'status', label: '状态', width: 100, slot: 'status' },
  { prop: 'duration', label: '执行时长', width: 120, slot: 'duration' },
]

let refreshTimer: number | null = null

// 加载执行详情
const loadExecution = async () => {
  loading.value = true
  try {
    const response = await executionService.getExecution(executionId)
    if (response.data) {
      execution.value = response.data
      
      // 如果还在运行，设置自动刷新
      if (response.data.status === 'running') {
        startAutoRefresh()
      } else {
        stopAutoRefresh()
      }
    }
  } catch (err) {
    error('加载执行详情失败')
    console.error('Load execution error:', err)
  } finally {
    loading.value = false
  }
}

// 加载日志
const loadLogs = async () => {
  try {
    const response = await executionService.getExecutionLogs(executionId)
    if (response.data) {
      logs.value = response.data
    }
  } catch (err) {
    console.error('Load logs error:', err)
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

// 格式化 JSON
const formatJSON = (obj: Record<string, unknown>): string => {
  return JSON.stringify(obj, null, 2)
}

// 自动刷新
const startAutoRefresh = () => {
  if (refreshTimer) return
  refreshTimer = window.setInterval(() => {
    if (execution.value?.status === 'running') {
      loadExecution()
      loadLogs()
    }
  }, 3000) // 每3秒刷新一次
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// 返回
const handleBack = () => {
  stopAutoRefresh()
  router.back()
}

// 刷新
const handleRefresh = async () => {
  refreshing.value = true
  try {
    await Promise.all([loadExecution(), loadLogs()])
  } finally {
    refreshing.value = false
  }
}

// 取消执行
const handleCancel = async () => {
  try {
    await executionService.cancelExecution(executionId)
    success('已取消执行')
    await loadExecution()
  } catch (err) {
    error('取消执行失败')
    console.error('Cancel execution error:', err)
  }
}

// 重试
const handleRetry = async () => {
  retrying.value = true
  try {
    const response = await executionService.retryExecution(executionId)
    if (response.data) {
      success('已重新执行')
      router.push(`/workflow/execution/${response.data.id}`)
    }
  } catch (err) {
    error('重试执行失败')
    console.error('Retry execution error:', err)
  } finally {
    retrying.value = false
  }
}

// 跳转到工作流
const goToWorkflow = () => {
  if (execution.value?.workflowId) {
    router.push(`/workflow/detail/${execution.value.workflowId}`)
  }
}

// 初始化
onMounted(() => {
  loadExecution()
  loadLogs()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped lang="scss">
.workflow-execution-page {
  padding: 24px;
  background-color: var(--el-bg-color-page);
  min-height: 100%;

  .execution-header {
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

      .execution-info {
        flex: 1;

        .execution-title {
          margin: 0 0 12px 0;
          font-size: 24px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }

        .execution-meta {
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

  .execution-content {
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

            &.error-text {
              color: var(--el-color-danger);
            }
          }
        }
      }
    }

    .logs-card {
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .log-filters {
        display: flex;
        gap: 12px;
      }

      .logs-container {
        max-height: 600px;
        overflow-y: auto;
        background-color: var(--el-bg-color);
        border: 1px solid var(--el-border-color-light);
        border-radius: 4px;
        padding: 16px;
        font-family: 'Courier New', monospace;
        font-size: 12px;
        line-height: 1.6;

        .log-item {
          display: flex;
          gap: 12px;
          padding: 4px 0;
          border-bottom: 1px solid var(--el-border-color-lighter);

          &:last-child {
            border-bottom: none;
          }

          .log-time {
            color: var(--el-text-color-secondary);
            min-width: 180px;
          }

          .log-level {
            font-weight: 600;
            min-width: 60px;

            &.log-debug {
              color: var(--el-text-color-secondary);
            }

            &.log-info {
              color: var(--el-color-info);
            }

            &.log-warn {
              color: var(--el-color-warning);
            }

            &.log-error {
              color: var(--el-color-danger);
            }
          }

          .log-node {
            color: var(--el-color-primary);
            min-width: 100px;
          }

          .log-message {
            flex: 1;
            color: var(--el-text-color-primary);
          }

          &.log-debug .log-level {
            color: var(--el-text-color-secondary);
          }

          &.log-info .log-level {
            color: var(--el-color-info);
          }

          &.log-warn .log-level {
            color: var(--el-color-warning);
          }

          &.log-error .log-level {
            color: var(--el-color-danger);
          }
        }
      }
    }

    .result-card {
      .result-content {
        .result-json {
          margin: 0;
          padding: 16px;
          background-color: var(--el-bg-color);
          border: 1px solid var(--el-border-color-light);
          border-radius: 4px;
          font-family: 'Courier New', monospace;
          font-size: 12px;
          line-height: 1.6;
          color: var(--el-text-color-primary);
          overflow-x: auto;
        }
      }
    }
  }
}
</style>
