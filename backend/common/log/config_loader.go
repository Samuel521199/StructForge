package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// LoadConfigFromFile 从文件加载配置
func LoadConfigFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config

	// 根据文件扩展名选择解析方式
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return Config{}, fmt.Errorf("YAML支持需要安装依赖: go get gopkg.in/yaml.v3。当前仅支持JSON格式")
	} else if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, &config); err != nil {
			return Config{}, fmt.Errorf("解析JSON配置失败: %w", err)
		}
	} else {
		return Config{}, fmt.Errorf("不支持的配置文件格式，支持 .yaml, .yml, .json")
	}

	// 应用环境变量覆盖
	config = applyEnvOverrides(config)

	// 验证配置
	if err := validateConfig(config); err != nil {
		return Config{}, fmt.Errorf("配置验证失败: %w", err)
	}

	return config, nil
}

// LoadConfigFromBytes 从字节数组加载配置
func LoadConfigFromBytes(data []byte, format string) (Config, error) {
	var config Config

	switch strings.ToLower(format) {
	case "yaml", "yml":
		return Config{}, fmt.Errorf("YAML支持需要安装依赖: go get gopkg.in/yaml.v3。当前仅支持JSON格式")
	case "json":
		if err := json.Unmarshal(data, &config); err != nil {
			return Config{}, fmt.Errorf("解析JSON配置失败: %w", err)
		}
	default:
		return Config{}, fmt.Errorf("不支持的配置格式: %s", format)
	}

	// 应用环境变量覆盖
	config = applyEnvOverrides(config)

	// 验证配置
	if err := validateConfig(config); err != nil {
		return Config{}, fmt.Errorf("配置验证失败: %w", err)
	}

	return config, nil
}

// applyEnvOverrides 应用环境变量覆盖
func applyEnvOverrides(config Config) Config {
	// 服务名称
	if serviceName := os.Getenv("LOG_SERVICE_NAME"); serviceName != "" {
		config.ServiceName = serviceName
	}
	if serviceID := os.Getenv("LOG_SERVICE_ID"); serviceID != "" {
		config.ServiceID = serviceID
	}
	if instanceID := os.Getenv("LOG_INSTANCE_ID"); instanceID != "" {
		config.InstanceID = instanceID
	}

	// 日志级别
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		config.Level = ParseLevel(levelStr)
	}

	// 启用颜色
	if colorStr := os.Getenv("LOG_ENABLE_COLOR"); colorStr != "" {
		config.EnableColor = colorStr == "true" || colorStr == "1"
	}

	// 文件路径
	if filePath := os.Getenv("LOG_FILE_PATH"); filePath != "" {
		config.File.Path = filePath
	}

	return config
}

// validateConfig 验证配置
func validateConfig(config Config) error {
	// 验证日志级别
	if config.Level < DebugLevel || config.Level > FatalLevel {
		return fmt.Errorf("无效的日志级别: %d", config.Level)
	}

	// 验证采样比例
	if config.Sampling.Enabled {
		if config.Sampling.Ratio < 0 || config.Sampling.Ratio > 1 {
			return fmt.Errorf("采样比例必须在0-1之间: %f", config.Sampling.Ratio)
		}
	}

	// 验证异步配置
	if config.Async.Enabled {
		if config.Async.QueueSize <= 0 {
			return fmt.Errorf("队列大小必须大于0: %d", config.Async.QueueSize)
		}
		if config.Async.BatchSize <= 0 {
			return fmt.Errorf("批量大小必须大于0: %d", config.Async.BatchSize)
		}
		if config.Async.FlushInterval <= 0 {
			return fmt.Errorf("刷新间隔必须大于0: %v", config.Async.FlushInterval)
		}
	}

	// 验证脱敏配置
	if config.Mask.Enabled {
		if config.Mask.KeepHead < 0 {
			return fmt.Errorf("保留头部字符数不能为负数: %d", config.Mask.KeepHead)
		}
		if config.Mask.KeepTail < 0 {
			return fmt.Errorf("保留尾部字符数不能为负数: %d", config.Mask.KeepTail)
		}
	}

	return nil
}

// SaveConfigToFile 保存配置到文件
func SaveConfigToFile(config Config, path string) error {
	var data []byte
	var err error

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return fmt.Errorf("YAML支持需要安装依赖: go get gopkg.in/yaml.v3。当前仅支持JSON格式")
	} else if strings.HasSuffix(path, ".json") {
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("序列化JSON配置失败: %w", err)
		}
	} else {
		return fmt.Errorf("不支持的配置文件格式，支持 .yaml, .yml, .json")
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
