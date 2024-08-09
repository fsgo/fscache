// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/17

package fsfreecache

import (
	"github.com/fsgo/fscache"
)

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
