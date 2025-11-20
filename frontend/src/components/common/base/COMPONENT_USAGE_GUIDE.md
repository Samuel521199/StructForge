# 通用组件使用指南

## 📚 概述

本项目采用统一的组件化架构，所有 UI 组件都封装在 `@/components/common/base` 中。这些组件是对 Element Plus 的封装，提供了统一的 API 和更好的可维护性。

## 🎯 设计原则

1. **统一性**: 所有页面必须使用通用组件，禁止直接使用 Element Plus 组件
2. **可扩展性**: 组件封装便于后续扩展和定制
3. **类型安全**: 所有组件都提供完整的 TypeScript 类型定义
4. **编辑器友好**: 组件设计考虑未来可视化编辑器的需求

## 📦 可用组件

### 基础组件

- **Button**: 按钮组件
- **Input**: 输入框组件
- **Form**: 表单组件
- **FormItem**: 表单项组件
- **Select**: 选择器组件
- **Checkbox**: 复选框组件
- **Link**: 链接组件
- **Icon**: 图标组件
- **Card**: 卡片组件
- **Dialog**: 对话框组件
- **Table**: 表格组件
- **Loading**: 加载组件
- **Empty**: 空状态组件
- **Message**: 消息提示（success, error, warning, info）

## 🚀 使用方法

### 1. 导入组件

```typescript
// ✅ 正确：从通用组件库导入
import { Button, Input, Form, FormItem, Checkbox, Link } from '@/components/common/base'
import { success, error } from '@/components/common/base/Message'

// ❌ 错误：直接使用 Element Plus
import { ElButton, ElInput } from 'element-plus'
```

### 2. 使用组件

```vue
<template>
  <Form :model="form" :rules="rules">
    <FormItem label="用户名" prop="username">
      <Input v-model="form.username" placeholder="请输入用户名" />
    </FormItem>
    
    <FormItem>
      <Checkbox v-model="agreeTerms">我已阅读并同意</Checkbox>
    </FormItem>
    
    <FormItem>
      <Button type="primary" @click="handleSubmit">提交</Button>
      <Link type="primary" @click="goToLogin">立即登录</Link>
    </FormItem>
  </Form>
</template>

<script setup lang="ts">
import { Form, FormItem, Input, Button, Checkbox, Link } from '@/components/common/base'
import { success, error } from '@/components/common/base/Message'
</script>
```

### 3. 图标使用

```vue
<template>
  <!-- 方式1: 直接导入图标（推荐用于 Input 的 prefix-icon） -->
  <Input :prefix-icon="User" />
  
  <!-- 方式2: 使用 Icon 组件（推荐用于独立图标显示） -->
  <Icon :icon="LoadingIcon" :size="40" :is-loading="true" />
</template>

<script setup lang="ts">
import { User, Lock } from '@element-plus/icons-vue'
import { Icon } from '@/components/common/base'
</script>
```

## 📋 组件迁移清单

### ✅ 已迁移的页面

- ✅ `Register.vue` - 注册页面
- ✅ `Login.vue` - 登录页面
- ✅ `ForgotPassword.vue` - 忘记密码页面
- ✅ `ResetPassword.vue` - 重置密码页面
- ✅ `VerifyEmail.vue` - 邮箱验证页面

### 🔄 迁移步骤

1. **替换导入语句**
   ```typescript
   // 旧代码
   import { ElButton, ElInput } from 'element-plus'
   
   // 新代码
   import { Button, Input } from '@/components/common/base'
   ```

2. **替换组件标签**
   ```vue
   <!-- 旧代码 -->
   <el-button>按钮</el-button>
   <el-input v-model="value" />
   
   <!-- 新代码 -->
   <Button>按钮</Button>
   <Input v-model="value" />
   ```

3. **更新样式选择器**
   ```scss
   // 旧代码
   .el-button { }
   .el-input { }
   
   // 新代码（使用 :deep 或全局样式）
   :deep(.el-button) { }
   :deep(.el-input) { }
   ```

## 🎨 样式规范

### 1. 避免直接使用 Element Plus 类名

```scss
// ❌ 不推荐：直接使用 Element Plus 类名
.el-button--primary {
  background-color: #00FF00;
}

// ✅ 推荐：使用组件类名或 :deep
:deep(.el-button--primary) {
  background-color: #00FF00;
}
```

### 2. 使用组件提供的样式变量

```scss
// 未来可以添加主题变量
$primary-color: #00FF00;
$border-color: rgba(0, 255, 0, 0.3);
```

## 🔧 扩展组件

### 添加新组件

1. 在 `frontend/src/components/common/base/` 下创建组件目录
2. 创建 `ComponentName.vue`, `types.ts`, `index.ts`, `README.md`
3. 在 `base/index.ts` 中导出新组件

### 组件结构示例

```
ComponentName/
├── ComponentName.vue    # 组件实现
├── types.ts             # TypeScript 类型定义
├── index.ts             # 导出文件
└── README.md            # 组件文档
```

## 📝 类型定义

所有组件都提供完整的 TypeScript 类型：

```typescript
import type { ButtonProps, ButtonEmits } from '@/components/common/base'
import type { InputProps, InputEmits } from '@/components/common/base'
```

## 🚫 禁止事项

1. ❌ **禁止直接使用 Element Plus 组件**
   ```vue
   <!-- ❌ 错误 -->
   <el-button>按钮</el-button>
   ```

2. ❌ **禁止在模板中使用 Element Plus 组件**
   ```vue
   <!-- ❌ 错误 -->
   <el-checkbox v-model="value" />
   ```

3. ❌ **禁止直接导入 Element Plus 组件**
   ```typescript
   // ❌ 错误
   import { ElButton } from 'element-plus'
   ```

## ✅ 检查清单

在提交代码前，请确认：

- [ ] 所有组件都从 `@/components/common/base` 导入
- [ ] 没有直接使用 `el-*` 标签
- [ ] 没有直接导入 Element Plus 组件
- [ ] 样式使用 `:deep()` 或全局样式
- [ ] 类型定义已正确导入

## 🔮 未来规划

1. **可视化编辑器支持**: 组件设计考虑未来可视化编辑器的需求
2. **主题系统**: 统一的主题配置和切换
3. **组件文档**: 自动生成组件文档和示例
4. **组件测试**: 完善的单元测试和集成测试

## 📚 相关文档

- [组件库架构说明](./COMPONENT_LIBRARY_GUIDE.md)
- [使用示例](./USAGE_EXAMPLES.md)
- [迁移总结](./MIGRATION_SUMMARY.md)

## 💡 最佳实践

1. **优先使用通用组件**: 新功能开发时优先使用通用组件
2. **保持一致性**: 相同功能的组件使用方式保持一致
3. **及时迁移**: 发现直接使用 Element Plus 的地方及时迁移
4. **文档更新**: 新增组件时及时更新文档

---

**最后更新**: 2025-11-20  
**维护者**: StructForge 开发团队

