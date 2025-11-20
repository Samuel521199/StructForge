package discovery

import (
	"context"
	"fmt"
	"sync"

	"StructForge/backend/common/log"
)

// Instance 服务实例
type Instance struct {
	// 实例ID
	ID string `json:"id"`
	// 主机地址
	Host string `json:"host"`
	// 端口
	Port int `json:"port"`
	// 权重（用于负载均衡）
	Weight int `json:"weight"`
	// 健康状态
	Healthy bool `json:"healthy"`
	// 元数据
	Metadata map[string]string `json:"metadata"`
}

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	// GetInstances 获取服务实例列表
	GetInstances(ctx context.Context, serviceName string) ([]Instance, error)
	// Watch 监听服务实例变化
	Watch(serviceName string, callback func([]Instance)) error
	// Register 注册服务实例（用于服务自注册）
	Register(ctx context.Context, serviceName string, instance Instance) error
	// Deregister 注销服务实例
	Deregister(ctx context.Context, serviceName string, instanceID string) error
}

// StaticDiscovery 静态服务发现（基于配置文件）
type StaticDiscovery struct {
	services map[string][]Instance
	watchers map[string][]func([]Instance)
	mu       sync.RWMutex
}

// NewStaticDiscovery 创建静态服务发现
func NewStaticDiscovery() *StaticDiscovery {
	return &StaticDiscovery{
		services: make(map[string][]Instance),
		watchers: make(map[string][]func([]Instance)),
	}
}

// RegisterService 注册服务（从配置文件加载）
func (d *StaticDiscovery) RegisterService(serviceName string, instances []Instance) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.services[serviceName] = instances

	// 通知监听者
	if watchers, exists := d.watchers[serviceName]; exists {
		for _, watcher := range watchers {
			watcher(instances)
		}
	}

	log.Info(context.Background(), "服务已注册（静态）",
		log.String("service", serviceName),
		log.Int("instances", len(instances)),
	)
}

// GetInstances 获取服务实例
func (d *StaticDiscovery) GetInstances(ctx context.Context, serviceName string) ([]Instance, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	instances, exists := d.services[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务 %s 未找到", serviceName)
	}

	// 只返回健康的实例
	healthyInstances := make([]Instance, 0)
	for _, instance := range instances {
		if instance.Healthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	if len(healthyInstances) == 0 {
		return nil, fmt.Errorf("服务 %s 没有健康的实例", serviceName)
	}

	return healthyInstances, nil
}

// Watch 监听服务实例变化
func (d *StaticDiscovery) Watch(serviceName string, callback func([]Instance)) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.watchers[serviceName] == nil {
		d.watchers[serviceName] = make([]func([]Instance), 0)
	}
	d.watchers[serviceName] = append(d.watchers[serviceName], callback)

	return nil
}

// Register 注册服务实例（用于服务自注册）
func (d *StaticDiscovery) Register(ctx context.Context, serviceName string, instance Instance) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.services[serviceName] == nil {
		d.services[serviceName] = make([]Instance, 0)
	}

	// 检查是否已存在
	for i, existing := range d.services[serviceName] {
		if existing.ID == instance.ID {
			d.services[serviceName][i] = instance
			log.Info(ctx, "服务实例已更新",
				log.String("service", serviceName),
				log.String("instance", instance.ID),
			)
			return nil
		}
	}

	d.services[serviceName] = append(d.services[serviceName], instance)

	// 通知监听者
	if watchers, exists := d.watchers[serviceName]; exists {
		for _, watcher := range watchers {
			watcher(d.services[serviceName])
		}
	}

	log.Info(ctx, "服务实例已注册",
		log.String("service", serviceName),
		log.String("instance", instance.ID),
	)

	return nil
}

// Deregister 注销服务实例
func (d *StaticDiscovery) Deregister(ctx context.Context, serviceName string, instanceID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	instances, exists := d.services[serviceName]
	if !exists {
		return fmt.Errorf("服务 %s 未找到", serviceName)
	}

	newInstances := make([]Instance, 0)
	for _, instance := range instances {
		if instance.ID != instanceID {
			newInstances = append(newInstances, instance)
		}
	}

	d.services[serviceName] = newInstances

	// 通知监听者
	if watchers, exists := d.watchers[serviceName]; exists {
		for _, watcher := range watchers {
			watcher(newInstances)
		}
	}

	log.Info(ctx, "服务实例已注销",
		log.String("service", serviceName),
		log.String("instance", instanceID),
	)

	return nil
}
