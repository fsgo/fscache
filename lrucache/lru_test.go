// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/23

package lrucache

import (
	"context"
	"testing"
	"time"

	"github.com/fsgo/fscache/cachetest"
)

func TestLRUCache(t *testing.T) {
	c, _ := New(&Option{
		Capacity: 100,
	})

	cachetest.CacheTest(t, c, "lruCache")
}

func TestLRUCache2(t *testing.T) {
	sc, _ := NewSCache(&Option{
		Capacity: 10,
	})

	for i := 0; i < 12; i++ {
		ret := sc.Set(context.Background(), i, i, 1*time.Hour)
		if err := ret.Err(); err != nil {
			t.Fatalf("set with error:%v", err)
		}
	}
}

func TestNewWithError(t *testing.T) {
	_, err := New(&Option{
		Capacity: 0,
	})
	if err == nil {
		t.Fatalf("expect has error")
	}

}
