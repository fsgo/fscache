# fscache


[![Build Status](https://travis-ci.org/fsgo/fscache.png?branch=master)](https://travis-ci.org/fsgo/fscache)
[![GoCover](http://gocover.io/_badge/github.com/fsgo/fscache)](http://gocover.io/github.com/fsgo/fscache)
[![GoDoc](https://godoc.org/github.com/fsgo/fscache?status.svg)](https://godoc.org/github.com/fsgo/fscache)

## 1.缓存接口定义
```go
// ICache 缓存API
type ICache interface {
	ISCache
	IMCache
}
```

```go
// ISCache 普通的单个缓存
type ISCache interface {
	Get(ctx context.Context, key interface{}) GetResult
	Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) SetResult
	Has(ctx context.Context, key interface{}) HasResult
	Delete(ctx context.Context, key interface{}) DeleteResult
}
```
注：为了将批请求结果和单个处理结果尽量保持一致，操作结果均返回一个值。可以使用对应的`Err()`方法来判断是否有异常

```go
// IMCache 缓存-批处理接口
type IMCache interface {
	MGet(ctx context.Context, keys []interface{}) MGetResult
	MSet(ctx context.Context, kvs KVData, ttl time.Duration) MSetResult
	MDelete(ctx context.Context, keys []interface{}) MDeleteResult
	MHas(ctx context.Context, keys []interface{}) MHasResult
}
```

## 2.文件缓存
```go
opt:=filecache.Option{
    Dir: "./testdata/cache_dir/",
}
fc,err:=filecache.New(opt)

# 1.读取缓存
retGet := fc.Get(context.Background(), "abc")
if err := retGet.Err(); err != nil {
    t.Fatalf("retGet.Value with error:%v", err)
}

# 2.获取真实内容
var num int
if has, err := retGet.Value(&num); err != nil {
    t.Fatalf("retGet.Value with error:%v", err)
}
```