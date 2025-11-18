# StructForge 启动指南

## 🚀 快速启动

### 方式一：一键启动（推荐）⭐

双击运行 `start-all.bat`，会自动启动前端和后端服务。

**特点：**
- ✅ 自动检查环境（Go、Node.js）
- ✅ 自动安装前端依赖（如果缺失）
- ✅ 自动安装 Vue Flow（如果缺失）
- ✅ 前端和后端在独立窗口中运行
- ✅ 关闭窗口即可停止对应服务
- ✅ 显示详细启动信息

### 方式二：分别启动

#### 启动后端

双击运行 `start-backend.bat`

**说明：**
- 启动 Gateway 服务
- 使用 local 环境配置
- 配置文件：`backend/configs/local/gateway.yaml`

#### 启动前端

双击运行 `start-frontend.bat`

**说明：**
- 启动前端开发服务器
- 自动检查并安装依赖
- 默认地址：http://localhost:5173

## 📋 前置要求

### 后端要求

1. **Go 1.21+**
   - 下载：https://golang.org/dl/
   - 验证：`go version`

2. **安装后端依赖**
   ```bash
   cd backend
   go mod tidy
   ```

3. **配置文件**
   - 确保 `backend/configs/local/gateway.yaml` 存在
   - 根据实际情况修改配置

### 前端要求

1. **Node.js 18+**
   - 下载：https://nodejs.org/
   - 验证：`node --version`

2. **npm**
   - 通常随 Node.js 一起安装
   - 验证：`npm --version`

3. **安装前端依赖**
   ```bash
   cd frontend
   npm install
   ```

4. **Vue Flow（必需）**
   ```bash
   cd frontend
   npm install @vue-flow/core
   ```

## 🔧 手动启动命令

### 后端手动启动

```bash
cd backend/apps/gateway/cmd/gateway
go run main.go wire_gen.go app.go logger.go -env local
```

### 前端手动启动

```bash
cd frontend
npm run dev
```

## ⚙️ 配置说明

### 后端配置

- **环境变量**：通过 `-env` 参数指定（local/test/prod）
- **配置文件**：`backend/configs/{env}/gateway.yaml`
- **默认环境**：local

### 前端配置

- **开发服务器端口**：默认 5173（Vite）
- **API 地址**：在 `frontend/src/api/client/index.ts` 中配置
- **环境变量**：支持 `.env` 文件

## 🐛 常见问题

### 1. 后端启动失败

**问题：找不到 Go 环境**
- 解决：安装 Go 1.21+ 并添加到 PATH

**问题：配置文件不存在**
- 解决：检查 `backend/configs/local/gateway.yaml` 是否存在

**问题：依赖缺失**
- 解决：运行 `cd backend && go mod tidy`

### 2. 前端启动失败

**问题：找不到 Node.js**
- 解决：安装 Node.js 18+ 并添加到 PATH

**问题：依赖安装失败**
- 解决：检查网络连接，或使用国内镜像：
  ```bash
  npm config set registry https://registry.npmmirror.com
  ```

**问题：Vue Flow 未安装**
- 解决：运行 `npm install @vue-flow/core`

### 3. 端口冲突

**后端端口冲突**
- 修改 `backend/configs/local/gateway.yaml` 中的端口配置

**前端端口冲突**
- Vite 会自动尝试下一个可用端口
- 或修改 `frontend/vite.config.ts` 中的端口配置

## 📝 服务地址

启动成功后：

- **前端**：http://localhost:5173
  - 会自动在浏览器中打开
  - 如果端口被占用，Vite 会自动使用下一个可用端口
  
- **后端 Gateway**：根据配置文件中的端口
  - 默认配置：查看 `backend/configs/local/gateway.yaml`
  - 常见端口：8000、8080

## 📂 文件说明

### 启动脚本

- **start-all.bat** - 一键启动前后端（推荐）
- **start-backend.bat** - 仅启动后端服务
- **start-frontend.bat** - 仅启动前端服务
- **stop-all.bat** - 停止所有服务

### 使用建议

1. **首次使用**：运行 `start-all.bat`，会自动处理依赖安装
2. **日常开发**：运行 `start-all.bat` 或分别运行对应的 bat 文件
3. **停止服务**：关闭对应窗口，或运行 `stop-all.bat`

## 🛑 停止服务

### 使用批处理文件启动的

#### 方式一：关闭窗口（推荐）
- 关闭对应的命令行窗口即可停止服务
- 前端窗口：关闭前端服务
- 后端窗口：关闭后端服务

#### 方式二：使用停止脚本
- 双击运行 `stop-all.bat` 可以停止所有服务
- 注意：此方式会强制终止进程，可能导致数据未保存

### 手动启动的

- 在命令行窗口按 `Ctrl+C` 停止服务

## 📚 更多信息

- 开发文档：查看 `docs/` 目录
- 前端架构：`frontend/FRONTEND_ARCHITECTURE.md`
- 后端架构：`backend/apps/gateway/GATEWAY_ARCHITECTURE.md`
- 测试指南：`frontend/TESTING_GUIDE.md`

