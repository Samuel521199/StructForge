package cors

import (
	"strconv"
	"strings"

	"StructForge/backend/common/log"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// CORSOptions CORS配置选项
type CORSOptions struct {
	// 允许的源（支持通配符 *）
	AllowedOrigins []string
	// 允许的方法
	AllowedMethods []string
	// 允许的请求头
	AllowedHeaders []string
	// 暴露的响应头
	ExposedHeaders []string
	// 是否允许携带凭证
	AllowCredentials bool
	// 预检请求的缓存时间（秒）
	MaxAge int
}

// DefaultCORSOptions 默认CORS配置
func DefaultCORSOptions() *CORSOptions {
	return &CORSOptions{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           86400, // 24小时
	}
}

// CORSHandler CORS处理器
type CORSHandler struct {
	options *CORSOptions
}

// IsOriginAllowed 检查源是否被允许（公开方法，供 Filter 使用）
func (h *CORSHandler) IsOriginAllowed(origin string) bool {
	return h.isOriginAllowed(origin)
}

// GetAllowedOrigin 获取允许的源（公开方法，供 Filter 使用）
func (h *CORSHandler) GetAllowedOrigin(origin string) string {
	return h.getAllowedOrigin(origin)
}

// AllowCredentials 是否允许携带凭证（公开方法，供 Filter 使用）
func (h *CORSHandler) AllowCredentials() bool {
	return h.options.AllowCredentials
}

// GetAllowedMethods 获取允许的方法（公开方法，供 Filter 使用）
func (h *CORSHandler) GetAllowedMethods() string {
	if len(h.options.AllowedMethods) > 0 {
		return strings.Join(h.options.AllowedMethods, ", ")
	}
	return "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD"
}

// GetAllowedHeaders 获取允许的请求头（公开方法，供 Filter 使用）
func (h *CORSHandler) GetAllowedHeaders() string {
	if len(h.options.AllowedHeaders) > 0 {
		return strings.Join(h.options.AllowedHeaders, ", ")
	}
	return "*"
}

// IsHeaderAllowed 检查请求头是否被允许（公开方法，供 Filter 使用）
func (h *CORSHandler) IsHeaderAllowed(requestedHeaders string) bool {
	return h.isHeaderAllowed(requestedHeaders)
}

// GetMaxAge 获取预检请求缓存时间（公开方法，供 Filter 使用）
func (h *CORSHandler) GetMaxAge() int {
	return h.options.MaxAge
}

// GetMaxAgeString 获取预检请求缓存时间（字符串格式，供 Filter 使用）
func (h *CORSHandler) GetMaxAgeString() string {
	return strconv.Itoa(h.options.MaxAge)
}

// NewCORSHandler 创建CORS处理器
func NewCORSHandler(options *CORSOptions) *CORSHandler {
	if options == nil {
		options = DefaultCORSOptions()
	}
	return &CORSHandler{
		options: options,
	}
}

// HandleCORS 处理CORS请求
func (h *CORSHandler) HandleCORS(ctx http.Context) error {
	req := ctx.Request()
	resp := ctx.Response()

	origin := req.Header.Get("Origin")

	// 添加调试日志
	log.Info(req.Context(), "HandleCORS 被调用",
		log.String("method", req.Method),
		log.String("path", req.URL.Path),
		log.String("origin", origin),
		log.String("allowed_origins", strings.Join(h.options.AllowedOrigins, ", ")),
	)

	if origin == "" {
		// 没有Origin头，不是CORS请求
		log.Debug(req.Context(), "没有 Origin 头，跳过 CORS 处理")
		return nil
	}

	// 检查是否允许该源
	allowed := h.isOriginAllowed(origin)
	log.Info(req.Context(), "源检查结果",
		log.String("origin", origin),
		log.Bool("allowed", allowed),
	)
	if !allowed {
		// 对于 OPTIONS 预检请求，如果源不被允许，返回 403
		if req.Method == "OPTIONS" {
			log.Warn(req.Context(), "CORS 预检请求被拒绝：源不被允许",
				log.String("origin", origin),
				log.String("allowed_origins", strings.Join(h.options.AllowedOrigins, ", ")),
			)
			resp.WriteHeader(403)
			return nil
		}
		// 对于普通请求，不设置 CORS 头，让浏览器处理
		log.Warn(req.Context(), "CORS 请求被拒绝：源不被允许",
			log.String("origin", origin),
			log.String("allowed_origins", strings.Join(h.options.AllowedOrigins, ", ")),
		)
		return nil
	}

	// 设置CORS响应头（必须在处理 OPTIONS 之前设置）
	allowedOrigin := h.getAllowedOrigin(origin)
	resp.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	log.Info(req.Context(), "设置 Access-Control-Allow-Origin",
		log.String("origin", allowedOrigin),
	)

	if h.options.AllowCredentials {
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		log.Debug(req.Context(), "设置 Access-Control-Allow-Credentials", log.String("value", "true"))
	}

	// 处理预检请求（OPTIONS）
	if req.Method == "OPTIONS" {
		log.Info(req.Context(), "处理 OPTIONS 预检请求，设置 CORS 响应头")

		// 允许的方法
		if len(h.options.AllowedMethods) > 0 {
			methods := strings.Join(h.options.AllowedMethods, ", ")
			resp.Header().Set("Access-Control-Allow-Methods", methods)
			log.Debug(req.Context(), "设置 Access-Control-Allow-Methods",
				log.String("methods", methods),
			)
		}

		// 允许的请求头
		requestedHeaders := req.Header.Get("Access-Control-Request-Headers")
		if requestedHeaders != "" {
			if h.isHeaderAllowed(requestedHeaders) {
				resp.Header().Set("Access-Control-Allow-Headers", requestedHeaders)
				log.Debug(req.Context(), "设置 Access-Control-Allow-Headers（来自请求）",
					log.String("headers", requestedHeaders),
				)
			} else if len(h.options.AllowedHeaders) > 0 {
				headers := strings.Join(h.options.AllowedHeaders, ", ")
				resp.Header().Set("Access-Control-Allow-Headers", headers)
				log.Debug(req.Context(), "设置 Access-Control-Allow-Headers（来自配置）",
					log.String("headers", headers),
				)
			}
		} else if len(h.options.AllowedHeaders) > 0 {
			headers := strings.Join(h.options.AllowedHeaders, ", ")
			resp.Header().Set("Access-Control-Allow-Headers", headers)
			log.Debug(req.Context(), "设置 Access-Control-Allow-Headers（默认）",
				log.String("headers", headers),
			)
		}

		// 暴露的响应头
		if len(h.options.ExposedHeaders) > 0 {
			exposedHeaders := strings.Join(h.options.ExposedHeaders, ", ")
			resp.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
		}

		// 预检请求缓存时间
		if h.options.MaxAge > 0 {
			resp.Header().Set("Access-Control-Max-Age", strconv.Itoa(h.options.MaxAge))
		}

		// 记录所有设置的响应头
		log.Info(req.Context(), "OPTIONS 预检请求处理完成，准备返回 204",
			log.String("Access-Control-Allow-Origin", resp.Header().Get("Access-Control-Allow-Origin")),
			log.String("Access-Control-Allow-Methods", resp.Header().Get("Access-Control-Allow-Methods")),
			log.String("Access-Control-Allow-Headers", resp.Header().Get("Access-Control-Allow-Headers")),
		)

		// 直接返回，不继续处理
		resp.WriteHeader(204)
		return nil
	}

	// 普通请求，设置暴露的响应头
	if len(h.options.ExposedHeaders) > 0 {
		exposedHeaders := strings.Join(h.options.ExposedHeaders, ", ")
		resp.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
	}

	return nil
}

// isOriginAllowed 检查源是否被允许
func (h *CORSHandler) isOriginAllowed(origin string) bool {
	if len(h.options.AllowedOrigins) == 0 {
		return false
	}

	for _, allowedOrigin := range h.options.AllowedOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
		// 支持通配符匹配（如 *.example.com）
		if strings.Contains(allowedOrigin, "*") {
			if h.matchWildcard(origin, allowedOrigin) {
				return true
			}
		}
	}

	return false
}

// getAllowedOrigin 获取允许的源（处理通配符）
func (h *CORSHandler) getAllowedOrigin(origin string) string {
	for _, allowedOrigin := range h.options.AllowedOrigins {
		if allowedOrigin == "*" {
			if h.options.AllowCredentials {
				// 如果允许凭证，不能使用通配符，必须返回具体源
				return origin
			}
			return "*"
		}
		if allowedOrigin == origin {
			return origin
		}
		// 支持通配符匹配
		if strings.Contains(allowedOrigin, "*") {
			if h.matchWildcard(origin, allowedOrigin) {
				return origin
			}
		}
	}

	return origin
}

// isHeaderAllowed 检查请求头是否被允许
func (h *CORSHandler) isHeaderAllowed(requestedHeaders string) bool {
	if len(h.options.AllowedHeaders) == 0 {
		return false
	}

	headers := strings.Split(requestedHeaders, ",")
	for _, header := range headers {
		header = strings.TrimSpace(header)
		allowed := false

		for _, allowedHeader := range h.options.AllowedHeaders {
			if allowedHeader == "*" {
				allowed = true
				break
			}
			if strings.EqualFold(header, allowedHeader) {
				allowed = true
				break
			}
		}

		if !allowed {
			return false
		}
	}

	return true
}

// matchWildcard 通配符匹配
func (h *CORSHandler) matchWildcard(origin, pattern string) bool {
	// 简单的通配符匹配：*.example.com
	if strings.HasPrefix(pattern, "*.") {
		domain := strings.TrimPrefix(pattern, "*.")
		return strings.HasSuffix(origin, domain)
	}

	return false
}

// 注意：CORS 处理现在在 handler 层面完成（gateway.go 的 Proxy 方法）
// 这里保留 CORSHandlerMiddleware 函数签名以保持兼容性，但实际不使用
// 如果需要中间件层面的 CORS 处理，需要使用 Kratos 的 HTTP Filter 机制
