/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/10
 */

package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// IMCache 缓存-批处理接口
type IMCache interface {
	MGet(ctx context.Context, keys []interface{}) MGetResult
	MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult
	MDelete(ctx context.Context, keys []interface{}) MDeleteResult
	MHas(ctx context.Context, keys []interface{}) MHasResult
}

// KVData  k-v pairs
type KVData map[interface{}]interface{}

// MGetResult 批量查询MGet接口的结果
// 若key不存在，是不存在
type MGetResult map[interface{}]GetResult

// MSetResult 批量设置MSet接口的结果
type MSetResult map[interface{}]SetResult

// HasError 是否有异常
func (mr MSetResult) HasError() bool {
	for _, ret := range mr {
		if err := ret.Err(); err != nil {
			return true
		}
	}
	return false
}

// MDeleteResult 批量删除MDelete接口的结果
type MDeleteResult map[interface{}]DeleteResult

// Num 删除的条数
func (md MDeleteResult) Num() int {
	var result int
	for _, ret := range md {
		result += ret.Num()
	}
	return result
}

// MHasResult 批量判断是否存在MHas接口的结果
type MHasResult map[interface{}]HasResult

// NewMCacheBySCache 创建一个MCacheBySCache实例
func NewMCacheBySCache(sCache ISCache, concurrent bool) IMCache {
	return &mCacheBySCache{
		sCache:     sCache,
		concurrent: concurrent,
	}
}

// mCacheBySCache 通过对sCache简单封装获取到的批量查询缓存实例
type mCacheBySCache struct {
	sCache     ISCache
	concurrent bool
}

func (m *mCacheBySCache) MGet(ctx context.Context, keys []interface{}) MGetResult {
	result := make(MGetResult, len(keys))
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, k := range keys {
		k := k
		wg.Add(1)
		m.runFn(func() {
			var val GetResult
			defer func() {
				if re := recover(); re != nil {
					val = NewGetResult(nil, fmt.Errorf("panic:%v", re), nil)
				}
				lock.Lock()
				result[k] = val
				lock.Unlock()

				wg.Done()
			}()
			val = m.sCache.Get(ctx, k)
		})
	}
	wg.Wait()
	return result
}

func (m *mCacheBySCache) MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult {
	result := make(MSetResult, len(kvs))
	var wg sync.WaitGroup
	var lock sync.Mutex
	for k, v := range kvs {
		wg.Add(1)
		k := k
		v := v
		m.runFn(func() {
			var val SetResult
			defer func() {
				if re := recover(); re != nil {
					val = NewSetResult(fmt.Errorf("panic:%v", re))
				}

				lock.Lock()
				result[k] = val
				lock.Unlock()

				wg.Done()
			}()
			val = m.sCache.Set(ctx, k, v, ttl)
		})
	}
	wg.Wait()
	return result
}

func (m *mCacheBySCache) MDelete(ctx context.Context, keys []interface{}) MDeleteResult {
	result := make(MDeleteResult, len(keys))
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, k := range keys {
		k := k
		wg.Add(1)
		m.runFn(func() {
			var val DeleteResult
			defer func() {
				if re := recover(); re != nil {
					val = NewDeleteResult(fmt.Errorf("panic:%v", re), 0)
				}
				lock.Lock()
				result[k] = val
				lock.Unlock()

				wg.Done()
			}()
			val = m.sCache.Delete(ctx, k)
		})
	}
	wg.Wait()
	return result
}

func (m *mCacheBySCache) MHas(ctx context.Context, keys []interface{}) MHasResult {
	result := make(MHasResult, len(keys))
	var wg sync.WaitGroup
	var lock sync.Mutex

	for _, k := range keys {
		k := k
		wg.Add(1)
		m.runFn(func() {
			var val HasResult
			defer func() {
				if re := recover(); re != nil {
					val = NewHasResult(fmt.Errorf("panic:%v", re), false)
				}
				lock.Lock()
				result[k] = val
				lock.Unlock()

				wg.Done()
			}()
			val = m.sCache.Has(ctx, k)
		})
	}
	wg.Wait()
	return result
}

func (m *mCacheBySCache) runFn(fn func()) {
	if m.concurrent {
		go fn()
	} else {
		fn()
	}
}

var _ IMCache = (*mCacheBySCache)(nil)
