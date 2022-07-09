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

type chanCtxKey uint8

const (
	chanCtxKeyGet chanCtxKey = iota
)

// WithValueFunc 给 context 中补充 生成空 value 的方法
//
// 当 Get 方法，从更低级的缓存中读取出来后，使用这个值解析，然后写入到更高一级的缓存中去
func WithValueFunc(ctx context.Context, fn func() any) context.Context {
	return context.WithValue(ctx, chanCtxKeyGet, fn)
}

func (c *sChains) Get(ctx context.Context, key any) (result fscache.GetResult) {
	var index int
	for i, subCache := range c.caches {
		if result = subCache.Cache.Get(ctx, key); result.Has() {
			index = i
			break
		}
	}

	if index == 0 {
		return result
	}

	vk := ctx.Value(chanCtxKeyGet)

	if vk != nil {
		vkv := vk.(func() any)()
		if _, err := result.Value(&vkv); err != nil {
			return result
		}
		for i := 0; i < index; i++ {
			subCache := c.caches[i].Cache
			ttlFn := c.caches[i].SetTTLFn
			subCache.Set(ctx, key, vkv, ttlFn(0))
		}
	}
	return result
}

func (c *sChains) Set(ctx context.Context, key any, value any, ttl time.Duration) (result fscache.SetResult) {
	for _, subCache := range c.caches {
		result = subCache.Cache.Set(ctx, key, value, subCache.SetTTLFn(ttl))
	}
	return
}

func (c *sChains) Has(ctx context.Context, key any) (result fscache.HasResult) {
	for _, subCache := range c.caches {
		if result = subCache.Cache.Has(ctx, key); result.Has() {
			return
		}
	}
	return
}

func (c *sChains) Delete(ctx context.Context, key any) (result fscache.DeleteResult) {
	for _, subCache := range c.caches {
		result = subCache.Cache.Delete(ctx, key)
	}
	return
}

func (c *sChains) Reset(ctx context.Context) (err error) {
	for _, subCache := range c.caches {
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
