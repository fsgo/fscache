// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/30

package chains

import (
	"context"
	"time"

	"github.com/fsgo/fscache"
)

// SetTTLFn 设置缓存的 ttl,参数 ttl 可能为空
type SetTTLFn func(ttl time.Duration) time.Duration

// New 创建一个链式缓存
func New(caches ...*Cache) fscache.Cache {
	if len(caches) == 0 {
		panic("no caches")
	}
	sc := &sChains{
		caches: caches,
	}
	return fscache.NewTemplate(sc, true)
}

type sChains struct {
	caches []*Cache
}

// Cache New 的参数
type Cache struct {
	Cache    fscache.SCache
	SetTTLFn SetTTLFn
}

func (c *Cache) getTTL(ttl time.Duration) time.Duration {
	if c.SetTTLFn == nil {
		return ttl
	}
	return c.SetTTLFn(ttl)
}

func (c *sChains) Get(ctx context.Context, key any) (result fscache.GetResult) {
	for i := 0; i < len(c.caches); i++ {
		subCache := c.caches[i]
		if result = subCache.Cache.Get(ctx, key); result.Has() {
			break
		}
	}
	// todo 将从后面cache 查询的结果赋值到前面的 cache 中
	return result
}

func (c *sChains) Set(ctx context.Context, key any, value any, ttl time.Duration) (result fscache.SetResult) {
	for i := 0; i < len(c.caches); i++ {
		subCache := c.caches[i]
		result = subCache.Cache.Set(ctx, key, value, subCache.getTTL(ttl))
	}
	return result
}

func (c *sChains) Has(ctx context.Context, key any) (result fscache.HasResult) {
	for i := 0; i < len(c.caches); i++ {
		result = c.caches[i].Cache.Has(ctx, key)
		if result.Has() {
			return result
		}
	}
	return result
}

func (c *sChains) Delete(ctx context.Context, key any) (result fscache.DeleteResult) {
	for i := 0; i < len(c.caches); i++ {
		subCache := c.caches[i]
		result = subCache.Cache.Delete(ctx, key)
	}
	return
}

func (c *sChains) Reset(ctx context.Context) (err error) {
	for i := 0; i < len(c.caches); i++ {
		subCache := c.caches[i]
		if sc, ok := subCache.Cache.(fscache.Reseter); ok {
			if e := sc.Reset(ctx); e != nil {
				err = e
			}
		}
	}
	return err
}

var _ fscache.SCache = (*sChains)(nil)
var _ fscache.Reseter = (*sChains)(nil)
