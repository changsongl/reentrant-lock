package reentrant_lock

import (
	"sync"
)

// Lock for reentrant
type Lock struct {
	lock     sync.Mutex
	ch       chan struct{}
	times    uint64
	gid      uint64 // goroutine id
	goIDFunc func() uint64
}

// Lock reentrant lock
func (l *Lock) Lock() {
	gid := l.getGID()
	if l.reentrant(gid) {
		return
	}

	l.ch <- struct{}{}

	l.lock.Lock()
	defer l.lock.Unlock()

	if l.gid != 0 || l.times != 0 {
		panic("reentrant Lock: wrong state")
	}

	l.gid = gid
	l.times = 1
}

// reentrant check if the current goroutine is locking it again. If yes, return true.
// Otherwise return false. Init the channel if it is nil.
func (l *Lock) reentrant(gid uint64) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.ch == nil {
		l.ch = make(chan struct{}, 1)
	}

	if l.gid == gid {
		l.times++
		return true
	}

	return false
}

// Unlock reentrant lock
func (l *Lock) Unlock() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.times == 0 {
		panic("reentrant Lock: unlock of unlocked mutex")
	}

	l.times--
	if l.times > 0 {
		return
	}

	l.gid = 0
	select {
	case <-l.ch:
	default:
		panic("reentrant Lock: select went wrong")
	}
}

// getGID return a gid for current goroutine
func (l *Lock) getGID() uint64 {
	if l.goIDFunc == nil {
		return defaultGetGoIDFunc()
	}

	return l.goIDFunc()
}

// SetGidFunc set diy get gid function
func (l *Lock) SetGidFunc(f func() uint64) {
	l.goIDFunc = f
}
