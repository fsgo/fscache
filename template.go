// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

import (
	"context"
	"time"
)

// Template 缓存模板类
type Template struct {
	SCache SCache
	MCache MCache
}

// NewTemplate 利用一个简单的缓存类，创建一个包含批量接口的缓存类
func NewTemplate(sCache SCache, concurrent bool) Cache {
	return &Template{
		SCache: sCache,
		MCache: NewMCacheBySCache(sCache, concurrent),
	}
}

// Get 读取
func (c *Template) Get(ctx context.Context, key interface{}) GetResult {
	return c.SCache.Get(ctx, key)
}

// Set 写入
func (c *Template) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) SetResult {
	return c.SCache.Set(ctx, key, value, ttl)
}

// Has 是否存在
func (c *Template) Has(ctx context.Context, key interface{}) HasResult {
	return c.SCache.Has(ctx, key)
}

// Delete  删除
func (c *Template) Delete(ctx context.Context, key interface{}) DeleteResult {
	return c.SCache.Delete(ctx, key)
}

// MGet 批量获取
func (c *Template) MGet(ctx context.Context, keys []interface{}) MGetResult {
	return c.MCache.MGet(ctx, keys)
}

// MSet 批量写入
func (c *Template) MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult {
	return c.MCache.MSet(ctx, kvs, ttl)
}

// MDelete 批量删除
func (c Template) MDelete(ctx context.Context, keys []interface{}) MDeleteResult {
	return c.MCache.MDelete(ctx, keys)
}

// MHas 批量判断是否存在
func (c *Template) MHas(ctx context.Context, keys []interface{}) MHasResult {
	return c.MCache.MHas(ctx, keys)
}

// Reset 重置缓存
func (c *Template) Reset(ctx context.Context) error {
	return c.SCache.Reset(ctx)
}

var _ Cache = (*Template)(nil)
