// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

// Option 配置
type Option struct {
	// Codec 编解码器，可选，默认为 msgpack
	Codec Codec
}

// GetCodec 获取编解码器，若没有设置，会返回默认值(msgpack)
func (o *Option) GetCodec() Codec {
	if o == nil || o.Codec == nil {
		return defaultCodec
	}
	return o.Codec
}
