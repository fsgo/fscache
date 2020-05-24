/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/24
 */

package fscache

import (
	"github.com/fsgo/fscache/cache"
	"github.com/fsgo/fscache/filecache"
	"github.com/fsgo/fscache/lrucache"
)

// ICache 统一的缓存接口，包含单个接口和批量接口
type ICache = cache.ICache

// ISCache 缓存接口，普通单个接口
// 如 Get(xxx,key)
type ISCache = cache.ISCache

// IMCache 缓存接口，批量接口
// 如 MGet(xxx,keys)
type IMCache = cache.IMCache

// Template 缓存的模块类
// 如 NewFileCache 返回的缓存实例，实际就是一个Template
type Template = cache.Template

// KVData k-v结构
// map[interface{}]interface{} 的别名
// 如用于MSet
type KVData = cache.KVData

// GetResult Get方法的返回值
type GetResult = cache.GetResult

// MGetResult MGet方法的返回值
type MGetResult = cache.MGetResult

// SetResult Set 方法的返回值
type SetResult = cache.SetResult

// MSetResult MSet方法的返回值
type MSetResult = cache.MSetResult

// DeleteResult Delete方法的返回值
type DeleteResult = cache.DeleteResult

// MDeleteResult MDelete方法的返回值
type MDeleteResult = cache.MDeleteResult

// HasResult Has接口的返回值
type HasResult = cache.HasResult

// MHasResult MHas接口的返回值
type MHasResult = cache.MHasResult

// Marshaler 数据序列化方法
type Marshaler = cache.Marshaler

// Unmarshaler 数据反序列化方法
type Unmarshaler = cache.Unmarshaler

// --------------------------------------------------------
// 文件 缓存

// FileIOption 文件缓存配置项接口定义
type FileIOption = filecache.IOption

// FileOption 文件缓存配置项
type FileOption = filecache.Option

// NewFileCache 创建文件缓存
func NewFileCache(opt FileIOption) (ICache, error) {
	return filecache.New(opt)
}

// --------------------------------------------------------
// 内存 LRU缓存

// LRUIOption LRU缓存配置项接口定义
type LRUIOption = lrucache.IOption

// LRUOption LRU缓存配置项
type LRUOption = lrucache.Option

// NewLRUCache 创建内存lru缓存
func NewLRUCache(opt LRUIOption) (ICache, error) {
	return lrucache.New(opt)
}

// --------------------------------------------------------
