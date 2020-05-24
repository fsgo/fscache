# fscache

统一封装的缓存接口，目前已包含文件缓存(FileCache)、内存LRU缓存。

[![Build Status](https://travis-ci.org/fsgo/fscache.png?branch=master)](https://travis-ci.org/fsgo/fscache)
[![GoCover](http://gocover.io/_badge/github.com/fsgo/fscache)](http://gocover.io/github.com/fsgo/fscache)
[![GoDoc](https://godoc.org/github.com/fsgo/fscache?status.svg)](https://godoc.org/github.com/fsgo/fscache)

## 1.缓存接口定义
```go
// ICache 缓存API
type ICache interface {
    Get(ctx context.Context, key interface{}) GetResult
    Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) SetResult
    Has(ctx context.Context, key interface{}) HasResult
    Delete(ctx context.Context, key interface{}) DeleteResult
    
    // 以下为批处理接口
    MGet(ctx context.Context, keys []interface{}) MGetResult
    MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult
    MDelete(ctx context.Context, keys []interface{}) MDeleteResult
    MHas(ctx context.Context, keys []interface{}) MHasResult
}
```
注：为了将批请求结果和单个处理结果尽量保持一致，操作结果均返回一个值。可以使用对应的`Err()`方法来判断是否有异常

  
已支持的缓存：
```go
// NewFileCache 创建文件缓存
func NewFileCache(opt FileIOption) (ICache, error)

// NewLRUCache 创建内存lru缓存
func NewLRUCache(opt LRUIOption) (ICache, error)
```

## 2.使用示例
```go
import (
    "github.com/fsgo/fscache"
)

opt:=fscache.FileOption{
    Dir: "./testdata/cache_dir/",
}

fc,err:=fscache.NewFileCache(opt)
if err != nil {
    log.Fatalf("init lru cache failed: %v", err)
}

// 1.读取缓存
retGet := fc.Get(context.Background(), "abc")
if err := retGet.Err(); err != nil {
    t.Fatalf("retGet.Value with error:%v", err)
}

// 2.获取缓存值
var num int
if has, err := retGet.Value(&num); err != nil {
    t.Fatalf("retGet.Value with error:%v", err)
}
```