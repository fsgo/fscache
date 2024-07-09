// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/6/3

package mapcache

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/fsgo/fst"
)

func TestMapCache(t *testing.T) {
	newFunc := func(ctx context.Context, key any) (any, error) {
		id := key.(int)
		if id == 1000 {
			return 0, errors.New("invalid id")
		}
		return id + 10, nil
	}

	count := func(mp *sync.Map) int {
		var num int
		mp.Range(func(key, value any) bool {
			num++
			return true
		})
		return num
	}

	check := func(t *testing.T, mc *MapCache) {
		for i := 0; i < 10000; i++ {
			val, err := mc.Get(i)
			if i == 1000 {
				fst.Error(t, err)
				fst.Equal(t, 0, val)
			} else {
				fst.NoError(t, err)
				fst.Equal(t, i+10, val.(int))
			}
		}
	}

	t.Run("No FailTTL", func(t *testing.T) {
		mc := &MapCache{
			New:     newFunc,
			Caption: 100,
		}
		check(t, mc)
		fst.LessOrEqual(t, count(&mc.values), 100)
	})

	t.Run("Has FailTTL", func(t *testing.T) {
		mc := &MapCache{
			New:     newFunc,
			FailTTL: 10 * time.Millisecond,
			Caption: 100,
		}
		check(t, mc)
		fst.LessOrEqual(t, count(&mc.values), 100)
	})
}

func BenchmarkMapCache(b *testing.B) {
	mc := &MapCache{
		New: func(ctx context.Context, key any) (any, error) {
			id := key.(int)
			if id == 1000 {
				return 0, errors.New("invalid id")
			}
			return id + 10, nil
		},
		Caption: 1000,
		FailTTL: time.Minute,
	}
	for i := 0; i < b.N; i++ {
		_, _ = mc.Get(i % 1000)
	}
}
