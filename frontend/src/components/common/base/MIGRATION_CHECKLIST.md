# 组件迁移检查清单

## 📋 迁移前检查

### 1. 导入检查

- [ ] 检查是否直接导入 Element Plus 组件
  ```typescript
  // ❌ 需要替换
  import { ElButton, ElInput, ElForm } from 'element-plus'
  
  // ✅ 替换为
  import { Button, Input, Form } from '@/components/common/base'
  ```

- [ ] 检查是否使用 `el-*` 标签
  ```vue
  <!-- ❌ 需要替换 -->
  <el-button>按钮</el-button>
  <el-input v-model="value" />
  <el-checkbox v-model="checked" />
  <el-link>链接</el-link>
  
  <!-- ✅ 替换为 -->
  <Button>按钮</Button>
  <Input v-model="value" />
  <Checkbox v-model="checked" />
  <Link>链接</Link>
  ```

### 2. 组件替换映射表

| Element Plus | 通用组件 | 说明 |
|-------------|---------|------|
| `el-button` | `Button` | 按钮组件 |
| `el-input` | `Input` | 输入框组件 |
| `el-form` | `Form` | 表单组件 |
| `el-form-item` | `FormItem` | 表单项组件 |
| `el-select` | `Select` | 选择器组件 |
| `el-checkbox` | `Checkbox` | 复选框组件 |
| `el-link` | `Link` | 链接组件 |
| `el-icon` | `Icon` | 图标组件（可选） |
| `el-card` | `Card` | 卡片组件 |
| `el-dialog` | `Dialog` | 对话框组件 |
| `el-table` | `Table` | 表格组件 |
| `el-loading` | `Loading` | 加载组件 |
| `el-empty` | `Empty` | 空状态组件 |
| `ElMessage` | `Message` | 消息提示 |

### 3. 样式检查

- [ ] 检查是否使用 `.el-*` 类名选择器
  ```scss
  // ❌ 需要更新
  .el-button { }
  .el-input { }
  
  // ✅ 更新为
  :deep(.el-button) { }
  :deep(.el-input) { }
  // 或使用全局样式
  ```

- [ ] 检查是否使用 Element Plus 的 CSS 变量
  ```scss
  // 如果使用了 Element Plus 的 CSS 变量，需要检查兼容性
  color: var(--el-color-primary);
  ```

### 4. 类型检查

- [ ] 检查是否使用 Element Plus 的类型
  ```typescript
  // ❌ 需要替换
  import type { FormInstance } from 'element-plus'
  
  // ✅ 可以继续使用（FormInstance 是 Element Plus 的类型，通用组件内部使用）
  // 但建议使用通用组件提供的类型
  ```

## 🔄 迁移步骤

### Step 1: 更新导入

```typescript
// 旧代码
import { ElButton, ElInput, ElForm, ElFormItem } from 'element-plus'
import { ElCheckbox, ElLink } from 'element-plus'

// 新代码
import { 
  Button, 
  Input, 
  Form, 
  FormItem, 
  Checkbox, 
  Link 
} from '@/components/common/base'
```

### Step 2: 更新模板

```vue
<!-- 旧代码 -->
<template>
  <el-form :model="form" :rules="rules">
    <el-form-item label="用户名" prop="username">
      <el-input v-model="form.username" />
    </el-form-item>
    <el-form-item>
      <el-checkbox v-model="agree">同意</el-checkbox>
      <el-link>链接</el-link>
    </el-form-item>
    <el-form-item>
      <el-button type="primary">提交</el-button>
    </el-form-item>
  </el-form>
</template>

<!-- 新代码 -->
<template>
  <Form :model="form" :rules="rules">
    <FormItem label="用户名" prop="username">
      <Input v-model="form.username" />
    </FormItem>
    <FormItem>
      <Checkbox v-model="agree">同意</Checkbox>
      <Link>链接</Link>
    </FormItem>
    <FormItem>
      <Button type="primary">提交</Button>
    </FormItem>
  </Form>
</template>
```

### Step 3: 更新样式

```scss
// 旧代码
.el-form-item__label {
  color: #fff;
}

.el-button--primary {
  background-color: #00FF00;
}

// 新代码（使用 :deep）
:deep(.el-form-item__label) {
  color: #fff;
}

:deep(.el-button--primary) {
  background-color: #00FF00;
}
```

### Step 4: 更新消息提示

```typescript
// 旧代码
import { ElMessage } from 'element-plus'
ElMessage.success('成功')
ElMessage.error('错误')

// 新代码
import { success, error } from '@/components/common/base/Message'
success('成功')
error('错误')
```

## ✅ 迁移后验证

### 功能验证

- [ ] 所有表单功能正常
- [ ] 所有按钮点击事件正常
- [ ] 所有输入框双向绑定正常
- [ ] 所有消息提示正常显示
- [ ] 所有样式显示正常

### 代码检查

- [ ] 没有 `el-*` 标签残留
- [ ] 没有直接导入 Element Plus 组件
- [ ] 所有组件都从 `@/components/common/base` 导入
- [ ] TypeScript 类型检查通过
- [ ] ESLint 检查通过

### 样式检查

- [ ] 所有样式正常显示
- [ ] 响应式布局正常
- [ ] 主题色正常
- [ ] 动画效果正常

## 📝 迁移记录模板

```markdown
## 迁移记录

**文件**: `src/views/example/Example.vue`
**日期**: 2025-11-20
**迁移内容**:
- ✅ 替换 `el-button` → `Button`
- ✅ 替换 `el-input` → `Input`
- ✅ 替换 `el-form` → `Form`
- ✅ 更新导入语句
- ✅ 更新样式选择器

**验证结果**:
- ✅ 功能正常
- ✅ 样式正常
- ✅ 类型检查通过
```

## 🚨 常见问题

### Q1: Icon 组件如何使用？

**A**: Icon 组件主要用于独立图标显示，Input 的 prefix-icon 可以直接使用图标组件：

```vue
<template>
  <!-- Input 的 prefix-icon 直接使用图标组件 -->
  <Input :prefix-icon="User" />
  
  <!-- 独立图标使用 Icon 组件 -->
  <Icon :icon="LoadingIcon" :size="40" :is-loading="true" />
</template>

<script setup lang="ts">
import { User, Loading } from '@element-plus/icons-vue'
import { Icon } from '@/components/common/base'
</script>
```

### Q2: FormInstance 类型从哪里导入？

**A**: FormInstance 是 Element Plus 的类型，通用组件内部使用，可以继续从 `element-plus` 导入：

```typescript
import type { FormInstance, FormRules } from 'element-plus'
```

### Q3: 样式不生效怎么办？

**A**: 使用 `:deep()` 选择器或全局样式：

```scss
// 使用 :deep()
:deep(.el-button--primary) {
  background-color: #00FF00;
}

// 或使用全局样式（在 <style> 中不使用 scoped）
<style>
.el-button--primary {
  background-color: #00FF00;
}
</style>
```

---

**最后更新**: 2025-11-20

