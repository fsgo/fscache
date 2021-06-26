/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/29
 */

package cache

import (
	"context"
	"time"
)

type Nil struct{}

var nilGetRet = NewGetResult(nil, ErrNotExists, nil)

func (n Nil) Get(ctx context.Context, key interface{}) GetResult {
	return nilGetRet
}

var nilSetRet = NewSetResult(nil)

func (n Nil) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) SetResult {
	return nilSetRet
}

var nilHasRet = NewHasResult(ErrNotExists, false)

func (n Nil) Has(ctx context.Context, key interface{}) HasResult {
	return nilHasRet
}

var nilDeleteRet = NewDeleteResult(nil, 0)

func (n Nil) Delete(ctx context.Context, key interface{}) DeleteResult {
	return nilDeleteRet
}

func (n Nil) Reset(ctx context.Context) error {
	return nil
}

func (n Nil) MGet(ctx context.Context, keys []interface{}) MGetResult {
	ret := make(MGetResult, len(keys))
	for _, k := range keys {
		ret[k] = nilGetRet
	}
	return ret
}

func (n *Nil) MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult {
	ret := make(MSetResult, len(kvs))
	for k := range kvs {
		ret[k] = nilSetRet
	}
	return ret
}

func (n *Nil) MDelete(ctx context.Context, keys []interface{}) MDeleteResult {
	ret := make(MDeleteResult, len(keys))
	for _, k := range keys {
		ret[k] = nilDeleteRet
	}
	return ret
}

func (n *Nil) MHas(ctx context.Context, keys []interface{}) MHasResult {
	ret := make(MHasResult, len(keys))
	for _, k := range keys {
		ret[k] = nilHasRet
	}
	return ret
}

var _ ICache = (*Nil)(nil)
