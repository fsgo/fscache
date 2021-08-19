// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/8/19

package filecache

import (
	"time"
)

var timeNow = time.Now

type canWriteString interface {
	WriteString(s string) (int, error)
}

func writeStrings(bf canWriteString, strs ...string) error {
	for _, str := range strs {
		_, err := bf.WriteString(str)
		if err != nil {
			return err
		}
	}
	return nil
}
