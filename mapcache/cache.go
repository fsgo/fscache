// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/6/2

package mapcache

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// MapCache 一个简单的，使用 sync.Map 作为存储的缓存
type MapCache struct {
	// New 创建新值的函数,
	New func(ctx context.Context, key any) (any, error)

	// TTL 缓存有效期，当为 0 时，默认值为 1 分钟
	TTL time.Duration

	// FailTTL 当 New 方法创建对象失败的时候，缓存的有效期，默认为 0
	FailTTL time.Duration

	// Caption 容量，当为 0 是，默认值为 100000
	// 这个值是一个近似值
	Caption int64

	count int64

	values sync.Map
}

func (mc *MapCache) getCaption() int64 {
	if mc.Caption == 0 {
		return 100000
	}
	return mc.Caption
}

func (mc *MapCache) getTTL() time.Duration {
	if mc.TTL == 0 {
		return time.Minute
	}
	return mc.TTL
}

// Get 读取一个值
func (mc *MapCache) Get(key any) (any, error) {
	return mc.GetContext(context.Background(), key)
}

// GetContext 读取一个值
func (mc *MapCache) GetContext(ctx context.Context, key any) (any, error) {
	cv, has := mc.values.Load(key)
	if has {
		if vv := cv.(*value); vv.IsOK() {
			return vv.payload, vv.err
		}
	}
	nv, err := mc.New(ctx, key)
	if err == nil {
		mc.store(key, nv, nil, mc.getTTL())
	} else {
		if mc.FailTTL > 0 {
			mc.store(key, nv, err, mc.FailTTL)
		}
	}
	return nv, err
}

func (mc *MapCache) store(key any, nv any, err error, ttl time.Duration) {
	cv := &value{
		payload: nv,
		err:     err,
		expired: time.Now().Add(ttl),
	}
	_, hasOld := mc.values.LoadOrStore(key, cv)
	if hasOld {
		return
	}
	num := atomic.AddInt64(&mc.count, 1)
	if del := num - mc.getCaption(); del > 0 {
		mc.clear(key, int(del))
	}
}

func (mc *MapCache) clear(notKey any, needDel int) {
	delKeys := make([]any, 0, 5)
	var loop int
	mc.values.Range(func(k, v any) bool {
		if k == notKey {
			return true
		}
		loop++

		// 当超出 Caption 10个以上的时候，直接删除一些
		if needDel > 10 && loop < 5 {
			delKeys = append(delKeys, k)
			return true
		}

		if loop >= 5 {
			// 当查找了几次没有找到过期数据的时候，直接删除一项
			delKeys = append(delKeys, k)
			return false
		}

		cv := v.(*value)
		if cv.IsOK() {
			return true
		}
		delKeys = append(delKeys, k)
		return false
	})

	for i := 0; i < len(delKeys); i++ {
		mc.Delete(delKeys[i])
	}
}

// Delete 删除值
func (mc *MapCache) Delete(key any) int {
	_, ok := mc.values.LoadAndDelete(key)
	if ok {
		atomic.AddInt64(&mc.count, -1)
		return 1
	}
	return 0
}

type value struct {
	index   int
	payload any
	err     error
	expired time.Time
}

func (v *value) IsOK() bool {
	return time.Now().Before(v.expired)
}
