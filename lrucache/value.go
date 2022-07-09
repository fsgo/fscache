// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/23

package lrucache

import (
	"time"
)

type value struct {
	Key      any
	Data     any
	ExpireAt time.Time
}

// Expired 是否已过期
func (v *value) Expired() bool {
	return time.Now().After(v.ExpireAt)
}
