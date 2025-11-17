# GoLand Gateway 运行配置说明

本文档说明如何在 Windows 下的 GoLand IDE 中配置 gateway 服务的运行配置。

## 运行配置设置

### 1. 创建运行配置

1. 在 GoLand 中，点击右上角的运行配置下拉菜单
2. 选择 `Edit Configurations...`
3. 点击 `+` 号，选择 `Go Build`

### 2. 基本配置

#### 配置名称
- **Name**: `Gateway (Local)`

#### 运行类型
- **Run kind**: `Package` 或 `File`
  - 如果选择 `Package`，填写：`StructForge/backend/apps/gateway/cmd/gateway`
  - 如果选择 `File`，填写：`backend/apps/gateway/cmd/gateway/main.go`

#### 工作目录（Working directory）
```
D:\StructForge\backend\apps\gateway\cmd\gateway
```
或者使用相对路径（相对于项目根目录）：
```
$PROJECT_DIR$/backend/apps/gateway/cmd/gateway
```

### 3. 程序参数（Program arguments）

#### 方式一：使用环境参数（推荐）
```
-env local
```

#### 方式二：指定配置文件路径
```
-conf ../../../../configs/local/gateway.yaml
```

#### 方式三：不指定参数（使用默认值）
```
（留空，程序会默认使用 local 环境）
```

### 4. 环境变量（Environment variables）

可选设置：
```
APP_ENV=local
```

如果需要其他环境变量，可以添加：
```
APP_ENV=local
POD_NAME=gateway-local-001
HOSTNAME=localhost
```

### 5. 完整配置示例

#### 配置 1：本地开发环境（Local）

- **Name**: `Gateway (Local)`
- **Run kind**: `Package`
- **Package path**: `StructForge/backend/apps/gateway/cmd/gateway`
- **Working directory**: `$PROJECT_DIR$/backend/apps/gateway/cmd/gateway`
- **Program arguments**: `-env local`
- **Environment variables**: `APP_ENV=local`

#### 配置 2：测试环境（Test）

- **Name**: `Gateway (Test)`
- **Run kind**: `Package`
- **Package path**: `StructForge/backend/apps/gateway/cmd/gateway`
- **Working directory**: `$PROJECT_DIR$/backend/apps/gateway/cmd/gateway`
- **Program arguments**: `-env test`
- **Environment variables**: `APP_ENV=test`

#### 配置 3：生产环境（Prod）

- **Name**: `Gateway (Prod)`
- **Run kind**: `Package`
- **Package path**: `StructForge/backend/apps/gateway/cmd/gateway`
- **Working directory**: `$PROJECT_DIR$/backend/apps/gateway/cmd/gateway`
- **Program arguments**: `-env prod`
- **Environment variables**: `APP_ENV=prod`

### 6. 构建参数（Build flags）

如果需要添加构建标签或构建参数，可以在 **Build flags** 中设置：

```
-tags wireinject
```

注意：`wire.go` 文件使用了 `//go:build wireinject` 构建标签，但在运行 `main.go` 时不需要这个标签。

### 7. 运行前检查清单

在运行 gateway 服务之前，请确保：

- [ ] 配置文件存在：`backend/configs/local/gateway.yaml`
- [ ] Nacos 服务已启动（如果使用 Nacos 配置中心）
- [ ] Redis 服务已启动（如果配置了 Redis）
- [ ] 已运行 `wire` 命令生成依赖注入代码（如果修改了 wire.go）
- [ ] 已运行 `make proto` 生成 protobuf 代码（如果修改了 .proto 文件）

### 8. 常见问题

#### 问题 1：找不到配置文件

**错误信息**：
```
Config file not found: ../../../../configs/local/gateway.yaml
```

**解决方法**：
- 确保工作目录设置为 `backend/apps/gateway/cmd/gateway`
- 检查配置文件是否存在：`backend/configs/local/gateway.yaml`
- 或者使用 `-conf` 参数指定绝对路径

#### 问题 2：Nacos 连接失败

**错误信息**：
```
创建 Nacos 配置客户端失败
```

**解决方法**：
- 确保 Nacos 服务已启动
- 检查 `gateway.yaml` 中的 Nacos 配置是否正确
- 检查 Nacos 服务器地址和端口

#### 问题 3：Wire 依赖注入错误

**错误信息**：
```
gateway 服务初始化失败
```

**解决方法**：
- 运行 `wire` 命令生成依赖注入代码：
  ```bash
  cd backend/apps/gateway/cmd/gateway
  wire
  ```

### 9. 调试配置

如果需要调试，可以：

1. 在运行配置中勾选 **Debug**
2. 设置断点
3. 点击调试按钮（或按 `Shift+F9`）

调试配置与运行配置相同，只需要将运行类型改为调试模式。

### 10. 快速启动命令

如果不想使用 GoLand 运行配置，也可以在终端中运行：

```bash
# 进入工作目录
cd backend/apps/gateway/cmd/gateway

# 运行服务（使用默认 local 环境）
go run main.go wire_gen.go

# 或者指定环境
go run main.go wire_gen.go -env local

# 或者指定配置文件
go run main.go wire_gen.go -conf ../../../../configs/local/gateway.yaml
```

---

**最后更新**: 2024年

