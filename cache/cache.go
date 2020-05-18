/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/10
 */

package cache

import (
	"fmt"
)

// ErrNotExists 缓存数据不存在
var ErrNotExists = fmt.Errorf("cache not exists")

// ICache 缓存API
type ICache interface {
	ISCache
	IMCache
}
