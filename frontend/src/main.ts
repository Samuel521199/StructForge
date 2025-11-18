/**
 * 应用入口文件
 */

import { createApp } from 'vue'
import { pinia } from './stores'
import router from './router'
import App from './App.vue'

// 样式
import './assets/styles/index.scss'

// Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

// 创建应用实例
const app = createApp(App)

// 使用插件
app.use(pinia)
app.use(router)
app.use(ElementPlus)

// 挂载应用
app.mount('#app')

