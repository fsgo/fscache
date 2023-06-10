// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/23

package filecache

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsgo/fscache"
	"github.com/fsgo/fscache/internal"
)

// NewSCache 创建普通的缓存实例
func NewSCache(opt *Option) (fscache.SCache, error) {
	if err := opt.Check(); err != nil {
		return nil, err
	}
	codec := opt.GetCodec()
	return &SCache{
		opt:    opt,
		encode: codec.Marshal,
		decode: codec.Unmarshal,
	}, nil
}

// SCache 普通(非批量)缓存
type SCache struct {
	opt *Option

	decode fscache.UnmarshalFunc
	encode fscache.MarshalFunc
	gcTime int64

	gcRunning atomic.Bool
}

// Get 获取
func (f *SCache) Get(ctx context.Context, key any) fscache.GetResult {
	defer f.autoGC()

	expire, data, err := f.readByKey(key, true)
	if err != nil {
		return fscache.GetResult{Err: err}
	}
	if expire {
		_, _ = f.delete(ctx, key)
		return internal.GetRetNotExists
	}
	return fscache.GetResult{Payload: data, UnmarshalFunc: f.decode}
}

// Set 写入
func (f *SCache) Set(ctx context.Context, key any, value any, ttl time.Duration) fscache.SetResult {
	defer f.autoGC()

	fp := f.opt.CachePath(key)
	dir := filepath.Dir(fp)
	if !fileExists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fscache.SetResult{Err: err}
		}
	}

	msg, err := f.encode(value)
	if err != nil {
		return fscache.SetResult{Err: err}
	}

	expireAt := timeNow().Add(ttl)

	file, err := os.CreateTemp(dir, filepath.Base(fp))
	if err != nil {
		return fscache.SetResult{Err: err}
	}

	defer func() {
		_, err2 := unlink(file.Name())
		if err2 != nil {
			log.Printf("[fileCache.Set] unlink(%q) with error:%v\n", file.Name(), err2)
		}
	}()

	// 写 cache 文件：
	writer := bufio.NewWriter(file)
	err = writeStrings(writer,
		// 第1行是缓存有效期，格式:etime=1590235951234907000
		"etime=",
		strconv.FormatInt(expireAt.UnixNano(), 10),
		"\n",

		// 第2行是创建时间：格式： ctime=1590235951
		"ctime=",
		strconv.FormatInt(timeNow().Unix(), 10),
		"\n",
	)

	if err == nil {
		_, err = writer.Write(msg)
	}
	if err != nil {
		return fscache.SetResult{Err: err}
	}

	if err = writer.Flush(); err != nil {
		return fscache.SetResult{Err: err}
	}
	if err = file.Close(); err != nil {
		return fscache.SetResult{Err: err}
	}
	if err = os.Rename(file.Name(), fp); err != nil {
		return fscache.SetResult{Err: err}
	}
	return internal.SetRetSuc
}

func (f *SCache) readByKey(key any, needData bool) (expire bool, data []byte, err error) {
	fp := f.opt.CachePath(key)
	return f.readByPath(fp, needData)
}

func (f *SCache) readByPath(fp string, needData bool) (expire bool, data []byte, err error) {
	if !fileExists(fp) {
		return false, nil, fscache.ErrNotExists
	}

	content, err := os.ReadFile(fp)
	if err != nil {
		return true, nil, err
	}
	lines := bytes.SplitN(content, []byte("\n"), 3)
	if len(lines) < 2 {
		return true, nil, errors.New("invalid cache file")
	}
	// 第一行为过期时间，格式为：etime=UnixNano()
	first := lines[0]
	if len(first) < len("etime=") {
		return true, nil, fmt.Errorf("not valid cache line, expect etime=\\d+, got=%q", first)
	}
	expireAt, err := strconv.ParseInt(string(first[len("etime="):]), 10, 64)
	if err != nil {
		return true, nil, err
	}
	expire = expireAt < timeNow().UnixNano()
	if !needData {
		return expire, nil, nil
	}
	// 第二行为创建时间，格式为：ctime=unix时间戳
	return expire, lines[2], err
}

// Has 判断是否存在
func (f *SCache) Has(ctx context.Context, key any) fscache.HasResult {
	defer f.autoGC()

	expire, _, err := f.readByKey(key, false)
	if err != nil {
		return fscache.HasResult{Err: err}
	}
	if !expire {
		return internal.HasRetYes
	}
	return internal.HasRetNot
}

// Delete 删除
func (f *SCache) Delete(ctx context.Context, key any) fscache.DeleteResult {
	num, err := f.delete(ctx, key)
	return fscache.DeleteResult{Deleted: num, Err: err}
}

func (f *SCache) delete(ctx context.Context, key any) (int, error) {
	fp := f.opt.CachePath(key)
	return unlink(fp)
}

// Reset  重置
func (f *SCache) Reset(ctx context.Context) error {
	return filepath.Walk(f.opt.CacheDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), cacheFileExt) {
			err1 := os.Remove(path)
			if err1 != nil {
				log.Printf("[fileCache][warn] remove %q failed, %s\n", path, err1.Error())
			}
		}
		return nil
	})
}

func (f *SCache) autoGC() {
	lastGc := atomic.LoadInt64(&f.gcTime)
	newVal := timeNow().UnixNano()
	if newVal-lastGc < int64(f.opt.GetGCInterval()) {
		return
	}

	if !atomic.CompareAndSwapInt64(&f.gcTime, lastGc, newVal) {
		return
	}

	go func() {
		defer func() {
			if re := recover(); re != nil {
				log.Printf("[fileCache][warn] autoGC panic:%v\n", re)
			}
		}()
		f.gc()
	}()
}

func (f *SCache) gc() {
	if !f.gcRunning.CompareAndSwap(false, true) {
		return
	}
	defer f.gcRunning.Store(false)

	err := filepath.Walk(f.opt.CacheDir(), func(path string, info os.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}
		if !info.IsDir() {
			if err1 := f.checkFile(path); err1 != nil {
				log.Printf("[fileCache][warn] checkFile %q failed, %s\n", path, err1.Error())
			}
		}
		return nil
	})
	if err != nil {
		log.Println("[fileCache.gc] filepath.Walk with error:", err)
	}
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

var _ fscache.SCache = (*SCache)(nil)
var _ fscache.ReSetter = (*SCache)(nil)

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
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	return 0, nil
}
