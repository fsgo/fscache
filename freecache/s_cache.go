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

// NewSCache 创建普通的缓存实例
func NewSCache(opt *Option) (fscache.SCache, error) {
	c := freecache.NewCache(opt.GetMemSize())
	codec := opt.GetCodec()
	return &sCache{
		opt:    opt,
		cache:  c,
		encode: codec.Marshal,
		decode: codec.Unmarshal,
	}, nil
}

// sCache 普通缓存
type sCache struct {
	opt    *Option
	cache  *freecache.Cache
	decode fscache.UnmarshalFunc
	encode fscache.MarshalFunc
}

func (s *sCache) Get(ctx context.Context, key any) fscache.GetResult {
	kb, err := s.opt.GetCodec().Marshal(key)
	if err != nil {
		return fscache.GetResult{Err: err}
	}
	vb, err := s.cache.Get(kb)
	if err != nil {
		if err == freecache.ErrNotFound {
			return internal.GetRetNotExists
		}
		return fscache.GetResult{Err: err}
	}
	return fscache.GetResult{Payload: vb, UnmarshalFunc: s.decode}
}

func (s *sCache) Set(ctx context.Context, key any, value any, ttl time.Duration) fscache.SetResult {
	kb, err := s.encode(key)
	if err != nil {
		return fscache.SetResult{Err: fmt.Errorf("encode key with error:%w", err)}
	}
	vb, err := s.encode(value)
	if err != nil {
		return fscache.SetResult{Err: fmt.Errorf("encode value with error:%w", err)}
	}
	errSet := s.cache.Set(kb, vb, int(ttl.Seconds()))
	return fscache.SetResult{Err: errSet}
}

func (s *sCache) Has(ctx context.Context, key any) fscache.HasResult {
	kb, err := s.encode(key)
	if err != nil {
		return fscache.HasResult{Err: fmt.Errorf("encode key with error:%w", err)}
	}
	_, errGet := s.cache.Get(kb)
	if errGet == nil {
		return internal.HasRetYes
	} else if errGet == freecache.ErrNotFound {
		return internal.HasRetNot
	} else {
		return fscache.HasResult{Err: errGet}
	}
}

func (s *sCache) Delete(ctx context.Context, key any) fscache.DeleteResult {
	kb, err := s.encode(key)
	if err != nil {
		return fscache.DeleteResult{Err: fmt.Errorf("encode key with error:%w", err)}
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
