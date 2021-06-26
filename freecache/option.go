/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/17
 */

package freecache

import (
	"github.com/fsgo/fscache/cache"
)

const cacheFileExt = ".cache"

// IOption filecache 选项接口
type IOption interface {
	cache.IOption

	// GetMemSize
	GetMemSize() int
}

// Option 配置选型
type Option struct {
	MemSize int
	cache.Option
}

// CacheDir 缓存根目录
func (o Option) GetMemSize() int {
	return o.MemSize
}

// Check 检查是否正确
func (o Option) Check() error {
	if err := o.Option.Check(); err != nil {
		return err
	}
	return nil
}

var _ IOption = (*Option)(nil)
