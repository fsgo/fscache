// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/17

package filecache

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsgo/fscache"
)

const cacheFileExt = ".cache"

// Option 配置选型
type Option struct {
	fscache.Option

	// Dir 缓存文件存储目录，必填
	Dir string

	// GCInterval 触发过期缓存清理的间隔时间，可选
	// 若为 0，会使用默认值 300秒
	GCInterval time.Duration
}

// GetGCInterval 获取自动 gc 的最小间隔
func (o *Option) GetGCInterval() time.Duration {
	if o.GCInterval == 0 {
		return 300 * time.Second
	}
	return o.GCInterval
}

// CacheDir 缓存根目录
func (o *Option) CacheDir() string {
	return o.Dir
}

// CachePath 获取缓存文件地址
func (o *Option) CachePath(key any) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprint(key)))
	s := hex.EncodeToString(h.Sum(nil))
	fp := filepath.Join(o.CacheDir(), s[:3], s[3:6], s[6:9], s[9:12], s[12:15], s[16:])
	return strings.Join([]string{fp, cacheFileExt}, "")
}

// Check 检查是否正确
func (o *Option) Check() error {
	if len(o.Dir) == 0 {
		return errors.New("cache dir is empty")
	}
	return nil
}
