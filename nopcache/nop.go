// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package nopcache

import (
	"context"
	"time"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/internal"
)

// Nop 黑洞 cache 实例
var Nop fscache.Cache = &nopCache{}

// nopCache 黑洞，可以写成功，但是读取的时候总是不存在
type nopCache struct{}

// Get 查询，总是返回 key 不存在
func (n *nopCache) Get(ctx context.Context, key any) fscache.GetResult {
	return internal.GetRetNotExists
}

// Set 写入，总是返回写成功
func (n *nopCache) Set(ctx context.Context, key any, value any, ttl time.Duration) fscache.SetResult {
	return internal.SetRetSuc
}

// Has 判断是否存在，总是返回 key 不存在
func (n *nopCache) Has(ctx context.Context, key any) fscache.HasResult {
	return internal.HasRetNot
}

// Delete 删除，总是成功，返回删除 0 条
func (n *nopCache) Delete(ctx context.Context, key any) fscache.DeleteResult {
	return internal.DeleteRetSucHas0
}

// Reset 重置
func (n *nopCache) Reset(ctx context.Context) error {
	return nil
}

// MGet 批量获取，总是返回 key 不存在
func (n *nopCache) MGet(ctx context.Context, keys []any) fscache.MGetResult {
	return nil
}

// MSet 批量写入，总是返回写成功
func (n *nopCache) MSet(ctx context.Context, kvs fscache.KVData, ttl time.Duration) fscache.MSetResult {
	return nil
}

// MDelete 批量删除，总是删除成功
func (n *nopCache) MDelete(ctx context.Context, keys []any) fscache.MDeleteResult {
	return nil
}

// MHas 判断是否存在，总是不存在
func (n *nopCache) MHas(ctx context.Context, keys []any) fscache.MHasResult {
	return nil
}

var _ fscache.Cache = (*nopCache)(nil)
