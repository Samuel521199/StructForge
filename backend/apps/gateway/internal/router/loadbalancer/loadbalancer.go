package loadbalancer

import (
	"math/rand"
	"sync"

	"StructForge/backend/apps/gateway/internal/router/discovery"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	// Select 选择服务实例
	Select(instances []discovery.Instance) *discovery.Instance
	// UpdateInstances 更新服务实例列表
	UpdateInstances(instances []discovery.Instance)
}

// RoundRobinLoadBalancer 轮询负载均衡器
type RoundRobinLoadBalancer struct {
	instances []discovery.Instance
	current   int
	mu        sync.Mutex
}

// NewRoundRobinLoadBalancer 创建轮询负载均衡器
func NewRoundRobinLoadBalancer() *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{
		instances: make([]discovery.Instance, 0),
		current:   0,
	}
}

// Select 选择服务实例（轮询）
func (lb *RoundRobinLoadBalancer) Select(instances []discovery.Instance) *discovery.Instance {
	if len(instances) == 0 {
		return nil
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 过滤健康的实例
	healthyInstances := make([]discovery.Instance, 0)
	for _, instance := range instances {
		if instance.Healthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	if len(healthyInstances) == 0 {
		return nil
	}

	instance := &healthyInstances[lb.current%len(healthyInstances)]
	lb.current++

	return instance
}

// UpdateInstances 更新服务实例列表
func (lb *RoundRobinLoadBalancer) UpdateInstances(instances []discovery.Instance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.instances = instances
}

// RandomLoadBalancer 随机负载均衡器
type RandomLoadBalancer struct {
	instances []discovery.Instance
	mu        sync.RWMutex
}

// NewRandomLoadBalancer 创建随机负载均衡器
func NewRandomLoadBalancer() *RandomLoadBalancer {
	return &RandomLoadBalancer{
		instances: make([]discovery.Instance, 0),
	}
}

// Select 选择服务实例（随机）
func (lb *RandomLoadBalancer) Select(instances []discovery.Instance) *discovery.Instance {
	if len(instances) == 0 {
		return nil
	}

	// 过滤健康的实例
	healthyInstances := make([]discovery.Instance, 0)
	for _, instance := range instances {
		if instance.Healthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	if len(healthyInstances) == 0 {
		return nil
	}

	index := rand.Intn(len(healthyInstances))
	return &healthyInstances[index]
}

// UpdateInstances 更新服务实例列表
func (rb *RandomLoadBalancer) UpdateInstances(instances []discovery.Instance) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.instances = instances
}

// LeastConnectionsLoadBalancer 最少连接负载均衡器
type LeastConnectionsLoadBalancer struct {
	instances   []discovery.Instance
	connections map[string]int
	mu          sync.RWMutex
}

// NewLeastConnectionsLoadBalancer 创建最少连接负载均衡器
func NewLeastConnectionsLoadBalancer() *LeastConnectionsLoadBalancer {
	return &LeastConnectionsLoadBalancer{
		instances:   make([]discovery.Instance, 0),
		connections: make(map[string]int),
	}
}

// Select 选择服务实例（最少连接）
func (lb *LeastConnectionsLoadBalancer) Select(instances []discovery.Instance) *discovery.Instance {
	if len(instances) == 0 {
		return nil
	}

	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 过滤健康的实例
	healthyInstances := make([]discovery.Instance, 0)
	for _, instance := range instances {
		if instance.Healthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	if len(healthyInstances) == 0 {
		return nil
	}

	// 选择连接数最少的实例
	var selected *discovery.Instance
	minConnections := -1

	for i := range healthyInstances {
		instanceID := healthyInstances[i].ID
		connections := lb.connections[instanceID]
		if minConnections == -1 || connections < minConnections {
			minConnections = connections
			selected = &healthyInstances[i]
		}
	}

	// 增加连接数
	if selected != nil {
		lb.connections[selected.ID]++
	}

	return selected
}

// UpdateInstances 更新服务实例列表
func (lb *LeastConnectionsLoadBalancer) UpdateInstances(instances []discovery.Instance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.instances = instances
}

// NewLoadBalancer 根据策略创建负载均衡器
func NewLoadBalancer(strategy string) LoadBalancer {
	switch strategy {
	case "random":
		return NewRandomLoadBalancer()
	case "least_connections":
		return NewLeastConnectionsLoadBalancer()
	case "round_robin":
		fallthrough
	default:
		return NewRoundRobinLoadBalancer()
	}
}
