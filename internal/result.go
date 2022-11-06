// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/27

package internal

import (
	"github.com/fsgo/fscache"
)

// GetRetNotExists get key 不存在
var GetRetNotExists = fscache.GetResult{Err: fscache.ErrNotExists}

// SetRetSuc set 成功
var SetRetSuc = fscache.SetResult{}

// DeleteRetSucHas0 delete 成功，删除0条
var DeleteRetSucHas0 = fscache.DeleteResult{}

// DeleteRetSucHas1 delete 成功，删除1条
var DeleteRetSucHas1 = fscache.DeleteResult{Deleted: 1}

// HasRetNot Has 成功判断，不存在
var HasRetNot = fscache.HasResult{Err: fscache.ErrNotExists}

// HasRetYes Has 成功判断，存在
var HasRetYes = fscache.HasResult{Has: true}
