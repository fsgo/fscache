// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/30

package freecache

import (
	"context"
	"fmt"
	"time"

	"github.com/coocood/freecache"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/internal"
)

// sCache 普通缓存
type sCache struct {
	opt   OptionType
	cache freecache.Cache
}

func (s *sCache) Get(ctx context.Context, key interface{}) fscache.GetResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return fscache.NewGetResult(nil, err, nil)
	}
	vb, err := s.cache.Get(kb)
	if err != nil {
		if err == freecache.ErrNotFound {
			return internal.GetRetNotExists
		}
		return fscache.NewGetResult(nil, err, nil)
	}
	return fscache.NewGetResult(vb, nil, s.opt.Unmarshaler())
}

func (s *sCache) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) fscache.SetResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return fscache.NewSetResult(fmt.Errorf("encode key with error:%w", err))
	}
	vb, err := s.opt.Marshaler()(value)
	if err != nil {
		return fscache.NewSetResult(fmt.Errorf("encode value with error:%w", err))
	}
	errSet := s.cache.Set(kb, vb, int(ttl.Seconds()))
	return fscache.NewSetResult(errSet)
}

func (s *sCache) Has(ctx context.Context, key interface{}) fscache.HasResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return fscache.NewHasResult(fmt.Errorf("encode key with error:%w", err), false)
	}
	_, errGet := s.cache.Get(kb)
	if errGet == nil {
		return internal.HasRetYes
	} else if errGet == freecache.ErrNotFound {
		return internal.HasRetNot
	} else {
		return fscache.NewHasResult(errGet, false)
	}
}

func (s *sCache) Delete(ctx context.Context, key interface{}) fscache.DeleteResult {
	kb, err := s.opt.Marshaler()(key)
	if err != nil {
		return fscache.NewDeleteResult(fmt.Errorf("encode key with error:%w", err), 0)
	}
	if ok := s.cache.Del(kb); ok {
		return internal.DeleteRetSucHas1
	}
	return internal.DeleteRetSucHas0

}

func (s *sCache) Reset(ctx context.Context) error {
	s.cache.Clear()
	return nil
}

var _ fscache.SCache = (*sCache)(nil)
