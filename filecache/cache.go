// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/10

package filecache

import (
	"github.com/fsgo/fscache"
)

// New 创建新缓存实例
func New(opt *Option) (fscache.Cache, error) {
	sc, err := NewSCache(opt)
	if err != nil {
		return nil, err
	}
	return fscache.NewTemplate(sc, false), nil
}
