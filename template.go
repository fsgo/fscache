// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

import (
	"context"
	"errors"
	"time"
)

// Template 缓存模板类
type Template struct {
	SCache SCache
	MCache MCache
}

// NewTemplate 利用一个简单的缓存类，创建一个包含批量接口的缓存类
func NewTemplate(sc SCache, concurrent bool) Cache {
	return &Template{
		SCache: sc,
		MCache: NewMCacheBySCache(sc, concurrent),
	}
}

// Get 读取
func (ct *Template) Get(ctx context.Context, key any) GetResult {
	return ct.SCache.Get(ctx, key)
}

// Set 写入
func (ct *Template) Set(ctx context.Context, key any, value any, ttl time.Duration) SetResult {
	return ct.SCache.Set(ctx, key, value, ttl)
}

// Has 是否存在
func (ct *Template) Has(ctx context.Context, key any) HasResult {
	return ct.SCache.Has(ctx, key)
}

// Delete  删除
func (ct *Template) Delete(ctx context.Context, key any) DeleteResult {
	return ct.SCache.Delete(ctx, key)
}

// MGet 批量获取
func (ct *Template) MGet(ctx context.Context, keys []any) MGetResult {
	return ct.MCache.MGet(ctx, keys)
}

// MSet 批量写入
func (ct *Template) MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult {
	return ct.MCache.MSet(ctx, kvs, ttl)
}

// MDelete 批量删除
func (ct *Template) MDelete(ctx context.Context, keys []any) MDeleteResult {
	return ct.MCache.MDelete(ctx, keys)
}

// MHas 批量判断是否存在
func (ct *Template) MHas(ctx context.Context, keys []any) MHasResult {
	return ct.MCache.MHas(ctx, keys)
}

// Reset 重置缓存
func (ct *Template) Reset(ctx context.Context) error {
	if rc, ok := ct.SCache.(ReSetter); ok {
		return rc.Reset(ctx)
	}
	return errors.New("not implemented ReSetter")
}

var _ Cache = (*Template)(nil)
var _ ReSetter = (*Template)(nil)
