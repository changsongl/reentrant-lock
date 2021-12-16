# reentrant-lock

### 介绍
重入锁，谨慎用于生产环境。此方法通过拿到goroutine id作为锁重入的判断，会因升级go版本而导致不可用。
假如发现此默认获取goroutine id的方法在新的go版本不可以使用，你可以注入自己的获取方法, 通过`OptionGetGoIDFunc`。

### 示例
````go
func TestNew(t *testing.T) {
	m := make(map[string]struct{})
	wg := sync.WaitGroup{}
	wg.Add(3)
	setFunc := func(s string) {
		t.Logf("%s entrant", s)
		m[s] = struct{}{}
		wg.Done()
	}

	l := New()

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
````
````shell
=== RUN   TestNew
    lock_test.go:15: lock1 entrant
    lock_test.go:15: lock2 entrant
    lock_test.go:15: lock3 entrant
--- PASS: TestNew (0.00s)
PASS
````
