package logging

import (
	"bytes"
	"io"
	"time"

	"StructForge/backend/common/log"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// RequestLogger 请求日志记录器
type RequestLogger struct {
	// 是否记录请求体
	LogRequestBody bool
	// 是否记录响应体
	LogResponseBody bool
	// 请求体最大大小（字节）
	MaxRequestBodySize int
	// 响应体最大大小（字节）
	MaxResponseBodySize int
}

// NewRequestLogger 创建请求日志记录器
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{
		LogRequestBody:      true,
		LogResponseBody:     true,
		MaxRequestBodySize:  1024 * 10, // 10KB
		MaxResponseBodySize: 1024 * 10, // 10KB
	}
}

// LogRequest 记录请求信息
func (l *RequestLogger) LogRequest(ctx http.Context, startTime time.Time) {
	req := ctx.Request()
	duration := time.Since(startTime)

	// 提取请求信息
	method := req.Method
	path := req.URL.Path
	query := req.URL.RawQuery
	remoteAddr := req.RemoteAddr
	userAgent := req.UserAgent()

	// 记录基本信息（不记录状态码，因为此时响应可能还未完成）
	log.Info(ctx, "HTTP请求",
		log.String("method", method),
		log.String("path", path),
		log.String("query", query),
		log.String("remote_addr", remoteAddr),
		log.String("user_agent", userAgent),
		log.Duration("duration", duration),
	)

	// 记录请求头（可选，避免日志过多）
	// log.Info(ctx, "请求头", log.Any("headers", req.Header))

	// 记录请求体（如果启用且大小合适）
	if l.LogRequestBody && req.Body != nil {
		bodyBytes, err := l.readRequestBody(req)
		if err == nil && len(bodyBytes) > 0 {
			bodyStr := string(bodyBytes)
			if len(bodyStr) > l.MaxRequestBodySize {
				bodyStr = bodyStr[:l.MaxRequestBodySize] + "...(truncated)"
			}
			log.Info(ctx, "请求体",
				log.String("body", bodyStr),
			)
		}
	}
}

// LogResponse 记录响应信息
func (l *RequestLogger) LogResponse(ctx http.Context, startTime time.Time, responseBody []byte) {
	duration := time.Since(startTime)

	// 记录响应信息（状态码需要在调用时传入）
	log.Info(ctx, "HTTP响应",
		log.Duration("duration", duration),
	)

	// 记录响应体（如果启用且大小合适）
	if l.LogResponseBody && len(responseBody) > 0 {
		bodyStr := string(responseBody)
		if len(bodyStr) > l.MaxResponseBodySize {
			bodyStr = bodyStr[:l.MaxResponseBodySize] + "...(truncated)"
		}
		log.Info(ctx, "响应体",
			log.String("body", bodyStr),
		)
	}
}

// readRequestBody 读取请求体（不消耗原始流）
func (l *RequestLogger) readRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}

	// 读取请求体
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// 重新创建请求体，以便后续处理
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}

// LogError 记录错误信息
func (l *RequestLogger) LogError(ctx http.Context, err error, startTime time.Time) {
	duration := time.Since(startTime)

	log.Error(ctx, "请求处理失败",
		log.ErrorField(err),
		log.Duration("duration", duration),
		log.String("path", ctx.Request().URL.Path),
		log.String("method", ctx.Request().Method),
	)
}
