// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

import (
	"context"
	"errors"
)

// ErrNotExists 缓存数据不存在
var ErrNotExists = errors.New("cache not exists")

// Cache 缓存API
type Cache interface {
	SCache
	MCache
}

// ReSetter 重置缓存
// 如本地文件缓存，可以实现该接口
type ReSetter interface {
	Reset(ctx context.Context) error
}
