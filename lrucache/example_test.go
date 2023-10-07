// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/10/7

package lrucache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/fsgo/fscache/lrucache"
)

func ExampleNew() {
	cache, err := lrucache.New(&lrucache.Option{Capacity: 1000})
	fmt.Println("err_is_nil=", err == nil)

	ret := cache.Set(context.Background(), "k1", "v1", time.Second)
	fmt.Println("set_success=", ret.Err == nil)

	// Output:
	// err_is_nil= true
	// set_success= true
}
