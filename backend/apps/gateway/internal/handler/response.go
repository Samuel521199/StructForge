package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	stdHttp "net/http"
	"strings"
	"time"

	"StructForge/backend/common/log"
)

// StandardResponse 标准响应结构
type StandardResponse struct {
	Code      int         `json:"code"`                // 状态码
	Message   string      `json:"message"`             // 消息
	Data      interface{} `json:"data,omitempty"`      // 数据（可选）
	Error     string      `json:"error,omitempty"`     // 错误详情（可选）
	TraceID   string      `json:"trace_id,omitempty"`  // 追踪ID（可选）
	Timestamp string      `json:"timestamp,omitempty"` // 时间戳（可选）
}

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeNetwork      ErrorType = "network"       // 网络错误
	ErrorTypeTimeout      ErrorType = "timeout"       // 超时错误
	ErrorTypeBusiness     ErrorType = "business"      // 业务错误
	ErrorTypeConfig       ErrorType = "config"        // 配置错误
	ErrorTypeAuth         ErrorType = "auth"          // 认证错误
	ErrorTypeRateLimit    ErrorType = "rate_limit"    // 限流错误
	ErrorTypeCircuitBreak ErrorType = "circuit_break" // 熔断错误
	ErrorTypeNotFound     ErrorType = "not_found"     // 未找到错误
	ErrorTypeInternal     ErrorType = "internal"      // 内部错误
)

// ErrorCode 错误码定义
const (
	// 通用错误码 (1000-1999)
	CodeSuccess            = 200
	CodeBadRequest         = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodeMethodNotAllowed   = 405
	CodeConflict           = 409
	CodeRateLimit          = 429
	CodeInternalError      = 500
	CodeBadGateway         = 502
	CodeServiceUnavailable = 503
	CodeGatewayTimeout     = 504

	// Gateway 特定错误码 (2000-2999)
	CodeRouteNotFound      = 2001 // 路由不存在
	CodeNoServiceInstance  = 2002 // 没有可用服务实例
	CodeCircuitBreakerOpen = 2003 // 熔断器已打开
	CodeRequestTimeout     = 2004 // 请求超时
	CodeDownstreamError    = 2005 // 下游服务错误
	CodeInvalidAuthFormat  = 2006 // 无效的认证格式
	CodeInvalidToken       = 2007 // 无效或过期的令牌
	CodeCacheError         = 2008 // 缓存错误
	CodeConfigError        = 2009 // 配置错误
)

// generateTraceID 生成追踪ID
func generateTraceID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// getTraceIDFromContext 从 Context 中获取 TraceID
func getTraceIDFromContext(ctx context.Context) string {
	if traceID := ctx.Value(log.CtxTraceID); traceID != nil {
		if id, ok := traceID.(string); ok && id != "" {
			return id
		}
	}
	return ""
}

// getTraceIDFromRequest 从请求头中获取 TraceID
func getTraceIDFromRequest(req *stdHttp.Request) string {
	// 优先从 X-Trace-ID 获取
	if traceID := req.Header.Get("X-Trace-ID"); traceID != "" {
		return traceID
	}
	// 其次从 X-Request-ID 获取
	if requestID := req.Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}
	// 最后从 Traceparent (W3C Trace Context) 获取
	if traceparent := req.Header.Get("Traceparent"); traceparent != "" {
		// 简单提取：Traceparent 格式为 "version-trace_id-parent_id-flags"
		parts := strings.Split(traceparent, "-")
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

// SuccessResponse 成功响应
func SuccessResponse(ctx context.Context, data interface{}) *StandardResponse {
	resp := &StandardResponse{
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 从 Context 中提取 TraceID
	if traceID := getTraceIDFromContext(ctx); traceID != "" {
		resp.TraceID = traceID
	}

	return resp
}

// ErrorResponse 错误响应（带 TraceID 和错误类型）
func ErrorResponse(ctx context.Context, code int, message string, err error, errorType ErrorType) *StandardResponse {
	resp := &StandardResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 从 Context 中提取 TraceID
	if traceID := getTraceIDFromContext(ctx); traceID != "" {
		resp.TraceID = traceID
	}

	// 添加错误详情
	if err != nil {
		// 根据错误类型格式化错误信息
		resp.Error = formatError(err, errorType)
	}

	return resp
}

// formatError 格式化错误信息
func formatError(err error, errorType ErrorType) string {
	if err == nil {
		return ""
	}

	errMsg := err.Error()

	// 根据错误类型添加前缀
	switch errorType {
	case ErrorTypeNetwork:
		return fmt.Sprintf("网络错误: %s", errMsg)
	case ErrorTypeTimeout:
		return fmt.Sprintf("请求超时: %s", errMsg)
	case ErrorTypeCircuitBreak:
		return fmt.Sprintf("服务熔断: %s", errMsg)
	case ErrorTypeBusiness:
		return fmt.Sprintf("业务错误: %s", errMsg)
	case ErrorTypeConfig:
		return fmt.Sprintf("配置错误: %s", errMsg)
	case ErrorTypeAuth:
		return fmt.Sprintf("认证错误: %s", errMsg)
	case ErrorTypeRateLimit:
		return fmt.Sprintf("限流错误: %s", errMsg)
	case ErrorTypeNotFound:
		return fmt.Sprintf("资源不存在: %s", errMsg)
	default:
		return errMsg
	}
}

// classifyError 分类错误类型
func classifyError(err error) ErrorType {
	if err == nil {
		return ErrorTypeInternal
	}

	errMsg := strings.ToLower(err.Error())

	// 网络错误
	if strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "network") ||
		strings.Contains(errMsg, "refused") ||
		strings.Contains(errMsg, "reset") ||
		strings.Contains(errMsg, "unreachable") {
		return ErrorTypeNetwork
	}

	// 超时错误
	if strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "deadline exceeded") ||
		strings.Contains(errMsg, "context deadline") {
		return ErrorTypeTimeout
	}

	// 熔断错误
	if strings.Contains(errMsg, "circuit breaker") ||
		strings.Contains(errMsg, "熔断器") {
		return ErrorTypeCircuitBreak
	}

	// 认证错误
	if strings.Contains(errMsg, "token") ||
		strings.Contains(errMsg, "auth") ||
		strings.Contains(errMsg, "unauthorized") ||
		strings.Contains(errMsg, "认证") {
		return ErrorTypeAuth
	}

	// 限流错误
	if strings.Contains(errMsg, "rate limit") ||
		strings.Contains(errMsg, "too many") ||
		strings.Contains(errMsg, "限流") {
		return ErrorTypeRateLimit
	}

	// 未找到错误
	if strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "不存在") {
		return ErrorTypeNotFound
	}

	// 客户端错误（4xx）
	if strings.Contains(errMsg, "client error") {
		return ErrorTypeBusiness
	}

	// 服务器错误（5xx）
	if strings.Contains(errMsg, "server error") {
		return ErrorTypeInternal
	}

	return ErrorTypeInternal
}

// getErrorCode 根据错误类型和 HTTP 状态码获取业务错误码
func getErrorCode(httpCode int, errorType ErrorType) int {
	switch errorType {
	case ErrorTypeNotFound:
		return CodeRouteNotFound
	case ErrorTypeCircuitBreak:
		return CodeCircuitBreakerOpen
	case ErrorTypeTimeout:
		return CodeRequestTimeout
	case ErrorTypeNetwork:
		return CodeBadGateway
	case ErrorTypeAuth:
		if httpCode == 401 {
			return CodeUnauthorized
		}
		return CodeInvalidToken
	case ErrorTypeRateLimit:
		return CodeRateLimit
	case ErrorTypeBusiness:
		return CodeDownstreamError
	default:
		return httpCode
	}
}

// NewErrorResponse 创建错误响应（自动分类错误）
func NewErrorResponse(ctx context.Context, httpCode int, message string, err error) *StandardResponse {
	errorType := classifyError(err)
	code := getErrorCode(httpCode, errorType)
	return ErrorResponse(ctx, code, message, err, errorType)
}

// 预定义的错误响应（使用函数，支持 TraceID）
func ErrNotFound(ctx context.Context) *StandardResponse {
	return ErrorResponse(ctx, CodeRouteNotFound, "路由不存在", nil, ErrorTypeNotFound)
}

func ErrUnauthorized(ctx context.Context) *StandardResponse {
	return ErrorResponse(ctx, CodeUnauthorized, "需要认证", nil, ErrorTypeAuth)
}

func ErrInvalidAuth(ctx context.Context) *StandardResponse {
	return ErrorResponse(ctx, CodeInvalidAuthFormat, "无效的认证格式", nil, ErrorTypeAuth)
}

func ErrInvalidToken(ctx context.Context, err error) *StandardResponse {
	return ErrorResponse(ctx, CodeInvalidToken, "无效或过期的令牌", err, ErrorTypeAuth)
}

func ErrRateLimit(ctx context.Context) *StandardResponse {
	return ErrorResponse(ctx, CodeRateLimit, "请求过于频繁，请稍后再试", nil, ErrorTypeRateLimit)
}

func ErrServiceUnavailable(ctx context.Context, err error) *StandardResponse {
	return ErrorResponse(ctx, CodeServiceUnavailable, "服务暂时不可用", err, ErrorTypeInternal)
}

func ErrCircuitBreakerOpen(ctx context.Context, service string) *StandardResponse {
	err := fmt.Errorf("服务 %s 的熔断器已打开", service)
	return ErrorResponse(ctx, CodeCircuitBreakerOpen, "服务暂时不可用（熔断器已打开）", err, ErrorTypeCircuitBreak)
}

func ErrNoServiceInstance(ctx context.Context, service string) *StandardResponse {
	err := fmt.Errorf("服务 %s 没有可用实例", service)
	return ErrorResponse(ctx, CodeNoServiceInstance, "服务暂时不可用", err, ErrorTypeInternal)
}

func ErrRequestTimeout(ctx context.Context, timeout time.Duration) *StandardResponse {
	err := fmt.Errorf("请求超时（%v）", timeout)
	return ErrorResponse(ctx, CodeRequestTimeout, "请求超时", err, ErrorTypeTimeout)
}

func ErrDownstreamError(ctx context.Context, service string, statusCode int, err error) *StandardResponse {
	msg := fmt.Sprintf("下游服务 %s 返回错误（状态码: %d）", service, statusCode)
	return ErrorResponse(ctx, CodeDownstreamError, msg, err, ErrorTypeBusiness)
}

// isNetworkError 判断是否为网络错误
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否为网络错误类型
	if _, ok := err.(net.Error); ok {
		return true
	}

	errMsg := strings.ToLower(err.Error())
	networkErrors := []string{
		"connection refused",
		"connection reset",
		"no such host",
		"network is unreachable",
		"timeout",
		"deadline exceeded",
	}

	for _, ne := range networkErrors {
		if strings.Contains(errMsg, ne) {
			return true
		}
	}

	return false
}

// isTimeoutError 判断是否为超时错误
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "deadline exceeded") ||
		strings.Contains(errMsg, "context deadline")
}
