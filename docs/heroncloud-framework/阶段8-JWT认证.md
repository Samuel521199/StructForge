# 阶段8：JWT认证

## 📚 学习目标
- 理解JWT的原理和结构
- 掌握Token的生成和解析
- 学习JWT在认证中的应用
- 理解Token刷新机制

---

## 🎯 什么是JWT？

### JWT简介

**JWT（JSON Web Token）** 是一种开放标准（RFC 7519），用于在各方之间安全地传输信息。

**特点：**
- 无状态：服务器不需要存储Session
- 自包含：Token中包含用户信息
- 可验证：使用签名防止篡改

### JWT结构

JWT由三部分组成，用 `.` 分隔：

```
Header.Payload.Signature
```

**示例：**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoi5byg5LiJIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

---

## 📝 核心代码：JWT实现

### 完整代码（手写练习）

```go
package jwt

import (
    commonpb "HeronGame/api/common/v1/pb"
    "HeronGame/common/log"
    gameErrors "HeronGame/game/errors"
    "context"
    "errors"
    "time"

    "go.uber.org/zap"
    "github.com/golang-jwt/jwt/v5"
)

// JWTClaims 自定义的 JWT Claims
type JWTClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar"`
    jwt.RegisteredClaims  // 标准Claims（过期时间等）
}

// JWTWebsocketClaims WebSocket专用Claims
type JWTWebsocketClaims struct {
    UserID   int64  `json:"user_id"`
    Nickname string `json:"nickname"`
    RoomID   int64  `json:"room_id"`
    GameID   int32  `json:"game_id"`
    jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(
    secret string,
    expire int64,
    signingMethod string,
    userID int64,
    username string,
    nickname string,
    avatar string,
) (string, error) {
    // 1. 选择签名算法
    var method jwt.SigningMethod
    switch signingMethod {
    case "HS256":
        method = jwt.SigningMethodHS256
    case "HS384":
        method = jwt.SigningMethodHS384
    case "HS512":
        method = jwt.SigningMethodHS512
    default:
        return "", errors.New("unsupported signing method")
    }

    // 2. 创建Claims
    claims := JWTClaims{
        UserID:   userID,
        Username: username,
        Nickname: nickname,
        Avatar:   avatar,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    // 3. 创建Token
    token := jwt.NewWithClaims(method, claims)
    
    // 4. 签名并生成Token字符串
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", gameErrors.New(int32(commonpb.ErrorCode_TOKEN_INVALID), tokenString)
    }

    return tokenString, nil
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string, secret string, signingMethod string) (*JWTClaims, error) {
    // 1. 选择签名算法
    var method jwt.SigningMethod
    switch signingMethod {
    case "HS256":
        method = jwt.SigningMethodHS256
    case "HS384":
        method = jwt.SigningMethodHS384
    case "HS512":
        method = jwt.SigningMethodHS512
    default:
        return nil, errors.New("unsupported signing method")
    }

    // 2. 解析Token
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        // 验证签名算法
        if token.Method != method {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secret), nil
    })

    // 3. 检查解析错误
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, gameErrors.New(int32(commonpb.ErrorCode_TOKEN_EXPIRED), tokenString)
        }
        return nil, gameErrors.New(int32(commonpb.ErrorCode_TOKEN_INVALID), tokenString)
    }

    // 4. 验证Token并提取Claims
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, gameErrors.New(int32(commonpb.ErrorCode_TOKEN_INVALID), tokenString)
}

// GenerateGameToken 生成游戏 JWT token（WebSocket用）
func GenerateGameToken(
    secret string,
    expire int64,
    signingMethod string,
    userID int64,
    nickname string,
    roomID int64,
    gameID int32,
) (string, error) {
    var method jwt.SigningMethod
    switch signingMethod {
    case "HS256":
        method = jwt.SigningMethodHS256
    case "HS384":
        method = jwt.SigningMethodHS384
    case "HS512":
        method = jwt.SigningMethodHS512
    default:
        return "", errors.New("unsupported signing method")
    }

    claims := JWTWebsocketClaims{
        UserID:   userID,
        Nickname: nickname,
        RoomID:   roomID,
        GameID:   gameID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(method, claims)
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        log.Error(context.Background(), "[JWT token生成] 解析token失败", zap.Any("err", err))
        return "", gameErrors.New(int32(commonpb.ErrorCode_TOKEN_INVALID), tokenString)
    }

    log.Info(context.Background(), "[JWT token生成成功] ",
        zap.Int64("user_id", userID),
        zap.Any("token", tokenString))
    
    return tokenString, nil
}

// ParseWebsocketToken 解析websocket token
func ParseWebsocketToken(tokenString string, secret string, signingMethod string) (*JWTWebsocketClaims, error) {
    var method jwt.SigningMethod
    switch signingMethod {
    case "HS256":
        method = jwt.SigningMethodHS256
    case "HS384":
        method = jwt.SigningMethodHS384
    case "HS512":
        method = jwt.SigningMethodHS512
    default:
        return nil, errors.New("unsupported signing method")
    }

    token, err := jwt.ParseWithClaims(tokenString, &JWTWebsocketClaims{}, func(token *jwt.Token) (interface{}, error) {
        if token.Method != method {
            return nil, gameErrors.New(int32(commonpb.ErrorCode_TOKEN_INVALID), tokenString)
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*JWTWebsocketClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, gameErrors.New(int32(commonpb.ErrorCode_TOKEN_INVALID), tokenString)
}
```

---

## 📖 JWT详解

### 1. JWT的三部分

#### Header（头部）

```json
{
  "alg": "HS256",  // 签名算法
  "typ": "JWT"     // Token类型
}
```

**Base64编码后：**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
```

#### Payload（载荷）

```json
{
  "user_id": 123,
  "username": "zhangsan",
  "nickname": "张三",
  "exp": 1704067200,  // 过期时间
  "iat": 1704063600   // 签发时间
}
```

**Base64编码后：**
```
eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoiemhhbmdzYW4ifQ
```

#### Signature（签名）

```
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret
)
```

**作用：**
- 验证Token是否被篡改
- 验证Token是否由服务器签发

### 2. 签名算法

**HS256（HMAC-SHA256）：**
- 使用对称密钥
- 速度快
- 适合单服务器场景

**RS256（RSA-SHA256）：**
- 使用非对称密钥
- 更安全
- 适合多服务器场景

### 3. Claims（声明）

**标准Claims：**
- `exp`：过期时间
- `iat`：签发时间
- `iss`：签发者
- `sub`：主题（用户ID）

**自定义Claims：**
- `user_id`：用户ID
- `username`：用户名
- `nickname`：昵称

### 4. Token生成流程

```
1. 创建Claims（包含用户信息）
   ↓
2. 选择签名算法（HS256/HS384/HS512）
   ↓
3. 使用密钥签名
   ↓
4. 生成Token字符串
```

### 5. Token解析流程

```
1. 分割Token（Header.Payload.Signature）
   ↓
2. 验证签名算法
   ↓
3. 使用密钥验证签名
   ↓
4. 检查过期时间
   ↓
5. 提取Claims
```

---

## 🔐 安全考虑

### 1. 密钥管理

```go
// ✅ 好的做法：从环境变量或配置中心获取
secret := os.Getenv("JWT_SECRET")
if secret == "" {
    panic("JWT_SECRET not set")
}

// ❌ 不好的做法：硬编码密钥
secret := "my-secret-key"  // 不安全！
```

### 2. Token过期时间

```go
// ✅ 好的做法：设置合理的过期时间
expire := 24 * 60 * 60  // 24小时

// ❌ 不好的做法：过期时间过长
expire := 365 * 24 * 60 * 60  // 1年，太长了！
```

### 3. HTTPS传输

```go
// ✅ 好的做法：使用HTTPS传输Token
// 防止Token被中间人攻击窃取

// ❌ 不好的做法：HTTP传输Token
// Token可能被窃取
```

### 4. Token刷新机制

```go
// 生成AccessToken（短期，15分钟）
accessToken, _ := GenerateToken(secret, 15*60, ...)

// 生成RefreshToken（长期，7天）
refreshToken, _ := GenerateToken(secret, 7*24*60*60, ...)
```

**刷新流程：**
```
1. 客户端使用AccessToken访问
   ↓
2. AccessToken过期
   ↓
3. 使用RefreshToken获取新的AccessToken
   ↓
4. 继续使用新的AccessToken
```

---

## 🎓 JWT vs Session

| 特性 | JWT | Session |
|------|-----|---------|
| 存储位置 | 客户端 | 服务器 |
| 扩展性 | 好（无状态） | 差（需要共享存储） |
| 性能 | 高（无需查询） | 低（需要查询） |
| 安全性 | 中（可被窃取） | 高（服务器存储） |
| 撤销 | 难（需黑名单） | 易（删除Session） |

---

## 💡 实践练习

### 练习1：理解JWT结构

使用在线工具（如 jwt.io）解析一个JWT Token：
1. 查看Header部分
2. 查看Payload部分
3. 验证Signature

### 练习2：手写Token生成

编写一个简单的Token生成函数：

```go
func GenerateSimpleToken(userID int64) (string, error) {
    // TODO: 实现简单的Token生成
}
```

### 练习3：Token刷新机制

设计一个Token刷新机制：
1. AccessToken过期时间：15分钟
2. RefreshToken过期时间：7天
3. 实现刷新逻辑

---

## 📌 下一阶段预告

**阶段9：WebSocket实现**
- WebSocket连接管理
- 消息处理流程
- 心跳机制
- 连接断开处理

---

## ❓ 思考题

1. JWT的Signature部分有什么作用？如果Signature被篡改会怎样？
2. 为什么JWT适合无状态认证？有什么缺点？
3. Token过期时间应该如何设置？过长或过短会有什么问题？
4. 如何实现Token的撤销机制？（提示：使用黑名单）

---

**完成本阶段后，请继续学习阶段9！** 🚀

