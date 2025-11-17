package nacos

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NacosNamingClient Nacos 服务发现客户端封装
type NacosNamingClient struct {
	client naming_client.INamingClient
	config *NacosConfig
}

// NewNacosNamingClient 创建 Nacos 服务发现客户端
func NewNacosNamingClient(startupConfig *StartupConfig) (*NacosNamingClient, error) {
	if startupConfig == nil {
		return nil, fmt.Errorf("startup config is nil")
	}

	nacosConfig := &startupConfig.Nacos

	// 转换服务器配置
	serverConfigs := make([]constant.ServerConfig, 0, len(nacosConfig.ServerConfigs))
	for _, sc := range nacosConfig.ServerConfigs {
		serverConfig := constant.ServerConfig{
			IpAddr:      sc.IpAddr,
			Port:        sc.Port,
			ContextPath: sc.ContextPath,
			Scheme:      sc.Scheme,
		}
		if serverConfig.ContextPath == "" {
			serverConfig.ContextPath = "/nacos"
		}
		if serverConfig.Scheme == "" {
			serverConfig.Scheme = "http"
		}
		serverConfigs = append(serverConfigs, serverConfig)
	}

	// 如果服务器配置为空，使用默认配置
	if len(serverConfigs) == 0 {
		serverConfigs = []constant.ServerConfig{
			*constant.NewServerConfig("127.0.0.1", 8848),
		}
	}

	// 转换客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosConfig.ClientConfig.NamespaceId,
		TimeoutMs:           nacosConfig.ClientConfig.TimeoutMs,
		NotLoadCacheAtStart: nacosConfig.ClientConfig.NotLoadCacheAtStart,
		LogDir:              nacosConfig.ClientConfig.LogDir,
		CacheDir:            nacosConfig.ClientConfig.CacheDir,
		LogLevel:            nacosConfig.ClientConfig.LogLevel,
		Username:            nacosConfig.ClientConfig.Username,
		Password:            nacosConfig.ClientConfig.Password,
		AccessKey:           nacosConfig.ClientConfig.AccessKey,
		SecretKey:           nacosConfig.ClientConfig.SecretKey,
	}

	// 设置默认值
	if clientConfig.NamespaceId == "" {
		clientConfig.NamespaceId = "public"
	}
	if clientConfig.TimeoutMs == 0 {
		clientConfig.TimeoutMs = 5000
	}
	if clientConfig.LogDir == "" {
		clientConfig.LogDir = "tmp/nacos/log"
	}
	if clientConfig.CacheDir == "" {
		clientConfig.CacheDir = "tmp/nacos/cache"
	}
	if clientConfig.LogLevel == "" {
		clientConfig.LogLevel = "info"
	}

	// 创建服务发现客户端
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create nacos naming client: %w", err)
	}

	return &NacosNamingClient{
		client: namingClient,
		config: nacosConfig,
	}, nil
}

// RegisterInstance 注册服务实例
func (c *NacosNamingClient) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	success, err := c.client.RegisterInstance(param)
	if err != nil {
		return false, fmt.Errorf("failed to register instance: %w", err)
	}

	return success, nil
}

// DeregisterInstance 注销服务实例
func (c *NacosNamingClient) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	success, err := c.client.DeregisterInstance(param)
	if err != nil {
		return false, fmt.Errorf("failed to deregister instance: %w", err)
	}

	return success, nil
}

// GetService 获取服务实例列表
func (c *NacosNamingClient) GetService(param vo.GetServiceParam) (model.Service, error) {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	service, err := c.client.GetService(param)
	if err != nil {
		return model.Service{}, fmt.Errorf("failed to get service: %w", err)
	}

	return service, nil
}

// SelectAllInstances 选择所有服务实例
func (c *NacosNamingClient) SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error) {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	instances, err := c.client.SelectAllInstances(param)
	if err != nil {
		return nil, fmt.Errorf("failed to select all instances: %w", err)
	}

	return instances, nil
}

// SelectInstances 选择服务实例（支持健康检查）
func (c *NacosNamingClient) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	instances, err := c.client.SelectInstances(param)
	if err != nil {
		return nil, fmt.Errorf("failed to select instances: %w", err)
	}

	return instances, nil
}

// SelectOneHealthyInstance 选择一个健康的服务实例
func (c *NacosNamingClient) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error) {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	instance, err := c.client.SelectOneHealthyInstance(param)
	if err != nil {
		return nil, fmt.Errorf("failed to select one healthy instance: %w", err)
	}

	return instance, nil
}

// Subscribe 订阅服务变化
func (c *NacosNamingClient) Subscribe(param *vo.SubscribeParam) error {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	err := c.client.Subscribe(param)
	if err != nil {
		return fmt.Errorf("failed to subscribe service: %w", err)
	}

	return nil
}

// Unsubscribe 取消订阅服务变化
func (c *NacosNamingClient) Unsubscribe(param *vo.SubscribeParam) error {
	if param.GroupName == "" {
		param.GroupName = "DEFAULT_GROUP"
	}

	err := c.client.Unsubscribe(param)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe service: %w", err)
	}

	return nil
}

// GetClient 获取底层服务发现客户端
func (c *NacosNamingClient) GetClient() naming_client.INamingClient {
	return c.client
}
