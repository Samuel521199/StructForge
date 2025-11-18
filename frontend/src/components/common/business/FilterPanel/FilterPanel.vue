<template>
  <div class="sf-filter-panel" :class="{ 'is-collapsed': collapsed }">
    <div v-if="collapsible" class="sf-filter-panel__header" @click="toggleCollapse">
      <span class="sf-filter-panel__title">筛选条件</span>
      <el-icon class="sf-filter-panel__icon" :class="{ 'is-rotated': !collapsed }">
        <ArrowDown />
      </el-icon>
    </div>

    <div v-show="!collapsed" class="sf-filter-panel__content">
      <div class="sf-filter-panel__filters">
        <div
          v-for="filter in filters"
          :key="filter.prop"
          class="sf-filter-panel__filter-item"
        >
          <label class="sf-filter-panel__filter-label">
            {{ filter.label }}
            <span v-if="filter.required" class="sf-filter-panel__required">*</span>
          </label>
          
          <!-- 选择器类型 -->
          <Select
            v-if="filter.type === 'select'"
            v-model="localValues[filter.prop]"
            :options="filter.options || []"
            :placeholder="filter.placeholder || '请选择'"
            clearable
            :multiple="filter.multiple"
            style="width: 100%;"
            @change="handleChange"
          />
          
          <!-- 输入框类型 -->
          <Input
            v-else-if="filter.type === 'input'"
            v-model="localValues[filter.prop]"
            :placeholder="filter.placeholder || '请输入'"
            clearable
            @input="handleChange"
          />
          
          <!-- 日期选择器类型 -->
          <el-date-picker
            v-else-if="filter.type === 'date'"
            v-model="localValues[filter.prop]"
            type="date"
            :placeholder="filter.placeholder || '选择日期'"
            clearable
            style="width: 100%;"
            @change="handleChange"
          />
          
          <!-- 日期范围选择器类型 -->
          <el-date-picker
            v-else-if="filter.type === 'daterange'"
            v-model="localValues[filter.prop]"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            clearable
            style="width: 100%;"
            @change="handleChange"
          />
        </div>
      </div>

      <div class="sf-filter-panel__actions">
        <Button @click="handleReset">重置</Button>
        <Button type="primary" @click="handleApply">应用</Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { Input, Select, Button } from '@/components/common/base'

export interface FilterOption {
  label: string
  prop: string
  type?: 'select' | 'input' | 'date' | 'daterange'
  options?: Array<{ label: string; value: any }>
  placeholder?: string
  required?: boolean
  multiple?: boolean
}

interface FilterPanelProps {
  filters: FilterOption[]
  modelValue: Record<string, any>
  collapsible?: boolean
}

const props = withDefaults(defineProps<FilterPanelProps>(), {
  collapsible: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
  'change': [value: Record<string, any>]
  'reset': []
  'apply': [value: Record<string, any>]
}>()

const collapsed = ref(!props.collapsible ? false : true)

// 本地筛选值
const localValues = ref<Record<string, any>>({})

// 初始化筛选值
watch(
  () => props.filters,
  (newFilters) => {
    if (newFilters) {
      const initialValues: Record<string, any> = {}
      newFilters.forEach(filter => {
        initialValues[filter.prop] = props.modelValue[filter.prop] || undefined
      })
      localValues.value = initialValues
    }
  },
  { immediate: true }
)

// 同步外部值
watch(
  () => props.modelValue,
  (newValue) => {
    localValues.value = { ...newValue }
  },
  { deep: true }
)

const toggleCollapse = () => {
  if (props.collapsible) {
    collapsed.value = !collapsed.value
  }
}

const handleChange = () => {
  emit('update:modelValue', { ...localValues.value })
  emit('change', { ...localValues.value })
}

const handleReset = () => {
  const resetValues: Record<string, any> = {}
  props.filters.forEach(filter => {
    resetValues[filter.prop] = undefined
  })
  localValues.value = resetValues
  emit('update:modelValue', {})
  emit('reset')
}

const handleApply = () => {
  emit('apply', { ...localValues.value })
  if (props.collapsible) {
    collapsed.value = true
  }
}
</script>

<style scoped lang="scss">
.sf-filter-panel {
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
  background: var(--el-bg-color);

  &__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    cursor: pointer;
    user-select: none;
    border-bottom: 1px solid var(--el-border-color-lighter);
    transition: background-color 0.2s;

    &:hover {
      background-color: var(--el-bg-color-page);
    }
  }

  &__title {
    font-size: 14px;
    font-weight: 500;
    color: var(--el-text-color-primary);
  }

  &__icon {
    transition: transform 0.3s;

    &.is-rotated {
      transform: rotate(180deg);
    }
  }

  &__content {
    padding: 16px;
  }

  &__filters {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 16px;
    margin-bottom: 16px;
  }

  &__filter-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  &__filter-label {
    font-size: 14px;
    color: var(--el-text-color-regular);
  }

  &__required {
    color: var(--el-color-danger);
    margin-left: 2px;
  }

  &__actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding-top: 16px;
    border-top: 1px solid var(--el-border-color-lighter);
  }
}
</style>

