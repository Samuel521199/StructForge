# StructForge 钱包系统设计文档

> **版本**: v1.0  
> **设计者**: AI 用户平台架构师  
> **日期**: 2025-01-18

## 1. 系统概述

### 1.1 业务定位

StructForge 钱包系统是一个独立的微服务，负责管理平台的所有虚拟货币和付费相关功能，包括：
- 金币系统（充值、消费、退款）
- VIP 系统（购买、续费、到期管理）
- 订阅系统（购买、续费、取消）
- 交易记录管理
- 支付集成

### 1.2 核心设计原则

1. **服务独立**：钱包服务完全独立，通过 gRPC 与其他服务通信
2. **数据一致性**：所有交易使用事务，保证数据一致性
3. **幂等性**：支付回调、消费操作支持幂等性
4. **可扩展性**：支持未来扩展更多付费模式和货币类型
5. **安全性**：交易记录完整、支付验证严格

### 1.3 与用户服务的关系

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

## 2. 业务模块划分

### 2.1 模块架构

```
Wallet Service (微服务)
├── 钱包模块 (Wallet)
│   ├── 钱包信息查询
│   ├── 余额管理
│   └── 钱包状态管理
│
├── 金币模块 (Coins)
│   ├── 金币充值
│   ├── 金币消费
│   ├── 金币退款
│   ├── 金币奖励
│   └── 交易记录
│
├── VIP 模块 (VIP)
│   ├── VIP 配置管理
│   ├── VIP 购买
│   ├── VIP 续费
│   ├── VIP 到期管理
│   └── VIP 购买记录
│
├── 订阅模块 (Subscription)
│   ├── 订阅套餐管理
│   ├── 订阅购买
│   ├── 订阅续费
│   ├── 订阅取消
│   └── 订阅到期管理
│
└── 支付模块 (Payment)
    ├── 支付订单创建
    ├── 支付回调处理
    ├── 支付方式集成
    └── 支付状态管理
```

### 2.2 模块职责

#### 钱包模块 (Wallet)
- **钱包信息**：查询用户钱包余额、VIP状态、订阅状态
- **余额管理**：实时余额计算、余额锁定（预留）
- **钱包状态**：钱包激活、冻结、解冻

#### 金币模块 (Coins)
- **充值**：支持多种支付方式充值金币
- **消费**：工作流执行、AI调用等场景消费金币
- **退款**：订单取消、服务异常等情况退款
- **奖励**：活动奖励、邀请奖励等
- **交易记录**：完整的交易流水记录

#### VIP 模块 (VIP)
- **配置管理**：VIP 等级配置、价格配置、权益配置
- **购买流程**：选择VIP等级 → 创建订单 → 支付 → 激活
- **续费流程**：VIP到期前续费，延长到期时间
- **到期管理**：定时任务检查到期VIP，自动降级
- **购买记录**：完整的VIP购买历史

#### 订阅模块 (Subscription)
- **套餐管理**：订阅套餐配置、价格配置、权益配置
- **购买流程**：选择套餐 → 创建订单 → 支付 → 激活
- **续费流程**：自动续费或手动续费
- **取消订阅**：取消自动续费，到期后停止服务
- **订阅记录**：完整的订阅历史

#### 支付模块 (Payment)
- **订单创建**：创建支付订单，生成订单号
- **回调处理**：处理支付平台回调，验证签名
- **支付集成**：支持支付宝、微信、Stripe等
- **状态管理**：订单状态流转（pending → paid → completed）

## 3. 数据库设计

### 3.1 核心表设计

#### 3.1.1 wallets（钱包表）

```sql
CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL,                  -- 用户ID（关联用户服务）
    coins DECIMAL(15, 2) DEFAULT 0.00,              -- 金币余额
    frozen_coins DECIMAL(15, 2) DEFAULT 0.00,      -- 冻结金币（预留）
    vip_level INT DEFAULT 0,                        -- VIP等级（0-5，0表示非VIP）
    vip_expires_at TIMESTAMP,                       -- VIP到期时间
    subscription_type VARCHAR(50),                  -- 订阅类型：basic, pro, enterprise
    subscription_expires_at TIMESTAMP,               -- 订阅到期时间
    status VARCHAR(20) DEFAULT 'active',             -- 状态：active, frozen, closed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_vip_level (vip_level),
    INDEX idx_subscription_type (subscription_type),
    INDEX idx_vip_expires_at (vip_expires_at),
    INDEX idx_subscription_expires_at (subscription_expires_at)
);
```

#### 3.1.2 coin_transactions（金币交易记录表）

```sql
CREATE TABLE coin_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                        -- 用户ID
    wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,                      -- 类型：recharge, consume, refund, reward, transfer
    amount DECIMAL(15, 2) NOT NULL,                 -- 金额（正数表示增加，负数表示减少）
    balance_before DECIMAL(15, 2) NOT NULL,        -- 交易前余额
    balance_after DECIMAL(15, 2) NOT NULL,        -- 交易后余额
    description VARCHAR(500),                      -- 交易描述
    related_id VARCHAR(100),                        -- 关联ID（如工作流ID、订单ID等）
    related_type VARCHAR(50),                       -- 关联类型（如 workflow, order 等）
    order_id VARCHAR(100),                          -- 订单ID（如果是充值）
    status VARCHAR(20) DEFAULT 'completed',        -- 状态：pending, completed, failed, cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_wallet_id (wallet_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_related (related_type, related_id),
    INDEX idx_order_id (order_id)
);
```

#### 3.1.3 payment_orders（支付订单表）

```sql
CREATE TABLE payment_orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(100) UNIQUE NOT NULL,           -- 订单号（唯一）
    user_id BIGINT NOT NULL,                        -- 用户ID
    wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    order_type VARCHAR(50) NOT NULL,                -- 订单类型：coin_recharge, vip_purchase, subscription_purchase
    amount DECIMAL(10, 2) NOT NULL,                -- 订单金额
    currency VARCHAR(10) DEFAULT 'CNY',            -- 货币类型
    payment_method VARCHAR(50),                     -- 支付方式：alipay, wechat, stripe
    payment_platform_order_id VARCHAR(200),         -- 支付平台订单ID
    status VARCHAR(20) DEFAULT 'pending',          -- 状态：pending, paid, failed, cancelled, refunded
    paid_at TIMESTAMP,                              -- 支付时间
    expires_at TIMESTAMP,                          -- 订单过期时间
    callback_data JSONB,                           -- 回调数据（JSON格式）
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_order_no (order_no),
    INDEX idx_user_id (user_id),
    INDEX idx_wallet_id (wallet_id),
    INDEX idx_order_type (order_type),
    INDEX idx_status (status),
    INDEX idx_payment_platform_order_id (payment_platform_order_id),
    INDEX idx_expires_at (expires_at)
);
```

#### 3.1.4 vip_configs（VIP配置表）

```sql
CREATE TABLE vip_configs (
    id BIGSERIAL PRIMARY KEY,
    vip_level INT UNIQUE NOT NULL,                  -- VIP等级（1-5）
    name VARCHAR(50) NOT NULL,                      -- VIP名称（如：VIP1、黄金会员）
    price_per_month DECIMAL(10, 2) NOT NULL,       -- 月价格
    price_per_year DECIMAL(10, 2) NOT NULL,        -- 年价格
    features JSONB,                                 -- 权益配置（JSON格式）
    -- features 示例：
    -- {
    --   "max_workflows": 100,
    --   "max_ai_calls_per_day": 1000,
    --   "priority_support": true,
    --   "advanced_features": ["workflow_template", "api_access"]
    -- }
    description TEXT,                               -- 描述
    is_active BOOLEAN DEFAULT TRUE,                -- 是否启用
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_vip_level (vip_level)
);
```

#### 3.1.5 vip_purchases（VIP购买记录表）

```sql
CREATE TABLE vip_purchases (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                        -- 用户ID
    wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    vip_level INT NOT NULL,                          -- 购买的VIP等级
    duration_days INT NOT NULL,                     -- 购买时长（天）
    price DECIMAL(10, 2) NOT NULL,                  -- 价格
    payment_order_id BIGINT REFERENCES payment_orders(id), -- 支付订单ID
    status VARCHAR(20) DEFAULT 'pending',           -- 状态：pending, paid, expired, refunded
    starts_at TIMESTAMP,                            -- VIP开始时间
    expires_at TIMESTAMP,                           -- VIP到期时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_wallet_id (wallet_id),
    INDEX idx_vip_level (vip_level),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_payment_order_id (payment_order_id)
);
```

#### 3.1.6 subscription_plans（订阅套餐表）

```sql
CREATE TABLE subscription_plans (
    id BIGSERIAL PRIMARY KEY,
    plan_type VARCHAR(50) UNIQUE NOT NULL,          -- 套餐类型：basic, pro, enterprise
    name VARCHAR(100) NOT NULL,                     -- 套餐名称
    description TEXT,                               -- 套餐描述
    price_per_month DECIMAL(10, 2) NOT NULL,      -- 月价格
    price_per_year DECIMAL(10, 2) NOT NULL,       -- 年价格
    features JSONB,                                 -- 权益配置（JSON格式）
    -- features 示例：
    -- {
    --   "max_workflows": 1000,
    --   "api_access": true,
    --   "priority_support": true,
    --   "custom_integrations": true
    -- }
    is_active BOOLEAN DEFAULT TRUE,                -- 是否启用
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_plan_type (plan_type)
);
```

#### 3.1.7 subscriptions（订阅记录表）

```sql
CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,                        -- 用户ID
    wallet_id BIGINT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    plan_type VARCHAR(50) NOT NULL,                -- 订阅套餐类型
    price DECIMAL(10, 2) NOT NULL,                 -- 价格
    billing_cycle VARCHAR(20) NOT NULL,            -- 计费周期：monthly, yearly
    payment_order_id BIGINT REFERENCES payment_orders(id), -- 支付订单ID
    status VARCHAR(20) DEFAULT 'active',           -- 状态：active, expired, cancelled, suspended
    starts_at TIMESTAMP,                            -- 订阅开始时间
    expires_at TIMESTAMP,                           -- 订阅到期时间
    auto_renew BOOLEAN DEFAULT TRUE,               -- 是否自动续费
    next_billing_date TIMESTAMP,                    -- 下次计费日期
    cancelled_at TIMESTAMP,                         -- 取消时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_user_id (user_id),
    INDEX idx_wallet_id (wallet_id),
    INDEX idx_plan_type (plan_type),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_next_billing_date (next_billing_date),
    INDEX idx_payment_order_id (payment_order_id)
);
```

### 3.2 表关系图

```
wallets (1) ──┬── (N) coin_transactions
              ├── (N) payment_orders
              ├── (N) vip_purchases
              └── (N) subscriptions

vip_configs (1) ── (N) vip_purchases (通过 vip_level 关联)
subscription_plans (1) ── (N) subscriptions (通过 plan_type 关联)
payment_orders (1) ── (N) vip_purchases
payment_orders (1) ── (N) subscriptions
```

## 4. API 设计

### 4.1 gRPC API 定义

```protobuf
// 钱包服务
service WalletService {
    // 获取钱包信息
    rpc GetWallet(GetWalletRequest) returns (GetWalletResponse);
    
    // 创建钱包（用户注册时调用）
    rpc CreateWallet(CreateWalletRequest) returns (CreateWalletResponse);
    
    // 冻结钱包
    rpc FreezeWallet(FreezeWalletRequest) returns (FreezeWalletResponse);
    
    // 解冻钱包
    rpc UnfreezeWallet(UnfreezeWalletRequest) returns (UnfreezeWalletResponse);
}

// 金币服务
service CoinService {
    // 充值金币（创建支付订单）
    rpc RechargeCoins(RechargeCoinsRequest) returns (RechargeCoinsResponse);
    
    // 消费金币
    rpc ConsumeCoins(ConsumeCoinsRequest) returns (ConsumeCoinsResponse);
    
    // 退款金币
    rpc RefundCoins(RefundCoinsRequest) returns (RefundCoinsResponse);
    
    // 奖励金币
    rpc RewardCoins(RewardCoinsRequest) returns (RewardCoinsResponse);
    
    // 获取金币交易记录
    rpc GetCoinTransactions(GetCoinTransactionsRequest) returns (GetCoinTransactionsResponse);
    
    // 获取金币余额
    rpc GetCoinBalance(GetCoinBalanceRequest) returns (GetCoinBalanceResponse);
}

// VIP 服务
service VIPService {
    // 获取 VIP 配置列表
    rpc GetVIPConfigs(GetVIPConfigsRequest) returns (GetVIPConfigsResponse);
    
    // 购买 VIP（创建支付订单）
    rpc PurchaseVIP(PurchaseVIPRequest) returns (PurchaseVIPResponse);
    
    // 获取用户 VIP 信息
    rpc GetUserVIP(GetUserVIPRequest) returns (GetUserVIPResponse);
    
    // 获取 VIP 购买记录
    rpc GetVIPPurchases(GetVIPPurchasesRequest) returns (GetVIPPurchasesResponse);
    
    // 续费 VIP
    rpc RenewVIP(RenewVIPRequest) returns (RenewVIPResponse);
}

// 订阅服务
service SubscriptionService {
    // 获取订阅套餐列表
    rpc GetSubscriptionPlans(GetSubscriptionPlansRequest) returns (GetSubscriptionPlansResponse);
    
    // 购买订阅（创建支付订单）
    rpc PurchaseSubscription(PurchaseSubscriptionRequest) returns (PurchaseSubscriptionResponse);
    
    // 获取用户订阅信息
    rpc GetUserSubscription(GetUserSubscriptionRequest) returns (GetUserSubscriptionResponse);
    
    // 取消订阅
    rpc CancelSubscription(CancelSubscriptionRequest) returns (CancelSubscriptionResponse);
    
    // 续费订阅
    rpc RenewSubscription(RenewSubscriptionRequest) returns (RenewSubscriptionResponse);
}

// 支付服务
service PaymentService {
    // 创建支付订单
    rpc CreatePaymentOrder(CreatePaymentOrderRequest) returns (CreatePaymentOrderResponse);
    
    // 查询支付订单
    rpc GetPaymentOrder(GetPaymentOrderRequest) returns (GetPaymentOrderResponse);
    
    // 处理支付回调
    rpc HandlePaymentCallback(HandlePaymentCallbackRequest) returns (HandlePaymentCallbackResponse);
    
    // 取消支付订单
    rpc CancelPaymentOrder(CancelPaymentOrderRequest) returns (CancelPaymentOrderResponse);
}
```

### 4.2 HTTP API（通过 Gateway）

```
# 钱包相关
GET    /api/v1/wallet/info                    # 获取钱包信息
POST   /api/v1/wallet/create                  # 创建钱包（内部调用）

# 金币相关
POST   /api/v1/coins/recharge                 # 充值金币
POST   /api/v1/coins/consume                   # 消费金币
POST   /api/v1/coins/refund                   # 退款金币
GET    /api/v1/coins/balance                  # 获取金币余额
GET    /api/v1/coins/transactions             # 获取交易记录

# VIP 相关
GET    /api/v1/vip/configs                     # 获取 VIP 配置
POST   /api/v1/vip/purchase                    # 购买 VIP
GET    /api/v1/vip/info                        # 获取用户 VIP 信息
GET    /api/v1/vip/purchases                   # 获取 VIP 购买记录
POST   /api/v1/vip/renew                       # 续费 VIP

# 订阅相关
GET    /api/v1/subscriptions/plans             # 获取订阅套餐
POST   /api/v1/subscriptions/purchase          # 购买订阅
GET    /api/v1/subscriptions/info              # 获取用户订阅信息
POST   /api/v1/subscriptions/cancel            # 取消订阅
POST   /api/v1/subscriptions/renew             # 续费订阅

# 支付相关
POST   /api/v1/payments/orders                 # 创建支付订单
GET    /api/v1/payments/orders/:id              # 查询支付订单
POST   /api/v1/payments/callback/:platform     # 支付回调（支付宝、微信等）
POST   /api/v1/payments/orders/:id/cancel       # 取消支付订单
```

## 5. 业务流程设计

### 5.1 钱包创建流程

```
1. 用户服务：用户注册成功
   ↓
2. 用户服务调用钱包服务：CreateWallet(user_id)
   ↓
3. 钱包服务创建钱包记录
   - coins = 0
   - vip_level = 0
   - status = 'active'
   ↓
4. 返回钱包ID
   ↓
5. 用户服务记录钱包ID（可选，或通过 user_id 查询）
```

### 5.2 金币充值流程

```
1. 用户选择充值金额
   ↓
2. 调用充值接口：RechargeCoins(amount, payment_method)
   ↓
3. 创建支付订单
   - 生成订单号
   - 记录订单信息
   - 设置订单过期时间（30分钟）
   ↓
4. 调用支付平台创建支付
   ↓
5. 返回支付链接/二维码
   ↓
6. 用户完成支付
   ↓
7. 支付平台回调：HandlePaymentCallback()
   ↓
8. 验证回调签名
   ↓
9. 开启数据库事务
   ↓
10. 更新支付订单状态（pending → paid）
   ↓
11. 增加用户金币
   ↓
12. 记录金币交易流水
   ↓
13. 提交事务
   ↓
14. 通知用户服务（可选，通过消息队列）
```

### 5.3 金币消费流程

```
1. 业务服务调用消费接口：ConsumeCoins(amount, description, related_id, related_type)
   ↓
2. 验证用户钱包状态（是否激活、是否冻结）
   ↓
3. 开启数据库事务
   ↓
4. 锁定钱包记录（SELECT ... FOR UPDATE）
   ↓
5. 检查余额是否充足
   ↓
6. 扣除金币（原子操作）
   ↓
7. 记录交易流水
   ↓
8. 提交事务
   ↓
9. 返回消费结果
```

### 5.4 VIP 购买流程

```
1. 用户选择 VIP 等级和时长
   ↓
2. 调用购买接口：PurchaseVIP(vip_level, duration_days)
   ↓
3. 查询 VIP 配置，计算价格
   ↓
4. 创建支付订单
   ↓
5. 创建 VIP 购买记录（status = 'pending'）
   ↓
6. 调用支付平台创建支付
   ↓
7. 返回支付链接/二维码
   ↓
8. 用户完成支付
   ↓
9. 支付回调处理
   ↓
10. 开启数据库事务
   ↓
11. 更新支付订单状态
   ↓
12. 更新 VIP 购买记录（status = 'paid'）
   ↓
13. 更新钱包 VIP 信息
    - 如果用户已有 VIP，计算续费时间
    - 如果用户没有 VIP，设置开始时间
   ↓
14. 提交事务
   ↓
15. 清除用户服务权限缓存（通过消息队列或直接调用）
```

### 5.5 订阅购买流程

```
1. 用户选择订阅套餐和计费周期
   ↓
2. 调用购买接口：PurchaseSubscription(plan_type, billing_cycle)
   ↓
3. 查询订阅套餐配置，计算价格
   ↓
4. 创建支付订单
   ↓
5. 创建订阅记录（status = 'active', auto_renew = true）
   ↓
6. 调用支付平台创建支付
   ↓
7. 返回支付链接/二维码
   ↓
8. 用户完成支付
   ↓
9. 支付回调处理
   ↓
10. 开启数据库事务
   ↓
11. 更新支付订单状态
   ↓
12. 更新订阅记录
   ↓
13. 更新钱包订阅信息
   ↓
14. 设置下次计费日期
   ↓
15. 提交事务
   ↓
16. 清除用户服务权限缓存
```

### 5.6 VIP/订阅到期管理

```
定时任务（每天凌晨执行）：
1. 查询即将到期的 VIP（expires_at < NOW() + 7 DAYS）
   ↓
2. 发送到期提醒（通过消息队列）
   ↓
3. 查询已到期的 VIP（expires_at < NOW()）
   ↓
4. 开启数据库事务
   ↓
5. 更新钱包 VIP 信息（vip_level = 0, vip_expires_at = NULL）
   ↓
6. 更新 VIP 购买记录（status = 'expired'）
   ↓
7. 提交事务
   ↓
8. 清除用户服务权限缓存
   ↓
9. 发送到期通知（通过消息队列）

订阅到期管理类似
```

## 6. 关键设计决策

### 6.1 服务间通信

**决策**：钱包服务与用户服务通过 gRPC 通信

**通信场景**：
- 用户服务 → 钱包服务：创建钱包、查询钱包信息
- 钱包服务 → 用户服务：VIP/订阅变更时通知（可选，通过消息队列）

### 6.2 幂等性保证

**决策**：所有消费操作支持幂等性

**实现方式**：
- 消费请求携带唯一ID（如：工作流执行ID）
- 检查是否已存在相同 related_id 和 related_type 的交易
- 如果存在，直接返回已存在的交易结果

### 6.3 并发控制

**决策**：使用数据库行锁（SELECT ... FOR UPDATE）防止并发问题

**场景**：
- 金币消费时锁定钱包记录
- VIP 购买时锁定钱包记录
- 订阅购买时锁定钱包记录

### 6.4 支付回调处理

**决策**：支付回调支持重试，保证最终一致性

**实现方式**：
- 回调处理前检查订单状态
- 如果订单已处理，直接返回成功
- 使用事务保证订单状态和钱包余额的一致性

### 6.5 订单过期处理

**决策**：支付订单30分钟过期，过期后自动取消

**实现方式**：
- 定时任务检查过期订单
- 自动取消过期订单
- 释放相关资源（如：VIP 购买记录）

## 7. 技术选型

### 7.1 数据库

- **主数据库**：PostgreSQL（使用 GORM 抽象层）
- **缓存**：Redis（钱包信息缓存、VIP配置缓存）

### 7.2 支付集成

- **国内**：支付宝、微信支付
- **国际**：Stripe、PayPal
- **支付SDK**：根据选择的支付方式使用对应SDK

### 7.3 消息队列（可选）

- **用途**：通知用户服务VIP/订阅变更
- **选型**：RabbitMQ、Kafka（根据项目需求）

## 8. 安全考虑

### 8.1 支付安全

- 支付回调验证签名
- 订单金额验证（防止篡改）
- 订单状态验证（防止重复处理）

### 8.2 交易安全

- 所有交易使用事务
- 关键操作记录审计日志
- 异常交易监控和告警

### 8.3 防刷机制

- 充值限流：同一用户 1小时内最多充值5次
- 消费限流：同一用户 1分钟内最多消费10次

## 9. 监控和日志

### 9.1 关键指标

- 充值成功率
- 消费成功率
- VIP 购买量
- 订阅购买量
- 支付订单完成率

### 9.2 日志记录

- 所有交易记录
- 支付回调日志
- 异常交易日志
- 性能日志

## 10. 扩展性考虑

### 10.1 未来扩展

- 多货币支持（CNY、USD等）
- 积分系统（独立于金币）
- 优惠券系统
- 会员卡系统

### 10.2 国际化

- 多货币支持
- 多支付方式支持
- 多语言支持

---

**设计完成时间**：2025-01-18  
**下一步**：根据此设计文档开始实现代码

