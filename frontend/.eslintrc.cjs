module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:vue/vue3-essential',
    'plugin:vue/vue3-strongly-recommended',
    'plugin:vue/vue3-recommended',
    'plugin:@typescript-eslint/recommended',
  ],
  parser: 'vue-eslint-parser',
  parserOptions: {
    ecmaVersion: 'latest',
    parser: '@typescript-eslint/parser',
    sourceType: 'module',
  },
  plugins: [
    'vue',
    '@typescript-eslint',
  ],
  rules: {
    // Vue 相关规则
    'vue/multi-word-component-names': 'off',
    'vue/no-v-html': 'warn',
    
    // TypeScript 相关规则
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    '@typescript-eslint/no-explicit-any': 'warn',
    
    // ========== 组件使用规范检查 ==========
    
    // 1. 禁止直接导入 Element Plus 组件（除了类型和图标）
    'no-restricted-imports': [
      'error',
      {
        paths: [
          {
            name: 'element-plus',
            importNames: [
              'ElButton', 'ElInput', 'ElSelect', 'ElDialog', 'ElTable',
              'ElForm', 'ElFormItem', 'ElCard', 'ElLoading', 'ElEmpty',
              'ElCheckbox', 'ElLink', 'ElIcon', 'ElMessage', 'ElNotification',
              'ElBadge', 'ElTag', 'ElTooltip', 'ElPopover', 'ElDropdown',
              'ElMenu', 'ElTabs', 'ElPagination', 'ElSwitch', 'ElRadio',
              'ElDatePicker', 'ElTimePicker', 'ElUpload', 'ElProgress',
            ],
            message: '❌ 禁止直接导入 Element Plus 组件。\n✅ 请使用 @/components/common/base 中的通用组件。\n例如：import { Button } from \'@/components/common/base\'',
          },
        ],
        patterns: [
          {
            group: ['element-plus/es/components/*'],
            message: '❌ 禁止直接导入 Element Plus 组件。\n✅ 请使用 @/components/common/base 中的通用组件。',
          },
        ],
      },
    ],
    
    // 2. 注意：Vue ESLint 没有直接检查模板标签名的规则
    // 但可以通过代码审查和 TypeScript 类型检查来确保规范
    // 如果需要，可以在 CI/CD 中使用脚本检查
  },
  overrides: [
    {
      files: ['*.vue'],
      rules: {
        // 在 Vue 文件中可以添加特定规则
      },
    },
    {
      // 在通用组件库中允许使用 Element Plus（因为它们是封装层）
      files: [
        'src/components/common/base/**/*.vue',
        'src/components/common/base/**/*.ts',
        'src/plugins/**/*.ts', // 插件文件（如 element-plus.ts）
      ],
      rules: {
        'no-restricted-imports': 'off',
      },
    },
  ],
}
