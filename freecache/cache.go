// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/19

package freecache

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
