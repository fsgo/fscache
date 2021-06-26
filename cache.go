// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

import (
	"fmt"
)

// ErrNotExists 缓存数据不存在
var ErrNotExists = fmt.Errorf("cache not exists")

// Cache 缓存API
type Cache interface {
	SCache
	MCache
}
