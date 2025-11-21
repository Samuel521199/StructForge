<template>
  <div class="dashboard-page">
    <!-- 页面标题 -->
    <div class="dashboard-header">
      <h1 class="page-title">仪表盘</h1>
      <div class="header-actions">
        <Button type="primary" size="large" @click="handleCreateWorkflow" class="create-btn">
          <template #icon>
            <Icon :icon="Plus" />
          </template>
          创建新工作流
        </Button>
        <Button @click="handleRefresh" :loading="refreshing">
          <template #icon>
            <Icon :icon="Refresh" />
          </template>
          刷新
        </Button>
      </div>
    </div>

    <!-- 快速操作卡片区域 -->
    <div class="quick-actions-grid">
      <Card class="quick-action-card action-create" @click="handleCreateWorkflow">
        <div class="action-content">
          <div class="action-icon-wrapper icon-green">
            <Icon :icon="Plus" :size="32" />
          </div>
          <h3 class="action-title">创建新工作流</h3>
          <p class="action-desc">空白画布，从头开始设计</p>
        </div>
      </Card>
      <Card class="quick-action-card action-template" @click="handleCreateFromTemplate">
        <div class="action-content">
          <div class="action-icon-wrapper icon-blue">
            <Icon :icon="Files" :size="32" />
          </div>
          <h3 class="action-title">从模板创建</h3>
          <p class="action-desc">选择预设模板快速开始</p>
        </div>
      </Card>
      <Card class="quick-action-card action-import" @click="handleImportWorkflow">
        <div class="action-content">
          <div class="action-icon-wrapper icon-purple">
            <Icon :icon="UploadFilled" :size="32" />
          </div>
          <h3 class="action-title">导入工作流</h3>
          <p class="action-desc">从 JSON 文件导入工作流</p>
        </div>
      </Card>
      <Card class="quick-action-card action-list" @click="goToWorkflowList">
        <div class="action-content">
          <div class="action-icon-wrapper icon-orange">
            <Icon :icon="Menu" :size="32" />
          </div>
          <h3 class="action-title">查看所有工作流</h3>
          <p class="action-desc">管理所有工作流</p>
        </div>
      </Card>
    </div>

    <!-- 执行状态卡片 -->
    <div class="stats-grid">
      <Card class="stat-card">
        <Statistic
          title="正在运行"
          :value="stats.running"
          value-color="#00FF00"
          :icon="VideoPlay"
          description="当前正在执行的工作流"
        />
      </Card>
      <Card class="stat-card">
        <Statistic
          title="排队等待"
          :value="stats.queued"
          value-color="#e6a23c"
          :icon="Clock"
          description="等待执行的工作流"
        />
      </Card>
      <Card class="stat-card">
        <Statistic
          title="总执行次数"
          :value="stats.totalExecutions"
          suffix="次"
          value-color="#409eff"
          :icon="DataAnalysis"
          description="最近 7 天"
        />
      </Card>
      <Card class="stat-card">
        <Statistic
          title="失败率"
          :value="stats.failureRate"
          suffix="%"
          :precision="1"
          :value-color="stats.failureRate > 5 ? '#f56c6c' : '#67c23a'"
          :icon="Warning"
          description="最近 7 天"
        />
      </Card>
    </div>

    <!-- 最近工作流快捷入口 -->
    <Card class="recent-workflows-card" v-if="recentWorkflows.length > 0">
      <template #header>
        <div class="card-header">
          <span>最近工作流</span>
          <Link type="primary" @click="goToWorkflowList">查看全部</Link>
        </div>
      </template>
      <div class="recent-workflows-list">
        <div
          v-for="workflow in recentWorkflows"
          :key="workflow.id"
          class="recent-workflow-item"
          @click="goToWorkflowDetail(workflow.id)"
        >
          <div class="workflow-icon">
            <Icon :icon="Folder" :size="20" />
          </div>
          <div class="workflow-info">
            <div class="workflow-name">{{ workflow.name }}</div>
            <div class="workflow-time">最后编辑：{{ formatTimeAgo(workflow.updatedAt) }}</div>
          </div>
          <div class="workflow-status">
            <StatusTag :status="workflow.status" />
          </div>
        </div>
      </div>
    </Card>

    <!-- 主要内容区域 -->
    <div class="dashboard-content">
      <!-- 左侧：执行历史列表 -->
      <div class="content-left">
        <Card>
          <template #header>
            <div class="card-header">
              <span>最近执行记录</span>
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

      <!-- 右侧：执行成功率统计 -->
      <div class="content-right">
        <Card>
          <template #header>
            <div class="card-header">
              <span>执行成功率（最近 24 小时）</span>
            </div>
          </template>
          <Table
            :data="successRateList"
            :columns="successRateColumns"
            :loading="loadingSuccessRate"
            stripe
          >
            <template #successRate="{ row }">
              <Progress
                :percentage="row.successRate"
                :status="row.successRate >= 99 ? 'success' : row.successRate >= 95 ? 'warning' : 'exception'"
                :text="`${row.successRate.toFixed(1)}%`"
              />
            </template>
            <template #failed="{ row }">
              <span :style="{ color: row.failed > 0 ? '#f56c6c' : '#67c23a' }">
                {{ row.failed }}
              </span>
            </template>
          </Table>
        </Card>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <!-- 错误趋势图 -->
      <Card>
        <template #header>
          <div class="card-header">
            <span>工作流错误趋势（最近 24 小时）</span>
          </div>
        </template>
        <Chart
          v-if="errorTrendOption"
          :option="errorTrendOption as any"
          height="300px"
          theme="dark"
        />
        <Empty v-else description="暂无数据" />
      </Card>

      <!-- 平均执行时长图 -->
      <Card>
        <template #header>
          <div class="card-header">
            <span>平均执行时长</span>
          </div>
        </template>
        <Chart
          v-if="durationOption"
          :option="durationOption as any"
          height="300px"
          theme="dark"
        />
        <Empty v-else description="暂无数据" />
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Refresh,
  VideoPlay,
  Clock,
  DataAnalysis,
  Warning,
  Plus,
  Files,
  UploadFilled,
  Menu,
  Folder,
} from '@element-plus/icons-vue'
import {
  Card,
  Button,
  Icon,
  Statistic,
  Table,
  Pagination,
  Progress,
  Chart,
  Empty,
  Link,
} from '@/components/common/base'
import { StatusTag } from '@/components/common/business'
import { dashboardService } from '@/api/services/dashboard.service'
import type {
  DashboardStats,
  ExecutionListItem,
  SuccessRateItem,
} from '@/api/types/dashboard.types'
import type { TableColumn } from '@/components/common/base'
import type { EChartsOption } from 'echarts'

const router = useRouter()

// 状态
const refreshing = ref(false)
const loadingExecutions = ref(false)
const loadingSuccessRate = ref(false)

const stats = reactive<DashboardStats>({
  running: 0,
  queued: 0,
  totalExecutions: 0,
  failedExecutions: 0,
  failureRate: 0,
  avgExecutionTime: 0,
  totalWorkflows: 0,
  activeWorkflows: 0,
})

const executionList = ref<ExecutionListItem[]>([])
const executionPagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const successRateList = ref<SuccessRateItem[]>([])

// 最近工作流
interface RecentWorkflow {
  id: string
  name: string
  status: 'active' | 'inactive' | 'error'
  updatedAt: string
}
const recentWorkflows = ref<RecentWorkflow[]>([])

// 表格列定义
const executionColumns: TableColumn[] = [
  { prop: 'workflowName', label: '工作流名称', width: 200 },
  { prop: 'startTime', label: '开始时间', width: 180 },
  { prop: 'status', label: '状态', width: 100, slot: 'status' },
  { prop: 'duration', label: '执行时长', width: 120, slot: 'duration' },
]

const successRateColumns: TableColumn[] = [
  { prop: 'workflowName', label: '工作流名称', minWidth: 200 },
  { prop: 'totalExecutions', label: '总执行次数', width: 120 },
  { prop: 'successful', label: '成功', width: 100 },
  { prop: 'failed', label: '失败', width: 100, slot: 'failed' },
  { prop: 'successRate', label: '成功率', width: 200, slot: 'successRate' },
]

// 图表配置
const errorTrendOption = ref<EChartsOption | null>(null)
const durationOption = ref<EChartsOption | null>(null)

// 类型断言辅助函数（用于解决 ECharts 类型兼容性问题）
// 类型断言辅助函数（用于解决 ECharts 类型兼容性问题）
const asEChartsOption = (option: Record<string, unknown>): EChartsOption => option as EChartsOption

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await dashboardService.getStats()
    if (response.data) {
      Object.assign(stats, response.data)
    }
  } catch (error) {
    console.error('Load stats error:', error)
  }
}

// 加载执行历史
const loadExecutions = async () => {
  loadingExecutions.value = true
  try {
    const response = await dashboardService.getExecutions({
      page: executionPagination.page,
      pageSize: executionPagination.pageSize,
    })
    if (response.data) {
      // 确保 list 是数组
      executionList.value = Array.isArray(response.data.list) ? response.data.list : []
      executionPagination.total = response.data.total || 0
    } else {
      executionList.value = []
      executionPagination.total = 0
    }
  } catch (error) {
    console.error('Load executions error:', error)
    // 出错时确保是空数组
    executionList.value = []
    executionPagination.total = 0
  } finally {
    loadingExecutions.value = false
  }
}

// 加载成功率
const loadSuccessRate = async () => {
  loadingSuccessRate.value = true
  try {
    const response = await dashboardService.getSuccessRate({ period: '24h' })
    // 确保返回的数据是数组
    if (response.data && Array.isArray(response.data)) {
      successRateList.value = response.data
    } else {
      successRateList.value = []
    }
  } catch (error) {
    console.error('Load success rate error:', error)
    // 出错时确保是空数组
    successRateList.value = []
  } finally {
    loadingSuccessRate.value = false
  }
}

// 加载错误趋势
const loadErrorTrend = async () => {
  try {
    const response = await dashboardService.getErrorTrend({ period: '24h' })
    if (response.data && response.data.length > 0) {
      errorTrendOption.value = asEChartsOption({
        backgroundColor: 'transparent',
        textStyle: {
          color: '#00FF00',
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true,
        },
        xAxis: {
          type: 'category',
          data: response.data.map((item) => item.time),
          axisLine: {
            lineStyle: {
              color: '#00FF00',
            },
          },
        },
        yAxis: {
          type: 'value',
          axisLine: {
            lineStyle: {
              color: '#00FF00',
            },
          },
        },
        series: [
          {
            name: '错误数量',
            type: 'line',
            data: response.data.map((item) => item.count),
            smooth: true,
            lineStyle: {
              color: '#00FF00',
            },
            itemStyle: {
              color: '#00FF00',
            },
            areaStyle: {
              color: {
                type: 'linear',
                x: 0,
                y: 0,
                x2: 0,
                y2: 1,
                colorStops: [
                  { offset: 0, color: 'rgba(0, 255, 0, 0.3)' },
                  { offset: 1, color: 'rgba(0, 255, 0, 0.05)' },
                ],
              },
            },
          },
        ],
      })
    }
  } catch (error) {
    console.error('Load error trend error:', error)
  }
}

// 加载执行时长
const loadExecutionDuration = async () => {
  try {
    const response = await dashboardService.getExecutionDuration()
    if (response.data && response.data.length > 0) {
      const data = response.data.slice(0, 10) // 只显示前10个
      durationOption.value = asEChartsOption({
        backgroundColor: 'transparent',
        textStyle: {
          color: '#00FF00',
        },
        grid: {
          left: '20%',
          right: '10%',
          bottom: '10%',
          containLabel: true,
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow',
          },
        },
        legend: {
          data: ['平均时长', '最小时长', '最大时长'],
          textStyle: {
            color: '#00FF00',
          },
        },
        xAxis: {
          type: 'value',
          axisLine: {
            lineStyle: {
              color: '#00FF00',
            },
          },
        },
        yAxis: {
          type: 'category',
          data: data.map((item) => item.workflowName),
          axisLine: {
            lineStyle: {
              color: '#00FF00',
            },
          },
        },
        series: [
          {
            name: '平均时长',
            type: 'bar',
            data: data.map((item) => item.avgDuration / 1000), // 转换为秒
            itemStyle: {
              color: '#00FF00',
            },
          },
          {
            name: '最小时长',
            type: 'bar',
            data: data.map((item) => item.minDuration / 1000),
            itemStyle: {
              color: '#67c23a',
            },
          },
          {
            name: '最大时长',
            type: 'bar',
            data: data.map((item) => item.maxDuration / 1000),
            itemStyle: {
              color: '#409eff',
            },
          },
        ],
      })
    }
  } catch (error) {
    console.error('Load execution duration error:', error)
  }
}

// 刷新所有数据
const handleRefresh = async () => {
  refreshing.value = true
  try {
    await Promise.all([
      loadStats(),
      loadExecutions(),
      loadSuccessRate(),
      loadErrorTrend(),
      loadExecutionDuration(),
    ])
  } finally {
    refreshing.value = false
  }
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

// 执行点击
const handleExecutionClick = (row: unknown) => {
  const execution = row as ExecutionListItem
  router.push(`/workflow/execution/${execution.id}`)
}

// 跳转到执行列表
const goToExecutions = () => {
  router.push('/workflow/executions')
}

// 创建工作流
const handleCreateWorkflow = () => {
  router.push('/workflow/editor')
}

// 从模板创建
const handleCreateFromTemplate = () => {
  // TODO: 实现模板选择功能
  router.push('/workflow/templates')
}

// 导入工作流
const handleImportWorkflow = () => {
  // TODO: 实现导入功能
  console.log('Import workflow')
}

// 跳转到工作流列表
const goToWorkflowList = () => {
  router.push('/workflow/list')
}

// 跳转到工作流详情
const goToWorkflowDetail = (id: string) => {
  router.push(`/workflow/detail/${id}`)
}

// 格式化相对时间
const formatTimeAgo = (time: string): string => {
  const now = new Date()
  const past = new Date(time)
  const diff = now.getTime() - past.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 0) return `${days}天前`
  if (hours > 0) return `${hours}小时前`
  if (minutes > 0) return `${minutes}分钟前`
  return '刚刚'
}

// 加载最近工作流
const loadRecentWorkflows = async () => {
  try {
    // TODO: 调用 API 获取最近工作流
    // 暂时使用模拟数据
    recentWorkflows.value = []
  } catch (error) {
    console.error('Load recent workflows error:', error)
    recentWorkflows.value = []
  }
}

// 初始化
onMounted(() => {
  handleRefresh()
  loadRecentWorkflows()
})
</script>

<style scoped lang="scss">
@use '@/assets/styles/glassmorphism' as *;

.dashboard-page {
  padding: 32px;
  min-height: 100%;
  position: relative;

  // 赛博朋克网格背景和光效
  &::before {
    content: '';
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: 
      // 细网格（20px）
      repeating-linear-gradient(
        0deg,
        transparent 0px,
        transparent 19px,
        rgba(0, 212, 255, 0.05) 20px,
        rgba(0, 212, 255, 0.05) 21px
      ),
      repeating-linear-gradient(
        90deg,
        transparent 0px,
        transparent 19px,
        rgba(0, 212, 255, 0.05) 20px,
        rgba(0, 212, 255, 0.05) 21px
      ),
      // 粗网格（100px）
      repeating-linear-gradient(
        0deg,
        transparent 0px,
        transparent 99px,
        rgba(0, 212, 255, 0.08) 100px,
        rgba(0, 212, 255, 0.08) 101px
      ),
      repeating-linear-gradient(
        90deg,
        transparent 0px,
        transparent 99px,
        rgba(0, 212, 255, 0.08) 100px,
        rgba(0, 212, 255, 0.08) 101px
      );
    background: 
      // 动态光效
      radial-gradient(circle at 10% 20%, rgba(0, 212, 255, 0.12) 0%, transparent 40%),
      radial-gradient(circle at 90% 80%, rgba(0, 255, 136, 0.1) 0%, transparent 40%),
      radial-gradient(circle at 50% 50%, rgba(183, 148, 246, 0.08) 0%, transparent 50%);
    pointer-events: none;
    z-index: 0;
    animation: lightMove 18s ease-in-out infinite;
  }

  // 扫描线效果（多层）
  &::after {
    content: '';
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, 
      transparent 0%,
      rgba(0, 212, 255, 0.9) 20%,
      rgba(0, 255, 136, 0.9) 50%,
      rgba(183, 148, 246, 0.9) 80%,
      transparent 100%
    );
    box-shadow: 
      0 0 20px rgba(0, 212, 255, 0.8),
      0 0 40px rgba(0, 212, 255, 0.4);
    pointer-events: none;
    z-index: 1;
    animation: scanLine 4s linear infinite;
  }

  > * {
    position: relative;
    z-index: 1;
  }

  .dashboard-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 32px;
    padding: 20px 0;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    position: relative;

    &::before {
      content: '';
      position: absolute;
      bottom: -1px;
      left: 0;
      width: 200px;
      height: 1px;
      background: linear-gradient(90deg, #409eff 0%, transparent 100%);
    }

    .page-title {
      margin: 0;
      font-size: 32px;
      font-weight: 800;
      color: #ffffff;
      text-shadow: 
        0 0 20px rgba(255, 255, 255, 0.5),
        0 2px 10px rgba(0, 0, 0, 0.8),
        0 0 40px rgba(64, 158, 255, 0.4);
      letter-spacing: 2px;
      position: relative;
      display: inline-block;

      &::after {
        content: '';
        position: absolute;
        bottom: -4px;
        left: 0;
        width: 100%;
        height: 3px;
        background: linear-gradient(90deg, #409eff 0%, #67c23a 50%, #9333ea 100%);
        border-radius: 2px;
        box-shadow: 0 0 10px rgba(64, 158, 255, 0.6);
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;
      align-items: center;

      .create-btn {
        background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
        border: none;
        color: #ffffff;
        box-shadow: 0 4px 20px rgba(64, 158, 255, 0.4);
        transition: all 0.3s ease;
        position: relative;
        overflow: hidden;

        &::before {
          content: '';
          position: absolute;
          top: 0;
          left: -100%;
          width: 100%;
          height: 100%;
          background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
          transition: left 0.5s ease;
        }

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 6px 30px rgba(64, 158, 255, 0.6);

          &::before {
            left: 100%;
          }
        }
      }
    }
  }

  // 快速操作卡片区域
  .quick-actions-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 24px;
    margin-bottom: 40px;
    position: relative;

      .quick-action-card {
      cursor: pointer;
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      border: 2px solid transparent;
      background: rgba(20, 20, 30, 0.9);
      backdrop-filter: blur(15px);
      overflow: hidden;
      position: relative;
      border-radius: 16px;

      &::before {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        opacity: 0;
        transition: opacity 0.3s ease;
        z-index: 0;
      }

      &:hover {
        transform: translateY(-8px) scale(1.02);
        border-color: currentColor;

        &::before {
          opacity: 1;
        }

        .action-icon-wrapper {
          transform: scale(1.15) rotate(8deg);

          &::after {
            opacity: 1;
          }
        }
      }

      .action-content {
        position: relative;
        z-index: 2;
        padding: 24px 16px;
        text-align: center;
        min-height: 160px;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
      }

      .action-icon-wrapper {
        width: 72px;
        height: 72px;
        border-radius: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 20px;
        transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
        box-shadow: 0 8px 25px rgba(0, 0, 0, 0.4);
        position: relative;

        &::after {
          content: '';
          position: absolute;
          inset: -2px;
          border-radius: 20px;
          padding: 2px;
          background: linear-gradient(135deg, rgba(255, 255, 255, 0.3), transparent);
          -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          -webkit-mask-composite: xor;
          mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          mask-composite: exclude;
          opacity: 0;
          transition: opacity 0.3s ease;
        }
      }

      .action-title {
        margin: 0 0 8px 0;
        font-size: 16px;
        font-weight: 600;
        color: #ffffff;
        text-shadow: 0 0 10px rgba(255, 255, 255, 0.5), 0 2px 4px rgba(0, 0, 0, 0.5);
      }

      .action-desc {
        margin: 0;
        font-size: 12px;
        color: rgba(255, 255, 255, 0.85);
        text-shadow: 0 0 8px rgba(255, 255, 255, 0.3), 0 1px 2px rgba(0, 0, 0, 0.5);
      }

      // 不同颜色的卡片
      &.action-create {
        border: 2px solid rgba(0, 255, 136, 0.5);
        box-shadow: 
          0 8px 32px rgba(0, 0, 0, 0.5),
          0 0 20px rgba(0, 255, 136, 0.4),
          0 0 40px rgba(0, 255, 136, 0.2),
          0 0 60px rgba(0, 255, 136, 0.1),
          0 0 30px rgba(0, 255, 136, 0.2) inset;
        position: relative;
        
        // 霓虹边框动画
        &::before {
          content: '';
          position: absolute;
          inset: -2px;
          border-radius: 16px;
          padding: 2px;
          background: linear-gradient(135deg, 
            rgba(0, 255, 136, 0.8) 0%,
            rgba(0, 255, 136, 0.4) 50%,
            transparent 100%
          );
          -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          -webkit-mask-composite: xor;
          mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          mask-composite: exclude;
          z-index: -1;
          animation: borderGlow 3s ease-in-out infinite;
        }
        
        // 内部光效
        &::after {
          content: '';
          position: absolute;
          inset: 0;
          background: linear-gradient(135deg, rgba(0, 255, 136, 0.1) 0%, rgba(103, 194, 58, 0.1) 100%);
          border-radius: 16px;
          z-index: 0;
        }

        .action-icon-wrapper.icon-green {
          background: linear-gradient(135deg, #00ff88 0%, #67c23a 100%);
          color: #fff;
          box-shadow: 
            0 0 20px rgba(0, 255, 136, 0.6),
            0 0 40px rgba(0, 255, 136, 0.4),
            inset 0 0 20px rgba(0, 255, 136, 0.3);
        }

        &:hover {
          border-color: rgba(0, 255, 136, 0.8);
          box-shadow: 
            0 12px 50px rgba(0, 0, 0, 0.6),
            0 0 30px rgba(0, 255, 136, 0.6),
            0 0 60px rgba(0, 255, 136, 0.4),
            0 0 90px rgba(0, 255, 136, 0.2),
            0 0 50px rgba(0, 255, 136, 0.3) inset;
          transform: translateY(-8px);
        }
      }

      &.action-template {
        border: 2px solid rgba(0, 212, 255, 0.5);
        box-shadow: 
          0 8px 32px rgba(0, 0, 0, 0.5),
          0 0 20px rgba(0, 212, 255, 0.4),
          0 0 40px rgba(0, 212, 255, 0.2),
          0 0 60px rgba(0, 212, 255, 0.1),
          0 0 30px rgba(0, 212, 255, 0.2) inset;
        position: relative;
        
        &::before {
          content: '';
          position: absolute;
          inset: -2px;
          border-radius: 16px;
          padding: 2px;
          background: linear-gradient(135deg, 
            rgba(0, 212, 255, 0.8) 0%,
            rgba(0, 212, 255, 0.4) 50%,
            transparent 100%
          );
          -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          -webkit-mask-composite: xor;
          mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          mask-composite: exclude;
          z-index: -1;
          animation: borderGlow 3s ease-in-out infinite;
        }
        
        &::after {
          content: '';
          position: absolute;
          inset: 0;
          background: linear-gradient(135deg, rgba(0, 212, 255, 0.1) 0%, rgba(64, 158, 255, 0.1) 100%);
          border-radius: 16px;
          z-index: 0;
        }

        .action-icon-wrapper.icon-blue {
          background: linear-gradient(135deg, #00d4ff 0%, #409eff 100%);
          color: #fff;
          box-shadow: 
            0 0 20px rgba(0, 212, 255, 0.6),
            0 0 40px rgba(0, 212, 255, 0.4),
            inset 0 0 20px rgba(0, 212, 255, 0.3);
        }

        &:hover {
          border-color: rgba(0, 212, 255, 0.8);
          box-shadow: 
            0 12px 50px rgba(0, 0, 0, 0.6),
            0 0 30px rgba(0, 212, 255, 0.6),
            0 0 60px rgba(0, 212, 255, 0.4),
            0 0 90px rgba(0, 212, 255, 0.2),
            0 0 50px rgba(0, 212, 255, 0.3) inset;
          transform: translateY(-8px);
        }
      }

      &.action-import {
        border: 2px solid rgba(183, 148, 246, 0.5);
        box-shadow: 
          0 8px 32px rgba(0, 0, 0, 0.5),
          0 0 20px rgba(183, 148, 246, 0.4),
          0 0 40px rgba(183, 148, 246, 0.2),
          0 0 60px rgba(183, 148, 246, 0.1),
          0 0 30px rgba(183, 148, 246, 0.2) inset;
        position: relative;
        
        &::before {
          content: '';
          position: absolute;
          inset: -2px;
          border-radius: 16px;
          padding: 2px;
          background: linear-gradient(135deg, 
            rgba(183, 148, 246, 0.8) 0%,
            rgba(183, 148, 246, 0.4) 50%,
            transparent 100%
          );
          -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          -webkit-mask-composite: xor;
          mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          mask-composite: exclude;
          z-index: -1;
          animation: borderGlow 3s ease-in-out infinite;
        }
        
        &::after {
          content: '';
          position: absolute;
          inset: 0;
          background: linear-gradient(135deg, rgba(183, 148, 246, 0.1) 0%, rgba(219, 39, 119, 0.1) 100%);
          border-radius: 16px;
          z-index: 0;
        }

        .action-icon-wrapper.icon-purple {
          background: linear-gradient(135deg, #b794f6 0%, #db2777 100%);
          color: #fff;
          box-shadow: 
            0 0 20px rgba(183, 148, 246, 0.6),
            0 0 40px rgba(183, 148, 246, 0.4),
            inset 0 0 20px rgba(183, 148, 246, 0.3);
        }

        &:hover {
          border-color: rgba(183, 148, 246, 0.8);
          box-shadow: 
            0 12px 50px rgba(0, 0, 0, 0.6),
            0 0 30px rgba(183, 148, 246, 0.6),
            0 0 60px rgba(183, 148, 246, 0.4),
            0 0 90px rgba(183, 148, 246, 0.2),
            0 0 50px rgba(183, 148, 246, 0.3) inset;
          transform: translateY(-8px);
        }
      }

      &.action-list {
        border: 2px solid rgba(255, 184, 77, 0.5);
        box-shadow: 
          0 8px 32px rgba(0, 0, 0, 0.5),
          0 0 20px rgba(255, 184, 77, 0.4),
          0 0 40px rgba(255, 184, 77, 0.2),
          0 0 60px rgba(255, 184, 77, 0.1),
          0 0 30px rgba(255, 184, 77, 0.2) inset;
        position: relative;
        
        &::before {
          content: '';
          position: absolute;
          inset: -2px;
          border-radius: 16px;
          padding: 2px;
          background: linear-gradient(135deg, 
            rgba(255, 184, 77, 0.8) 0%,
            rgba(255, 184, 77, 0.4) 50%,
            transparent 100%
          );
          -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          -webkit-mask-composite: xor;
          mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
          mask-composite: exclude;
          z-index: -1;
          animation: borderGlow 3s ease-in-out infinite;
        }
        
        &::after {
          content: '';
          position: absolute;
          inset: 0;
          background: linear-gradient(135deg, rgba(255, 184, 77, 0.1) 0%, rgba(239, 68, 68, 0.1) 100%);
          border-radius: 16px;
          z-index: 0;
        }

        .action-icon-wrapper.icon-orange {
          background: linear-gradient(135deg, #ffb84d 0%, #ef4444 100%);
          color: #fff;
          box-shadow: 
            0 0 20px rgba(255, 184, 77, 0.6),
            0 0 40px rgba(255, 184, 77, 0.4),
            inset 0 0 20px rgba(255, 184, 77, 0.3);
        }

        &:hover {
          border-color: rgba(255, 184, 77, 0.8);
          box-shadow: 
            0 12px 50px rgba(0, 0, 0, 0.6),
            0 0 30px rgba(255, 184, 77, 0.6),
            0 0 60px rgba(255, 184, 77, 0.4),
            0 0 90px rgba(255, 184, 77, 0.2),
            0 0 50px rgba(255, 184, 77, 0.3) inset;
          transform: translateY(-8px);
        }
      }
    }
  }

  // 最近工作流卡片
  .recent-workflows-card {
    margin-bottom: 32px;
    background: rgba(20, 20, 30, 0.9);
    border: 1px solid rgba(64, 158, 255, 0.3);
    box-shadow: 
      0 0 30px rgba(64, 158, 255, 0.2),
      inset 0 0 10px rgba(64, 158, 255, 0.1);
    backdrop-filter: blur(15px);
    border-radius: 16px;
    transition: all 0.3s ease;

    &:hover {
      border-color: rgba(64, 158, 255, 0.5);
      box-shadow: 
        0 0 40px rgba(64, 158, 255, 0.3),
        inset 0 0 15px rgba(64, 158, 255, 0.15);
    }

    .recent-workflows-list {
      display: flex;
      gap: 16px;
      flex-wrap: wrap;

        .recent-workflow-item {
        flex: 1;
        min-width: 200px;
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 16px;
        border-radius: 12px;
        background: rgba(64, 158, 255, 0.05);
        border: 1px solid rgba(64, 158, 255, 0.2);
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

        &:hover {
          background: rgba(64, 158, 255, 0.1);
          border-color: rgba(64, 158, 255, 0.4);
          transform: translateX(8px) scale(1.02);
          box-shadow: 0 6px 20px rgba(64, 158, 255, 0.3);
        }

        .workflow-icon {
          width: 44px;
          height: 44px;
          border-radius: 12px;
          background: linear-gradient(135deg, rgba(64, 158, 255, 0.2) 0%, rgba(103, 194, 58, 0.2) 100%);
          display: flex;
          align-items: center;
          justify-content: center;
          color: #409eff;
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
        }

        .workflow-info {
          flex: 1;
          min-width: 0;

          .workflow-name {
            font-size: 14px;
            font-weight: 600;
            color: var(--el-text-color-primary);
            margin-bottom: 4px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .workflow-time {
            font-size: 12px;
            color: var(--el-text-color-secondary);
          }
        }

        .workflow-status {
          flex-shrink: 0;
        }
      }
    }
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 24px;
    margin-bottom: 32px;

    .stat-card {
      background: rgba(20, 20, 30, 0.9);
      border: 1px solid rgba(64, 158, 255, 0.3);
      box-shadow: 
        0 0 30px rgba(64, 158, 255, 0.2),
        inset 0 0 10px rgba(64, 158, 255, 0.1);
      backdrop-filter: blur(15px);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      position: relative;
      overflow: hidden;
      border-radius: 16px;

      &::before {
        content: '';
        position: absolute;
        top: -50%;
        left: -50%;
        width: 200%;
        height: 200%;
        background: radial-gradient(circle, rgba(64, 158, 255, 0.15) 0%, transparent 70%);
        opacity: 0;
        transition: opacity 0.4s ease;
      }

      &:hover {
        transform: translateY(-8px) scale(1.02);
        border-color: rgba(64, 158, 255, 0.6);
        box-shadow: 
          0 12px 60px rgba(64, 158, 255, 0.4),
          inset 0 0 20px rgba(64, 158, 255, 0.2);

        &::before {
          opacity: 1;
        }
      }
    }
  }

  .dashboard-content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
    margin-bottom: 24px;

    @media (max-width: 1200px) {
      grid-template-columns: 1fr;
    }

    .content-left,
    .content-right {
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }

      .table-pagination {
        margin-top: 16px;
        display: flex;
        justify-content: flex-end;
      }
    }
  }

  .charts-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;

    @media (max-width: 1200px) {
      grid-template-columns: 1fr;
    }

    :deep(.el-card) {
      background: rgba(20, 20, 30, 0.9);
      border: 1px solid rgba(64, 158, 255, 0.3);
      box-shadow: 
        0 0 30px rgba(64, 158, 255, 0.2),
        inset 0 0 10px rgba(64, 158, 255, 0.1);
      backdrop-filter: blur(15px);
      transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
      border-radius: 16px;

      &:hover {
        border-color: rgba(64, 158, 255, 0.5);
        box-shadow: 
          0 0 40px rgba(64, 158, 255, 0.3),
          inset 0 0 15px rgba(64, 158, 255, 0.15);
      }
    }
  }

  // 增强卡片样式（赛博朋克风格）- 覆盖通用卡片样式
  :deep(.el-card) {
    background: rgba(20, 20, 30, 0.85);
    border: 2px solid rgba(0, 212, 255, 0.4);
    box-shadow: 
      0 8px 32px rgba(0, 0, 0, 0.5),
      0 0 20px rgba(0, 212, 255, 0.3),
      0 0 40px rgba(0, 212, 255, 0.15),
      0 0 60px rgba(0, 212, 255, 0.08),
      inset 0 0 20px rgba(0, 212, 255, 0.1);
    backdrop-filter: blur(20px) saturate(180%);
    -webkit-backdrop-filter: blur(20px) saturate(180%);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 16px;
    position: relative;
    overflow: visible;
    
    // 左上角装饰
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      width: 50px;
      height: 50px;
      border-top: 3px solid rgba(0, 212, 255, 0.8);
      border-left: 3px solid rgba(0, 212, 255, 0.8);
      border-radius: 16px 0 0 0;
      box-shadow: 
        -3px -3px 15px rgba(0, 212, 255, 0.6),
        -3px -3px 30px rgba(0, 212, 255, 0.4);
      pointer-events: none;
      z-index: 1;
      animation: cornerGlow 2s ease-in-out infinite;
    }

    &:hover {
      border-color: rgba(0, 212, 255, 0.7);
      box-shadow: 
        0 12px 50px rgba(0, 0, 0, 0.6),
        0 0 30px rgba(0, 212, 255, 0.5),
        0 0 60px rgba(0, 212, 255, 0.3),
        0 0 90px rgba(0, 212, 255, 0.15),
        inset 0 0 30px rgba(0, 212, 255, 0.2);
      transform: translateY(-4px);
      
      &::before {
        box-shadow: 
          -3px -3px 25px rgba(0, 212, 255, 0.9),
          -3px -3px 50px rgba(0, 212, 255, 0.6);
      }
    }

    .el-card__header {
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.08) 0%, rgba(0, 255, 136, 0.05) 100%);
      border-bottom: 1px solid rgba(0, 212, 255, 0.3);
      color: rgba(0, 212, 255, 0.95);
      font-weight: 600;
      padding: 16px 20px;
      position: relative;
      z-index: 1;
      text-shadow: 0 0 10px rgba(0, 212, 255, 0.3);
    }

    .el-card__body {
      color: rgba(255, 255, 255, 0.9);
      padding: 20px;
      position: relative;
      z-index: 1;
    }
  }
}

// 扫描线动画
@keyframes scanLine {
  0% {
    top: 0;
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    top: 100%;
    opacity: 0;
  }
}

// 边框光晕动画
@keyframes borderGlow {
  0%, 100% {
    opacity: 0.3;
    transform: rotate(0deg);
  }
  50% {
    opacity: 0.8;
    transform: rotate(180deg);
  }
}

// 装饰角光晕动画
@keyframes cornerGlow {
  0%, 100% {
    opacity: 0.6;
    box-shadow: 
      -3px -3px 15px rgba(0, 212, 255, 0.6),
      -3px -3px 30px rgba(0, 212, 255, 0.4);
  }
  50% {
    opacity: 1;
    box-shadow: 
      -3px -3px 25px rgba(0, 212, 255, 0.9),
      -3px -3px 50px rgba(0, 212, 255, 0.6);
  }
}
</style>
