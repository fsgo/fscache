// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/24

package lrucache

import (
	"fmt"
)

// Option LRU缓存的配置
type Option struct {
	Capacity int
}

// Check 检查配置是否正常
func (o *Option) Check() error {
	if o.Capacity < 1 {
		return fmt.Errorf("option.Capacity=%d, expect >= 1", o.Capacity)
	}
	return nil
}

// GetCapacity 获取容量
func (o *Option) GetCapacity() int {
	return o.Capacity
}
