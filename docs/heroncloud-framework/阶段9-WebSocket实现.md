# 阶段9：WebSocket实现

## 📚 学习目标
- 理解WebSocket的工作原理
- 掌握WebSocket服务器的实现
- 学习连接管理和消息处理
- 理解心跳机制和断开处理

---

## 🎯 什么是WebSocket？

### WebSocket vs HTTP

**HTTP：**
- 请求-响应模式
- 客户端发起请求，服务器响应
- 无法主动推送消息

**WebSocket：**
- 全双工通信
- 建立连接后，双方都可以主动发送消息
- 适合实时通信（聊天、游戏、推送等）

### WebSocket连接流程

```
1. 客户端发起HTTP请求（带Upgrade头）
   GET /ws HTTP/1.1
   Upgrade: websocket
   Connection: Upgrade
   ↓
2. 服务器响应101 Switching Protocols
   HTTP/1.1 101 Switching Protocols
   Upgrade: websocket
   Connection: Upgrade
   ↓
3. 连接升级为WebSocket
   ↓
4. 双方可以自由发送消息
```

---

## 📝 核心代码：WebSocket服务器

### 完整代码（手写练习）

```go
package server

import (
    "HeronGame/api/common/v1/pb"
    gatewayv1 "HeronGame/api/gateway/v1/pb"
    "HeronGame/apps/gateway/internal/conf"
    "HeronGame/apps/gateway/internal/handler"
    "HeronGame/apps/gateway/internal/manager"
    "HeronGame/apps/gateway/pkg/pack"
    "HeronGame/common/log"
    "HeronGame/common/middleware/mdx"
    "HeronGame/common/selector"
    "HeronGame/common/websocket/auth"
    "HeronGame/game/errors"
    "context"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
    "github.com/go-kratos/kratos/v2/registry"
    "github.com/go-kratos/kratos/v2/transport"
    "github.com/gorilla/websocket"
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
    "google.golang.org/protobuf/types/known/timestamppb"
)

var _ transport.Server = (*WebSocketServer)(nil)

// WebSocketServer WebSocket服务器
type WebSocketServer struct {
    config   *conf.Bootstrap
    upgrader websocket.Upgrader
    server   *http.Server
    rdb      redis.UniversalClient

    // 管理器
    connectionManager *manager.ConnectionManager
    messageHandler    *handler.MessageHandler

    // 服务注册
    registry *nacos.Registry
    instance *registry.ServiceInstance
}

// NewWebSocketServer 创建WebSocket服务器
func NewWebSocketServer(
    c *conf.Bootstrap,
    messageHandler *handler.MessageHandler,
    connectionManager *manager.ConnectionManager,
    rdb redis.UniversalClient,
    serverRegistry *nacos.Registry,
) *WebSocketServer {
    // 1. 配置Upgrader
    upgrader := websocket.Upgrader{
        ReadBufferSize:  int(c.Server.Websocket.ReadBufferSize),
        WriteBufferSize: int(c.Server.Websocket.WriteBufferSize),
        CheckOrigin: func(r *http.Request) bool {
            return !c.Server.Websocket.CheckOrigin
        },
    }

    // 2. 构建服务实例信息
    instance := &registry.ServiceInstance{
        ID:      c.Server.Id,
        Name:    c.Server.Name,
        Version: c.Server.Version,
        Metadata: map[string]string{
            "kind":        "websocket",
            "instance_id": c.Server.Id,
            "client_addr": c.Server.Websocket.ClientWsAddr,
        },
        Endpoints: []string{
            fmt.Sprintf("%s://%s/ws", c.Server.Websocket.Network, c.Server.Websocket.Addr),
        },
    }

    // 3. 创建WebSocket服务器
    ws := &WebSocketServer{
        config:            c,
        upgrader:          upgrader,
        connectionManager: connectionManager,
        messageHandler:    messageHandler,
        rdb:               rdb,
        registry:          serverRegistry,
        instance:          instance,
    }

    // 4. 创建HTTP服务器用于WebSocket升级
    mux := http.NewServeMux()
    mux.HandleFunc("/ws", ws.handleWebSocket)

    ws.server = &http.Server{
        Addr:    c.Server.Websocket.Addr,
        Handler: mux,
        ReadTimeout: func() time.Duration {
            if c.Server.Websocket.Timeout != nil {
                return c.Server.Websocket.Timeout.AsDuration()
            }
            return 30 * time.Second
        }(),
        WriteTimeout: func() time.Duration {
            if c.Server.Websocket.Timeout != nil {
                return c.Server.Websocket.Timeout.AsDuration()
            }
            return 30 * time.Second
        }(),
    }

    return ws
}

// Start 启动WebSocket服务器
func (ws *WebSocketServer) Start(ctx context.Context) error {
    log.Info(ctx, "[gateway] WebSocket 服务器启动 "+ws.config.Server.Websocket.Addr)
    return ws.server.ListenAndServe()
}

// Stop 停止WebSocket服务器
func (ws *WebSocketServer) Stop(ctx context.Context) error {
    log.Info(ctx, "[gateway] 正在停止 WebSocket 服务器")
    return ws.server.Shutdown(ctx)
}

// handleWebSocket 处理WebSocket连接
func (ws *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    httpCtx := r.Context()

    log.Info(httpCtx, "[gateway] 收到 WebSocket 连接请求",
        zap.String("remote_addr", r.RemoteAddr),
        zap.String("client_ip", clientIP(r)),
        zap.String("user_agent", r.UserAgent()))

    // 1. 升级HTTP连接为WebSocket
    conn, err := ws.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Error(httpCtx, "[gateway] WebSocket 升级失败",
            zap.String("remote_addr", r.RemoteAddr),
            zap.String("client_ip", clientIP(r)),
            zap.Error(err))
        return
    }

    // 2. 在goroutine中处理连接
    go func() {
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()
        defer conn.Close()

        // 3. 处理连接认证和后续逻辑
        ws.handleConnection(ctx, conn, r)
    }()
}

// handleConnection 处理WebSocket连接认证和后续逻辑
func (ws *WebSocketServer) handleConnection(ctx context.Context, conn *websocket.Conn, r *http.Request) {
    // 1. JWT认证检查
    log.Info(ctx, "[gateway] WebSocket 连接升级成功,开始认证",
        zap.String("remote_addr", r.RemoteAddr),
        zap.String("client_ip", clientIP(r)))

    claims, err := auth.AuthenticateWebSocketRequest(
        ctx, r,
        ws.config.Auth.Jwt.Enabled,
        ws.config.Auth.Jwt.TestAssetsUrl,
        ws.config.Auth.Jwt.Secret,
        ws.config.Auth.Jwt.SigningMethod,
        ws.rdb,
    )
    if err != nil {
        log.Error(ctx, "[gateway] WebSocket 认证失败",
            zap.String("remote_addr", r.RemoteAddr),
            zap.String("client_ip", clientIP(r)),
            zap.Error(err))
        ws.connectionManager.ServerDisconnect(ctx, conn, 0, websocket.CloseUnsupportedData, "认证失败", "authentication_failed")
        return
    }

    // 2. 连接数限制检查
    if !ws.connectionManager.CheckConnectionLimit() {
        log.Warn(ctx, "[gateway] WebSocket 连接数达到上限",
            zap.String("remote_addr", r.RemoteAddr),
            zap.String("client_ip", clientIP(r)),
            zap.Int64("user_id", claims.UserID),
            zap.Int64("current", ws.connectionManager.GetConnectionCount()),
            zap.Int32("max", ws.config.Server.Websocket.MaxConnections))
        ws.connectionManager.ServerDisconnect(ctx, conn, 0, websocket.CloseServiceRestart, "服务器连接数已满", "connection_limit_exceeded")
        return
    }

    // 3. 检查是否有旧连接
    oldConn, exist := ws.connectionManager.GetConnection(claims.UserID)
    if exist {
        // 优雅关闭旧连接
        done := ws.closeOldConnection(ctx, oldConn)
        <-done  // 等待旧连接完全清理
        log.Info(ctx, "[gateway] 旧连接已完全清理, 准备建立新连接",
            zap.String("remote_addr", r.RemoteAddr),
            zap.Int64("user_id", claims.UserID))
    }

    // 4. 创建新连接
    connection := ws.connectionManager.AddConnection(
        ctx, conn,
        claims.UserID,
        claims.RoomID,
        claims.GameID,
        clientIP(r),
    )

    defer ws.setupPanicRecovery(ctx, claims.UserID)()

    log.Info(ctx, "[gateway] WebSocket 连接建立成功",
        zap.String("remote_addr", r.RemoteAddr),
        zap.Int64("user_id", claims.UserID),
        zap.Int64("room_id", claims.RoomID),
        zap.Int32("game_id", claims.GameID))

    // 5. 注入用户信息到context
    userCtx := ws.injectUserInfoToContext(ctx, claims.UserID, claims.RoomID, claims.GameID)

    // 6. 发送连接成功消息
    connectedMsg := &gatewayv1.ConnectedMessage{
        ConnectionId: fmt.Sprintf("ws_%d", claims.UserID),
        UserId:       claims.UserID,
        ServerTime:   timestamppb.Now(),
    }
    if err := connection.SendPacket(userCtx, uint32(pb.MessageType_GATEWAY_CONNECTED), 0, connectedMsg); err != nil {
        log.Error(ctx, "[gateway] 发送连接成功消息失败",
            zap.Int64("user_id", claims.UserID),
            zap.Error(err))
    }

    // 7. 开始消息处理循环
    ws.handleInboundMessage(userCtx, connection)
}

// handleInboundMessage 处理WebSocket连接消息循环
func (ws *WebSocketServer) handleInboundMessage(ctx context.Context, connection *manager.Connection) {
    userID := connection.GetUserID()
    roomID := connection.GetRoomID()
    gameID := connection.GetGameID()
    clientIP := connection.GetClientIP()

    log.Info(ctx, "[gateway] WebSocket 开始处理连接消息循环",
        zap.Int64("user_id", userID),
        zap.Int64("room_id", roomID),
        zap.Int32("game_id", gameID))

    // 消息处理循环
    for {
        select {
        case <-ctx.Done():
            // 外部Context被取消，优雅退出
            log.Info(ctx, "[gateway] WebSocket 连接外部context被取消",
                zap.Int64("user_id", userID))
            ws.connectionManager.ServerDisconnect(ctx, nil, userID, 0, "", "context_cancelled")
            return

        case <-connection.GetContext().Done():
            // Connection内部context被取消（优雅关闭触发）
            log.Info(ctx, "[gateway] Connection内部context被取消, 消息循环退出",
                zap.Int64("user_id", userID))
            return

        default:
            // 读取并处理WebSocket消息
            if err := ws.readAndHandleMessage(ctx, connection); err != nil {
                // 检查是否是连接已关闭的错误
                if ws.isConnectionClosedError(err) {
                    log.Info(ctx, "[gateway] WebSocket 连接已关闭, 退出消息循环",
                        zap.Int64("user_id", userID),
                        zap.Error(err))
                    return
                }

                // 处理玩家离线事件
                if offlineErr := ws.messageHandler.HandleClientOffline(ctx, userID, roomID, gameID); offlineErr != nil {
                    log.Error(ctx, "[gateway] 处理玩家离线事件失败",
                        zap.Int64("user_id", userID),
                        zap.Error(offlineErr))
                }

                // 读取消息异常，断开连接
                ws.connectionManager.ClientDisconnect(ctx, userID, roomID, gameID, clientIP, err)
                return
            }

            // 更新最后活跃时间
            connection.UpdateLastPing()
        }
    }
}

// readAndHandleMessage 读取并处理单个WebSocket消息
func (ws *WebSocketServer) readAndHandleMessage(ctx context.Context, conn *manager.Connection) error {
    // 1. 读取WebSocket消息
    messageType, msgData, err := conn.Conn.ReadMessage()
    if err != nil {
        return errors.Wrap(err, "[gateway] 读取WebSocket消息失败")
    }

    // 2. 只支持二进制包协议
    if messageType != websocket.BinaryMessage {
        return nil
    }

    // 3. 解析数据包
    header, payload, err := pack.DecodePacket(msgData)
    if err != nil {
        log.Error(ctx, "[gateway] 数据包解析失败",
            zap.Int64("user_id", conn.GetUserID()),
            zap.Int("data_len", len(msgData)),
            zap.Error(err))
        return nil
    }

    // 4. 记录接收日志
    log.Info(ctx, "[gateway] 接收到二进制包",
        zap.Int64("user_id", conn.GetUserID()),
        zap.Uint32("msg_type", header.MsgType),
        zap.String("msg_type_name", pb.MessageType_name[int32(header.MsgType)]),
        zap.Uint32("request_id", header.RequestId),
        zap.Int("payload_size", len(payload)))

    // 5. 使用MessageHandler处理包
    return ws.messageHandler.HandlePacket(ctx, conn, header, payload)
}

// injectUserInfoToContext 注入用户信息到context
func (ws *WebSocketServer) injectUserInfoToContext(ctx context.Context, userID, roomID int64, gameID int32) context.Context {
    // 注入用户信息
    userCtx := context.WithValue(ctx, mdx.CtxUserID, userID)
    userCtx = context.WithValue(userCtx, mdx.CtxRoomID, roomID)
    userCtx = context.WithValue(userCtx, mdx.CtxGameID, gameID)

    // 用于服务亲和计算
    userCtx = context.WithValue(userCtx, mdx.CtxGatewayServiceID, ws.config.Server.Id)

    // 注入到选择器中，后续这个房间的玩家请求，会路由到固定的一台游戏服务实例上
    userCtx = selector.WithHashKey(userCtx, gameID, roomID)

    log.Info(userCtx, "[gateway] 用户信息已注入context",
        zap.Int64(string(mdx.CtxUserID), userID),
        zap.Int64(string(mdx.CtxRoomID), roomID),
        zap.Int32(string(mdx.CtxGameID), gameID))

    return userCtx
}

// setupPanicRecovery 设置 panic 恢复机制
func (ws *WebSocketServer) setupPanicRecovery(ctx context.Context, userID int64) func() {
    return func() {
        if r := recover(); r != nil {
            ws.connectionManager.ServerDisconnect(ctx, nil, userID, 0, "", "server_panic")
            log.Error(ctx, "[gateway] WebSocket 连接处理发生panic, 已清理资源",
                zap.Int64("user_id", userID),
                zap.Any("panic", r))
            panic(r)  // 重新抛出panic
        }
    }
}

// isConnectionClosedError 检查错误是否是连接已关闭的错误
func (ws *WebSocketServer) isConnectionClosedError(err error) bool {
    if err == nil {
        return false
    }

    errorMsg := err.Error()
    closedErrorPatterns := []string{
        "use of closed network connection",
        "connection reset by peer",
        "broken pipe",
        "websocket: close sent",
        "websocket: connection closed",
    }

    for _, pattern := range closedErrorPatterns {
        if strings.Contains(errorMsg, pattern) {
            return true
        }
    }

    return false
}

// clientIP 获取客户端IP
func clientIP(r *http.Request) string {
    ip := r.Header.Get("X-Real-IP")
    if ip == "" {
        ip = r.Header.Get("X-Forwarded-For")
    }
    return ip
}
```

---

## 📖 WebSocket详解

### 1. 连接升级

```go
conn, err := ws.upgrader.Upgrade(w, r, nil)
```

**过程：**
1. 客户端发送HTTP请求，带 `Upgrade: websocket` 头
2. 服务器验证请求
3. 返回 `101 Switching Protocols` 响应
4. 连接升级为WebSocket

### 2. 连接管理

**连接状态：**
- 已连接：在连接管理器中
- 已断开：从连接管理器中移除
- 重连：关闭旧连接，创建新连接

**连接限制：**
```go
if !ws.connectionManager.CheckConnectionLimit() {
    // 连接数达到上限，拒绝连接
    return
}
```

### 3. 消息处理循环

```go
for {
    select {
    case <-ctx.Done():
        // Context取消，退出循环
        return
    default:
        // 读取消息
        if err := ws.readAndHandleMessage(ctx, connection); err != nil {
            // 处理错误
            return
        }
    }
}
```

**为什么用select？**
- 可以同时监听多个channel
- 可以处理超时和取消
- 非阻塞读取

### 4. 心跳机制

```go
// 更新最后活跃时间
connection.UpdateLastPing()
```

**作用：**
- 检测连接是否存活
- 定期发送ping/pong
- 超时自动断开

### 5. 优雅关闭

```go
// 关闭旧连接
done := ws.closeOldConnection(ctx, oldConn)
<-done  // 等待完全清理
```

**步骤：**
1. 通知客户端即将关闭
2. 停止接收新消息
3. 处理完现有消息
4. 关闭连接
5. 清理资源

---

## 🎓 WebSocket最佳实践

### 1. 错误处理

```go
// ✅ 好的做法：区分不同类型的错误
if ws.isConnectionClosedError(err) {
    // 连接已关闭，正常退出
    return
}
// 其他错误，记录并断开
log.Error(ctx, "处理消息失败", zap.Error(err))
ws.connectionManager.ClientDisconnect(...)
```

### 2. 并发安全

```go
// ✅ 好的做法：每个连接在独立的goroutine中处理
go func() {
    ws.handleConnection(ctx, conn, r)
}()

// ❌ 不好的做法：在主goroutine中阻塞处理
ws.handleConnection(ctx, conn, r)  // 会阻塞其他连接
```

### 3. 资源清理

```go
// ✅ 好的做法：使用defer确保清理
defer conn.Close()
defer ws.setupPanicRecovery(ctx, userID)()

// ❌ 不好的做法：忘记关闭连接
// conn.Close()  // 可能忘记调用
```

### 4. 消息格式

```go
// ✅ 好的做法：使用二进制协议
if messageType != websocket.BinaryMessage {
    return nil  // 只处理二进制消息
}

// ❌ 不好的做法：处理所有类型的消息
// 可能导致安全问题
```

---

## 💡 实践练习

### 练习1：理解连接流程

画出WebSocket连接的完整流程：
```
HTTP请求 → 升级 → 认证 → 创建连接 → 消息循环 → 断开
```

### 练习2：实现心跳检测

编写一个心跳检测函数：

```go
func CheckHeartbeat(connection *Connection, timeout time.Duration) bool {
    // TODO: 检查最后活跃时间是否超过timeout
}
```

### 练习3：优雅关闭

实现一个优雅关闭连接的函数：

```go
func GracefulClose(connection *Connection) error {
    // TODO: 实现优雅关闭逻辑
}
```

---

## 📌 总结

恭喜你完成了所有9个阶段的学习！现在你已经：

✅ 理解了微服务架构设计
✅ 掌握了Kratos框架的使用
✅ 学会了依赖注入（Wire）
✅ 理解了HTTP/gRPC/WebSocket服务器
✅ 掌握了日志、中间件、认证等核心功能
✅ 能够手写核心代码并理解每一行

---

## ❓ 思考题

1. WebSocket和HTTP长轮询有什么区别？各有什么优缺点？
2. 为什么每个WebSocket连接要在独立的goroutine中处理？
3. 心跳机制的作用是什么？如何设计一个高效的心跳机制？
4. 如何实现WebSocket连接的负载均衡？

---

## 🎉 下一步

现在你可以：
1. 尝试搭建自己的微服务项目
2. 扩展现有功能
3. 优化性能
4. 深入学习Go语言和微服务架构

**继续加油！** 💪🚀

