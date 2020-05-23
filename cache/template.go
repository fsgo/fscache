/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/23
 */

package cache

import (
	"context"
	"time"
)

// Template 缓存模板
type Template struct {
	SCache ISCache
	MCache IMCache
}

func NewTemplate(sCache ISCache, concurrent bool) ICache {
	return &Template{
		SCache: sCache,
		MCache: NewMCacheBySCache(sCache, concurrent),
	}
}

func (c *Template) Get(ctx context.Context, key interface{}) GetResult {
	return c.SCache.Get(ctx, key)
}

func (c *Template) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) SetResult {
	return c.SCache.Set(ctx, key, value, ttl)
}

func (c *Template) Has(ctx context.Context, key interface{}) HasResult {
	return c.SCache.Has(ctx, key)
}

func (c *Template) Delete(ctx context.Context, key interface{}) DeleteResult {
	return c.SCache.Delete(ctx, key)
}

func (c *Template) MGet(ctx context.Context, keys []interface{}) MGetResult {
	return c.MCache.MGet(ctx, keys)
}

func (c *Template) MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult {
	return c.MCache.MSet(ctx, kvs, ttl)
}

func (c Template) MDelete(ctx context.Context, keys []interface{}) MDeleteResult {
	return c.MCache.MDelete(ctx, keys)
}

func (c *Template) MHas(ctx context.Context, keys []interface{}) MHasResult {
	return c.MCache.MHas(ctx, keys)
}

var _ ICache = (*Template)(nil)
