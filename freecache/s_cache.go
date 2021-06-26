/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/30
 */

package freecache

import (
	"context"
	"fmt"
	"time"

	"github.com/coocood/freecache"

	"github.com/fsgo/fscache/cache"
)

type SCache struct {
	opt   IOption
	cache freecache.Cache
}

func (s *SCache) Get(ctx context.Context, key interface{}) cache.GetResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return cache.NewGetResult(nil, err, nil)
	}
	if vb, err := s.cache.Get(kb); err != nil {
		if err == freecache.ErrNotFound {
			return cache.NewGetResult(nil, cache.ErrNotExists, nil)
		}
		return cache.NewGetResult(nil, err, nil)
	} else {
		return cache.NewGetResult(vb, nil, s.opt.Unmarshaler())
	}
}

func (s *SCache) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) cache.SetResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return cache.NewSetResult(fmt.Errorf("encode key with error:%w", err))
	}
	vb, err := s.opt.Marshaler()(value)
	if err != nil {
		return cache.NewSetResult(fmt.Errorf("encode value with error:%w", err))
	}
	errSet := s.cache.Set(kb, vb, int(ttl.Seconds()))
	return cache.NewSetResult(errSet)
}

func (s *SCache) Has(ctx context.Context, key interface{}) cache.HasResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return cache.NewHasResult(fmt.Errorf("encode key with error:%w", err), false)
	}
	_, errGet := s.cache.Get(kb)
	if errGet == nil {
		return cache.NewHasResult(nil, true)
	} else if errGet == freecache.ErrNotFound {
		return cache.NewHasResult(cache.ErrNotExists, false)
	} else {
		return cache.NewHasResult(errGet, false)
	}
}

func (s *SCache) Delete(ctx context.Context, key interface{}) cache.DeleteResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return cache.NewDeleteResult(fmt.Errorf("encode key with error:%w", err), 0)
	}
	if ok := s.cache.Del(kb); ok {
		return cache.NewDeleteResult(nil, 1)
	}
	return cache.NewDeleteResult(nil, 0)

}

func (s *SCache) Reset(ctx context.Context) error {
	panic("implement me")
}

var _ cache.ISCache = (*SCache)(nil)
