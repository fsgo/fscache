/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/23
 */

package lrucache

import (
	"github.com/fsgo/fscache/cache"
)

// New 创建新的lru缓存实例
func New(opt IOption) (cache.ICache, error) {
	sc, err := NewSCache(opt)
	if err != nil {
		return nil, err
	}
	return cache.NewTemplate(sc, false), nil
}
