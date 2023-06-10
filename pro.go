// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/6/10

package fscache

import (
	"context"
	"time"
)

type ProS[K any, V any] struct {
	// SCache 必填
	SCache SCache
}

// Get 查询单个
func (pc *ProS[K, V]) Get(ctx context.Context, key K) (value V, err error) {
	ret := pc.SCache.Get(ctx, key)
	if ret.Err != nil {
		return value, ret.Err
	}
	_, err = ret.Value(&value)
	return value, err
}

// Set 设置并附带有效期
func (pc *ProS[K, V]) Set(ctx context.Context, key K, value V, ttl time.Duration) error {
	ret := pc.SCache.Set(ctx, key, value, ttl)
	return ret.Err
}

// Has 判断是否存在
func (pc *ProS[K, V]) Has(ctx context.Context, key K) (has bool, err error) {
	ret := pc.SCache.Has(ctx, key)
	return ret.Has, ret.Err
}

// Delete 删除指定的 key
func (pc *ProS[K, V]) Delete(ctx context.Context, key K) (deleted int, err error) {
	ret := pc.SCache.Delete(ctx, key)
	return ret.Deleted, ret.Err
}
