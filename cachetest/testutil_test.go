/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/23
 */

package cachetest

import (
	"testing"

	"github.com/fsgo/fscache/lrucache"
)

func TestCache(t *testing.T) {
	c, _ := lrucache.New(lrucache.Option{
		Capacity: 100,
	})
	CacheTest(t, c, "test")
}
