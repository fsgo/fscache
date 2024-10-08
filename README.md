# fscache

统一封装的缓存接口，目前已包含文件缓存(FileCache)、内存LRU缓存。

[![GoDoc](https://pkg.go.dev/badge/github.com/fsgo/fscache)](https://pkg.go.dev/github.com/fsgo/fscache)

## 1.缓存接口定义
```go
// Cache 缓存API
type Cache interface {
    Get(ctx context.Context, key any) GetResult
    Set(ctx context.Context, key any, value any, ttl time.Duration) SetResult
    Has(ctx context.Context, key any) HasResult
    Delete(ctx context.Context, key any) DeleteResult
    
    // 以下为批处理接口
    MGet(ctx context.Context, keys []any) MGetResult
    MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult
    MDelete(ctx context.Context, keys []any) MDeleteResult
    MHas(ctx context.Context, keys []any) MHasResult
}
```
注：为了将批请求结果和单个处理结果尽量保持一致，操作结果均返回一个值。可以使用对应的`Err()`方法来判断是否有异常


## 2.使用示例
```go
import (
    "github.com/fsgo/fscache/filecache"
)

opt:=&filecache.Option{
    Dir: "./testdata/cache_dir/",
}

fc,err:=filecache.New(opt)
if err != nil {
    log.Fatalf("init cache failed: %v", err)
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