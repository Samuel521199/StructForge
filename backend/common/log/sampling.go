package log

import (
	"hash/fnv"
	"sync"
)

// sampler 采样器
type sampler struct {
	config SamplingConfig
	mu     sync.RWMutex
}

// shouldSample 判断是否应该采样
func (s *sampler) shouldSample(level Level) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 如果未启用采样，全部记录
	if !s.config.Enabled {
		return true
	}

	// 如果指定了级别列表，检查当前级别是否在列表中
	if len(s.config.Levels) > 0 {
		found := false
		for _, l := range s.config.Levels {
			if l == level {
				found = true
				break
			}
		}
		// 如果当前级别不在采样列表中，全部记录（不采样）
		if !found {
			return true
		}
		// 如果当前级别在采样列表中，继续采样判断
	}

	// 使用哈希算法决定是否采样（保证相同trace_id的日志采样一致性）
	// 这里使用简单的哈希，实际可以使用trace_id等更稳定的值
	hash := fnv.New32a()
	hash.Write([]byte{byte(level)})
	hashValue := hash.Sum32()

	// 将hash值映射到0-1之间
	ratio := float64(hashValue%10000) / 10000.0

	// 如果hash值小于采样比例，则采样
	return ratio < s.config.Ratio
}

// shouldSampleWithKey 使用指定key进行采样（保证相同key的日志采样一致性）
func (s *sampler) shouldSampleWithKey(level Level, key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 如果未启用采样，全部记录
	if !s.config.Enabled {
		return true
	}

	// 如果指定了级别列表，检查当前级别是否在列表中
	if len(s.config.Levels) > 0 {
		found := false
		for _, l := range s.config.Levels {
			if l == level {
				found = true
				break
			}
		}
		// 如果当前级别不在采样列表中，全部记录（不采样）
		if !found {
			return true
		}
		// 如果当前级别在采样列表中，继续采样判断
	}

	// 使用key进行哈希，保证相同key的日志采样一致性
	hash := fnv.New32a()
	hash.Write([]byte(key))
	hashValue := hash.Sum32()

	// 将hash值映射到0-1之间
	ratio := float64(hashValue%10000) / 10000.0

	// 如果hash值小于采样比例，则采样
	return ratio < s.config.Ratio
}
