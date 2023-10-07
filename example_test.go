// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/10/7

package fscache_test

import (
	"context"
	"fmt"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/nopcache"
)

func ExampleProS_Get() {
	ps := &fscache.ProS[string, string]{
		SCache: nopcache.Nop,
	}
	got, err := ps.Get(context.Background(), "hello")
	fmt.Println("got=", got, ", err=", err)
	// Output:
	// got=  , err= cache not exists
}
