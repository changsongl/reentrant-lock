package reentrant_lock

import (
	"bytes"
	"runtime"
	"strconv"
)

// opt option function wrapper
type opt func(*lock)

// apply function
func (o opt) apply(l *lock) {
	o(l)
}

// Option interface
type Option interface {
	apply(*lock)
}

// OptionGetGoIDFunc to set a get goroutine id function
func OptionGetGoIDFunc(f func() uint64) Option {
	return opt(func(l *lock) {
		l.goIDFunc = f
	})
}

// defaultGetGoIDFunc return go id, this function may not work for your go version
func defaultGetGoIDFunc() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
