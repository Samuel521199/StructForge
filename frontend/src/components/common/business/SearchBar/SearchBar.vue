<template>
  <div class="sf-search-bar">
    <div class="sf-search-bar__content">
      <!-- 搜索输入框 -->
      <Input
        v-if="searchable"
        :model-value="modelValue"
        :placeholder="placeholder"
        clearable
        class="sf-search-bar__input"
        @update:model-value="handleSearchUpdate"
        @keyup.enter="handleSearch"
        @clear="handleClear"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </Input>

      <!-- 筛选按钮 -->
      <Button
        v-if="filterable && filters && filters.length > 0"
        :type="hasActiveFilters ? 'primary' : 'default'"
        class="sf-search-bar__filter-btn"
        @click="showFilterPanel = !showFilterPanel"
      >
        <el-icon><Filter /></el-icon>
        <span v-if="hasActiveFilters" class="sf-search-bar__filter-count">
          {{ activeFilterCount }}
        </span>
        筛选
      </Button>

      <!-- 搜索按钮 -->
      <Button
        v-if="searchable"
        type="primary"
        class="sf-search-bar__search-btn"
        @click="handleSearch"
      >
        搜索
      </Button>
    </div>

    <!-- 筛选面板 -->
    <div
      v-if="filterable && filters && filters.length > 0 && showFilterPanel"
      class="sf-search-bar__filter-panel"
    >
      <div class="sf-search-bar__filter-content">
        <div
          v-for="filter in filters"
          :key="filter.prop"
          class="sf-search-bar__filter-item"
        >
          <label class="sf-search-bar__filter-label">{{ filter.label }}</label>
          <Select
            v-if="filter.type === 'select'"
            v-model="filterValues[filter.prop]"
            :options="filter.options || []"
            :placeholder="filter.placeholder || '请选择'"
            clearable
            style="width: 200px;"
            @change="handleFilterChange"
          />
          <Input
            v-else-if="filter.type === 'input'"
            v-model="filterValues[filter.prop]"
            :placeholder="filter.placeholder || '请输入'"
            clearable
            style="width: 200px;"
            @input="handleFilterChange"
          />
          <!-- TODO: 添加日期选择器支持 -->
        </div>
      </div>
      <div class="sf-search-bar__filter-actions">
        <Button @click="handleResetFilters">重置</Button>
        <Button type="primary" @click="handleApplyFilters">应用</Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Search, Filter } from '@element-plus/icons-vue'
import { Input, Button, Select } from '@/components/common/base'
import type { SearchBarProps, SearchBarEmits } from './types'

const props = withDefaults(defineProps<SearchBarProps>(), {
  placeholder: '请输入搜索关键词',
  searchable: true,
  filterable: false,
})

const emit = defineEmits<SearchBarEmits>()

// 筛选面板显示状态
const showFilterPanel = ref(false)

// 筛选值
const filterValues = ref<Record<string, any>>({})

// 是否有激活的筛选条件
const hasActiveFilters = computed(() => {
  return Object.values(filterValues.value).some(value => {
    if (Array.isArray(value)) {
      return value.length > 0
    }
    return value !== null && value !== undefined && value !== ''
  })
})

// 激活的筛选条件数量
const activeFilterCount = computed(() => {
  return Object.values(filterValues.value).filter(value => {
    if (Array.isArray(value)) {
      return value.length > 0
    }
    return value !== null && value !== undefined && value !== ''
  }).length
})

// 监听筛选值变化
watch(
  () => props.filters,
  (newFilters) => {
    if (newFilters) {
      const initialValues: Record<string, any> = {}
      newFilters.forEach(filter => {
        initialValues[filter.prop] = undefined
      })
      filterValues.value = initialValues
    }
  },
  { immediate: true }
)

const handleSearchUpdate = (value: string) => {
  emit('update:modelValue', value)
}

const handleSearch = () => {
  emit('search', props.modelValue)
}

const handleClear = () => {
  emit('update:modelValue', '')
  emit('search', '')
}

const handleFilterChange = () => {
  // 筛选值变化时自动触发（可选）
  // emit('filter', filterValues.value)
}

const handleApplyFilters = () => {
  emit('filter', { ...filterValues.value })
  showFilterPanel.value = false
}

const handleResetFilters = () => {
  if (props.filters) {
    const resetValues: Record<string, any> = {}
    props.filters.forEach(filter => {
      resetValues[filter.prop] = undefined
    })
    filterValues.value = resetValues
    emit('filter', {})
  }
  showFilterPanel.value = false
}
</script>

<style scoped lang="scss">
.sf-search-bar {
  &__content {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  &__input {
    flex: 1;
    max-width: 400px;
  }

  &__filter-btn {
    position: relative;
  }

  &__filter-count {
    display: inline-block;
    min-width: 18px;
    height: 18px;
    line-height: 18px;
    padding: 0 4px;
    margin-right: 4px;
    background: var(--el-color-danger);
    color: white;
    border-radius: 9px;
    font-size: 12px;
    text-align: center;
  }

  &__search-btn {
    // 搜索按钮样式
  }

  &__filter-panel {
    margin-top: 16px;
    padding: 16px;
    background: var(--el-bg-color-page);
    border: 1px solid var(--el-border-color-light);
    border-radius: 4px;
  }

  &__filter-content {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 16px;
  }

  &__filter-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  &__filter-label {
    min-width: 80px;
    font-size: 14px;
    color: var(--el-text-color-regular);
  }

  &__filter-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding-top: 16px;
    border-top: 1px solid var(--el-border-color-lighter);
  }
}
</style>
