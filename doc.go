/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/23
 */

package fscache

/*
import(
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fsgo/fscache"
)

func main(){
	opt:=fscache.LRUOption{
		Capacity: 100,
	}
	cache, err := fscache.NewLRUCache(opt)
	if err != nil {
		log.Fatalf("init lru cache failed: %v", err)
	}

	key := "hello"
	value := "world"

	// 写缓存
	setRet := cache.Set(context.Background(), key, value, 1*time.Hour)
	if err := setRet.Err(); err != nil {
		log.Fatalf("Set has error: %v", err)
	}

	// 读取缓存
	getRet := cache.Get(context.Background(), key)
	if err := getRet.Err(); err != nil {
		log.Fatalf("Get has error: %v", err)
	}
	var got string
	has, _ := getRet.Value(&got)
	fmt.Println(has)
	fmt.Println(got)
}
*/
