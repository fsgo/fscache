// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package internal

import (
	"github.com/fsgo/fscache"
)

// GetRetNotExists get key 不存在
var GetRetNotExists = fscache.NewGetResult(nil, fscache.ErrNotExists, nil)

// SetRetSuc set 成功
var SetRetSuc = fscache.NewSetResult(nil)

// DeleteRetSucHas0 delete 成功，删除0条
var DeleteRetSucHas0 = fscache.NewDeleteResult(nil, 0)

// DeleteRetSucHas1 delete 成功，删除1条
var DeleteRetSucHas1 = fscache.NewDeleteResult(nil, 1)

// HasRetNot Has 成功判断，不存在
var HasRetNot = fscache.NewHasResult(fscache.ErrNotExists, false)

// HasRetYes Has 成功判断，存在
var HasRetYes = fscache.NewHasResult(nil, true)
