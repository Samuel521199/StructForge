# 用户服务快速开始指南

## 前置要求

1. **Go 1.21+** 已安装
2. **PostgreSQL** 数据库已启动
3. **Wire** 工具已安装（用于依赖注入）
4. **Protoc** 工具已安装（用于生成 Protobuf 代码）

## 安装工具

### 1. 安装 Wire

```bash
go install github.com/google/wire/cmd/wire@latest
```

### 2. 安装 Protoc 和相关插件

```bash
# 安装 protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 安装 protoc-gen-go-http（Kratos HTTP 路由生成）
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
```

## 快速启动步骤

### 步骤 1: 配置数据库

编辑 `backend/configs/local/user.yaml`，确保数据库配置正确：

```yaml
database:
  adapter_type: postgres
  postgres:
    host: localhost
    port: 5432
    user: postgres
    password: postgres
    dbname: structforge
    sslmode: disable
    timezone: Asia/Shanghai
```

### 步骤 2: 生成 Protobuf 代码

```bash
cd backend

# 如果有 Makefile
make proto

# 或者手动运行（需要先下载 googleapis）
protoc --proto_path=./api \
  --proto_path=./third_party \
  --go_out=paths=source_relative:./api \
  --go-grpc_out=paths=source_relative:./api \
  --go-http_out=paths=source_relative:./api \
  ./api/user/v1/user.proto
```

**注意**：如果缺少 `google/api/annotations.proto`，需要下载：

```bash
# 创建 third_party 目录
mkdir -p backend/third_party

# 下载 googleapis（如果还没有）
git clone https://github.com/googleapis/googleapis.git backend/third_party/googleapis
```

### 步骤 3: 生成 Wire 代码

```bash
cd backend/apps/user/cmd/user
wire
```

这会生成 `wire_gen.go` 文件。

### 步骤 4: 安装依赖

```bash
cd backend
go mod tidy
```

### 步骤 5: 启动服务

#### Windows

```bash
cd backend/apps/user/cmd/user
start-user.bat
```

#### Linux/Mac

```bash
cd backend/apps/user/cmd/user
chmod +x start-user.sh
./start-user.sh
```

#### 手动启动

```bash
cd backend/apps/user/cmd/user
go run main.go wire_gen.go
```

## 验证服务

### 1. 检查服务是否启动

服务启动后，应该看到类似以下日志：

```
[INFO] user 服务启动中
[INFO] 正在加载配置文件...
[INFO] 数据库系统初始化成功
[INFO] 应用实例初始化成功
[INFO] user 服务开始运行
```

### 2. 测试注册接口

```bash
curl -X POST http://localhost:8001/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. 测试登录接口

```bash
curl -X POST http://localhost:8001/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## 常见问题

### 问题 1: 找不到 wire_gen.go

**错误信息**：
```
wire_gen.go: no such file or directory
```

**解决方法**：
```bash
cd backend/apps/user/cmd/user
wire
```

### 问题 2: 找不到 Protobuf 生成的代码

**错误信息**：
```
could not import StructForge/backend/api/user/v1
```

**解决方法**：
```bash
cd backend
make proto
# 或手动运行 protoc 命令
```

### 问题 3: 数据库连接失败

**错误信息**：
```
数据库连接测试失败
```

**解决方法**：
1. 确保 PostgreSQL 已启动
2. 检查 `configs/local/user.yaml` 中的数据库配置
3. 确保数据库 `structforge` 已创建：
   ```sql
   CREATE DATABASE structforge;
   ```

### 问题 4: 端口被占用

**错误信息**：
```
bind: address already in use
```

**解决方法**：
1. 修改 `configs/local/user.yaml` 中的端口配置
2. 或关闭占用端口的进程

## 下一步

- [ ] 配置 Gateway 路由到用户服务
- [ ] 实现 JWT 认证中间件
- [ ] 实现邮件发送功能
- [ ] 添加单元测试和集成测试

