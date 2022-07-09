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
	Get(ctx context.Context, key any) GetResult
	Set(ctx context.Context, key any, value any, ttl time.Duration) SetResult
	Has(ctx context.Context, key any) HasResult
	Delete(ctx context.Context, key any) DeleteResult
}

// Reseter 重置缓存
type Reseter interface {
	Reset(ctx context.Context) error
}

// Unmarshaler 数据反序列化方法
type Unmarshaler func(bf []byte, obj any) error

// Marshaler 数据序列化方法
type Marshaler func(obj any) ([]byte, error)

// ResultError 所有结果都包含的错误信息
type ResultError interface {
	Err() error
}

type resultError struct {
	err error
}

func (r *resultError) Err() error {
	return r.err
}

// NewGetResult 创建 Get 方法的结果
func NewGetResult(bf []byte, err error, unmarshaler Unmarshaler) GetResult {
	return &getResult{
		err:         err,
		val:         bf,
		unmarshaler: unmarshaler,
	}
}

// GetResult Get 方法的结果的接口定义
type GetResult interface {
	ResultError

	// Value 获取缓存的值
	Value(obj any) (has bool, err error)

	// Has 是否有值
	Has() bool
}

var getRetNotExists = NewGetResult(nil, ErrNotExists, nil)

type getResult struct {
	err         error
	val         []byte
	unmarshaler Unmarshaler
}

func (g *getResult) Has() bool {
	if g.err == ErrNotExists {
		return false
	}
	return g.err == nil
}

func (g *getResult) Err() error {
	if g.err == ErrNotExists {
		return nil
	}
	return g.err
}

func (g *getResult) String() string {
	return fmt.Sprintf("err=%v; val=%q; unmarshaler=%v", g.err, g.val, g.unmarshaler)
}

func (g *getResult) Value(obj any) (has bool, err error) {
	if g.err == ErrNotExists {
		return false, nil
	}
	if g.err != nil {
		return false, g.err
	}

	if g.unmarshaler == nil {
		return false, errors.New("unmarshaler is nil")
	}

	err = g.unmarshaler(g.val, obj)
	if err == nil {
		return true, nil
	}
	return false, err
}

var _ GetResult = (*getResult)(nil)

// NewSetResult 创建一个Set接口的结果
func NewSetResult(err error) SetResult {
	return &resultError{
		err: err,
	}
}

var setRetSuc = NewSetResult(nil)

// SetResult Set接口的结果接口定义
type SetResult interface {
	ResultError
}

// NewDeleteResult 创建Delete方法的结果
func NewDeleteResult(err error, num int) DeleteResult {
	return &deleteResult{
		err: err,
		num: num,
	}
}

var deleteRetSucHas0 = NewDeleteResult(nil, 0)

// DeleteResult Delete方法的结果接口定义
type DeleteResult interface {
	ResultError

	// Deleted 删除的数量
	Deleted() int
}

type deleteResult struct {
	err error
	num int
}

func (d *deleteResult) Err() error {
	return d.err
}

func (d *deleteResult) Deleted() int {
	return d.num
}

var _ DeleteResult = (*deleteResult)(nil)

// NewHasResult 创建Has接口的结果
func NewHasResult(err error, has bool) HasResult {
	return &hasResult{
		err: err,
		has: has,
	}
}

var hasRetNot = NewHasResult(ErrNotExists, false)

// HasResult Has接口的结果接口定义
type HasResult interface {
	ResultError
	Has() bool
}

type hasResult struct {
	err error
	has bool
}

func (h *hasResult) Err() error {
	if h.err == ErrNotExists {
		return nil
	}
	return h.err
}

func (h *hasResult) Has() bool {
	return h.has
}

var _ HasResult = (*hasResult)(nil)
