/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/10
 */

package filecache

import (
	"bufio"
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fsgo/fscache/cache"
)

// New 创建新缓存实例
func New(opt IOption) (cache.ICache, error) {
	if err := opt.Check(); err != nil {
		return nil, err
	}

	sc := &sCache{
		opt: opt,
	}
	c := &fileCache{
		sCache:  sc,
		IMCache: cache.NewMCacheBySCache(sc,false),
	}
	return c, nil
}

type fileCache struct {
	*sCache
	cache.IMCache
}

var _ cache.ICache = (*fileCache)(nil)

type sCache struct {
	opt IOption
}

func (f *sCache) Get(ctx context.Context, key interface{}) cache.GetResult {
	expireAt, data, err := f.read(key, true)
	if err != nil {
		return cache.NewGetResult(nil, err, nil)
	}
	if expireAt < time.Now().Unix() {
		f.Delete(ctx, key)
		return cache.NewGetResult(nil, cache.ErrNotExists, nil)
	}
	return cache.NewGetResult(data, nil, f.opt.Unmarshaler())
}

func (f *sCache) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) cache.SetResult {
	fp := f.opt.CachePath(key)
	dir := filepath.Dir(fp)
	if !fileExists(dir) {
		os.MkdirAll(dir, 0777)
	}

	msg, err := f.opt.Marshaler()(value)
	if err != nil {
		return cache.NewSetResult(err)
	}

	expireAt := time.Now().Add(ttl).Unix()

	file, err := ioutil.TempFile(dir, filepath.Base(fp))
	if err != nil {
		return cache.NewSetResult(err)
	}

	defer unlink(file.Name())
	// cache文件格式：
	// 第一行为缓存有效期，为时间戳
	// 之后为缓存数据
	writer := bufio.NewWriter(file)
	writer.WriteString(strconv.FormatInt(expireAt, 10))
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

func (f *sCache) read(key interface{}, needData bool) (expireAt int64, data []byte, err error) {
	fp := f.opt.CachePath(key)
	if !fileExists(fp) {
		return 0, nil, cache.ErrNotExists
	}

	file, err := os.Open(fp)
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	first, _, err := reader.ReadLine()
	if err != nil {
		return 0, nil, err
	}
	expireAt, err = strconv.ParseInt(string(first), 10, 64)
	if err != nil {
		return 0, nil, err
	}
	if !needData {
		return expireAt, nil, nil
	}
	data, _, err = reader.ReadLine()
	return expireAt, data, err
}

func (f *sCache) Has(ctx context.Context, key interface{}) cache.HasResult {
	expireAt, _, err := f.read(key, false)
	if err != nil {
		return cache.NewHasResult(err, false)
	}
	if time.Now().Unix() < expireAt {
		return cache.NewHasResult(nil, true)
	}
	return cache.NewHasResult(nil, false)
}

func (f *sCache) Delete(ctx context.Context, key interface{}) cache.DeleteResult {
	fp := f.opt.CachePath(key)
	num, err := unlink(fp)
	return cache.NewDeleteResult(err, num)
}

var _ cache.ISCache = (*sCache)(nil)

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
