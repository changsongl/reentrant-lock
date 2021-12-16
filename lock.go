package reentrant_lock

import (
	"sync"
)

type lock struct {
	lock     sync.Mutex
	ch       chan struct{}
	times    uint64
	gid      uint64 // goroutine id
	goIDFunc func() uint64
}

func (l *lock) Lock() {
	if l.goIDFunc == nil {
		panic("reentrant lock: goIDFunc is nil")
	}

	gid := l.goIDFunc()
	if l.reentrant(gid) {
		return
	}

	l.ch <- struct{}{}

	l.lock.Lock()
	defer l.lock.Unlock()

	if l.gid != 0 || l.times != 0 {
		panic("reentrant lock: wrong state")
	}

	l.gid = gid
	l.times = 1
}

func (l *lock) reentrant(gid uint64) bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.gid == gid {
		l.times++
		return true
	}

	return false
}

func (l *lock) Unlock() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.times == 0 {
		panic("reentrant lock: unlock of unlocked mutex")
	}

	l.times--
	if l.times > 0 {
		return
	}

	l.gid = 0
	select {
	case <-l.ch:
	default:
		panic("reentrant lock: select went wrong")
	}
}

// New a reentrant lock
func New(options ...Option) sync.Locker {
	l := &lock{
		ch:       make(chan struct{}, 1),
		goIDFunc: defaultGetGoIDFunc,
	}

	for _, option := range options {
		option.apply(l)
	}

	return l
}
