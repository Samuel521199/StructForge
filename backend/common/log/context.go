package log

import (
	"context"
	"os"
)

// ContextKey 上下文键类型
type ContextKey string

const (
	CtxTraceID   ContextKey = "trace_id"
	CtxSpanID    ContextKey = "span_id"
	CtxUserID    ContextKey = "user_id"
	CtxRequestID ContextKey = "request_id"
)

// extractContextFields 从Context中提取字段
func extractContextFields(ctx context.Context) []Field {
	var fields []Field

	// 提取追踪信息
	if traceID := ctx.Value(CtxTraceID); traceID != nil {
		if id, ok := traceID.(string); ok && id != "" {
			fields = append(fields, String("trace_id", id))
		}
	}

	if spanID := ctx.Value(CtxSpanID); spanID != nil {
		if id, ok := spanID.(string); ok && id != "" {
			fields = append(fields, String("span_id", id))
		}
	}

	// 提取用户信息
	if userID := ctx.Value(CtxUserID); userID != nil {
		switch v := userID.(type) {
		case string:
			if v != "" {
				fields = append(fields, String("user_id", v))
			}
		case int64:
			fields = append(fields, Int64("user_id", v))
		case int:
			fields = append(fields, Int("user_id", v))
		}
	}

	// 提取请求ID
	if requestID := ctx.Value(CtxRequestID); requestID != nil {
		if id, ok := requestID.(string); ok && id != "" {
			fields = append(fields, String("request_id", id))
		}
	}

	return fields
}

// extractContainerInfo 从环境变量提取容器信息
func extractContainerInfo() []Field {
	var fields []Field

	if podName := os.Getenv("POD_NAME"); podName != "" {
		fields = append(fields, String("pod_name", podName))
	}

	if hostname := os.Getenv("HOSTNAME"); hostname != "" {
		fields = append(fields, String("hostname", hostname))
	}

	if namespace := os.Getenv("POD_NAMESPACE"); namespace != "" {
		fields = append(fields, String("namespace", namespace))
	}

	if podIP := os.Getenv("POD_IP"); podIP != "" {
		fields = append(fields, String("pod_ip", podIP))
	}

	if nodeName := os.Getenv("NODE_NAME"); nodeName != "" {
		fields = append(fields, String("node_name", nodeName))
	}

	return fields
}
