// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/6

package fscache

import (
	"github.com/vmihailenco/msgpack/v5"
)

// Codec 数据编解码器
type Codec interface {
	// Marshal 数据序列化方法
	Marshal(obj any) ([]byte, error)

	// Unmarshal 数据反序列化方法
	Unmarshal(bf []byte, obj any) error
}

type (
	// MarshalFunc 数据序列化方法
	MarshalFunc func(obj any) ([]byte, error)

	// UnmarshalFunc 数据反序列化方法
	UnmarshalFunc func(bf []byte, obj any) error
)

var defaultCodec = NewCodec(msgpack.Marshal, msgpack.Unmarshal)

// NewCodec 创建一个新的编解码器
func NewCodec(encode MarshalFunc, decode UnmarshalFunc) Codec {
	return &codec{
		encode: encode,
		decode: decode,
	}
}

var _ Codec = (*codec)(nil)

type codec struct {
	encode func(obj any) ([]byte, error)
	decode func(bf []byte, obj any) error
}

func (c *codec) Marshal(obj any) ([]byte, error) {
	return c.encode(obj)
}

func (c *codec) Unmarshal(bf []byte, obj any) error {
	return c.decode(bf, obj)
}
