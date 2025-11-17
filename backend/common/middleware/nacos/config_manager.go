package nacos

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	configClient  *NacosConfigClient
	startupConfig *StartupConfig
	bootstrap     interface{}
}

// NewConfigManager 创建配置管理器
// bootstrapPtr 是指向 Bootstrap 结构体的指针，用于扫描配置
func NewConfigManager(configClient *NacosConfigClient, startupConfig *StartupConfig, bootstrapPtr interface{}) (*ConfigManager, error) {
	if configClient == nil {
		return nil, fmt.Errorf("config client is nil")
	}
	if startupConfig == nil {
		return nil, fmt.Errorf("startup config is nil")
	}
	if bootstrapPtr == nil {
		return nil, fmt.Errorf("bootstrap pointer is nil")
	}

	cc := configClient.GetClient()
	if cc == nil {
		return nil, fmt.Errorf("nacos config client is nil")
	}

	// 获取配置中心配置
	configCenter := startupConfig.Nacos.ConfigCenter
	if !configCenter.Enabled {
		return nil, fmt.Errorf("config center is not enabled")
	}

	dataId := configCenter.DataId
	if dataId == "" {
		return nil, fmt.Errorf("config center data_id is empty")
	}

	group := configCenter.Group
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	// 从 Nacos 获取配置内容
	content, err := configClient.GetConfig(dataId, group)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from nacos: %w", err)
	}

	// 使用 YAML 解析配置内容到结构体
	if err := yaml.Unmarshal([]byte(content), bootstrapPtr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &ConfigManager{
		configClient:  configClient,
		startupConfig: startupConfig,
		bootstrap:     bootstrapPtr,
	}, nil
}

// GetBootstrap 获取 Bootstrap 配置
func (m *ConfigManager) GetBootstrap() interface{} {
	return m.bootstrap
}

// Reload 重新加载配置
func (m *ConfigManager) Reload() error {
	configCenter := m.startupConfig.Nacos.ConfigCenter
	dataId := configCenter.DataId
	group := configCenter.Group
	if group == "" {
		group = "DEFAULT_GROUP"
	}

	// 从 Nacos 重新获取配置内容
	content, err := m.configClient.GetConfig(dataId, group)
	if err != nil {
		return fmt.Errorf("failed to reload config from nacos: %w", err)
	}

	// 重新解析配置
	if err := yaml.Unmarshal([]byte(content), m.bootstrap); err != nil {
		return fmt.Errorf("failed to unmarshal reloaded config: %w", err)
	}

	return nil
}

// Close 关闭配置管理器
func (m *ConfigManager) Close() error {
	// 目前没有需要关闭的资源
	return nil
}
