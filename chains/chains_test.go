// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/7/9

package chains

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/lrucache"
)

func Test_sChains(t *testing.T) {
	lc1, _ := lrucache.New(&lrucache.Option{
		Capacity: 100,
	})
	c1 := &Cache{
		Cache: lc1,
	}
	lc2, _ := lrucache.New(&lrucache.Option{
		Capacity: 100,
	})

	c2 := &Cache{
		Cache: lc2,
	}
	ctx := context.Background()
	cc1 := New(c1, c2)
	key := "abc"
	value := "hello"
	got1 := cc1.Get(ctx, key)
	require.Error(t, got1.Err)

	ret2 := lc2.Set(ctx, key, value, time.Second)
	require.NoError(t, ret2.Err)

	checkHas := func(t *testing.T, got fscache.GetResult, want string) {
		var v2 string
		has, err := got.Value(&v2)
		require.True(t, has)
		require.NoError(t, err)
		require.Equal(t, want, v2)
	}

	checkHas(t, lc2.Get(ctx, key), value)

	checkHas(t, cc1.Get(ctx, key), value)

	ret3 := cc1.Delete(ctx, key)
	require.Equal(t, 1, ret3.Deleted)

	key2 := "world"
	got3 := cc1.Get(ctx, key2)
	require.Error(t, got3.Err)

	ret4 := cc1.Set(ctx, key2, key2, time.Second)
	require.NoError(t, ret4.Err)

	checkHas(t, cc1.Get(ctx, key2), key2)
	checkHas(t, lc1.Get(ctx, key2), key2)
	checkHas(t, lc2.Get(ctx, key2), key2)
}
