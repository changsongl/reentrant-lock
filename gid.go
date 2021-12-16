package reentrant_lock

import (
	"bytes"
	"runtime"
	"strconv"
)

// defaultGetGoIDFunc return go id, this function may not work for your go version
func defaultGetGoIDFunc() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
