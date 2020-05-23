/*
 * Copyright(C) 2020 github.com/hidu  All Rights Reserved.
 * Author: hidu (duv123+git@baidu.com)
 * Date: 2020/5/10
 */

package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ISCache 普通的单个缓存
type ISCache interface {
	Get(ctx context.Context, key interface{}) GetResult
	Set(ctx context.Context, key interface{}, value interface{}, ttl time.Duration) SetResult
	Has(ctx context.Context, key interface{}) HasResult
	Delete(ctx context.Context, key interface{}) DeleteResult

	// 缓存数据重置
	Reset(ctx context.Context) error
}

// Unmarshaler 数据反序列化方法
type Unmarshaler func(bf []byte, obj interface{}) error

// Marshaler 数据序列化方法
type Marshaler func(obj interface{}) ([]byte, error)

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

// NewGetResult 创建Get方法的结果
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
	Value(obj interface{}) (has bool, err error)
}

type getResult struct {
	err         error
	val         []byte
	unmarshaler Unmarshaler
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

func (g *getResult) Value(obj interface{}) (has bool, err error) {
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

// DeleteResult Delete方法的结果接口定义
type DeleteResult interface {
	ResultError
	Num() int
}

type deleteResult struct {
	err error
	num int
}

func (d *deleteResult) Err() error {
	return d.err
}

func (d *deleteResult) Num() int {
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
