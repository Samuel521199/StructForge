package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"StructForge/backend/common/log"
	nacosClient "StructForge/backend/common/middleware/nacos"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NacosDiscovery Nacos服务发现实现
type NacosDiscovery struct {
	nacosClient *nacosClient.NacosNamingClient // 使用封装的客户端
	instances   map[string][]Instance          // 服务实例缓存
	watchers    map[string][]func([]Instance)  // 服务监听器回调函数
	subscribes  map[string]*vo.SubscribeParam  // 已订阅的服务
	mu          sync.RWMutex
	stopCh      chan struct{}
}

// NewNacosDiscovery 创建Nacos服务发现
func NewNacosDiscovery(nacosConfig *nacosClient.StartupConfig) (*NacosDiscovery, error) {
	if nacosConfig == nil || len(nacosConfig.Nacos.ServerConfigs) == 0 {
		return nil, fmt.Errorf("Nacos配置为空")
	}

	// 使用 common/middleware/nacos 中的命名客户端
	nacosNamingClient, err := nacosClient.NewNacosNamingClient(nacosConfig)
	if err != nil {
		return nil, fmt.Errorf("创建Nacos命名客户端失败: %w", err)
	}

	discovery := &NacosDiscovery{
		nacosClient: nacosNamingClient,
		instances:   make(map[string][]Instance),
		watchers:    make(map[string][]func([]Instance)),
		subscribes:  make(map[string]*vo.SubscribeParam),
		stopCh:      make(chan struct{}),
	}

	// 启动服务监听
	go discovery.watchServices()

	return discovery, nil
}

// GetInstances 获取服务实例
func (d *NacosDiscovery) GetInstances(ctx context.Context, serviceName string) ([]Instance, error) {
	d.mu.RLock()
	instances, exists := d.instances[serviceName]
	d.mu.RUnlock()

	if exists && len(instances) > 0 {
		return instances, nil
	}

	// 如果缓存中没有，从Nacos获取
	return d.fetchInstancesFromNacos(ctx, serviceName)
}

// fetchInstancesFromNacos 从Nacos获取服务实例
func (d *NacosDiscovery) fetchInstancesFromNacos(ctx context.Context, serviceName string) ([]Instance, error) {
	// 使用 NacosNamingClient 获取服务实例
	instances, err := d.nacosClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   "DEFAULT_GROUP",
		HealthyOnly: true,
	})
	if err != nil {
		log.Warn(ctx, "从Nacos获取服务实例失败",
			log.String("service", serviceName),
			log.ErrorField(err),
		)
		return nil, err
	}

	// 转换为内部Instance格式
	result := make([]Instance, 0, len(instances))
	for _, s := range instances {
		// 转换元数据
		metadata := make(map[string]string)
		if s.Metadata != nil {
			metadata = s.Metadata
		}

		result = append(result, Instance{
			ID:       s.InstanceId,
			Host:     s.Ip,
			Port:     int(s.Port),
			Weight:   int(s.Weight),
			Healthy:  s.Healthy,
			Metadata: metadata,
		})
	}

	// 更新缓存
	d.mu.Lock()
	d.instances[serviceName] = result
	d.mu.Unlock()

	return result, nil
}

// RegisterService 注册服务（用于服务注册）
func (d *NacosDiscovery) RegisterService(serviceName string, instances []Instance) {
	// Nacos服务发现模式下，服务应该自己注册到Nacos
	// 这里只更新本地缓存
	d.mu.Lock()
	defer d.mu.Unlock()
	d.instances[serviceName] = instances
}

// Watch 监听服务实例变化（实现ServiceDiscovery接口）
// 使用 Nacos 真正的订阅机制，而不是轮询
func (d *NacosDiscovery) Watch(serviceName string, callback func([]Instance)) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 保存回调函数
	if d.watchers[serviceName] == nil {
		d.watchers[serviceName] = make([]func([]Instance), 0)
	}
	d.watchers[serviceName] = append(d.watchers[serviceName], callback)

	// 如果已经订阅过该服务，直接返回
	if _, exists := d.subscribes[serviceName]; exists {
		log.Info(context.Background(), "服务已订阅，添加新的监听器",
			log.String("service", serviceName),
		)
		return nil
	}

	// TODO: 实现真正的 Nacos 订阅
	// 注意：Nacos SDK 的 SubscribeCallback 类型签名与我们的接口不匹配
	// SubscribeCallback 期望的是 func(services []model.SubscribeService, err error)
	// 但我们使用的是 func(services []model.Instance, err error)
	// 暂时使用轮询方式，后续需要适配 Nacos SDK 的类型
	//
	// 创建订阅参数（暂时不设置 SubscribeCallback，避免类型不匹配）
	// subscribeParam := &vo.SubscribeParam{
	// 	ServiceName:       serviceName,
	// 	GroupName:         "DEFAULT_GROUP",
	// 	SubscribeCallback: d.createSubscribeCallback(serviceName),
	// }
	//
	// err := d.nacosClient.Subscribe(subscribeParam)
	// if err != nil {
	// 	log.Warn(context.Background(), "订阅服务失败，将使用轮询方式",
	// 		log.String("service", serviceName),
	// 		log.ErrorField(err),
	// 	)
	// 	return nil
	// }
	//
	// d.subscribes[serviceName] = subscribeParam

	log.Info(context.Background(), "已注册服务监听器（使用轮询方式）",
		log.String("service", serviceName),
	)

	return nil
}

// createSubscribeCallback 创建 Nacos 订阅回调函数
// TODO: 需要适配 Nacos SDK 的 SubscribeService 类型
// 当前 Nacos SDK 的 SubscribeCallback 期望的是 func(services []model.SubscribeService, err error)
// 但我们的实现使用的是 model.Instance，需要类型转换
// 暂时注释掉，使用轮询方式
// func (d *NacosDiscovery) createSubscribeCallback(serviceName string) func(services []model.SubscribeService, err error) {
// 	return func(services []model.SubscribeService, err error) {
// 		// 转换逻辑...
// 	}
// }

// Register 注册服务实例（实现ServiceDiscovery接口）
func (d *NacosDiscovery) Register(ctx context.Context, serviceName string, instance Instance) error {
	// 使用 NacosNamingClient 注册实例
	param := vo.RegisterInstanceParam{
		Ip:          instance.Host,
		Port:        uint64(instance.Port),
		ServiceName: serviceName,
		Weight:      float64(instance.Weight),
		Enable:      true,
		Healthy:     instance.Healthy,
		Metadata:    instance.Metadata,
	}

	_, err := d.nacosClient.RegisterInstance(param)
	if err != nil {
		return err
	}

	// 同时更新本地缓存
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.instances[serviceName] == nil {
		d.instances[serviceName] = make([]Instance, 0)
	}
	d.instances[serviceName] = append(d.instances[serviceName], instance)

	return nil
}

// Deregister 注销服务实例（实现ServiceDiscovery接口）
func (d *NacosDiscovery) Deregister(ctx context.Context, serviceName string, instanceID string) error {
	// 从本地缓存中查找实例信息
	d.mu.RLock()
	instances, exists := d.instances[serviceName]
	var instanceToDeregister *Instance
	for i := range instances {
		if instances[i].ID == instanceID {
			instanceToDeregister = &instances[i]
			break
		}
	}
	d.mu.RUnlock()

	if !exists || instanceToDeregister == nil {
		return fmt.Errorf("服务实例 %s 未找到", instanceID)
	}

	// 使用 NacosNamingClient 注销实例
	param := vo.DeregisterInstanceParam{
		Ip:          instanceToDeregister.Host,
		Port:        uint64(instanceToDeregister.Port),
		ServiceName: serviceName,
	}

	_, err := d.nacosClient.DeregisterInstance(param)
	if err != nil {
		return err
	}

	// 从本地缓存中移除
	d.mu.Lock()
	defer d.mu.Unlock()
	newInstances := make([]Instance, 0)
	for _, instance := range d.instances[serviceName] {
		if instance.ID != instanceID {
			newInstances = append(newInstances, instance)
		}
	}
	d.instances[serviceName] = newInstances

	return nil
}

// watchServices 监听服务变化
// 对于未订阅的服务，使用轮询方式作为降级方案
func (d *NacosDiscovery) watchServices() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒刷新一次（作为降级方案）
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 只刷新未订阅的服务（降级方案）
			d.refreshUnsubscribedServices()
		case <-d.stopCh:
			return
		}
	}
}

// refreshUnsubscribedServices 刷新未订阅的服务实例列表（降级方案）
func (d *NacosDiscovery) refreshUnsubscribedServices() {
	d.mu.RLock()
	serviceNames := make([]string, 0)
	watchersCopy := make(map[string][]func([]Instance))

	// 只处理未订阅的服务
	for serviceName := range d.instances {
		if _, subscribed := d.subscribes[serviceName]; !subscribed {
			serviceNames = append(serviceNames, serviceName)
		}
	}

	for serviceName, watchers := range d.watchers {
		if _, subscribed := d.subscribes[serviceName]; !subscribed {
			watchersCopy[serviceName] = make([]func([]Instance), len(watchers))
			copy(watchersCopy[serviceName], watchers)
		}
	}
	d.mu.RUnlock()

	ctx := context.Background()
	for _, serviceName := range serviceNames {
		instances, err := d.fetchInstancesFromNacos(ctx, serviceName)
		if err == nil {
			log.Debug(ctx, "刷新服务实例（轮询方式）",
				log.String("service", serviceName),
				log.Int("instances", len(instances)),
			)

			// 调用监听器回调
			if watchers, exists := watchersCopy[serviceName]; exists {
				for _, watcher := range watchers {
					if watcher != nil {
						watcher(instances)
					}
				}
			}
		}
	}
}

// Stop 停止服务发现
func (d *NacosDiscovery) Stop() {
	// 取消所有订阅
	d.mu.Lock()
	defer d.mu.Unlock()

	for serviceName, param := range d.subscribes {
		err := d.nacosClient.Unsubscribe(param)
		if err != nil {
			log.Warn(context.Background(), "取消订阅失败",
				log.String("service", serviceName),
				log.ErrorField(err),
			)
		} else {
			log.Info(context.Background(), "已取消订阅服务",
				log.String("service", serviceName),
			)
		}
	}

	// 清空订阅列表
	d.subscribes = make(map[string]*vo.SubscribeParam)

	// 停止轮询 goroutine
	close(d.stopCh)
}
