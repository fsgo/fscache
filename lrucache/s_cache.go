/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/23
 */

package lrucache

import (
	"container/list"
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/fsgo/fscache/cache"
)

func NewSCache(capacity int) (cache.ISCache, error) {
	sc := &SCache{
		capacity: capacity,
	}
	sc.Reset()
	return sc, nil
}

// SCache lru普通缓存
type SCache struct {
	capacity int
	data     map[interface{}]*list.Element
	list     *list.List
	lock     sync.Mutex
}

// Get 读取
func (L *SCache) Get(ctx context.Context, key interface{}) cache.GetResult {
	L.lock.Lock()
	defer L.lock.Unlock()
	el, has := L.data[key]
	if !has {
		return cache.NewGetResult([]byte("err1-"+fmt.Sprint(key)), cache.ErrNotExists, nil)
	}
	val := el.Value.(*value)

	if val.Expired() {
		L.list.Remove(el)
		delete(L.data, key)
		return cache.NewGetResult([]byte("err2"), cache.ErrNotExists, nil)
	}
	L.list.MoveToFront(el)
	return cache.NewGetResult(nil, nil, newUnmarshaler(val.Data))
}

// Set 设置
func (L *SCache) Set(ctx context.Context, key interface{}, val interface{}, ttl time.Duration) cache.SetResult {
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
		return cache.NewSetResult(nil)
	}
	elm := L.list.PushFront(cacheVal)
	L.data[key] = elm
	if L.list.Len() > L.capacity {
		L.weedOut()
	}
	return cache.NewSetResult(nil)
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
func (L *SCache) Has(ctx context.Context, key interface{}) cache.HasResult {
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
	}
	return cache.NewHasResult(nil, has)
}

// Delete 删除
func (L *SCache) Delete(ctx context.Context, key interface{}) cache.DeleteResult {
	L.lock.Lock()
	defer L.lock.Unlock()
	el, has := L.data[key]
	if !has {
		return cache.NewDeleteResult(nil, 0)
	}
	delete(L.data, key)
	L.list.Remove(el)
	return cache.NewDeleteResult(nil, 1)
}

// Reset 重置、清空所有缓存
func (L *SCache) Reset() {
	L.lock.Lock()
	defer L.lock.Unlock()
	L.data = make(map[interface{}]*list.Element, L.capacity)
	L.list = list.New()
}

var _ cache.ISCache = (*SCache)(nil)

func newUnmarshaler(val interface{}) cache.Unmarshaler {
	return func(_ []byte, obj interface{}) error {
		rv := reflect.ValueOf(obj).Elem()
		if !rv.CanSet() {
			return fmt.Errorf("cannot Unmarshal, %s cannot set", rv.String())
		}
		rv.Set(reflect.ValueOf(val))
		return nil
	}
}
