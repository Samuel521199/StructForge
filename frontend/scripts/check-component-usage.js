#!/usr/bin/env node

/**
 * 检查组件使用规范的脚本
 * 用于 CI/CD 或手动检查
 * 
 * 检查规则：
 * 1. 禁止在模板中使用 el-* 标签
 * 2. 禁止直接导入 Element Plus 组件（除了类型和图标）
 */

const fs = require('fs')
const path = require('path')
const { glob } = require('glob')

// Element Plus 组件列表
const ELEMENT_PLUS_COMPONENTS = [
  'el-button', 'el-input', 'el-select', 'el-dialog', 'el-table',
  'el-form', 'el-form-item', 'el-card', 'el-loading', 'el-empty',
  'el-checkbox', 'el-link', 'el-icon', 'el-message', 'el-notification',
]

// 允许使用 Element Plus 的目录
const ALLOWED_DIRS = [
  'src/components/common/base',
  'src/plugins',
]

/**
 * 检查文件是否在允许的目录中
 */
function isAllowedFile(filePath) {
  return ALLOWED_DIRS.some(dir => filePath.includes(dir))
}

/**
 * 检查 Vue 文件中的 el-* 标签
 */
function checkVueFile(filePath) {
  const content = fs.readFileSync(filePath, 'utf-8')
  const errors = []
  
  // 检查模板中的 el-* 标签
  const templateMatch = content.match(/<template[^>]*>([\s\S]*?)<\/template>/)
  if (templateMatch) {
    const template = templateMatch[1]
    
    ELEMENT_PLUS_COMPONENTS.forEach(component => {
      const regex = new RegExp(`<${component}[\\s>]`, 'g')
      if (regex.test(template)) {
        const lines = template.split('\n')
        lines.forEach((line, index) => {
          if (line.includes(`<${component}`)) {
            errors.push({
              file: filePath,
              line: index + 1,
              message: `❌ 禁止使用 ${component} 标签。请使用 @/components/common/base 中的通用组件。`,
            })
          }
        })
      }
    })
  }
  
  // 检查导入语句（除了类型导入和图标）
  const importRegex = /import\s+(?:type\s+)?\{[^}]*\}\s+from\s+['"]element-plus['"]/g
  const typeImportRegex = /import\s+type\s+\{[^}]*\}\s+from\s+['"]element-plus['"]/g
  const lines = content.split('\n')
  
  lines.forEach((line, index) => {
    if (importRegex.test(line) && !typeImportRegex.test(line) && !line.includes('@element-plus/icons-vue')) {
      // 检查是否导入了组件（而不是类型）
      const componentImports = ELEMENT_PLUS_COMPONENTS.map(c => c.replace('el-', 'El')).join('|')
      if (new RegExp(componentImports).test(line)) {
        errors.push({
          file: filePath,
          line: index + 1,
          message: '❌ 禁止直接导入 Element Plus 组件。请使用 @/components/common/base 中的通用组件。',
        })
      }
    }
  })
  
  return errors
}

/**
 * 主函数
 */
async function main() {
  console.log('🔍 开始检查组件使用规范...\n')
  
  const vueFiles = await glob('src/**/*.vue', {
    ignore: ['node_modules/**', 'dist/**'],
  })
  
  const allErrors = []
  
  for (const file of vueFiles) {
    if (isAllowedFile(file)) {
      continue // 跳过允许的目录
    }
    
    const errors = checkVueFile(file)
    if (errors.length > 0) {
      allErrors.push(...errors)
    }
  }
  
  if (allErrors.length > 0) {
    console.error('❌ 发现以下问题：\n')
    allErrors.forEach(error => {
      console.error(`  ${error.file}:${error.line}`)
      console.error(`  ${error.message}\n`)
    })
    console.error(`总共发现 ${allErrors.length} 个问题\n`)
    process.exit(1)
  } else {
    console.log('✅ 所有文件都符合组件使用规范！\n')
    process.exit(0)
  }
}

if (require.main === module) {
  main().catch(error => {
    console.error('检查过程中出现错误：', error)
    process.exit(1)
  })
}

module.exports = { checkVueFile, isAllowedFile }

