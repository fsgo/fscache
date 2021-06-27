// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package internal

import (
	"testing"
)

func TestResultVars(t *testing.T) {
	t.Run("GetRetNotExists", func(t *testing.T) {
		if err := GetRetNotExists.Err(); err != nil {
			t.Fatalf("got.Err()=%v want=%v", err, nil)
		}
		if has := GetRetNotExists.Has(); has {
			t.Fatalf("expect not")
		}
	})

	t.Run("SetRetSuc", func(t *testing.T) {
		if got := SetRetSuc.Err(); got != nil {
			t.Fatalf("got=%v want nil", got)
		}
	})
}
