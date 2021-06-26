// Copyright(C) 2020 github.com/hidu  All Rights Reserved.
// Author: hidu (duv123+git@baidu.com)
// Date: 2020/5/24

package lrucache

import (
	"reflect"
	"testing"
)

func Test_newUnmarshaler(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name    string
		args    args
		got     interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				val: true,
			},
			got:     false,
			want:    true,
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{
				val: 999,
			},
			got:     0,
			want:    999,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := newUnmarshaler(tt.args.val)
			err := fn(nil, &tt.got)
			if hasErr := err != nil; hasErr != tt.wantErr {
				t.Errorf("wantErr=%v, but not", tt.wantErr)
			}
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("newUnmarshaler() = %v, want %v", tt.got, tt.want)
			}
		})
	}

	t.Run("expect_err", func(t *testing.T) {
		fn := newUnmarshaler("abc")
		var n int
		err := fn(nil, n)
		if err == nil {
			t.Errorf("expect has error")
		}
	})
}
