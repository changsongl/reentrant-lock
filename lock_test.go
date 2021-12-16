package reentrant_lock

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := make(map[string]struct{})
	wg := sync.WaitGroup{}
	wg.Add(3)
	setFunc := func(s string) {
		t.Logf("%s entrant", s)
		m[s] = struct{}{}
		wg.Done()
	}

	l := Lock{}

	go func() {
		l.Lock()
		setFunc("lock1")
		l.Unlock()

		l.Lock()
		setFunc("lock2")
		l.Unlock()
	}()

	time.Sleep(time.Millisecond)
	l.Lock()
	setFunc("lock3")
	l.Unlock()

	wg.Wait()

	expected := map[string]struct{}{
		"lock1": {},
		"lock2": {},
		"lock3": {},
	}
	require.EqualValues(t, expected, m)
}
