/**
 * 应用入口文件
 */

import { createApp } from 'vue'
import { pinia } from './stores'
import router from './router'
import App from './App.vue'

// Element Plus（必须在自定义样式之前导入）
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

// 样式（在 Element Plus 之后导入，确保可以覆盖）
import './assets/styles/index.scss'
import './assets/styles/element-override.scss'

// 创建应用实例
const app = createApp(App)

// 使用插件
app.use(pinia)
app.use(router)
app.use(ElementPlus)

// 错误处理
app.config.errorHandler = (err, _instance, info) => {
  console.error('Vue 应用错误:', err, info)
}

// 挂载应用
try {
  app.mount('#app')
  console.log('应用已成功挂载')
} catch (error) {
  console.error('应用挂载失败:', error)
  document.body.innerHTML = `
    <div style="padding: 20px; font-family: Arial, sans-serif;">
      <h1>应用加载失败</h1>
      <p>错误信息: ${error instanceof Error ? error.message : String(error)}</p>
      <p>请检查浏览器控制台获取更多信息。</p>
    </div>
  `
}

