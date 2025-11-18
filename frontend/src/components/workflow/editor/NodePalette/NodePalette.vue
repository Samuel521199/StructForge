<template>
  <div class="node-palette">
    <div class="node-palette__header">
      <h3>节点面板</h3>
      <Input
        v-model="searchKeyword"
        placeholder="搜索节点"
        clearable
        size="small"
        class="node-palette__search"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </Input>
    </div>

    <div class="node-palette__content">
      <div
        v-for="category in filteredCategories"
        :key="category.name"
        class="node-palette__category"
      >
        <div class="node-palette__category-header" @click="toggleCategory(category.name)">
          <span>{{ category.label }}</span>
          <el-icon :class="{ 'is-expanded': expandedCategories.includes(category.name) }">
            <ArrowDown />
          </el-icon>
        </div>

        <div
          v-show="expandedCategories.includes(category.name)"
          class="node-palette__nodes"
        >
          <div
            v-for="node in category.nodes"
            :key="node.type"
            class="node-palette__node"
            draggable="true"
            @dragstart="handleDragStart($event, node)"
            @dragend="handleDragEnd"
          >
            <el-icon class="node-palette__node-icon">
              <Circle />
            </el-icon>
            <span class="node-palette__node-label">{{ node.label }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, ArrowDown, Circle } from '@element-plus/icons-vue'
import { Input } from '@/components/common/base'

interface NodeType {
  type: string
  label: string
  icon: any
  category: string
}

interface NodeCategory {
  name: string
  label: string
  nodes: NodeType[]
}

const searchKeyword = ref('')
const expandedCategories = ref<string[]>(['trigger', 'ai', 'data', 'control', 'integration', 'tool'])

// 节点类型定义
const nodeTypes: NodeType[] = [
  // 触发节点
  { type: 'webhook', label: 'Webhook', icon: 'Link', category: 'trigger' },
  { type: 'timer', label: '定时器', icon: 'Clock', category: 'trigger' },
  { type: 'manual', label: '手动触发', icon: 'Pointer', category: 'trigger' },
  
  // AI节点
  { type: 'text-generation', label: '文本生成', icon: 'Document', category: 'ai' },
  { type: 'image-generation', label: '图像生成', icon: 'Picture', category: 'ai' },
  { type: 'code-generation', label: '代码生成', icon: 'Code', category: 'ai' },
  
  // 数据处理节点
  { type: 'transform', label: '数据转换', icon: 'Refresh', category: 'data' },
  { type: 'filter', label: '数据过滤', icon: 'Filter', category: 'data' },
  { type: 'aggregate', label: '数据聚合', icon: 'DataAnalysis', category: 'data' },
  
  // 控制节点
  { type: 'condition', label: '条件判断', icon: 'Switch', category: 'control' },
  { type: 'loop', label: '循环', icon: 'RefreshRight', category: 'control' },
  { type: 'parallel', label: '并行', icon: 'Grid', category: 'control' },
  
  // 集成节点
  { type: 'http', label: 'HTTP请求', icon: 'Connection', category: 'integration' },
  { type: 'database', label: '数据库', icon: 'Database', category: 'integration' },
  { type: 'file', label: '文件操作', icon: 'Folder', category: 'integration' },
  
  // 工具节点
  { type: 'script', label: '脚本', icon: 'Document', category: 'tool' },
  { type: 'code-executor', label: '代码执行', icon: 'Cpu', category: 'tool' },
]

const categories: NodeCategory[] = [
  {
    name: 'trigger',
    label: '触发',
    nodes: nodeTypes.filter(n => n.category === 'trigger'),
  },
  {
    name: 'ai',
    label: 'AI',
    nodes: nodeTypes.filter(n => n.category === 'ai'),
  },
  {
    name: 'data',
    label: '数据处理',
    nodes: nodeTypes.filter(n => n.category === 'data'),
  },
  {
    name: 'control',
    label: '控制流',
    nodes: nodeTypes.filter(n => n.category === 'control'),
  },
  {
    name: 'integration',
    label: '集成',
    nodes: nodeTypes.filter(n => n.category === 'integration'),
  },
  {
    name: 'tool',
    label: '工具',
    nodes: nodeTypes.filter(n => n.category === 'tool'),
  },
]

const filteredCategories = computed(() => {
  if (!searchKeyword.value) {
    return categories
  }

  const keyword = searchKeyword.value.toLowerCase()
  return categories.map(category => ({
    ...category,
    nodes: category.nodes.filter(node =>
      node.label.toLowerCase().includes(keyword) ||
      node.type.toLowerCase().includes(keyword)
    ),
  })).filter(category => category.nodes.length > 0)
})

const toggleCategory = (categoryName: string) => {
  const index = expandedCategories.value.indexOf(categoryName)
  if (index > -1) {
    expandedCategories.value.splice(index, 1)
  } else {
    expandedCategories.value.push(categoryName)
  }
}

const handleDragStart = (event: DragEvent, node: NodeType) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData('application/json', JSON.stringify(node))
  }
}

const handleDragEnd = () => {
  // 拖拽结束处理
}
</script>

<style scoped lang="scss">
.node-palette {
  width: 250px;
  height: 100%;
  background: var(--el-bg-color);
  border-right: 1px solid var(--el-border-color-light);
  display: flex;
  flex-direction: column;

  &__header {
    padding: 16px;
    border-bottom: 1px solid var(--el-border-color-light);

    h3 {
      margin: 0 0 12px 0;
      font-size: 16px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }
  }

  &__search {
    width: 100%;
  }

  &__content {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }

  &__category {
    margin-bottom: 8px;

    &-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 8px 12px;
      font-size: 14px;
      font-weight: 500;
      color: var(--el-text-color-primary);
      background: var(--el-bg-color-page);
      border-radius: 4px;
      cursor: pointer;
      user-select: none;

      .el-icon {
        transition: transform 0.2s;

        &.is-expanded {
          transform: rotate(180deg);
        }
      }
    }
  }

  &__nodes {
    padding: 4px 0;
  }

  &__node {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    margin: 4px 0;
    background: var(--el-bg-color);
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 4px;
    cursor: grab;
    transition: all 0.2s;
    user-select: none;

    &:hover {
      background: var(--el-bg-color-page);
      border-color: var(--el-color-primary);
      transform: translateX(4px);
    }

    &:active {
      cursor: grabbing;
    }
  }

  &__node-icon {
    font-size: 16px;
    color: var(--el-color-primary);
  }

  &__node-label {
    font-size: 14px;
    color: var(--el-text-color-primary);
  }
}
</style>

