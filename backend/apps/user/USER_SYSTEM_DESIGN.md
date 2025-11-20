# StructForge 用户系统设计文档

> **版本**: v1.0  
> **设计者**: AI 用户平台架构师  
> **日期**: 2025-01-18

## 1. 系统概述

### 1.1 业务定位

StructForge 用户系统是一个独立的微服务，负责管理平台的所有用户相关功能，包括：
- 用户账户管理（注册、登录、认证）
- 用户资料管理（头像、个人信息）
- 成长体系（等级、经验值）
- 权限查询（从钱包服务获取VIP/订阅信息）

**注意**：付费体系（金币、VIP、订阅）已独立到钱包服务（Wallet Service），详见 `backend/apps/wallet/WALLET_SYSTEM_DESIGN.md`

### 1.2 与钱包服务的关系

```
┌─────────────────┐         ┌─────────────────┐
│   User Service  │         │  Wallet Service │
│                 │         │                 │
│ - 用户账户      │◄───────►│ - 钱包信息      │
│ - 用户资料      │  gRPC   │ - 金币交易      │
│ - 认证授权      │         │ - VIP/订阅      │
│ - 等级经验      │         │ - 支付集成      │
└─────────────────┘         └─────────────────┘
```

**职责划分**：
- **用户服务**：用户账户、资料、认证、等级经验
- **钱包服务**：金币、VIP、订阅、支付

### 1.3 核心设计原则

1. **微服务独立**：用户服务完全独立，通过 gRPC 与其他服务通信
2. **数据一致性**：关键操作使用事务，保证数据一致性
3. **可扩展性**：支持未来扩展更多等级体系和功能
4. **安全性**：密码加密、JWT 认证、邮件验证
5. **性能优化**：关键数据缓存、读写分离（未来）
6. **服务解耦**：与钱包服务解耦，通过 gRPC 查询VIP/订阅信息

## 2. 业务模块划分

### 2.1 模块架构

```
User Service (微服务)
├── 账户模块 (Account)
│   ├── 注册/登录
│   ├── 邮件验证
│   ├── 密码管理
│   └── JWT 认证
│
├── 用户资料模块 (Profile)
│   ├── 基本信息
│   ├── 头像管理
│   └── 个人设置
│
├── 成长模块 (Growth)
│   ├── 等级系统
│   ├── 经验值系统
│   └── 成就系统（可选）
│
└── 权限模块 (Permission)
    ├── VIP 权限
    ├── 等级权限
    └── 功能权限
```

### 2.2 模块职责

#### 账户模块 (Account)
- **注册流程**：用户名/邮箱注册 → 发送验证邮件 → 验证激活
- **登录流程**：用户名/邮箱+密码 → JWT Token 生成
- **认证流程**：JWT Token 验证 → 用户信息注入
- **密码管理**：密码加密（bcrypt）、密码重置、修改密码

#### 用户资料模块 (Profile)
- **基本信息**：昵称、邮箱、手机号、个人简介
- **头像管理**：头像上传、头像裁剪、头像CDN
- **个人设置**：通知设置、隐私设置、主题设置

#### 成长模块 (Growth)
- **等级系统**：
  - 等级定义（1-100级）
  - 等级升级规则
  - 等级权益配置
- **经验值系统**：
  - 经验值获取规则（登录、完成任务、使用功能等）
  - 经验值计算
  - 经验值历史记录

#### 权限模块 (Permission)
- **权限查询**：从钱包服务查询用户VIP/订阅信息
- **权限计算**：综合 VIP 等级和普通等级计算用户权限
- **权限缓存**：权限信息缓存到 Redis，提高性能
- **权限验证**：提供权限验证接口供其他服务调用

## 3. 数据库设计

### 3.1 核心表设计

#### 3.1.1 users（用户基础表）

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,           -- 用户名（唯一）
    email VARCHAR(255) UNIQUE NOT NULL,               -- 邮箱（唯一）
    password_hash VARCHAR(255) NOT NULL,             -- 密码哈希（bcrypt）
    email_verified BOOLEAN DEFAULT FALSE,           -- 邮箱是否已验证
    email_verified_at TIMESTAMP,                    -- 邮箱验证时间
    status VARCHAR(20) DEFAULT 'active',            -- 状态：active, inactive, banned
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,                         -- 最后登录时间
    last_login_ip VARCHAR(45),                       -- 最后登录IP
    
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status)
);
```

#### 3.1.2 user_profiles（用户资料表）

```sql
CREATE TABLE user_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nickname VARCHAR(100),                           -- 昵称
    avatar_url VARCHAR(500),                         -- 头像URL
    bio TEXT,                                        -- 个人简介
    phone VARCHAR(20),                               -- 手机号
    gender VARCHAR(10),                              -- 性别：male, female, other
    birthday DATE,                                   -- 生日
    location VARCHAR(100),                           -- 所在地
    website VARCHAR(255),                            -- 个人网站
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id)
);
```

#### 3.1.3 user_levels（用户等级表）

```sql
CREATE TABLE user_levels (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level INT DEFAULT 1,                             -- 普通等级（1-100）
    experience BIGINT DEFAULT 0,                     -- 经验值
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_level (level)
);
```

#### 3.1.4 experience_logs（经验值记录表）

```sql
CREATE TABLE experience_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL,               -- 动作类型：login, task_complete, workflow_run 等
    experience_points INT NOT NULL,                 -- 获得的经验值
    level_before INT NOT NULL,                      -- 升级前等级
    level_after INT NOT NULL,                       -- 升级后等级
    description VARCHAR(500),                        -- 描述
    related_id VARCHAR(100),                         -- 关联ID
    related_type VARCHAR(50),                        -- 关联类型
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_action_type (action_type),
    INDEX idx_created_at (created_at)
);
```

#### 3.1.8 email_verifications（邮箱验证表）

```sql
CREATE TABLE email_verifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,                     -- 待验证的邮箱
    token VARCHAR(100) UNIQUE NOT NULL,             -- 验证令牌
    type VARCHAR(20) NOT NULL,                      -- 类型：register, reset_password, change_email
    expires_at TIMESTAMP NOT NULL,                  -- 过期时间（通常24小时）
    used BOOLEAN DEFAULT FALSE,                      -- 是否已使用
    used_at TIMESTAMP,                              -- 使用时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_token (token),
    INDEX idx_email (email),
    INDEX idx_expires_at (expires_at)
);
```

#### 3.1.5 level_configs（等级配置表）

```sql
CREATE TABLE level_configs (
    id BIGSERIAL PRIMARY KEY,
    level INT UNIQUE NOT NULL,                      -- 等级（1-100）
    name VARCHAR(50) NOT NULL,                       -- 等级名称（如：新手、初级、中级）
    experience_required BIGINT NOT NULL,            -- 升级所需经验值
    features JSONB,                                 -- 等级权益配置（JSON格式）
    -- features 示例：
    -- {
    --   "max_workflows": 10,
    --   "max_ai_calls_per_day": 100,
    --   "unlock_features": ["workflow_share"]
    -- }
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_level (level)
);
```

### 3.2 表关系图

```
users (1) ──┬── (1) user_profiles
            ├── (1) user_levels
            ├── (N) experience_logs
            └── (N) email_verifications

level_configs (1) ── (N) user_levels (通过 level 字段关联)

注意：钱包相关表（wallets, coin_transactions, vip_purchases, subscriptions）在钱包服务中
```

## 4. API 设计

### 4.1 gRPC API 定义

#### 4.1.1 账户相关 API

```protobuf
// 用户服务
service UserService {
    // 注册
    rpc Register(RegisterRequest) returns (RegisterResponse);
    
    // 登录
    rpc Login(LoginRequest) returns (LoginResponse);
    
    // 验证邮箱
    rpc VerifyEmail(VerifyEmailRequest) returns (VerifyEmailResponse);
    
    // 重新发送验证邮件
    rpc ResendVerificationEmail(ResendVerificationEmailRequest) returns (ResendVerificationEmailResponse);
    
    // 获取用户信息
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    
    // 更新用户信息
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    
    // 修改密码
    rpc ChangePassword(ChangePasswordRequest) returns (ChangePasswordResponse);
    
    // 重置密码（发送邮件）
    rpc RequestPasswordReset(RequestPasswordResetRequest) returns (RequestPasswordResetResponse);
    
    // 重置密码（验证令牌）
    rpc ResetPassword(ResetPasswordRequest) returns (ResetPasswordResponse);
}
```

#### 4.1.2 成长相关 API

```protobuf
// 成长服务
service GrowthService {
    // 获取用户等级信息
    rpc GetUserLevel(GetUserLevelRequest) returns (GetUserLevelResponse);
    
    // 增加经验值
    rpc AddExperience(AddExperienceRequest) returns (AddExperienceResponse);
    
    // 获取经验值记录
    rpc GetExperienceLogs(GetExperienceLogsRequest) returns (GetExperienceLogsResponse);
    
    // 获取等级配置列表
    rpc GetLevelConfigs(GetLevelConfigsRequest) returns (GetLevelConfigsResponse);
}
```

#### 4.1.3 权限相关 API

```protobuf
// 权限服务
service PermissionService {
    // 获取用户权限
    rpc GetUserPermissions(GetUserPermissionsRequest) returns (GetUserPermissionsResponse);
    
    // 检查权限
    rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
    
    // 获取用户功能限制
    rpc GetUserLimits(GetUserLimitsRequest) returns (GetUserLimitsResponse);
}
```

### 4.2 HTTP API（通过 Gateway）

```
# 账户相关
POST   /api/v1/users/register              # 注册
POST   /api/v1/users/login                 # 登录
POST   /api/v1/users/verify-email          # 验证邮箱
POST   /api/v1/users/resend-verification   # 重新发送验证邮件
GET    /api/v1/users/me                    # 获取当前用户信息
PUT    /api/v1/users/me                    # 更新用户信息
POST   /api/v1/users/change-password      # 修改密码
POST   /api/v1/users/request-password-reset # 请求重置密码
POST   /api/v1/users/reset-password        # 重置密码

# 成长相关
GET    /api/v1/users/level                 # 获取用户等级
POST   /api/v1/users/experience           # 增加经验值（内部调用）
GET    /api/v1/users/experience-logs       # 获取经验值记录
GET    /api/v1/levels/configs              # 获取等级配置

# 权限相关
GET    /api/v1/users/permissions           # 获取用户权限
POST   /api/v1/users/check-permission     # 检查权限
GET    /api/v1/users/limits               # 获取用户功能限制
```

## 5. 业务流程设计

### 5.1 用户注册流程

```
1. 用户提交注册信息（用户名、邮箱、密码）
   ↓
2. 验证用户名和邮箱唯一性
   ↓
3. 密码加密（bcrypt）
   ↓
4. 创建用户记录（status = 'inactive', email_verified = false）
   ↓
5. 创建用户钱包（初始金币=0，等级=1，经验值=0）
   ↓
6. 生成邮箱验证令牌
   ↓
7. 发送验证邮件（异步）
   ↓
8. 返回注册成功（提示用户查收邮件）
   ↓
9. 用户点击邮件中的验证链接
   ↓
10. 验证令牌有效性
   ↓
11. 更新用户状态（email_verified = true, status = 'active'）
   ↓
12. 标记验证令牌为已使用
   ↓
13. 给予初始奖励（如：100金币、100经验值）
```

### 5.2 用户登录流程

```
1. 用户提交登录信息（用户名/邮箱 + 密码）
   ↓
2. 查询用户（根据用户名或邮箱）
   ↓
3. 验证密码（bcrypt）
   ↓
4. 检查用户状态（是否激活、是否被封禁）
   ↓
5. 检查邮箱是否已验证（可选，根据业务需求）
   ↓
6. 生成 JWT Token（包含用户ID、用户名、VIP等级、普通等级）
   ↓
7. 更新最后登录时间和IP
   ↓
8. 记录登录日志（可选）
   ↓
9. 返回 Token 和用户信息
```

### 5.3 经验值增加流程

```
1. 用户执行某个动作（如：完成工作流、登录等）
   ↓
2. 业务服务调用增加经验值接口
   ↓
3. 查询经验值规则（根据动作类型）
   ↓
4. 开启数据库事务
   ↓
5. 增加经验值（原子操作）
   ↓
6. 检查是否升级（查询等级配置）
   ↓
7. 如果升级，更新等级
   ↓
8. 记录经验值日志
   ↓
9. 提交事务
   ↓
10. 清除权限缓存（如果升级了）
   ↓
11. 发送升级通知（可选）
```

### 5.4 权限计算流程

```
1. 获取用户信息（普通等级）
   ↓
2. 从钱包服务查询 VIP 信息（gRPC 调用）
   ↓
3. 从钱包服务查询订阅信息（gRPC 调用）
   ↓
4. 查询等级配置（获取等级权益）
   ↓
5. 查询 VIP 配置（从钱包服务或缓存获取）
   ↓
6. 查询订阅配置（从钱包服务或缓存获取）
   ↓
7. 合并权限（VIP 权限优先于等级权限，订阅权限叠加）
   ↓
8. 计算功能限制（如：最大工作流数、每日AI调用次数）
   ↓
9. 缓存权限信息（Redis，TTL = 1小时）
   ↓
10. 返回权限信息
```

## 6. 关键设计决策

### 6.1 服务间通信

**决策**：用户服务通过 gRPC 调用钱包服务获取VIP/订阅信息

**通信场景**：
- 权限计算时查询VIP/订阅信息
- 用户信息展示时查询钱包信息
- 缓存VIP/订阅信息，减少调用频率

### 6.2 经验值获取规则

**决策**：经验值获取规则配置化，存储在配置表或配置文件中

**建议规则**：
- 每日登录：+10 经验值
- 完成工作流：+5 经验值/次
- 执行AI调用：+1 经验值/次
- 分享工作流：+20 经验值/次
- 邀请好友：+50 经验值/人

### 6.3 权限合并策略

**决策**：VIP 权限优先于等级权限，订阅权限叠加，取最大值

**示例**：
- VIP3 允许 100 个工作流
- 等级 10 允许 20 个工作流
- Pro 订阅允许 500 个工作流
- 最终权限：500 个工作流（取最大值）

### 6.4 缓存策略

**决策**：关键数据使用 Redis 缓存

**缓存内容**：
- 用户基本信息（TTL = 1小时）
- 用户权限信息（TTL = 1小时，包含从钱包服务获取的VIP/订阅信息）
- 等级配置（TTL = 24小时）

**缓存更新**：
- 用户信息更新时，清除用户缓存
- 等级升级时，清除权限缓存
- 钱包服务通知VIP/订阅变更时，清除权限缓存（通过消息队列或直接调用）

### 6.6 邮件服务集成

**决策**：邮件服务作为独立模块，支持多种邮件服务商

**建议**：
- 支持 SMTP（通用）
- 支持 SendGrid、阿里云邮件推送等
- 邮件模板化，支持 HTML 模板
- 异步发送，使用消息队列（可选）

### 6.7 支付服务集成

**决策**：支付服务作为独立模块，支持多种支付方式

**建议**：
- 支持支付宝、微信支付
- 支持 Stripe、PayPal（国际）
- 支付回调验证
- 支付订单状态管理

## 7. 技术选型

### 7.1 数据库

- **主数据库**：PostgreSQL（使用 GORM 抽象层）
- **缓存**：Redis（用户信息、权限缓存）

### 7.2 认证

- **JWT**：使用 JWT 进行无状态认证
- **Token 刷新**：支持 Refresh Token 机制
- **Token 存储**：客户端存储（localStorage）

### 7.3 密码加密

- **算法**：bcrypt
- **成本因子**：10（平衡安全性和性能）

### 7.4 邮件服务

- **SMTP**：支持标准 SMTP 协议
- **模板引擎**：支持 HTML 邮件模板

### 7.5 文件存储

- **头像存储**：支持本地存储、OSS（阿里云/腾讯云）、S3
- **CDN**：头像通过 CDN 加速

## 8. 安全考虑

### 8.1 密码安全

- 密码使用 bcrypt 加密
- 密码强度验证（至少8位，包含字母和数字）
- 密码重置令牌有效期24小时

### 8.2 邮箱验证

- 验证令牌使用随机字符串（32位）
- 令牌有效期24小时
- 令牌只能使用一次

### 8.3 防刷机制

- 注册限流：同一IP 1小时内最多注册3次
- 登录限流：同一账户 5分钟内最多尝试5次
- 邮件发送限流：同一邮箱 1小时内最多发送3封

### 8.4 数据安全

- 敏感信息加密存储（如：支付订单号）
- 操作日志记录（重要操作记录审计日志）
- 数据备份（定期备份数据库）

## 9. 性能优化

### 9.1 数据库优化

- 关键字段建立索引
- 分页查询优化
- 读写分离（未来扩展）

### 9.2 缓存优化

- 热点数据缓存
- 缓存预热
- 缓存穿透防护

### 9.3 查询优化

- 避免 N+1 查询
- 使用批量查询
- 合理使用 JOIN

## 10. 监控和日志

### 10.1 关键指标

- 注册成功率
- 登录成功率
- 金币消费量
- VIP 购买量
- 用户活跃度

### 10.2 日志记录

- 用户操作日志
- 支付交易日志
- 错误日志
- 性能日志

## 11. 扩展性考虑

### 11.1 未来扩展

- 多租户支持（企业版）
- 积分系统（独立于金币）
- 成就系统
- 用户标签系统
- 推荐系统

### 11.2 国际化

- 多语言支持
- 多时区支持
- 多货币支持

## 12. 实施计划

### 阶段一：基础功能（MVP）

1. 用户注册/登录
2. 邮箱验证
3. 用户资料管理
4. 基础权限系统

### 阶段二：成长系统

1. 等级系统
2. 经验值系统
3. 权限合并

### 阶段三：高级功能

1. 权限缓存优化
2. 与钱包服务集成
3. 性能优化

**注意**：付费系统（金币、VIP、订阅）在钱包服务中实现

---

**设计完成时间**：2025-01-18  
**下一步**：根据此设计文档开始实现代码

