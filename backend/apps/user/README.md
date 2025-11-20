# User Service - 用户服务

## 概述

用户服务是 StructForge 平台的独立微服务，负责管理用户账户、资料、认证等功能。

## 功能特性

- ✅ 用户注册（支持邮箱验证）
- ✅ 用户登录（支持用户名/邮箱登录）
- ✅ JWT 认证
- ✅ 用户资料管理（昵称、头像等）
- ✅ 密码管理（修改密码、重置密码）
- ✅ 邮箱验证

## 快速开始

### 1. 生成 Protobuf 代码

```bash
cd backend
make proto
```

### 2. 生成 Wire 依赖注入代码

```bash
cd apps/user/cmd/user
wire
```

### 3. 配置数据库

编辑 `backend/configs/local/user.yaml`，配置数据库连接信息。

### 4. 启动服务

```bash
cd backend/apps/user/cmd/user
go run main.go wire_gen.go
```

或者使用启动脚本：

```bash
# Windows
start-user.bat

# Linux/Mac
./start-user.sh
```

## API 端点

### HTTP API（通过 Gateway）

```
POST   /api/v1/users/register              # 用户注册
POST   /api/v1/users/login                 # 用户登录
POST   /api/v1/users/verify-email          # 验证邮箱
POST   /api/v1/users/resend-verification   # 重新发送验证邮件
GET    /api/v1/users/me                   # 获取当前用户信息
PUT    /api/v1/users/me                   # 更新用户信息
POST   /api/v1/users/change-password      # 修改密码
POST   /api/v1/users/request-password-reset # 请求重置密码
POST   /api/v1/users/reset-password       # 重置密码
```

### gRPC API

服务定义在 `backend/api/user/v1/user.proto`。

## 配置说明

配置文件位置：`backend/configs/local/user.yaml`

主要配置项：
- `server`: 服务器配置（HTTP/gRPC 端口）
- `database`: 数据库配置
- `nacos`: Nacos 配置中心（可选）

## 数据库

服务使用 PostgreSQL 数据库，通过 `common/data/database` 抽象层访问。

数据库表：
- `users`: 用户基础表
- `user_profiles`: 用户资料表
- `email_verifications`: 邮箱验证表

## 开发说明

### 目录结构

```
apps/user/
├── cmd/user/          # 启动入口
│   ├── main.go        # 主程序
│   ├── wire.go        # Wire 配置
│   └── wire_gen.go    # Wire 生成的代码（运行 wire 后生成）
├── internal/
│   ├── biz/           # 业务逻辑层
│   ├── data/          # 数据访问层
│   ├── service/       # 服务层（gRPC 实现）
│   ├── server/        # 服务器配置
│   └── conf/          # 配置定义
└── USER_SYSTEM_DESIGN.md  # 设计文档
```

### 依赖注入

使用 Google Wire 进行依赖注入。修改依赖后需要重新运行 `wire` 命令。

### 日志

使用 `backend/common/log` 统一日志系统。

## 下一步

- [ ] 实现邮件发送功能
- [ ] 实现头像上传功能
- [ ] 实现用户等级和经验值系统
- [ ] 与钱包服务集成（查询VIP/订阅信息）

