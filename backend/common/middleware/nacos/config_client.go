package nacos

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NacosConfigClient Nacos 配置客户端封装
type NacosConfigClient struct {
	client config_client.IConfigClient
	config *NacosConfig
}

// NewNacosConfigClient 创建 Nacos 配置客户端
func NewNacosConfigClient(startupConfig *StartupConfig) (*NacosConfigClient, error) {
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

	// 创建配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create nacos config client: %w", err)
	}

	return &NacosConfigClient{
		client: configClient,
		config: nacosConfig,
	}, nil
}

// GetConfig 获取配置
func (c *NacosConfigClient) GetConfig(dataId, group string) (string, error) {
	if dataId == "" {
		dataId = c.config.ConfigCenter.DataId
	}
	if group == "" {
		group = c.config.ConfigCenter.Group
	}
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	content, err := c.client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get config from nacos: %w", err)
	}

	return content, nil
}

// PublishConfig 发布配置
func (c *NacosConfigClient) PublishConfig(dataId, group, content string) (bool, error) {
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	success, err := c.client.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	if err != nil {
		return false, fmt.Errorf("failed to publish config to nacos: %w", err)
	}

	return success, nil
}

// DeleteConfig 删除配置
func (c *NacosConfigClient) DeleteConfig(dataId, group string) (bool, error) {
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	success, err := c.client.DeleteConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return false, fmt.Errorf("failed to delete config from nacos: %w", err)
	}

	return success, nil
}

// ListenConfig 监听配置变化
func (c *NacosConfigClient) ListenConfig(dataId, group string, onChange func(namespace, group, dataId, data string)) error {
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	err := c.client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			if onChange != nil {
				onChange(namespace, group, dataId, data)
			}
		},
	})
	if err != nil {
		return fmt.Errorf("failed to listen config from nacos: %w", err)
	}

	return nil
}

// GetClient 获取底层配置客户端
func (c *NacosConfigClient) GetClient() config_client.IConfigClient {
	return c.client
}
