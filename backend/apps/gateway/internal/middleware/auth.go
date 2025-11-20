package middleware

import (
	"context"

	jwtMiddleware "StructForge/backend/apps/gateway/internal/middleware/jwt"

	"github.com/go-kratos/kratos/v2/middleware"
)

// UserIDKey Context 中存储用户ID的键
type userIDKey struct{}

// GetUserID 从 Context 中获取用户ID
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey{}).(int64)
	return userID, ok
}

// GetUsername 从 Context 中获取用户名
type usernameKey struct{}

func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(usernameKey{}).(string)
	return username, ok
}

// AuthMiddleware JWT 认证中间件（用于 Kratos HTTP 中间件链）
func AuthMiddleware(jwtManager *jwtMiddleware.Manager) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 注意：Kratos HTTP 中间件中，req 是 HTTP 请求对象
			// 这里暂时保留接口，实际认证逻辑在 handler 中实现
			return handler(ctx, req)
		}
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
