package router

import (
	"bytes"
	"io"
	"net/http"

	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// ResponseWriter 响应写入器（用于捕获响应内容）
type ResponseWriter struct {
	kratosHttp.Context
	statusCode int
	headers    http.Header
	body       *bytes.Buffer
	written    bool
}

// NewResponseWriter 创建响应写入器
func NewResponseWriter(ctx kratosHttp.Context) *ResponseWriter {
	return &ResponseWriter{
		Context:    ctx,
		statusCode: 200,
		headers:    make(http.Header),
		body:       &bytes.Buffer{},
		written:    false,
	}
}

// WriteHeader 写入状态码
func (w *ResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.statusCode = code
		w.Context.Response().WriteHeader(code)
		w.written = true
	}
}

// Write 写入响应体
func (w *ResponseWriter) Write(data []byte) (int, error) {
	if !w.written {
		w.WriteHeader(200)
	}

	// 同时写入到原始响应和缓冲区
	n1, err1 := w.Context.Response().Write(data)
	if err1 != nil {
		return n1, err1
	}

	n2, err2 := w.body.Write(data)
	if err2 != nil {
		return n2, err2
	}

	// 返回较小的写入数（确保一致性）
	if n1 < n2 {
		return n1, nil
	}
	return n2, nil
}

// Header 获取响应头
func (w *ResponseWriter) Header() http.Header {
	return w.Context.Response().Header()
}

// GetStatusCode 获取状态码
func (w *ResponseWriter) GetStatusCode() int {
	return w.statusCode
}

// GetBody 获取响应体
func (w *ResponseWriter) GetBody() []byte {
	return w.body.Bytes()
}

// GetHeaders 获取响应头
func (w *ResponseWriter) GetHeaders() http.Header {
	headers := make(http.Header)
	for key, values := range w.Context.Response().Header() {
		headers[key] = values
	}
	return headers
}

// CopyToOriginal 将响应复制到原始响应（用于缓存命中时）
func (w *ResponseWriter) CopyToOriginal(ctx kratosHttp.Context, statusCode int, headers http.Header, body []byte) error {
	// 复制响应头
	for key, values := range headers {
		for _, value := range values {
			ctx.Response().Header().Set(key, value)
		}
	}

	// 设置状态码
	ctx.Response().WriteHeader(statusCode)

	// 写入响应体
	_, err := ctx.Response().Write(body)
	return err
}

// TeeWriter 创建一个 TeeWriter，同时写入多个目标
type TeeWriter struct {
	writers []io.Writer
}

// NewTeeWriter 创建 TeeWriter
func NewTeeWriter(writers ...io.Writer) *TeeWriter {
	return &TeeWriter{
		writers: writers,
	}
}

// Write 写入数据到所有目标
func (t *TeeWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}
