package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

// Serializer 序列化接口
type Serializer interface {
	// Serialize 序列化对象为字节数组
	Serialize(v interface{}) ([]byte, error)

	// Deserialize 从字节数组反序列化对象
	Deserialize(data []byte, v interface{}) error
}

// JSONSerializer JSON序列化器
type JSONSerializer struct{}

// Serialize 序列化
func (s *JSONSerializer) Serialize(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, ErrInvalidValue
	}
	return json.Marshal(v)
}

// Deserialize 反序列化
func (s *JSONSerializer) Deserialize(data []byte, v interface{}) error {
	if len(data) == 0 {
		return ErrInvalidValue
	}
	return json.Unmarshal(data, v)
}

// GobSerializer Gob序列化器
type GobSerializer struct{}

// Serialize 序列化
func (s *GobSerializer) Serialize(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, ErrInvalidValue
	}
	// Gob序列化需要特殊处理，这里简化实现
	// 实际使用时可能需要更复杂的实现
	return nil, errors.New("gob serializer not fully implemented, use json instead")
}

// Deserialize 反序列化
func (s *GobSerializer) Deserialize(data []byte, v interface{}) error {
	if len(data) == 0 {
		return ErrInvalidValue
	}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(v)
}

// NewSerializer 创建序列化器
func NewSerializer(serializerType SerializerType) Serializer {
	switch serializerType {
	case SerializerJSON:
		return &JSONSerializer{}
	case SerializerGob:
		return &GobSerializer{}
	case SerializerMsgPack:
		// MsgPack需要第三方库，这里先返回JSON
		// 后续可以集成 github.com/vmihailenco/msgpack
		return &JSONSerializer{}
	default:
		return &JSONSerializer{}
	}
}
