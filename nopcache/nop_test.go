// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package nopcache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fsgo/fscache/internal"
)

func TestNop(t *testing.T) {
	key := "key"
	value := "value"
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		got := Nop.Get(ctx, key)
		require.Equal(t, internal.GetRetNotExists, got)
	})

	t.Run("Set", func(t *testing.T) {
		got := Nop.Set(ctx, key, value, time.Second)
		require.Equal(t, internal.SetRetSuc, got)
	})

	t.Run("Delete", func(t *testing.T) {
		got := Nop.Delete(ctx, key)
		require.Equal(t, internal.DeleteRetSucHas0, got)
	})
}

func Test_nopCache_MGet(t *testing.T) {
	ret := Nop.MGet(context.Background(), []any{"abc", "def"})
	got := ret.Get("abc")
	require.Equal(t, internal.GetRetNotExists, got)
}
