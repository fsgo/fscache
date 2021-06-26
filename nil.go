// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

import (
	"context"
	"time"
)

// Nil 黑洞 cache 实例
var Nil Cache = &nilCache{}

// nilCache 黑洞，可以写成功，但是读取的时候总是不存在
type nilCache struct{}

var nilGetRet = NewGetResult(nil, ErrNotExists, nil)

// Get 查询，总是返回 key 不存在
func (n *nilCache) Get(ctx context.Context, key interface{}) GetResult {
	return nilGetRet
}

var nilSetRet = NewSetResult(nil)

// Set 写入，总是返回写成功
func (n *nilCache) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) SetResult {
	return nilSetRet
}

var nilHasRet = NewHasResult(ErrNotExists, false)

// Has 判断是否存在，总是返回 key 不存在
func (n *nilCache) Has(ctx context.Context, key interface{}) HasResult {
	return nilHasRet
}

var nilDeleteRet = NewDeleteResult(nil, 0)

// Delete 删除，总是成功，返回删除 0 条
func (n *nilCache) Delete(ctx context.Context, key interface{}) DeleteResult {
	return nilDeleteRet
}

// Reset 重置
func (n *nilCache) Reset(ctx context.Context) error {
	return nil
}

// MGet 批量获取，总是返回 key 不存在
func (n *nilCache) MGet(ctx context.Context, keys []interface{}) MGetResult {
	ret := make(MGetResult, len(keys))
	for _, k := range keys {
		ret[k] = nilGetRet
	}
	return ret
}

// MSet 批量写入，总是返回写成功
func (n *nilCache) MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult {
	ret := make(MSetResult, len(kvs))
	for k := range kvs {
		ret[k] = nilSetRet
	}
	return ret
}

// MDelete 批量删除，总是删除成功
func (n *nilCache) MDelete(ctx context.Context, keys []interface{}) MDeleteResult {
	ret := make(MDeleteResult, len(keys))
	for _, k := range keys {
		ret[k] = nilDeleteRet
	}
	return ret
}

// MHas 判断是否存在，总是不存在
func (n *nilCache) MHas(ctx context.Context, keys []interface{}) MHasResult {
	ret := make(MHasResult, len(keys))
	for _, k := range keys {
		ret[k] = nilHasRet
	}
	return ret
}

var _ Cache = (*nilCache)(nil)
