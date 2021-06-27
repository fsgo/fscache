// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package nopcache

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/fsgo/fscache/internal"
)

func TestNop(t *testing.T) {
	key := "key"
	value := "value"
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		got := Nop.Get(ctx, key)
		if got != internal.GetRetNotExists {
			t.Fatalf("not eq")
		}
	})

	t.Run("Set", func(t *testing.T) {
		got := Nop.Set(ctx, key, value, time.Second)
		if got != internal.SetRetSuc {
			t.Fatalf("not eq,got=%v want=%v", got, internal.SetRetSuc)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		got := Nop.Delete(ctx, key)
		if got != internal.DeleteRetSucHas0 {
			t.Fatalf("not eq,got=%v want=%v", got, internal.DeleteRetSucHas0)
		}
	})
}

func Test_nopCache_MGet(t *testing.T) {
	ret := Nop.MGet(context.Background(), []interface{}{"abc", "def"})
	if got := ret.Get("abc"); !reflect.DeepEqual(got, internal.GetRetNotExists) {
		t.Fatalf("not eq")
	}
}
