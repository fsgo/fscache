/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/23
 */

package filecache

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fsgo/fscache/cache"
)

// NewSCache 创建普通的缓存实例
func NewSCache(opt IOption) (cache.ISCache, error) {
	if err := opt.Check(); err != nil {
		return nil, err
	}
	return &SCache{
		opt: opt,
	}, nil
}

type SCache struct {
	opt IOption
}

func (f *SCache) Get(ctx context.Context, key interface{}) cache.GetResult {
	expire, data, err := f.read(key, true)
	if err != nil {
		return cache.NewGetResult(nil, err, nil)
	}

	if expire {
		f.Delete(ctx, key)
		return cache.NewGetResult(nil, cache.ErrNotExists, nil)
	}
	return cache.NewGetResult(data, nil, f.opt.Unmarshaler())
}

func (f *SCache) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) cache.SetResult {
	fp := f.opt.CachePath(key)
	dir := filepath.Dir(fp)
	if !fileExists(dir) {
		os.MkdirAll(dir, 0777)
	}

	msg, err := f.opt.Marshaler()(value)
	if err != nil {
		return cache.NewSetResult(err)
	}

	expireAt := time.Now().Add(ttl)

	file, err := ioutil.TempFile(dir, filepath.Base(fp))
	if err != nil {
		return cache.NewSetResult(err)
	}

	defer unlink(file.Name())
	// cache文件格式：
	// 第一行为缓存有效期，格式:etime=1590235951234907000
	writer := bufio.NewWriter(file)
	writer.WriteString("etime=")
	writer.WriteString(strconv.FormatInt(expireAt.UnixNano(), 10))
	writer.WriteString("\n")
	// 第二行为创建时间：格式： ctime=1590235951
	writer.WriteString("ctime=")
	writer.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	writer.WriteString("\n")
	writer.Write(msg)
	if err := writer.Flush(); err != nil {
		return cache.NewSetResult(err)
	}
	if err := file.Close(); err != nil {
		return cache.NewSetResult(err)
	}
	if err := os.Rename(file.Name(), fp); err != nil {
		return cache.NewSetResult(err)
	}
	return cache.NewSetResult(nil)
}

func (f *SCache) read(key interface{}, needData bool) (expire bool, data []byte, err error) {
	fp := f.opt.CachePath(key)
	if !fileExists(fp) {
		return false, nil, cache.ErrNotExists
	}

	file, err := os.Open(fp)
	if err != nil {
		return true, nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	// 第一行为过期时间，格式为：etime=UnixNano()
	first, _, err := reader.ReadLine()
	if err != nil {
		return true, nil, err
	}

	if len(first) < len("etime=") {
		return true, nil, fmt.Errorf("not valid cache line, expect etime=\\d+, got=%q", first)
	}
	expireAt, err := strconv.ParseInt(string(first[len("etime="):]), 10, 64)
	if err != nil {
		return true, nil, err
	}
	expire = expireAt < time.Now().UnixNano()
	if !needData {
		return expire, nil, nil
	}
	// 第二行为创建时间，格式为：ctime=unix时间戳
	reader.ReadLine()

	data, _, err = reader.ReadLine()
	return expire, data, err
}

func (f *SCache) Has(ctx context.Context, key interface{}) cache.HasResult {
	expire, _, err := f.read(key, false)
	if err != nil {
		return cache.NewHasResult(err, false)
	}
	if !expire {
		return cache.NewHasResult(nil, true)
	}
	return cache.NewHasResult(nil, false)
}

func (f *SCache) Delete(ctx context.Context, key interface{}) cache.DeleteResult {
	fp := f.opt.CachePath(key)
	num, err := unlink(fp)
	return cache.NewDeleteResult(err, num)
}

var _ cache.ISCache = (*SCache)(nil)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func unlink(name string) (int, error) {
	if fileExists(name) {
		err := os.Remove(name)
		if err == nil {
			return 1, nil
		}
		return 0, err
	}
	return 0, nil
}
