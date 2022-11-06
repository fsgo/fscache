// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/6/26

package fscache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// SCache 普通的单个缓存
type SCache interface {
	// Get 查询单个
	Get(ctx context.Context, key any) GetResult

	// Set 设置并附带有效期
	Set(ctx context.Context, key any, value any, ttl time.Duration) SetResult

	// Has 判断是否存在
	Has(ctx context.Context, key any) HasResult

	// Delete 删除指定的 key
	Delete(ctx context.Context, key any) DeleteResult
}

var getRetNotExists = GetResult{Err: ErrNotExists}

// GetResult Get 方法的结果
type GetResult struct {
	Err           error
	UnmarshalFunc UnmarshalFunc
	Payload       []byte
}

func (g GetResult) String() string {
	return fmt.Sprintf("err=%v; Payload=%q; unmarshaler=%v", g.Err, g.Payload, g.UnmarshalFunc)
}

// Value 获取值
func (g GetResult) Value(obj any) (has bool, err error) {
	if g.Err == ErrNotExists {
		return false, nil
	}
	if g.Err != nil {
		return false, g.Err
	}

	if g.UnmarshalFunc == nil {
		return false, errors.New("unmarshaler is nil")
	}

	err = g.UnmarshalFunc(g.Payload, obj)
	if err == nil {
		return true, nil
	}
	return false, err
}

// SetResult Set 方法的结果
type SetResult struct {
	Err error
}

var setRetSuc = SetResult{}

var deleteRetSucHas0 = DeleteResult{Deleted: 0}

// DeleteResult Delete 方法的结果接口定义
type DeleteResult struct {
	Err     error
	Deleted int
}

var hasRetNot = HasResult{Err: ErrNotExists, Has: false}

// HasResult Has 方法的结果
type HasResult struct {
	Err error
	Has bool
}
