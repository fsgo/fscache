// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/23

package lrucache

import (
	"container/list"
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/internal"
)

// NewSCache 创建普通(非批量)
func NewSCache(opt OptionType) (fscache.SCache, error) {
	if err := opt.Check(); err != nil {
		return nil, err
	}
	sc := &SCache{
		opt: opt,
	}
	_ = sc.Reset(context.Background())
	return sc, nil
}

// SCache lru普通缓存
type SCache struct {
	opt  OptionType
	data map[interface{}]*list.Element
	list *list.List
	lock sync.Mutex
}

// Get 读取
func (L *SCache) Get(ctx context.Context, key interface{}) fscache.GetResult {
	L.lock.Lock()
	defer L.lock.Unlock()
	el, has := L.data[key]
	if !has {
		return fscache.NewGetResult([]byte("err1-"+fmt.Sprint(key)), fscache.ErrNotExists, nil)
	}
	val := el.Value.(*value)

	if val.Expired() {
		L.list.Remove(el)
		delete(L.data, key)
		return fscache.NewGetResult([]byte("err2"), fscache.ErrNotExists, nil)
	}
	L.list.MoveToFront(el)
	return fscache.NewGetResult(nil, nil, newUnmarshaler(val.Data))
}

// Set 设置
func (L *SCache) Set(ctx context.Context, key interface{}, val interface{}, ttl time.Duration) fscache.SetResult {
	cacheVal := &value{
		Key:      key,
		Data:     val,
		ExpireAt: time.Now().Add(ttl),
	}
	L.lock.Lock()
	defer L.lock.Unlock()
	el, has := L.data[key]
	if has {
		el.Value = cacheVal
		L.list.MoveToFront(el)
		return internal.SetRetSuc
	}
	elm := L.list.PushFront(cacheVal)
	L.data[key] = elm
	if L.list.Len() > L.opt.GetCapacity() {
		L.weedOut()
	}
	return internal.SetRetSuc
}

func (L *SCache) weedOut() {
	el := L.list.Back()
	if el == nil {
		return
	}
	v := el.Value.(*value)
	delete(L.data, v.Key)
	L.list.Remove(el)
}

// Has 判断是否存在
func (L *SCache) Has(ctx context.Context, key interface{}) fscache.HasResult {
	L.lock.Lock()
	el, has := L.data[key]
	L.lock.Unlock()
	if has {
		val := el.Value.(*value)
		if val.Expired() {
			L.Delete(ctx, key)
			has = false
		}
	}

	if has {
		L.lock.Lock()
		delete(L.data, key)
		L.list.Remove(el)
		L.lock.Unlock()
		return internal.HasRetYes
	}
	return internal.HasRetNot
}

// Delete 删除
func (L *SCache) Delete(ctx context.Context, key interface{}) fscache.DeleteResult {
	L.lock.Lock()
	defer L.lock.Unlock()
	el, has := L.data[key]
	if !has {
		return internal.DeleteRetSucHas0
	}
	delete(L.data, key)
	L.list.Remove(el)
	return internal.DeleteRetSucHas1
}

// Reset 重置、清空所有缓存
func (L *SCache) Reset(ctx context.Context) error {
	L.lock.Lock()
	L.data = make(map[interface{}]*list.Element, L.opt.GetCapacity())
	L.list = list.New()
	L.lock.Unlock()
	return nil
}

var _ fscache.SCache = (*SCache)(nil)

func newUnmarshaler(val interface{}) fscache.Unmarshaler {
	return func(_ []byte, obj interface{}) (err error) {
		defer func() {
			if re := recover(); re != nil {
				err = fmt.Errorf("panic:%v", re)
			}
		}()

		rv := reflect.ValueOf(obj).Elem()
		if !rv.CanSet() {
			return fmt.Errorf("cannot Unmarshal, %s cannot set", rv.String())
		}
		rv.Set(reflect.ValueOf(val))
		return nil
	}
}
