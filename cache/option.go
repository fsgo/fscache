/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/17
 */

package cache

import (
	"github.com/vmihailenco/msgpack/v4"
)

// Checker 检查是否有错误
type Checker interface {
	Check() error
}

// IOption 配置选项接口
type IOption interface {
	Unmarshaler() Unmarshaler
	Marshaler() Marshaler

	Checker
}

// Option 配置选型
type Option struct {
	// 数据反序列化方法
	UnmarshalFn Unmarshaler

	// 数据序列化方法
	MarshalFn Marshaler
}

// Unmarshaler 获取反序列化的方法
func (o Option) Unmarshaler() Unmarshaler {
	if o.UnmarshalFn == nil {
		return msgpack.Unmarshal
	}
	return o.UnmarshalFn
}

// Marshaler 获取序列化的方法
func (o Option) Marshaler() Marshaler {
	if o.MarshalFn == nil {
		return msgpack.Marshal
	}
	return o.MarshalFn
}

// Check 检查是否正确
func (o Option) Check() error {
	return nil
}

var _ IOption = (*Option)(nil)
