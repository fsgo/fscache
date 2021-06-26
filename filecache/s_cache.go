// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/23

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

	"github.com/fsgo/fscache"
)

// NewSCache 创建普通的缓存实例
func NewSCache(opt OptionType) (fscache.SCache, error) {
	if err := opt.Check(); err != nil {
		return nil, err
	}
	return &SCache{
		opt: opt,
	}, nil
}

// SCache 普通(非批量)缓存
type SCache struct {
	opt    OptionType
	gcTime int64

	gcRunning int32
}

// Get 获取
func (f *SCache) Get(ctx context.Context, key interface{}) fscache.GetResult {
	defer f.autoGC()

	expire, data, err := f.readByKey(key, true)
	if err != nil {
		return fscache.NewGetResult(nil, err, nil)
	}
	if expire {
		f.Delete(ctx, key)
		return fscache.NewGetResult(nil, fscache.ErrNotExists, nil)
	}
	return fscache.NewGetResult(data, nil, f.opt.Unmarshaler())
}

// Set 写入
func (f *SCache) Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) fscache.SetResult {
	defer f.autoGC()

	fp := f.opt.CachePath(key)
	dir := filepath.Dir(fp)
	if !fileExists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fscache.NewSetResult(err)
		}
	}

	msg, err := f.opt.Marshaler()(value)
	if err != nil {
		return fscache.NewSetResult(err)
	}

	expireAt := time.Now().Add(ttl)

	file, err := ioutil.TempFile(dir, filepath.Base(fp))
	if err != nil {
		return fscache.NewSetResult(err)
	}

	defer func() {
		_, err2 := unlink(file.Name())
		if err2 != nil {
			log.Printf("[filecache.Set] unlink(%q) with error:%v\n", file.Name(), err2)
		}
	}()

	// cache文件格式：
	// 第一行为缓存有效期，格式:etime=1590235951234907000
	writer := bufio.NewWriter(file)
	_, _ = writer.WriteString("etime=")
	_, _ = writer.WriteString(strconv.FormatInt(expireAt.UnixNano(), 10))
	_, _ = writer.WriteString("\n")
	// 第二行为创建时间：格式： ctime=1590235951
	_, _ = writer.WriteString("ctime=")
	_, _ = writer.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	_, _ = writer.WriteString("\n")
	_, _ = writer.Write(msg)
	if err = writer.Flush(); err != nil {
		return fscache.NewSetResult(err)
	}
	if err = file.Close(); err != nil {
		return fscache.NewSetResult(err)
	}
	if err = os.Rename(file.Name(), fp); err != nil {
		return fscache.NewSetResult(err)
	}
	return fscache.NewSetResult(nil)
}

func (f *SCache) readByKey(key interface{}, needData bool) (expire bool, data []byte, err error) {
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
		return true, nil, fmt.Errorf("invalid cache file")
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
	expire = expireAt < time.Now().UnixNano()
	if !needData {
		return expire, nil, nil
	}
	// 第二行为创建时间，格式为：ctime=unix时间戳
	return expire, lines[2], err
}

// Has 判断是否存在
func (f *SCache) Has(ctx context.Context, key interface{}) fscache.HasResult {
	defer f.autoGC()

	expire, _, err := f.readByKey(key, false)
	if err != nil {
		return fscache.NewHasResult(err, false)
	}
	if !expire {
		return fscache.NewHasResult(nil, true)
	}
	return fscache.NewHasResult(nil, false)
}

// Delete 删除
func (f *SCache) Delete(ctx context.Context, key interface{}) fscache.DeleteResult {
	fp := f.opt.CachePath(key)
	num, err := unlink(fp)
	return fscache.NewDeleteResult(err, num)
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
	newVal := time.Now().UnixNano()
	if newVal-lastGc < int64(f.opt.GetGCInterval()) {
		return
	}

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

	err := filepath.Walk(f.opt.CacheDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if err1 := f.checkFile(path); err1 != nil {
				log.Printf("[filecache][warn] checkFile %q failed, %s\n", path, err1.Error())
			}
		}
		return nil
	})
	if err != nil {
		log.Println("[filecache.gc] filepath.Walk with error:", err)
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
