package handler

import (
	"time"

	"StructForge/backend/common/log"

	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
)

// DashboardHandler 处理 Dashboard 相关请求
type DashboardHandler struct{}

// NewDashboardHandler 创建 Dashboard 处理器
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// GetStats 获取统计数据
func (h *DashboardHandler) GetStats(ctx kratosHttp.Context) error {
	log.Info(ctx.Request().Context(), "获取 Dashboard 统计数据")

	// 返回模拟数据
	stats := map[string]interface{}{
		"running":          0,
		"queued":           0,
		"totalExecutions":  0,
		"failedExecutions": 0,
		"failureRate":      0.0,
		"avgExecutionTime": 0.0,
		"totalWorkflows":   0,
		"activeWorkflows":  0,
	}

	return ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    stats,
	})
}

// GetExecutions 获取执行历史列表
func (h *DashboardHandler) GetExecutions(ctx kratosHttp.Context) error {
	log.Info(ctx.Request().Context(), "获取执行历史列表")

	// 获取查询参数
	page := 1
	pageSize := 20
	if pageStr := ctx.Request().URL.Query().Get("page"); pageStr != "" {
		if p := parseInt(pageStr); p > 0 {
			page = p
		}
	}
	if pageSizeStr := ctx.Request().URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps := parseInt(pageSizeStr); ps > 0 {
			pageSize = ps
		}
	}

	// 返回模拟数据
	response := map[string]interface{}{
		"list":     []interface{}{},
		"total":    0,
		"page":     page,
		"pageSize": pageSize,
	}

	return ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

// parseInt 简单的字符串转整数辅助函数
func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return result
}

// GetSuccessRate 获取执行成功率
func (h *DashboardHandler) GetSuccessRate(ctx kratosHttp.Context) error {
	log.Info(ctx.Request().Context(), "获取执行成功率")

	// 返回模拟数据
	response := []interface{}{}

	return ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

// GetErrorTrend 获取错误趋势
func (h *DashboardHandler) GetErrorTrend(ctx kratosHttp.Context) error {
	log.Info(ctx.Request().Context(), "获取错误趋势")

	// 返回模拟数据（最近24小时）
	now := time.Now()
	response := []map[string]interface{}{}
	for i := 23; i >= 0; i-- {
		t := now.Add(-time.Duration(i) * time.Hour)
		response = append(response, map[string]interface{}{
			"time":  t.Format("2006-01-02T15:04:05Z07:00"),
			"count": 0,
		})
	}

	return ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}

// GetExecutionDuration 获取平均执行时长
func (h *DashboardHandler) GetExecutionDuration(ctx kratosHttp.Context) error {
	log.Info(ctx.Request().Context(), "获取平均执行时长")

	// 返回模拟数据
	response := []interface{}{}

	return ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}
