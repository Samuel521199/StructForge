package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestCircuitBreakerStates 测试熔断器状态转换
func TestCircuitBreakerStates(t *testing.T) {
	cb := NewCircuitBreaker(&Config{
		FailureThreshold: 0.5, // 50% 失败率
		MinRequests:      10,  // 至少10个请求
		WindowSize:       60,  // 60秒窗口
		OpenDuration:     30,  // 打开状态持续30秒
		HalfOpenRequests: 3,   // 半开状态允许3个请求
		Timeout:          5,   // 5秒超时
	})

	// 初始状态应该是 Closed
	if cb.GetState() != StateClosed {
		t.Errorf("初始状态应该是 Closed，实际 %s", cb.GetState())
	}

	// 执行一些请求，但未达到最小请求数
	for i := 0; i < 5; i++ {
		err := cb.Execute(context.Background(), func() error {
			return errors.New("test error")
		})
		if err == nil {
			t.Error("应该返回错误")
		}
	}

	// 状态应该仍然是 Closed（未达到最小请求数）
	if cb.GetState() != StateClosed {
		t.Errorf("状态应该是 Closed，实际 %s", cb.GetState())
	}
}

// TestCircuitBreakerOpen 测试熔断器打开
func TestCircuitBreakerOpen(t *testing.T) {
	cb := NewCircuitBreaker(&Config{
		FailureThreshold: 0.5,
		MinRequests:      10,
		WindowSize:       60,
		OpenDuration:     1, // 1秒打开持续时间（测试用）
		HalfOpenRequests: 3,
		Timeout:          5,
	})

	// 执行足够的请求，其中一半失败
	for i := 0; i < 10; i++ {
		cb.Execute(context.Background(), func() error {
			if i < 5 {
				return errors.New("test error")
			}
			return nil
		})
	}

	// 等待状态更新
	time.Sleep(100 * time.Millisecond)

	// 熔断器应该打开
	if cb.GetState() != StateOpen {
		t.Errorf("熔断器应该打开，实际状态 %s", cb.GetState())
	}

	// 打开状态下，请求应该立即失败
	err := cb.Execute(context.Background(), func() error {
		return nil
	})
	if err == nil || !IsCircuitBreakerError(err) {
		t.Error("打开状态下应该返回熔断器错误")
	}
}

// TestCircuitBreakerHalfOpen 测试熔断器半开状态
func TestCircuitBreakerHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(&Config{
		FailureThreshold: 0.5,
		MinRequests:      10,
		WindowSize:       60,
		OpenDuration:     1, // 1秒后进入半开
		HalfOpenRequests: 3,
		Timeout:          5,
	})

	// 先触发打开状态
	for i := 0; i < 10; i++ {
		cb.Execute(context.Background(), func() error {
			return errors.New("test error")
		})
	}

	// 等待进入半开状态
	time.Sleep(1100 * time.Millisecond)

	// 状态应该是 HalfOpen
	if cb.GetState() != StateHalfOpen {
		t.Errorf("状态应该是 HalfOpen，实际 %s", cb.GetState())
	}

	// 半开状态下，允许少量请求通过
	successCount := 0
	for i := 0; i < 3; i++ {
		err := cb.Execute(context.Background(), func() error {
			return nil // 成功
		})
		if err == nil {
			successCount++
		}
	}

	if successCount == 0 {
		t.Error("半开状态下应该允许请求通过")
	}
}

// TestCircuitBreakerRecovery 测试熔断器恢复
func TestCircuitBreakerRecovery(t *testing.T) {
	cb := NewCircuitBreaker(&Config{
		FailureThreshold: 0.5,
		MinRequests:      10,
		WindowSize:       60,
		OpenDuration:     1,
		HalfOpenRequests: 3,
		Timeout:          5,
	})

	// 触发打开
	for i := 0; i < 10; i++ {
		cb.Execute(context.Background(), func() error {
			return errors.New("test error")
		})
	}

	// 等待进入半开
	time.Sleep(1100 * time.Millisecond)

	// 半开状态下，所有请求都成功，应该恢复到关闭状态
	for i := 0; i < 3; i++ {
		cb.Execute(context.Background(), func() error {
			return nil
		})
	}

	// 等待状态更新
	time.Sleep(100 * time.Millisecond)

	// 应该恢复到 Closed 状态
	if cb.GetState() != StateClosed {
		t.Errorf("应该恢复到 Closed 状态，实际 %s", cb.GetState())
	}
}

// TestCircuitBreakerManager 测试熔断器管理器
func TestCircuitBreakerManager(t *testing.T) {
	mgr := NewCircuitBreakerManager()

	config := &Config{
		FailureThreshold: 0.5,
		MinRequests:      10,
		WindowSize:       60,
		OpenDuration:     30,
		HalfOpenRequests: 3,
		Timeout:          5,
	}

	// 为不同服务创建熔断器
	service1 := "user-service"
	service2 := "order-service"

	// 执行请求
	err1 := mgr.Execute(context.Background(), service1, config, func() error {
		return errors.New("test error")
	})
	if err1 == nil {
		t.Error("应该返回错误")
	}

	err2 := mgr.Execute(context.Background(), service2, config, func() error {
		return nil
	})
	if err2 != nil {
		t.Error("不应该返回错误")
	}

	// 获取统计信息
	stats := mgr.GetBreakerStats()
	if len(stats) != 2 {
		t.Errorf("应该有2个服务的统计信息，实际 %d", len(stats))
	}
}
