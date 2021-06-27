// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

import (
	"github.com/vmihailenco/msgpack/v5"
)

// OptionType 配置选项接口
type OptionType interface {
	Unmarshaler() Unmarshaler
	Marshaler() Marshaler

	Check() error
}

// OptionDefault 默认全局的配置
var OptionDefault = &Option{
	UnmarshalFn: msgpack.Unmarshal,
	MarshalFn:   msgpack.Marshal,
}

// Option 配置选型
type Option struct {
	// 数据反序列化方法
	UnmarshalFn Unmarshaler

	// 数据序列化方法
	MarshalFn Marshaler
}

// Unmarshaler 获取反序列化的方法
func (o *Option) Unmarshaler() Unmarshaler {
	if o.UnmarshalFn == nil {
		return OptionDefault.UnmarshalFn
	}
	return o.UnmarshalFn
}

// Marshaler 获取序列化的方法
func (o *Option) Marshaler() Marshaler {
	if o.MarshalFn == nil {
		return OptionDefault.MarshalFn
	}
	return o.MarshalFn
}

// Check 检查配置是否正确
func (o *Option) Check() error {
	return nil
}

var _ OptionType = (*Option)(nil)
