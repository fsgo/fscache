// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/19

package filecache

import (
	"os"
	"testing"
)

func Test_fileExists(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "not exists",
			args: args{
				name: "not_exists.txt",
			},
			want: false,
		},
		{
			name: "exists",
			args: args{
				name: "s_cache.go",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.args.name); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unlink(t *testing.T) {
	type args struct {
		getName func() string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "not exists",
			args: args{
				getName: func() string {
					return "not_exists.txt"
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "new file",
			args: args{
				getName: func() string {
					name := "testdata/test_unlink_0.txt"
					err := os.WriteFile(name, []byte("test"), 0655)
					if err != nil {
						panic(err)
					}
					return name
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unlink(tt.args.getName())
			if (err != nil) != tt.wantErr {
				t.Errorf("unlink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unlink() got = %v, want %v", got, tt.want)
			}
		})
	}
}
