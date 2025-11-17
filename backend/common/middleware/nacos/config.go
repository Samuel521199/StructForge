package nacos

// NacosConfig Nacos 配置结构
type NacosConfig struct {
	// 服务器配置
	ServerConfigs []ServerConfig `yaml:"server_configs" json:"server_configs"`
	// 客户端配置
	ClientConfig ClientConfig `yaml:"client_config" json:"client_config"`
	// 配置中心相关
	ConfigCenter ConfigCenter `yaml:"config_center" json:"config_center"`
}

// ServerConfig Nacos 服务器配置
type ServerConfig struct {
	// 服务器地址
	IpAddr string `yaml:"ip_addr" json:"ip_addr"`
	// 服务器端口
	Port uint64 `yaml:"port" json:"port"`
	// 上下文路径
	ContextPath string `yaml:"context_path" json:"context_path"`
	// 是否使用 HTTPS
	Scheme string `yaml:"scheme" json:"scheme"`
}

// ClientConfig Nacos 客户端配置
type ClientConfig struct {
	// 命名空间 ID
	NamespaceId string `yaml:"namespace_id" json:"namespace_id"`
	// 超时时间（毫秒）
	TimeoutMs uint64 `yaml:"timeout_ms" json:"timeout_ms"`
	// 是否在启动时不加载本地缓存
	NotLoadCacheAtStart bool `yaml:"not_load_cache_at_start" json:"not_load_cache_at_start"`
	// 日志目录
	LogDir string `yaml:"log_dir" json:"log_dir"`
	// 缓存目录
	CacheDir string `yaml:"cache_dir" json:"cache_dir"`
	// 日志级别
	LogLevel string `yaml:"log_level" json:"log_level"`
	// 用户名（可选）
	Username string `yaml:"username" json:"username"`
	// 密码（可选）
	Password string `yaml:"password" json:"password"`
	// 访问令牌（可选）
	AccessKey string `yaml:"access_key" json:"access_key"`
	// 密钥（可选）
	SecretKey string `yaml:"secret_key" json:"secret_key"`
}

// ConfigCenter 配置中心配置
type ConfigCenter struct {
	// 是否启用配置中心
	Enabled bool `yaml:"enabled" json:"enabled"`
	// 数据 ID
	DataId string `yaml:"data_id" json:"data_id"`
	// 分组
	Group string `yaml:"group" json:"group"`
	// 命名空间
	Namespace string `yaml:"namespace" json:"namespace"`
}

// StartupConfig 启动配置（包含 Nacos 配置）
type StartupConfig struct {
	// Nacos 配置
	Nacos NacosConfig `yaml:"nacos" json:"nacos"`
}
