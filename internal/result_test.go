// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResultVars(t *testing.T) {
	t.Run("GetRetNotExists", func(t *testing.T) {
		require.Error(t, GetRetNotExists.Err)
	})

	t.Run("SetRetSuc", func(t *testing.T) {
		require.NoError(t, SetRetSuc.Err)
	})
}
