// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package fscache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMGetResult_Get(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var mr MGetResult
		got := mr.Get("key")
		require.Equal(t, getRetNotExists, got)
	})
}
