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
	fscache.Option

	// MemSize 内存大小，最小值 512 * 1024
	// 若为 0，使用默认值 8*1024*1024
	MemSize int
}

const defaultSize = 8 * 1024 * 1024

// GetMemSize 获取配置的内存大小
func (o *Option) GetMemSize() int {
	if o.MemSize == 0 {
		return defaultSize
	}
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
