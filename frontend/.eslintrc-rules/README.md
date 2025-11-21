# ESLint 自定义规则说明

## 为什么不需要自定义规则？

经过评估，**不需要自定义ESLint规则**，原因如下：

### 1. Vue ESLint 内置规则足够

Vue ESLint 插件已经提供了足够的规则来检查组件使用：
- `no-restricted-imports`: 检查导入语句
- `vue/no-restricted-v-bind`: 检查模板中的绑定

### 2. 自定义规则维护成本高

- 需要维护自定义规则代码
- 需要处理 Vue 模板解析的复杂性
- 需要随着 Vue/ESLint 版本更新而更新

### 3. 更简单的方案

使用 `.eslintrc.cjs` 中的配置即可：
- 使用 `no-restricted-imports` 禁止导入 Element Plus 组件
- 使用 `vue/no-restricted-v-bind` 限制模板中的使用
- 在通用组件库目录中允许使用（通过 `overrides`）

### 4. TypeScript 类型检查

TypeScript 编译器已经可以检查类型导入，不需要 ESLint 规则来检查 `FormInstance` 等类型。

## 当前方案

使用 `.eslintrc.cjs` 中的配置，简单有效：
- ✅ 禁止直接导入 Element Plus 组件
- ✅ 允许导入类型（`import type`）
- ✅ 允许在通用组件库中使用
- ✅ 配置简单，易于维护

## 如果未来需要更严格的检查

如果未来需要检查模板中的 `el-*` 标签，可以考虑：
1. 使用 `vue/no-restricted-component-names`（如果 Vue ESLint 支持）
2. 或者使用简单的脚本在 CI/CD 中检查
3. 或者使用 Prettier 插件

但目前，`.eslintrc.cjs` 的配置已经足够。

