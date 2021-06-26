// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/17

package freecache

import (
	"github.com/fsgo/fscache"
)

// OptionType filecache 选项接口
type OptionType interface {
	fscache.OptionType

	GetMemSize() int
}

// Option 配置
type Option struct {
	MemSize int
	fscache.Option
}

// GetMemSize p获取配置的内存大小
func (o *Option) GetMemSize() int {
	return o.MemSize
}

// Check 检查是否正确
func (o *Option) Check() error {
	if err := o.Option.Check(); err != nil {
		return err
	}
	return nil
}

var _ OptionType = (*Option)(nil)
