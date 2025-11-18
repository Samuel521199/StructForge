# 快速启动指南

## 🚀 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
```

**重要：** 如果工作流编辑器无法使用，请确保安装 Vue Flow：

```bash
npm install @vue-flow/core
```

### 2. 启动开发服务器

```bash
npm run dev
```

访问浏览器：`http://localhost:5173` (或终端显示的端口)

## 📋 快速测试流程

### 第一步：测试登录/注册

1. 访问 `http://localhost:5173/auth/login`
2. 点击"立即注册"测试注册页面
3. 返回登录页面（注意：后端未连接时登录会失败，这是正常的）

### 第二步：测试工作流列表

1. 如果已登录，访问 `http://localhost:5173/workflow/list`
2. 如果未登录，会自动跳转到登录页
3. 查看工作流列表（目前使用模拟数据）

### 第三步：测试工作流编辑器

1. 点击"创建工作流"或访问 `http://localhost:5173/workflow/editor`
2. 从左侧节点面板拖拽节点到画布
3. 点击节点查看属性面板
4. 尝试连接节点（拖拽连接点）

### 第四步：测试个人中心

1. 点击右上角用户头像（如果已登录）
2. 选择"个人中心"
3. 尝试上传头像和编辑信息

## ⚠️ 注意事项

### 后端 API 未连接
- 登录/注册会调用 API，如果后端未启动会失败
- 这是正常的，可以查看浏览器控制台的网络请求
- 后续连接后端后即可正常工作

### Vue Flow 依赖
如果工作流编辑器报错，请确保已安装：
```bash
npm install @vue-flow/core
```

### 图标显示
如果某些图标不显示，检查：
- Element Plus 图标是否正确导入
- 浏览器控制台是否有错误

## 🐛 常见问题

### 1. 页面空白
- 检查浏览器控制台错误
- 确认所有依赖已安装
- 检查路由配置

### 2. 样式异常
- 确认 Element Plus 样式已导入
- 检查 SCSS 文件是否存在

### 3. 路由跳转失败
- 检查路由守卫配置
- 查看浏览器控制台错误

### 4. 工作流编辑器无法显示
- 确认 `@vue-flow/core` 已安装
- 检查浏览器控制台错误
- 确认 Vue Flow 样式已导入

## 📝 测试检查清单

- [ ] 登录页面正常显示
- [ ] 注册页面正常显示
- [ ] Header 组件正常显示
- [ ] 工作流列表正常显示
- [ ] 工作流编辑器正常显示
- [ ] 节点面板可以拖拽节点
- [ ] 画布可以接收拖放的节点
- [ ] 节点可以连接
- [ ] 个人中心正常显示

## 🔍 调试工具

### 浏览器开发者工具
- 按 `F12` 打开开发者工具
- 查看 `Console` 标签页的错误
- 查看 `Network` 标签页的请求

### Vue DevTools
建议安装 Vue DevTools 浏览器扩展：
- Chrome: [Vue.js devtools](https://chrome.google.com/webstore/detail/vuejs-devtools)
- Firefox: [Vue.js devtools](https://addons.mozilla.org/firefox/addon/vue-js-devtools/)

## 📚 更多信息

详细测试指南请查看：[TESTING_GUIDE.md](./TESTING_GUIDE.md)

开发进度请查看：[DEVELOPMENT_PROGRESS.md](./DEVELOPMENT_PROGRESS.md)

