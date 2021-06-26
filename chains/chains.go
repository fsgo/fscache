/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/30
 */

package chains

import (
	"context"
	"time"

	"github.com/fsgo/fscache/cache"
)

type SetTTLFn func(ttl time.Duration) time.Duration

type IChains interface {
	cache.ICache
	AddCache(cache cache.ICache, ttlFn SetTTLFn)
}

type sChains struct {
	caches []chainsCache
}

type chainsCache struct {
	cache    cache.ISCache
	setTTLFn SetTTLFn
}

func (c *sChains) Get(ctx context.Context, key interface{}) (result cache.GetResult) {
	var index int
	for i, subCache := range c.caches {
		if result = subCache.cache.Get(ctx, key); result.Has() {
			index = i
			break
		}
	}
	for i := 0; i < index; i++ {
		subCache := c.caches[i].cache
		subCache.Set(ctx, key, nil, 0)
	}
	return result
}

func (c *sChains) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) (result cache.SetResult) {
	for _, subCache := range c.caches {
		result = subCache.cache.Set(ctx, key, value, subCache.setTTLFn(ttl))
	}
	return
}

func (c *sChains) Has(ctx context.Context, key interface{}) (result cache.HasResult) {
	for _, subCache := range c.caches {
		if result = subCache.cache.Has(ctx, key); result.Has() {
			return
		}
	}
	return
}

func (c *sChains) Delete(ctx context.Context, key interface{}) (result cache.DeleteResult) {
	for _, subCache := range c.caches {
		result = subCache.cache.Delete(ctx, key)
	}
	return
}

func (c *sChains) Reset(ctx context.Context) (err error) {
	for _, subCache := range c.caches {
		err = subCache.cache.Reset(ctx)
	}
	return
}

var _ cache.ISCache = (*sChains)(nil)
