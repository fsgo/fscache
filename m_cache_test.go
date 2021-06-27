// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package fscache

import (
	"testing"
)

func TestMGetResult_Get(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var mr MGetResult
		if got := mr.Get("key"); got != getRetNotExists {
			t.Fatalf("not eq")
		}
	})
}
