/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/10
 */

package filecache

import (
	"testing"

	"github.com/fsgo/fscache/cachetest"
)

func TestNew(t *testing.T) {
	c, err := New(Option{
		Dir: "./testdata/cache_dir/",
	})
	if err != nil {
		t.Fatalf("new cache with error:%v", err)
	}
	cachetest.CacheTest(t, c, "fileCache")
}

func TestNewWithErr(t *testing.T) {
	_, err := New(Option{
		Dir: "",
	})
	if err == nil {
		t.Fatalf("new cache expect error")
	}
}
