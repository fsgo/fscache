/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/24
 */

package fscache_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fsgo/fscache"
)

func ExampleNewLRUCache() {
	cache, err := fscache.NewLRUCache(fscache.LRUOption{
		Capacity: 100,
	})
	if err != nil {
		log.Fatalf("init lru cache failed: %v", err)
	}
	cacheWriteRead(cache)

	// OutPut:
	// true
	// world
}

func ExampleNewFileCache() {
	cache, err := fscache.NewFileCache(fscache.FileOption{
		Dir: "filecache/testdata/cache_dir/",
	})
	if err != nil {
		log.Fatalf("init lru cache failed: %v", err)
	}
	cacheWriteRead(cache)

	// OutPut:
	// true
	// world
}

func cacheWriteRead(cache fscache.ICache) {
	key := "hello"
	value := "world"

	setRet := cache.Set(context.Background(), key, value, 1*time.Hour)
	if err := setRet.Err(); err != nil {
		log.Fatalf("Set has error: %v", err)
	}

	getRet := cache.Get(context.Background(), key)
	if err := getRet.Err(); err != nil {
		log.Fatalf("Get has error: %v", err)
	}
	var got string
	has, _ := getRet.Value(&got)
	fmt.Println(has)
	fmt.Println(got)
}
