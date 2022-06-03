// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/6/3

package mapcache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMapCache(t *testing.T) {
	newFunc := func(ctx context.Context, key interface{}) (interface{}, error) {
		id := key.(int)
		if id == 1000 {
			return 0, fmt.Errorf("invalid id")
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
				require.Error(t, err)
				require.Equal(t, 0, val)
			} else {
				require.NoError(t, err)
				require.Equal(t, i+10, val)
			}
		}
	}

	t.Run("No FailTTL", func(t *testing.T) {
		mc := &MapCache{
			New:     newFunc,
			Caption: 100,
		}
		check(t, mc)
		require.LessOrEqual(t, count(&mc.values), 100)
	})

	t.Run("Has FailTTL", func(t *testing.T) {
		mc := &MapCache{
			New:     newFunc,
			FailTTL: 10 * time.Millisecond,
			Caption: 100,
		}
		check(t, mc)
		require.LessOrEqual(t, count(&mc.values), 100)
	})
}

func BenchmarkMapCache(b *testing.B) {
	mc := &MapCache{
		New: func(ctx context.Context, key any) (any, error) {
			id := key.(int)
			if id == 1000 {
				return 0, fmt.Errorf("invalid id")
			}
			return id + 10, nil
		},
		Caption: 200,
		FailTTL: time.Minute,
	}
	for i := 0; i < b.N; i++ {
		_, _ = mc.Get(i % 1000)
	}
}
