package router

import (
	"testing"

	"StructForge/backend/apps/gateway/internal/router/discovery"
	"StructForge/backend/apps/gateway/internal/router/loadbalancer"
)

// TestRouteMatchExact 测试精确匹配
func TestRouteMatchExact(t *testing.T) {
	discovery := discovery.NewStaticDiscovery()
	router := NewRouter(discovery)
	route := &Route{
		Path:      "/api/v1/users",
		MatchType: "exact",
		Service:   "user-service",
	}
	router.AddRoute(route)

	// 测试精确匹配
	found := router.FindRoute("/api/v1/users")
	if found == nil {
		t.Fatal("应该找到路由")
	}
	if found.Path != "/api/v1/users" {
		t.Errorf("期望路径 /api/v1/users，实际 %s", found.Path)
	}

	// 测试不匹配
	notFound := router.FindRoute("/api/v1/users/123")
	if notFound != nil {
		t.Error("不应该找到路由")
	}
}

// TestRouteMatchPrefix 测试前缀匹配
func TestRouteMatchPrefix(t *testing.T) {
	discovery := discovery.NewStaticDiscovery()
	router := NewRouter(discovery)
	route := &Route{
		Path:      "/api/v1/users",
		MatchType: "prefix",
		Service:   "user-service",
	}
	router.AddRoute(route)

	// 测试前缀匹配
	testCases := []struct {
		path     string
		expected bool
	}{
		{"/api/v1/users", true},
		{"/api/v1/users/123", true},
		{"/api/v1/users/123/profile", true},
		{"/api/v1/user", false},
		{"/api/v1/userss", false},
	}

	for _, tc := range testCases {
		found := router.FindRoute(tc.path)
		if tc.expected && found == nil {
			t.Errorf("路径 %s 应该匹配", tc.path)
		}
		if !tc.expected && found != nil {
			t.Errorf("路径 %s 不应该匹配", tc.path)
		}
	}
}

// TestRouteMatchRegex 测试正则匹配
func TestRouteMatchRegex(t *testing.T) {
	discovery := discovery.NewStaticDiscovery()
	router := NewRouter(discovery)
	route := &Route{
		Path:      `^/api/v1/users/\d+$`,
		MatchType: "regex",
		Service:   "user-service",
	}
	router.AddRoute(route)

	// 测试正则匹配
	testCases := []struct {
		path     string
		expected bool
	}{
		{"/api/v1/users/123", true},
		{"/api/v1/users/456", true},
		{"/api/v1/users/abc", false},
		{"/api/v1/users/123/profile", false},
	}

	for _, tc := range testCases {
		found := router.FindRoute(tc.path)
		if tc.expected && found == nil {
			t.Errorf("路径 %s 应该匹配正则", tc.path)
		}
		if !tc.expected && found != nil {
			t.Errorf("路径 %s 不应该匹配正则", tc.path)
		}
	}
}

// TestRoutePriority 测试路由优先级（第一个匹配的路由生效）
func TestRoutePriority(t *testing.T) {
	discovery := discovery.NewStaticDiscovery()
	router := NewRouter(discovery)

	// 添加更具体的路由
	route1 := &Route{
		Path:      "/api/v1/users/me",
		MatchType: "exact",
		Service:   "user-service",
	}
	router.AddRoute(route1)

	// 添加更通用的路由
	route2 := &Route{
		Path:      "/api/v1/users",
		MatchType: "prefix",
		Service:   "user-service-v2",
	}
	router.AddRoute(route2)

	// 精确匹配应该优先
	found := router.FindRoute("/api/v1/users/me")
	if found == nil || found.Service != "user-service" {
		t.Error("应该匹配到更具体的路由")
	}

	// 前缀匹配应该匹配到第二个路由
	found2 := router.FindRoute("/api/v1/users/123")
	if found2 == nil || found2.Service != "user-service-v2" {
		t.Error("应该匹配到前缀路由")
	}
}

// TestLoadBalancer 测试负载均衡器
func TestLoadBalancer(t *testing.T) {
	instances := []discovery.Instance{
		{ID: "1", Host: "localhost", Port: 8001, Weight: 100, Healthy: true},
		{ID: "2", Host: "localhost", Port: 8002, Weight: 100, Healthy: true},
		{ID: "3", Host: "localhost", Port: 8003, Weight: 100, Healthy: true},
	}

	// 测试轮询负载均衡
	lb := loadbalancer.NewLoadBalancer("round_robin")
	lb.UpdateInstances(instances)

	selected := make(map[string]int)
	for i := 0; i < 30; i++ {
		instance := lb.Select(instances)
		if instance == nil {
			t.Fatal("应该选择到实例")
		}
		selected[instance.ID]++
	}

	// 轮询应该均匀分布
	if len(selected) != 3 {
		t.Errorf("期望选择3个不同实例，实际 %d", len(selected))
	}
	for id, count := range selected {
		if count != 10 {
			t.Errorf("实例 %s 应该被选择10次，实际 %d", id, count)
		}
	}
}

// TestLoadBalancerWeight 测试权重负载均衡
func TestLoadBalancerWeight(t *testing.T) {
	instances := []discovery.Instance{
		{ID: "1", Host: "localhost", Port: 8001, Weight: 100, Healthy: true},
		{ID: "2", Host: "localhost", Port: 8002, Weight: 200, Healthy: true},
		{ID: "3", Host: "localhost", Port: 8003, Weight: 100, Healthy: true},
	}

	lb := loadbalancer.NewLoadBalancer("round_robin")
	lb.UpdateInstances(instances)

	selected := make(map[string]int)
	for i := 0; i < 400; i++ {
		instance := lb.Select(instances)
		if instance == nil {
			t.Fatal("应该选择到实例")
		}
		selected[instance.ID]++
	}

	// 权重为 200 的实例应该被选择更多次
	if selected["2"] <= selected["1"] || selected["2"] <= selected["3"] {
		t.Errorf("权重高的实例应该被选择更多次: %v", selected)
	}
}

// TestLoadBalancerUnhealthy 测试不健康实例过滤
func TestLoadBalancerUnhealthy(t *testing.T) {
	instances := []discovery.Instance{
		{ID: "1", Host: "localhost", Port: 8001, Weight: 100, Healthy: true},
		{ID: "2", Host: "localhost", Port: 8002, Weight: 100, Healthy: false},
		{ID: "3", Host: "localhost", Port: 8003, Weight: 100, Healthy: true},
	}

	lb := loadbalancer.NewLoadBalancer("round_robin")
	lb.UpdateInstances(instances)

	// 只选择健康的实例
	healthyInstances := make([]discovery.Instance, 0)
	for _, inst := range instances {
		if inst.Healthy {
			healthyInstances = append(healthyInstances, inst)
		}
	}

	for i := 0; i < 20; i++ {
		instance := lb.Select(healthyInstances)
		if instance == nil {
			t.Fatal("应该选择到实例")
		}
		if instance.ID == "2" {
			t.Error("不应该选择不健康的实例")
		}
	}
}
