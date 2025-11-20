# Gateway 服务快速启动指南

## 问题诊断

如果遇到 `ERR_CONNECTION_REFUSED` 错误，说明 Gateway 服务没有运行。

## 快速启动方法

### 方法 1：使用启动脚本（最简单）

在项目根目录运行：

```bash
start-backend.bat
```

这会同时启动 Gateway 和 User 服务。

### 方法 2：单独启动 Gateway 服务

1. 打开新的命令提示符（CMD）或 PowerShell 窗口
2. 切换到 Gateway 目录：
   ```bash
   cd backend\apps\gateway\cmd\gateway
   ```
3. 运行启动命令：
   ```bash
   go run main.go wire_gen.go app.go logger.go -env local
   ```

### 方法 3：使用新创建的启动脚本

```bash
cd backend\apps\gateway\cmd\gateway
start-gateway.bat
```

## 验证服务是否启动成功

启动后，在另一个终端窗口运行：

```bash
netstat -ano | findstr ":8000"
```

如果看到类似以下输出，说明 Gateway 服务已启动：
```
TCP    0.0.0.0:8000           0.0.0.0:0              LISTENING       <PID>
```

或者直接在浏览器访问：
```
http://localhost:8000/health
```

应该返回健康检查信息。

## 常见问题

### 1. 如果启动失败，检查以下内容：

- **配置文件是否存在**：`backend/configs/local/gateway.yaml`
- **依赖是否正确**：运行 `go mod tidy`
- **Wire 代码是否生成**：运行 `cd backend/apps/gateway/cmd/gateway && wire`

### 2. 如果端口被占用：

检查是否有其他程序占用 8000 端口：
```bash
netstat -ano | findstr ":8000"
```

如果端口被占用，可以：
- 修改 `backend/configs/local/gateway.yaml` 中的端口号
- 或者关闭占用端口的程序

### 3. 启动后立即退出：

查看 Gateway 服务的日志输出，通常在：
- 控制台输出
- `backend/apps/gateway/cmd/gateway/logs/gateway-*.log`

## 服务架构

```
前端 (localhost:5173)
    ↓ HTTP 请求
Gateway (localhost:8000) ← 需要启动
    ↓ 转发请求
User 服务 (localhost:8001) ← 已启动
```

## 启动顺序

1. **先启动 User 服务**（如果还没启动）
2. **再启动 Gateway 服务**
3. **最后启动前端服务**

Gateway 服务启动后，前端就可以正常发送请求了。

