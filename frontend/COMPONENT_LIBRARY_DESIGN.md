# 公共组件库设计文档

## 📋 目录

1. [设计理念](#1-设计理念)
2. [组件分类体系](#2-组件分类体系)
3. [核心组件详细设计](#3-核心组件详细设计)
4. [组件开发规范](#4-组件开发规范)
5. [组件使用指南](#5-组件使用指南)

---

## 1. 设计理念

### 1.1 设计原则

#### 🎯 原子化设计（Atomic Design）

```
原子（Atoms）
  ↓
分子（Molecules）
  ↓
组织（Organisms）
  ↓
模板（Templates）
  ↓
页面（Pages）
```

**我们的组件层级**：
- **基础组件（Base）**：原子级，最小UI单元
- **组合组件（Composite）**：分子级，基础组件的组合
- **业务组件（Business）**：组织级，包含业务逻辑

#### 🔄 可组合性

每个组件都应该：
- **独立**：可以单独使用
- **可组合**：可以与其他组件组合
- **可扩展**：支持通过props和slots扩展

#### 🎨 一致性

- **视觉一致性**：统一的颜色、字体、间距
- **交互一致性**：统一的交互模式
- **API一致性**：统一的props命名和结构

### 1.2 组件设计目标

1. **易用性**：API简单直观，开箱即用
2. **灵活性**：支持多种使用场景
3. **可维护性**：代码清晰，易于维护
4. **可测试性**：易于编写测试

---

## 2. 组件分类体系

### 2.1 组件分类图

```
components/common/
├── base/                    # 基础组件（原子级）
│   ├── Button              # 按钮
│   ├── Input               # 输入框
│   ├── Select              # 选择器
│   └── ...
│
├── data-display/           # 数据展示组件
│   ├── DataTable          # 数据表格
│   ├── DataList           # 数据列表
│   └── ...
│
├── feedback/              # 反馈组件
│   ├── Alert              # 警告
│   ├── Toast              # 提示
│   └── ...
│
├── form/                  # 表单组件
│   ├── FormField         # 表单项
│   ├── DatePicker        # 日期选择
│   └── ...
│
├── navigation/           # 导航组件
│   ├── Breadcrumb        # 面包屑
│   └── ...
│
├── layout/               # 布局组件
│   ├── Container         # 容器
│   └── ...
│
└── business/             # 业务通用组件
    ├── SearchBar         # 搜索栏
    └── ...
```

### 2.2 组件依赖关系

```
业务组件 (Business)
    ↓ 依赖
组合组件 (Composite: data-display, form, feedback)
    ↓ 依赖
基础组件 (Base: Button, Input, Select)
    ↓ 依赖
UI框架 (Element Plus)
```

---

## 3. 核心组件详细设计

### 3.1 基础组件（Base）

#### 3.1.1 Button 按钮

**功能**：基础按钮组件，支持多种类型和状态

**Props**：
```typescript
interface ButtonProps {
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'text'
  size?: 'large' | 'default' | 'small'
  disabled?: boolean
  loading?: boolean
  icon?: string
  round?: boolean
  circle?: boolean
  plain?: boolean
  nativeType?: 'button' | 'submit' | 'reset'
}
```

**Events**：
- `click`: 点击事件

**Slots**：
- `default`: 按钮内容
- `icon`: 图标插槽

**使用示例**：
```vue
<Button type="primary" size="large" :loading="isLoading" @click="handleClick">
  提交
</Button>
```

**设计要点**：
- 支持多种类型和尺寸
- 加载状态显示spinner
- 禁用状态视觉反馈
- 支持图标和文字组合

---

#### 3.1.2 Input 输入框

**功能**：文本输入组件

**Props**：
```typescript
interface InputProps {
  modelValue: string | number
  type?: 'text' | 'password' | 'number' | 'email' | 'url' | 'tel'
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  clearable?: boolean
  showPassword?: boolean
  prefixIcon?: string
  suffixIcon?: string
  maxlength?: number
  minlength?: number
  showWordLimit?: boolean
  validateEvent?: boolean
  size?: 'large' | 'default' | 'small'
}
```

**Events**：
- `update:modelValue`: v-model更新
- `focus`: 获得焦点
- `blur`: 失去焦点
- `clear`: 清空
- `input`: 输入事件

**Slots**：
- `prefix`: 前缀内容
- `suffix`: 后缀内容
- `prepend`: 前置内容
- `append`: 后置内容

**使用示例**：
```vue
<Input
  v-model="username"
  placeholder="请输入用户名"
  clearable
  :maxlength="20"
  @clear="handleClear"
/>
```

---

#### 3.1.3 Select 选择器

**功能**：下拉选择组件

**Props**：
```typescript
interface SelectProps {
  modelValue: string | number | Array<string | number>
  options: SelectOption[]
  placeholder?: string
  multiple?: boolean
  disabled?: boolean
  clearable?: boolean
  filterable?: boolean
  allowCreate?: boolean
  size?: 'large' | 'default' | 'small'
  loading?: boolean
}
```

**使用示例**：
```vue
<Select
  v-model="selectedValue"
  :options="options"
  placeholder="请选择"
  filterable
  clearable
/>
```

---

#### 3.1.4 Dialog 对话框

**功能**：模态对话框组件

**Props**：
```typescript
interface DialogProps {
  modelValue: boolean
  title?: string
  width?: string | number
  fullscreen?: boolean
  top?: string
  modal?: boolean
  closeOnClickModal?: boolean
  closeOnPressEscape?: boolean
  showClose?: boolean
  appendToBody?: boolean
  lockScroll?: boolean
}
```

**Events**：
- `update:modelValue`: v-model更新
- `open`: 打开
- `close`: 关闭
- `opened`: 打开后
- `closed`: 关闭后

**Slots**：
- `default`: 对话框内容
- `header`: 头部内容
- `footer`: 底部内容

**使用示例**：
```vue
<Dialog v-model="visible" title="确认删除" width="500px">
  <p>确定要删除这个工作流吗？</p>
  <template #footer>
    <Button @click="visible = false">取消</Button>
    <Button type="primary" @click="handleConfirm">确认</Button>
  </template>
</Dialog>
```

---

#### 3.1.5 Table 表格

**功能**：数据表格组件

**Props**：
```typescript
interface TableProps {
  data: any[]
  columns: TableColumn[]
  stripe?: boolean
  border?: boolean
  size?: 'large' | 'default' | 'small'
  showHeader?: boolean
  highlightCurrentRow?: boolean
  emptyText?: string
  loading?: boolean
  height?: string | number
  maxHeight?: string | number
}
```

**Events**：
- `selection-change`: 选择变化
- `row-click`: 行点击
- `sort-change`: 排序变化

**Slots**：
- `default`: 表格列（使用TableColumn组件）
- `empty`: 空数据内容

**使用示例**：
```vue
<Table :data="tableData" :columns="columns" stripe border>
  <template #empty>
    <Empty description="暂无数据" />
  </template>
</Table>
```

---

### 3.2 数据展示组件（Data Display）

#### 3.2.1 DataTable 数据表格（增强版）

**功能**：增强的数据表格，支持搜索、筛选、分页等

**Props**：
```typescript
interface DataTableProps {
  data: any[]
  columns: TableColumn[]
  loading?: boolean
  pagination?: PaginationConfig
  searchable?: boolean
  filterable?: boolean
  exportable?: boolean
  selectable?: boolean
}
```

**Features**：
- 内置搜索功能
- 内置筛选功能
- 内置分页功能
- 数据导出功能
- 行选择功能

**使用示例**：
```vue
<DataTable
  :data="workflows"
  :columns="columns"
  :pagination="{ page: 1, pageSize: 10, total: 100 }"
  searchable
  filterable
  exportable
/>
```

---

#### 3.2.2 Statistic 统计数字

**功能**：展示统计数字

**Props**：
```typescript
interface StatisticProps {
  title: string
  value: string | number
  prefix?: string
  suffix?: string
  precision?: number
  valueStyle?: CSSProperties
}
```

**使用示例**：
```vue
<Statistic
  title="总工作流数"
  :value="workflowCount"
  suffix="个"
/>
```

---

### 3.3 反馈组件（Feedback）

#### 3.3.1 Alert 警告提示

**功能**：页面内警告提示

**Props**：
```typescript
interface AlertProps {
  title?: string
  type?: 'success' | 'warning' | 'info' | 'error'
  description?: string
  closable?: boolean
  showIcon?: boolean
  center?: boolean
}
```

**使用示例**：
```vue
<Alert
  title="提示"
  type="warning"
  description="工作流执行失败，请检查配置"
  closable
  show-icon
/>
```

---

#### 3.3.2 Toast 轻提示

**功能**：全局轻提示（通过composable使用）

**API**：
```typescript
// 通过composable使用
const toast = useToast()

toast.success('操作成功')
toast.error('操作失败')
toast.warning('警告信息')
toast.info('提示信息')
```

---

### 3.4 表单组件（Form）

#### 3.4.1 FormField 表单项

**功能**：统一的表单项组件，包含标签、输入、错误提示

**Props**：
```typescript
interface FormFieldProps {
  label: string
  prop: string
  required?: boolean
  error?: string
  labelWidth?: string
}
```

**使用示例**：
```vue
<FormField label="工作流名称" prop="name" required :error="errors.name">
  <Input v-model="form.name" placeholder="请输入工作流名称" />
</FormField>
```

---

#### 3.4.2 DatePicker 日期选择器

**功能**：日期/时间选择组件

**Props**：
```typescript
interface DatePickerProps {
  modelValue: Date | string | number
  type?: 'date' | 'datetime' | 'daterange' | 'datetimerange'
  placeholder?: string
  format?: string
  valueFormat?: string
  disabled?: boolean
  clearable?: boolean
}
```

**使用示例**：
```vue
<DatePicker
  v-model="date"
  type="datetime"
  placeholder="选择日期时间"
  format="YYYY-MM-DD HH:mm:ss"
/>
```

---

### 3.5 业务通用组件（Business）

#### 3.5.1 SearchBar 搜索栏

**功能**：统一的搜索栏组件

**Props**：
```typescript
interface SearchBarProps {
  modelValue: string
  placeholder?: string
  searchable?: boolean
  filterable?: boolean
  filters?: FilterOption[]
  onSearch?: (keyword: string) => void
  onFilter?: (filters: Record<string, any>) => void
}
```

**使用示例**：
```vue
<SearchBar
  v-model="searchKeyword"
  placeholder="搜索工作流..."
  :filters="filterOptions"
  @search="handleSearch"
  @filter="handleFilter"
/>
```

---

#### 3.5.2 FilterPanel 筛选面板

**功能**：筛选条件面板

**Props**：
```typescript
interface FilterPanelProps {
  filters: FilterOption[]
  modelValue: Record<string, any>
  collapsible?: boolean
}
```

**使用示例**：
```vue
<FilterPanel
  :filters="filterOptions"
  v-model="filterValues"
  collapsible
/>
```

---

#### 3.5.3 ActionBar 操作栏

**功能**：统一的操作按钮栏

**Props**：
```typescript
interface ActionBarProps {
  actions: ActionItem[]
  align?: 'left' | 'right' | 'center'
}
```

**使用示例**：
```vue
<ActionBar
  :actions="[
    { label: '新建', type: 'primary', onClick: handleCreate },
    { label: '删除', type: 'danger', onClick: handleDelete }
  ]"
/>
```

---

#### 3.5.4 StatusTag 状态标签

**功能**：统一的状态标签组件

**Props**：
```typescript
interface StatusTagProps {
  status: string
  statusMap?: Record<string, { label: string; type: string }>
}
```

**使用示例**：
```vue
<StatusTag status="running" />
```

---

## 4. 组件开发规范

### 4.1 组件文件结构

```
ComponentName/
├── ComponentName.vue          # 组件主文件
├── ComponentName.test.ts      # 单元测试
├── types.ts                   # 类型定义
├── index.ts                   # 导出文件
└── README.md                  # 组件文档
```

### 4.2 组件代码规范

#### 4.2.1 组件命名

- **组件名**：PascalCase，如 `Button`, `DataTable`
- **文件名**：与组件名一致
- **Props**：camelCase，如 `modelValue`, `showHeader`

#### 4.2.2 Props定义

```typescript
// 使用interface定义Props
interface ComponentProps {
  // 必填属性
  requiredProp: string
  // 可选属性，提供默认值
  optionalProp?: number
  // 带默认值的属性
  defaultProp?: boolean
}

// 在组件中使用
const props = withDefaults(defineProps<ComponentProps>(), {
  optionalProp: 0,
  defaultProp: true
})
```

#### 4.2.3 Events定义

```typescript
// 定义Events
interface ComponentEmits {
  'update:modelValue': [value: string]
  'change': [value: string]
  'click': [event: MouseEvent]
}

// 在组件中使用
const emit = defineEmits<ComponentEmits>()
```

#### 4.2.4 Slots定义

```typescript
// 定义Slots
interface ComponentSlots {
  default(): any
  header(): any
  footer(): any
}

// 在组件中使用
defineSlots<ComponentSlots>()
```

### 4.3 组件文档规范

每个组件都应该有完整的文档：

```markdown
# ComponentName

## 功能说明
组件的功能描述

## 基础用法
代码示例

## API

### Props
| 参数 | 说明 | 类型 | 默认值 |
|------|------|------|--------|
| prop1 | 说明 | type | default |

### Events
| 事件名 | 说明 | 参数 |
|--------|------|------|
| event1 | 说明 | param |

### Slots
| 插槽名 | 说明 |
|--------|------|
| slot1 | 说明 |

## 示例
更多使用示例
```

---

## 5. 组件使用指南

### 5.1 组件导入

```typescript
// 方式1：按需导入（推荐）
import { Button, Input, Dialog } from '@/components/common'

// 方式2：全量导入
import * as CommonComponents from '@/components/common'
```

### 5.2 组件注册

```typescript
// 全局注册（在main.ts中）
import { Button, Input } from '@/components/common'

app.component('Button', Button)
app.component('Input', Input)

// 局部注册（在组件中）
import { Button } from '@/components/common'

export default {
  components: {
    Button
  }
}
```

### 5.3 组件组合使用

```vue
<template>
  <Dialog v-model="visible" title="编辑工作流">
    <Form>
      <FormField label="名称" prop="name" required>
        <Input v-model="form.name" />
      </FormField>
      <FormField label="类型" prop="type">
        <Select v-model="form.type" :options="typeOptions" />
      </FormField>
    </Form>
    <template #footer>
      <Button @click="visible = false">取消</Button>
      <Button type="primary" @click="handleSubmit">确认</Button>
    </template>
  </Dialog>
</template>
```

---

## 6. 组件开发优先级

### Phase 1: 核心基础组件（P0）

必须优先实现，所有其他组件都依赖这些组件：

1. Button
2. Input
3. Select
4. Dialog
5. Table
6. Form
7. Card
8. Loading
9. Message
10. Empty

### Phase 2: 重要组件（P1）

尽快实现，提升开发效率：

1. DataTable
2. FormField
3. DatePicker
4. SearchBar
5. FilterPanel
6. ActionBar
7. StatusTag
8. Alert
9. Toast
10. Pagination

### Phase 3: 增强组件（P2）

后续实现，增强用户体验：

1. DataList
2. Statistic
3. Timeline
4. Tree
5. Upload
6. Editor
7. CodeEditor
8. Breadcrumb
9. Steps

---

## 7. 组件测试规范

### 7.1 测试要求

- **单元测试**：所有公共组件必须有单元测试
- **覆盖率**：核心组件覆盖率>90%
- **测试工具**：使用Vitest

### 7.2 测试内容

1. **Props测试**：验证props是否正确传递
2. **Events测试**：验证事件是否正确触发
3. **Slots测试**：验证插槽是否正确渲染
4. **交互测试**：验证用户交互是否正确
5. **边界测试**：验证边界情况处理

---

## 总结

公共组件库是整个前端项目的基础，设计良好的组件库可以：

1. **提升开发效率**：复用组件，减少重复开发
2. **保证一致性**：统一的UI和交互体验
3. **降低维护成本**：集中维护，统一更新
4. **提升代码质量**：经过充分测试的组件更可靠

通过遵循本文档的设计规范，我们可以构建一个高质量、易用、可维护的公共组件库。

---

**最后更新**: 2024年

