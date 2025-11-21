# ESLint 组件使用规范检查

## 方案说明

经过评估，我们采用**简化方案**，而不是复杂的自定义 ESLint 插件。

### ✅ 当前方案（推荐）

使用 **ESLint 内置规则** + **TypeScript 类型检查**：

1. **`.eslintrc.cjs`** - 使用 ESLint 内置规则
   - `no-restricted-imports`: 禁止直接导入 Element Plus 组件
   - 允许类型导入（`import type`）
   - 允许在通用组件库中使用（通过 `overrides`）

2. **TypeScript 编译器** - 类型检查
   - 已经可以检查类型导入
   - 编译时发现类型错误

3. **可选：检查脚本** - 用于 CI/CD
   - `scripts/check-component-usage.js` - 检查模板中的 `el-*` 标签

### ❌ 为什么不使用自定义 ESLint 插件？

1. **维护成本高**
   - 需要维护自定义规则代码
   - 需要处理 Vue 模板解析的复杂性
   - 需要随着 Vue/ESLint 版本更新

2. **Vue ESLint 内置规则足够**
   - `no-restricted-imports` 已经可以检查导入
   - TypeScript 可以检查类型
   - 不需要额外的自定义规则

3. **更简单、更灵活**
   - 配置简单，易于理解
   - 不需要额外的插件加载
   - 更容易调试和维护

## 使用方法

### 1. 运行 ESLint 检查

```bash
npm run lint
```

### 2. 运行组件使用规范检查（可选）

```bash
node scripts/check-component-usage.js
```

### 3. 在 CI/CD 中使用

```yaml
# .github/workflows/ci.yml
- name: Check component usage
  run: |
    npm run lint
    node scripts/check-component-usage.js
```

## 检查规则

### ✅ 允许的操作

1. **导入类型**：
   ```typescript
   import type { FormInstance, FormRules } from 'element-plus'
   ```

2. **导入图标**：
   ```typescript
   import { User, Lock } from '@element-plus/icons-vue'
   ```

3. **在通用组件库中使用**：
   - `src/components/common/base/**/*` - 允许使用 Element Plus

### ❌ 禁止的操作

1. **直接导入组件**：
   ```typescript
   // ❌ 错误
   import { ElButton } from 'element-plus'
   
   // ✅ 正确
   import { Button } from '@/components/common/base'
   ```

2. **在模板中使用 el-* 标签**：
   ```vue
   <!-- ❌ 错误 -->
   <el-button>按钮</el-button>
   
   <!-- ✅ 正确 -->
   <Button>按钮</Button>
   ```

## 配置说明

### `.eslintrc.cjs`

- 使用 `no-restricted-imports` 规则禁止导入 Element Plus 组件
- 通过 `overrides` 在通用组件库中允许使用
- 配置简单，易于维护

### `scripts/check-component-usage.js`

- 可选脚本，用于检查模板中的 `el-*` 标签
- 可以在 CI/CD 中使用
- 不依赖 ESLint，可以独立运行

## 总结

**简化方案的优势**：
- ✅ 配置简单，易于维护
- ✅ 使用标准 ESLint 规则，兼容性好
- ✅ 不需要自定义插件，减少维护成本
- ✅ TypeScript 类型检查已经足够
- ✅ 可选脚本用于更严格的检查

**如果未来需要更严格的检查**：
- 可以考虑使用 Prettier 插件
- 或者使用更专业的代码检查工具
- 但目前这个方案已经足够

