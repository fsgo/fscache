// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/19

package fsfreecache

import (
	"testing"

	"github.com/fsgo/fscache/cachetest"
)

func TestNew(t *testing.T) {
	c, err := New(&Option{})
	if err != nil {
		t.Fatalf("new cache with error:%v", err)
	}
	cachetest.CacheTest(t, c, "freeCache")
}
