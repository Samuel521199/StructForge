# 核心组件使用示例

本文档展示如何在实际项目中使用核心公共组件。

## 📦 组件导入

```typescript
// 推荐：从统一入口导入
import {
  Button,
  Input,
  Select,
  Dialog,
  Table,
  Form,
  Card,
  Loading,
  Empty
} from '@/components/common'

// Message组件特殊导入（方法调用）
import { success, error, warning, info, useMessage } from '@/components/common/base/Message'
```

## 🎯 实际使用场景

### 1. 表单页面

```vue
<template>
  <Card header="创建工作流">
    <Form ref="formRef" :model="form" :rules="rules" label-width="120px">
      <FormItem label="工作流名称" prop="name" required>
        <Input
          v-model="form.name"
          placeholder="请输入工作流名称"
          clearable
          :maxlength="50"
        />
      </FormItem>
      
      <FormItem label="工作流类型" prop="type" required>
        <Select
          v-model="form.type"
          :options="typeOptions"
          placeholder="请选择类型"
          clearable
        />
      </FormItem>
      
      <FormItem label="描述" prop="description">
        <Input
          v-model="form.description"
          type="textarea"
          :rows="4"
          placeholder="请输入描述"
          :maxlength="200"
          show-word-limit
        />
      </FormItem>
      
      <FormItem>
        <Button type="primary" :loading="submitting" @click="handleSubmit">
          创建
        </Button>
        <Button @click="handleReset">重置</Button>
      </FormItem>
    </Form>
  </Card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { success, error } from '@/components/common/base/Message'

const formRef = ref()
const submitting = ref(false)

const form = ref({
  name: '',
  type: '',
  description: ''
})

const rules = {
  name: [
    { required: true, message: '请输入工作流名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择工作流类型', trigger: 'change' }
  ]
}

const typeOptions = [
  { label: '数据处理', value: 'data' },
  { label: 'AI生成', value: 'ai' },
  { label: '自动化', value: 'automation' }
]

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true
    
    // 提交表单
    await createWorkflow(form.value)
    
    success('创建工作流成功')
    handleReset()
  } catch (err: any) {
    if (err.fields) {
      // 表单验证错误
      return
    }
    error('创建工作流失败：' + err.message)
  } finally {
    submitting.value = false
  }
}

const handleReset = () => {
  formRef.value.resetFields()
}
</script>
```

### 2. 列表页面

```vue
<template>
  <Card>
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span>工作流列表</span>
        <Button type="primary" @click="handleCreate">新建工作流</Button>
      </div>
    </template>
    
    <Loading :loading="loading">
      <Table
        :data="tableData"
        :columns="columns"
        stripe
        border
        @row-click="handleRowClick"
      >
        <template #empty>
          <Empty description="暂无工作流">
            <Button type="primary" @click="handleCreate">创建第一个工作流</Button>
          </Empty>
        </template>
      </Table>
    </Loading>
  </Card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from '@/components/common/base/Message'

const router = useRouter()
const message = useMessage()

const loading = ref(false)
const tableData = ref([])

const columns = [
  { prop: 'name', label: '名称', width: 200, sortable: true },
  { prop: 'type', label: '类型', width: 120 },
  { prop: 'status', label: '状态', width: 100 },
  { prop: 'createTime', label: '创建时间', width: 180, sortable: true },
  { prop: 'action', label: '操作', width: 200, fixed: 'right' }
]

onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const data = await fetchWorkflows()
    tableData.value = data
  } catch (err: any) {
    message.error('加载失败：' + err.message)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  router.push('/workflow/editor')
}

const handleRowClick = (row: any) => {
  router.push(`/workflow/detail/${row.id}`)
}
</script>
```

### 3. 确认对话框

```vue
<template>
  <div>
    <Button type="danger" @click="showDeleteDialog">删除</Button>
    
    <Dialog v-model="deleteVisible" title="确认删除" width="400px">
      <p>确定要删除工作流 <strong>{{ currentWorkflow?.name }}</strong> 吗？</p>
      <p style="color: #f56c6c; margin-top: 10px;">此操作不可恢复，请谨慎操作。</p>
      
      <template #footer>
        <Button @click="deleteVisible = false">取消</Button>
        <Button type="danger" :loading="deleting" @click="handleDelete">
          确认删除
        </Button>
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { success, error } from '@/components/common/base/Message'

const deleteVisible = ref(false)
const deleting = ref(false)
const currentWorkflow = ref(null)

const showDeleteDialog = (workflow: any) => {
  currentWorkflow.value = workflow
  deleteVisible.value = true
}

const handleDelete = async () => {
  deleting.value = true
  try {
    await deleteWorkflow(currentWorkflow.value.id)
    success('删除成功')
    deleteVisible.value = false
    // 刷新列表
  } catch (err: any) {
    error('删除失败：' + err.message)
  } finally {
    deleting.value = false
  }
}
</script>
```

### 4. 搜索和筛选

```vue
<template>
  <Card>
    <div style="margin-bottom: 20px;">
      <Input
        v-model="searchKeyword"
        placeholder="搜索工作流名称..."
        clearable
        style="width: 300px; margin-right: 10px;"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </Input>
      
      <Select
        v-model="filterType"
        :options="typeOptions"
        placeholder="筛选类型"
        clearable
        style="width: 150px; margin-right: 10px;"
      />
      
      <Button type="primary" @click="handleSearch">搜索</Button>
      <Button @click="handleReset">重置</Button>
    </div>
    
    <Table :data="filteredData" :columns="columns" />
  </Card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const searchKeyword = ref('')
const filterType = ref('')

const typeOptions = [
  { label: '全部', value: '' },
  { label: '数据处理', value: 'data' },
  { label: 'AI生成', value: 'ai' }
]

const allData = ref([])

const filteredData = computed(() => {
  let result = allData.value
  
  if (searchKeyword.value) {
    result = result.filter(item => 
      item.name.includes(searchKeyword.value)
    )
  }
  
  if (filterType.value) {
    result = result.filter(item => item.type === filterType.value)
  }
  
  return result
})

const handleSearch = () => {
  // 搜索逻辑已在computed中处理
}

const handleReset = () => {
  searchKeyword.value = ''
  filterType.value = ''
}
</script>
```

### 5. 消息提示使用

```vue
<script setup lang="ts">
import { ref } from 'vue'
import { success, error, warning, info, useMessage } from '@/components/common/base/Message'

// 方式1：直接调用
const handleSuccess = () => {
  success('操作成功')
}

const handleError = () => {
  error('操作失败')
}

// 方式2：使用组合函数
const message = useMessage()

const handleSubmit = async () => {
  try {
    await submitForm()
    message.success('提交成功')
  } catch (err: any) {
    message.error('提交失败：' + err.message)
  }
}

// 自定义选项
const handleImportant = () => {
  warning('重要提示', {
    duration: 0, // 不自动关闭
    showClose: true,
    center: true
  })
}
</script>
```

### 6. 加载状态

```vue
<template>
  <!-- 方式1：使用Loading组件 -->
  <Loading :loading="loading" text="加载中...">
    <Table :data="tableData" :columns="columns" />
  </Loading>
  
  <!-- 方式2：使用Table的loading属性 -->
  <Table
    :data="tableData"
    :columns="columns"
    :loading="loading"
  />
  
  <!-- 方式3：使用v-loading指令 -->
  <div v-loading="loading" style="min-height: 200px;">
    内容区域
  </div>
</template>
```

### 7. 空状态

```vue
<template>
  <!-- 基础用法 -->
  <Empty description="暂无数据" />
  
  <!-- 自定义内容 -->
  <Empty description="还没有创建工作流">
    <Button type="primary" @click="handleCreate">创建第一个工作流</Button>
  </Empty>
  
  <!-- 在表格中使用 -->
  <Table :data="tableData" :columns="columns">
    <template #empty>
      <Empty description="暂无工作流" />
    </template>
  </Table>
</template>
```

## 🔗 组件组合使用

### 完整的CRUD页面示例

```vue
<template>
  <div class="workflow-list">
    <!-- 操作栏 -->
    <Card style="margin-bottom: 20px;">
      <div style="display: flex; justify-content: space-between;">
        <div>
          <Input
            v-model="searchKeyword"
            placeholder="搜索..."
            clearable
            style="width: 300px;"
          />
          <Select
            v-model="filterType"
            :options="typeOptions"
            placeholder="类型"
            clearable
            style="width: 150px; margin-left: 10px;"
          />
          <Button type="primary" style="margin-left: 10px;" @click="handleSearch">
            搜索
          </Button>
        </div>
        <Button type="primary" @click="handleCreate">新建工作流</Button>
      </div>
    </Card>
    
    <!-- 数据表格 -->
    <Card>
      <Loading :loading="loading">
        <Table
          :data="tableData"
          :columns="columns"
          stripe
          border
          @row-click="handleRowClick"
        >
          <template #empty>
            <Empty description="暂无工作流">
              <Button type="primary" @click="handleCreate">创建第一个工作流</Button>
            </Empty>
          </template>
        </Table>
      </Loading>
    </Card>
    
    <!-- 删除确认对话框 -->
    <Dialog v-model="deleteVisible" title="确认删除" width="400px">
      <p>确定要删除这个工作流吗？此操作不可恢复。</p>
      <template #footer>
        <Button @click="deleteVisible = false">取消</Button>
        <Button type="danger" :loading="deleting" @click="handleDelete">
          确认删除
        </Button>
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Button,
  Input,
  Select,
  Dialog,
  Table,
  Card,
  Loading,
  Empty
} from '@/components/common'
import { success, error, useMessage } from '@/components/common/base/Message'

const router = useRouter()
const message = useMessage()

// 数据
const loading = ref(false)
const tableData = ref([])
const searchKeyword = ref('')
const filterType = ref('')

// 对话框
const deleteVisible = ref(false)
const deleting = ref(false)
const currentWorkflow = ref(null)

// 表格列
const columns = [
  { prop: 'name', label: '名称', width: 200 },
  { prop: 'type', label: '类型', width: 120 },
  { prop: 'status', label: '状态', width: 100 },
  { prop: 'createTime', label: '创建时间', width: 180 }
]

// 选项
const typeOptions = [
  { label: '全部', value: '' },
  { label: '数据处理', value: 'data' },
  { label: 'AI生成', value: 'ai' }
]

// 方法
onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const data = await fetchWorkflows()
    tableData.value = data
  } catch (err: any) {
    message.error('加载失败：' + err.message)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  router.push('/workflow/editor')
}

const handleSearch = () => {
  loadData()
}

const handleRowClick = (row: any) => {
  router.push(`/workflow/detail/${row.id}`)
}

const showDeleteDialog = (workflow: any) => {
  currentWorkflow.value = workflow
  deleteVisible.value = true
}

const handleDelete = async () => {
  deleting.value = true
  try {
    await deleteWorkflow(currentWorkflow.value.id)
    success('删除成功')
    deleteVisible.value = false
    loadData()
  } catch (err: any) {
    error('删除失败：' + err.message)
  } finally {
    deleting.value = false
  }
}
</script>
```

## 💡 最佳实践

### 1. 组件导入

```typescript
// ✅ 推荐：从统一入口导入
import { Button, Input, Dialog } from '@/components/common'

// ❌ 不推荐：从Element Plus直接导入
import { ElButton, ElInput } from 'element-plus'
```

### 2. 消息提示

```typescript
// ✅ 推荐：使用封装的方法
import { success, error } from '@/components/common/base/Message'
success('操作成功')

// ❌ 不推荐：直接使用Element Plus
import { ElMessage } from 'element-plus'
ElMessage.success('操作成功')
```

### 3. 表单验证

```typescript
// ✅ 推荐：使用Form组件的validate方法
const formRef = ref()
await formRef.value.validate()

// ❌ 不推荐：手动验证
```

### 4. 加载状态

```vue
<!-- ✅ 推荐：使用组件自带的loading属性 -->
<Table :data="data" :columns="columns" :loading="loading" />

<!-- ✅ 也可以使用Loading组件包裹 -->
<Loading :loading="loading">
  <Table :data="data" :columns="columns" />
</Loading>
```

## 📚 相关文档

- [组件库设计文档](./COMPONENT_LIBRARY_DESIGN.md)
- [前端架构设计文档](../FRONTEND_ARCHITECTURE.md)
- 各组件README文档

---

**最后更新**: 2024年

