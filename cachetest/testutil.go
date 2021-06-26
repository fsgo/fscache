// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/23

package cachetest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fsgo/fscache"
)

// CacheTest  测试缓存
func CacheTest(t *testing.T, c fscache.Cache, prex string) {
	SCacheTest(t, c, prex+"_sCache")
	MCacheTest(t, c, prex+"_mCache")
}

// SCacheTest 测试SCache
func SCacheTest(t *testing.T, c fscache.SCache, prex string) {
	kv := map[interface{}]interface{}{
		123:   234,
		234:   456,
		"789": 789,
	}
	for k, v := range kv {
		name := fmt.Sprintf("%s case:key=%v, val=%v", prex, k, v)
		t.Run(name, func(t *testing.T) {
			t.Run("Set", func(t *testing.T) {
				for i := 0; i < 5; i++ {
					retSet := c.Set(context.Background(), k, v, 10*time.Second)
					checkNoErr(t, retSet, "Set")
				}
			})

			t.Run("Get_has", func(t *testing.T) {
				retGet := c.Get(context.Background(), k)
				checkNoErr(t, retGet, "Get")

				var num int
				has, err := retGet.Value(&num)
				if err != nil {
					t.Fatalf("retGet.Value with error:%v", err)
				}
				if !has {
					t.Fatalf("expect has cache")
				}
				if num != v {
					t.Fatalf("got=%v,want=%v", num, v)
				}
			})

			t.Run("Get_miss", func(t *testing.T) {
				keyMiss := fmt.Sprintf("miss_%v", k)
				retGet := c.Get(context.Background(), keyMiss)
				checkNoErr(t, retGet, "Get")

				var num int
				has, err := retGet.Value(&num)
				if err != nil {
					t.Fatalf("retGet.Value with error:%v", err)
				}
				if has {
					t.Fatalf("expect no cache")
				}
			})

			t.Run("Has", func(t *testing.T) {
				retHas := c.Has(context.Background(), k)
				checkNoErr(t, retHas, "Has")

				if !retHas.Has() {
					t.Errorf("expect has")
				}
			})
		})
	}

	t.Run("cache_expire", func(t *testing.T) {
		key1 := "key_expire"
		key2 := "key_expire"

		// set cache
		{
			setRet := c.Set(context.Background(), key1, 0, 1*time.Second)
			checkNoErr(t, setRet, "Set")

			c.Set(context.Background(), key2, 0, 1*time.Second)
		}

		// check has
		{
			getRet := c.Get(context.Background(), key1)
			checkNoErr(t, getRet, "Get")

			var n int
			if has, _ := getRet.Value(&n); !has {
				t.Fatalf("expect has cache")
			}
		}

		time.Sleep(1 * time.Second)

		// check no
		{
			getRet := c.Get(context.Background(), key1)
			var n int
			if has, _ := getRet.Value(&n); has {
				t.Fatalf("expect no cache")
			}
		}
	})

	t.Run("Delete_miss", func(t *testing.T) {
		delRet := c.Delete(context.Background(), "not_exists")
		checkNoErr(t, delRet, "Delete_miss")
		if num := delRet.Num(); num != 0 {
			t.Errorf("Num=%d, want=0", num)
		}
	})
}

// MCacheTest 测试MCache
func MCacheTest(t *testing.T, c fscache.MCache, prex string) {
	kv := map[interface{}]interface{}{
		12345:    234,
		23456:    456,
		"789000": 789,
	}
	var keys []interface{}

	for k := range kv {
		keys = append(keys, k)
	}
	mSetRet := c.MSet(context.Background(), kv, 10*time.Second)

	if mSetRet.HasError() {
		t.Fatalf("mset has error,ret=%v", mSetRet)
	}

	t.Run(prex+"_MGET", func(t *testing.T) {
		t.Logf("mget keys=%v", keys)

		retMGet := c.MGet(context.Background(), keys)
		if len(retMGet) != len(keys) {
			t.Fatalf("result.len=%d, want=%d", len(retMGet), len(keys))
		}
		t.Logf("retMGet=%v", retMGet)
		for k, v := range kv {
			t.Run(fmt.Sprintf("case key=%v,val=%v", k, v), func(t *testing.T) {
				ret := retMGet[k]
				checkNoErr(t, ret, "Get")

				var num int

				has, err := ret.Value(&num)
				if err != nil {
					t.Fatalf("Value() error=%v", err)
				}
				if !has {
					t.Fatalf("expect has")
				}
				if num != v {
					t.Fatalf("got=%v,want=%v", num, v)
				}
			})
		}
	})

	t.Run(prex+"_MDelete", func(t *testing.T) {
		retMDel := c.MDelete(context.Background(), keys)
		got := retMDel.Num()
		want := len(keys)
		if got != want {
			t.Errorf("MDelete ret.Num()=%d want=%d", got, want)
		}
	})
}

func checkNoErr(t *testing.T, ret fscache.ResultError, msg string) {
	if err := ret.Err(); err != nil {
		t.Fatalf("%s err=", err)
	}
}
