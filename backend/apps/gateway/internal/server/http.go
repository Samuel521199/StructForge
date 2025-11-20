package server

import (
	"net/http"

	"StructForge/backend/apps/gateway/internal/conf"
	corsMiddleware "StructForge/backend/apps/gateway/internal/middleware/cors"
	"StructForge/backend/common/log"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer 创建HTTP服务器
// 使用 HTTP Filter 在路由之前处理 OPTIONS 请求
func NewHTTPServer(c *conf.Bootstrap, corsHandler *corsMiddleware.CORSHandler) *kratosHttp.Server {
	var opts = []kratosHttp.ServerOption{}

	// 恢复中间件
	opts = append(opts, kratosHttp.Middleware(recovery.Recovery()))

	// 使用 HTTP Filter 在路由之前处理 OPTIONS 请求和 CORS
	// 这样可以确保所有 OPTIONS 请求都能被捕获，即使路由没有匹配
	if corsHandler != nil {
		opts = append(opts, kratosHttp.Filter(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				method := r.Method
				path := r.URL.Path
				origin := r.Header.Get("Origin")

				// 添加调试日志
				log.Info(ctx, "HTTP Filter 收到请求",
					log.String("method", method),
					log.String("path", path),
					log.String("origin", origin),
				)

				// 如果是 OPTIONS 请求，直接处理 CORS
				if method == "OPTIONS" {
					log.Info(ctx, "Filter 处理 OPTIONS 预检请求",
						log.String("path", path),
						log.String("origin", origin),
					)
					// 直接处理 CORS，不通过 kratosHttp.Context
					if origin != "" {
						// 检查源是否被允许
						allowed := corsHandler.IsOriginAllowed(origin)
						log.Info(ctx, "源检查结果",
							log.String("origin", origin),
							log.Bool("allowed", allowed),
						)
						if allowed {
							// 设置 CORS 响应头
							allowedOrigin := corsHandler.GetAllowedOrigin(origin)
							w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
							if corsHandler.AllowCredentials() {
								w.Header().Set("Access-Control-Allow-Credentials", "true")
							}
							// 设置允许的方法
							w.Header().Set("Access-Control-Allow-Methods", corsHandler.GetAllowedMethods())
							// 设置允许的请求头
							requestedHeaders := r.Header.Get("Access-Control-Request-Headers")
							if requestedHeaders != "" && corsHandler.IsHeaderAllowed(requestedHeaders) {
								w.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
							} else {
								w.Header().Set("Access-Control-Allow-Headers", corsHandler.GetAllowedHeaders())
							}
							// 设置预检请求缓存时间
							if corsHandler.GetMaxAge() > 0 {
								w.Header().Set("Access-Control-Max-Age", corsHandler.GetMaxAgeString())
							}
							log.Info(ctx, "OPTIONS 请求处理完成，返回 204",
								log.String("Access-Control-Allow-Origin", w.Header().Get("Access-Control-Allow-Origin")),
							)
							// 返回 204 No Content
							w.WriteHeader(204)
							return
						} else {
							// 源不被允许，返回 403
							log.Warn(ctx, "CORS 预检请求被拒绝：源不被允许",
								log.String("origin", origin),
							)
							w.WriteHeader(403)
							return
						}
					}
					// 没有 Origin 头，返回 204
					log.Info(ctx, "OPTIONS 请求没有 Origin 头，返回 204")
					w.WriteHeader(204)
					return
				}

				// 对于非 OPTIONS 请求，也设置 CORS 响应头（如果存在 Origin）
				if origin != "" {
					allowed := corsHandler.IsOriginAllowed(origin)
					if allowed {
						allowedOrigin := corsHandler.GetAllowedOrigin(origin)
						w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
						if corsHandler.AllowCredentials() {
							w.Header().Set("Access-Control-Allow-Credentials", "true")
						}
					}
				}

				// 继续处理其他请求
				log.Info(ctx, "Filter 继续处理请求",
					log.String("method", method),
					log.String("path", path),
				)
				next.ServeHTTP(w, r)
			})
		}))
	}

	// 设置默认地址（如果配置中没有指定）
	addr := ":8000"
	if c.Server != nil && c.Server.Id != "" {
		// 可以根据服务ID设置不同的端口
		addr = ":8000"
	}

	opts = append(opts, kratosHttp.Address(addr))

	return kratosHttp.NewServer(opts...)
}
