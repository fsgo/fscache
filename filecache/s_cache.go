/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/23
 */

package filecache

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
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

// SCache 普通(非批量)缓存
type SCache struct {
	opt    IOption
	gcTime int64

	gcRunning int32
}

// Get 获取
func (f *SCache) Get(ctx context.Context, key interface{}) cache.GetResult {
	defer f.autoGC()

	expire, data, err := f.readByKey(key, true)
	if err != nil {
		return cache.NewGetResult(nil, err, nil)
	}
	if expire {
		f.Delete(ctx, key)
		return cache.NewGetResult(nil, cache.ErrNotExists, nil)
	}
	return cache.NewGetResult(data, nil, f.opt.Unmarshaler())
}

// Set 写入
func (f *SCache) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) cache.SetResult {
	defer f.autoGC()

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
	if err = writer.Flush(); err != nil {
		return cache.NewSetResult(err)
	}
	if err = file.Close(); err != nil {
		return cache.NewSetResult(err)
	}
	if err = os.Rename(file.Name(), fp); err != nil {
		return cache.NewSetResult(err)
	}
	return cache.NewSetResult(nil)
}

func (f *SCache) readByKey(key interface{}, needData bool) (expire bool, data []byte, err error) {
	fp := f.opt.CachePath(key)
	return f.readByPath(fp, needData)
}

func (f *SCache) readByPath(fp string, needData bool) (expire bool, data []byte, err error) {
	if !fileExists(fp) {
		return false, nil, cache.ErrNotExists
	}

	content,err:=os.ReadFile(fp)
	if err!=nil{
		return true,nil,err
	}
	lines:=bytes.SplitN(content,[]byte("\n"),3)
	if len(lines)<2{
		return true,nil,fmt.Errorf("invalid cache file")
	}
	// 第一行为过期时间，格式为：etime=UnixNano()
	first:=lines[0]
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
	return expire, lines[2], err
}

// Has 判断是否存在
func (f *SCache) Has(ctx context.Context, key interface{}) cache.HasResult {
	defer f.autoGC()

	expire, _, err := f.readByKey(key, false)
	if err != nil {
		return cache.NewHasResult(err, false)
	}
	if !expire {
		return cache.NewHasResult(nil, true)
	}
	return cache.NewHasResult(nil, false)
}

// Delete 删除
func (f *SCache) Delete(ctx context.Context, key interface{}) cache.DeleteResult {
	fp := f.opt.CachePath(key)
	num, err := unlink(fp)
	return cache.NewDeleteResult(err, num)
}

// Reset  重置
func (f *SCache) Reset(ctx context.Context) error {
	return filepath.Walk(f.opt.CacheDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), cacheFileExt) {
			err := os.Remove(path)
			if err != nil {
				log.Printf("[filecache][warn] remove %q failed, %s\n", path, err.Error())
			}
		}
		return nil
	})
}

func (f *SCache) autoGC() {
	lastGc := atomic.LoadInt64(&f.gcTime)
	if time.Now().UnixNano()-lastGc < int64(f.opt.GetGCInterval()) {
		return
	}

	newVal := time.Now().UnixNano()
	if !atomic.CompareAndSwapInt64(&f.gcTime, lastGc, newVal) {
		return
	}

	go func() {
		defer func() {
			if re := recover(); re != nil {
				log.Printf("[filecache][warn] autoGC panic:%v\n", re)
			}
		}()
		f.gc()
	}()

}
func (f *SCache) gc() {
	if atomic.LoadInt32(&f.gcRunning) == 1 {
		return
	}
	atomic.StoreInt32(&f.gcRunning, 1)
	defer func() {
		atomic.StoreInt32(&f.gcRunning, 0)
	}()
	filepath.Walk(f.opt.CacheDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if err = f.checkFile(path); err != nil {
				log.Printf("[filecache][warn] checkFile %q failed, %s\n", path, err.Error())
			}
		}
		return nil
	})
}

func (f *SCache) checkFile(fp string) error {
	if strings.HasSuffix(fp, cacheFileExt) {
		return nil
	}
	expire, _, _ := f.readByPath(fp, false)
	if expire {
		return os.Remove(fp)
	}

	return nil
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
