# Gateway 服务 - 通用API网关

## 概述

Gateway 服务是 StructForge 平台的统一入口，提供灵活的、可配置的路由转发功能，支持客户端通过统一的网关连接到各个微服务。

## 核心特性

### 1. 动态路由系统
- **路径匹配**：支持前缀匹配（prefix）、精确匹配（exact）、正则匹配（regex）
- **路由规则**：可配置的路由规则，支持路径重写、超时设置、重试策略
- **灵活配置**：通过 YAML 配置文件定义路由规则

### 2. 服务发现
- **静态服务发现**：基于配置文件的服务实例注册（适合开发环境）
- **动态服务发现**：预留接口，支持集成 Nacos、Consul 等服务发现组件
- **健康检查**：自动过滤不健康的服务实例

### 3. 负载均衡
- **轮询（Round Robin）**：按顺序分发请求
- **随机（Random）**：随机选择服务实例
- **最少连接（Least Connections）**：选择连接数最少的实例

### 4. 认证授权
- **JWT 认证**：支持 JWT Token 验证
- **路由级认证**：可为每个路由单独配置是否需要认证
- **Token 验证**：自动验证 Bearer Token 并提取用户信息

### 5. 请求转发
- **HTTP 代理**：完整的 HTTP 请求/响应转发
- **请求头复制**：自动复制所有请求头到目标服务
- **超时控制**：可配置的请求超时时间
- **错误处理**：完善的错误处理和日志记录

## 配置说明

### 路由配置示例

```yaml
routes:
  # 用户服务路由（不需要认证）
  - path: "/api/v1/users"
    match_type: "prefix"
    service: "user-service"
    require_auth: false
    timeout: 30
    load_balance_strategy: "round_robin"
    rate_limit:
      qps: 100
      burst: 200

  # 需要认证的用户服务路由
  - path: "/api/v1/users/me"
    match_type: "exact"
    service: "user-service"
    require_auth: true
    timeout: 30
```

### 服务配置示例

```yaml
services:
  user-service:
    - id: "user-service-1"
      host: "localhost"
      port: 8001
      weight: 100
      healthy: true
      metadata:
        version: "v1.0.0"
```

## 架构设计

### 目录结构

```
internal/
├── handler/          # HTTP 处理器
│   ├── gateway.go    # 网关主处理器
│   └── provider.go   # 依赖注入提供者
├── router/           # 路由系统
│   ├── router.go     # 路由管理器
│   ├── discovery/    # 服务发现
│   ├── loadbalancer/ # 负载均衡
│   └── loader.go     # 配置加载器
├── middleware/       # 中间件
│   ├── auth.go       # 认证中间件
│   └── jwt/          # JWT 管理器
├── server/           # 服务器配置
└── conf/             # 配置结构
```

### 工作流程

1. **请求接收**：Gateway 接收客户端请求
2. **路由匹配**：根据请求路径匹配路由规则
3. **认证检查**：如果路由需要认证，验证 JWT Token
4. **服务发现**：从服务发现组件获取服务实例列表
5. **负载均衡**：使用负载均衡策略选择服务实例
6. **请求转发**：将请求转发到选定的服务实例
7. **响应返回**：将服务响应返回给客户端

## 使用示例

### 启动服务

```bash
cd backend/apps/gateway/cmd/gateway
go run main.go wire_gen.go app.go logger.go -env local
```

### 客户端请求

```bash
# 不需要认证的请求
curl http://localhost:8000/api/v1/users/register \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'

# 需要认证的请求
curl http://localhost:8000/api/v1/users/me \
  -H "Authorization: Bearer <your-jwt-token>"
```

## 扩展性

### 添加新的服务发现组件

1. 实现 `discovery.ServiceDiscovery` 接口
2. 在 `router/provider.go` 中注册新的服务发现提供者
3. 更新配置文件以使用新的服务发现组件

### 添加新的负载均衡策略

1. 在 `router/loadbalancer/` 中实现新的负载均衡器
2. 在 `loadbalancer.go` 的 `NewLoadBalancer` 函数中注册新策略
3. 在路由配置中使用新的策略名称

### 添加新的中间件

1. 在 `middleware/` 目录中创建新的中间件
2. 在 `handler/gateway.go` 中应用中间件
3. 根据需要配置中间件参数

## 待完善功能

- [ ] 实现动态服务发现（Nacos/Consul 集成）
- [ ] 实现限流中间件
- [ ] 实现请求重试机制
- [ ] 实现熔断器
- [ ] 实现请求/响应日志记录
- [ ] 实现监控和指标收集
- [ ] 支持 WebSocket 代理
- [ ] 支持 gRPC 代理

## 注意事项

1. **JWT Secret Key**：生产环境必须修改默认的 JWT Secret Key
2. **服务实例健康检查**：确保服务实例的健康检查端点正常工作
3. **超时配置**：根据实际业务需求调整超时时间
4. **负载均衡策略**：根据服务特性选择合适的负载均衡策略

