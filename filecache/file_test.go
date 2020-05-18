/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/10
 */

package filecache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c, err := New(Option{
		Dir: "./testdata/cache_dir/",
	})
	if err != nil {
		t.Fatalf("new cache with error:%v", err)
	}
	kv := map[interface{}]interface{}{
		123:   234,
		234:   456,
		"789": 789,
	}
	var keys []interface{}
	for k, v := range kv {
		keys = append(keys, k)
		t.Run(fmt.Sprintf("case:key=%v, val=%v", k, v), func(t *testing.T) {
			retSet := c.Set(context.Background(), k, v, 10*time.Second)
			if err := retSet.Err(); err != nil {
				t.Fatalf("Set with error:%v", err)
			}

			retGet := c.Get(context.Background(), k)
			if err := retGet.Err(); err != nil {
				t.Fatalf("Get with error:%v", err)
			}
			var num int
			if _, err := retGet.Value(&num); err != nil {
				t.Fatalf("retGet.Value with error:%v", err)
			}
			if num != v {
				t.Fatalf("got=%v,want=%v", num, v)
			}
			retHas := c.Has(context.Background(), k)

			if err := retHas.Err(); err != nil {
				t.Errorf("Has with error:%v", err)
			}
			if !retHas.Has() {
				t.Errorf("expect has")
			}
		})
	}

	t.Run("MGET", func(t *testing.T) {
		retMGet := c.MGet(context.Background(), keys)
		if len(retMGet) != len(keys) {
			t.Fatalf("result.len=%d, want=%d", len(retMGet), len(keys))
		}
		for k, v := range kv {
			ret := retMGet[k]
			if err := ret.Err(); err != nil {
				t.Errorf("has error:%v", err)
			}
			var num int
			if _, err := ret.Value(&num); err != nil {
				t.Errorf("Value() error=%v", err)
			}
			if num != v {
				t.Fatalf("got=%v,want=%v", num, v)
			}
		}
	})

	t.Run("MDelete", func(t *testing.T) {
		retMDel := c.MDelete(context.Background(), keys)
		got := retMDel.Num()
		want := len(keys)
		if got != want {
			t.Errorf("MDelete ret.Num()=%d want=%d", got, want)
		}
	})
}
