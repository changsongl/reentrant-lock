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

### 获取goroutine id方式
1. 从stack获得gid, 本包默认使用的。性能会较差10000次>50ms。
````go
package main

import (
    "bytes"
    "fmt"
    "runtime"
    "strconv"
)

func main() {
    fmt.Println(GetGID())
}

func GetGID() uint64 {
    b := make([]byte, 64)
    b = b[:runtime.Stack(b, false)]
    b = bytes.TrimPrefix(b, []byte("goroutine "))
    b = b[:bytes.IndexByte(b, ' ')]
    n, _ := strconv.ParseUint(string(b), 10, 64)
    return n
}
````

2. 修改go的代码 `src/runtime/runtime2.go`, 但以后只能在自己的go代码编译。
````go
func Goid() int64 {
    _g_ := getg()
    return _g_.goid
}
````

3. CGo去获得gid, 不影响移植和性能。但要开启cgo.

````cgo
// 文件id.c:
#include "runtime.h"

int64 ·Id(void) {
    return g-&gt;goid;
}
````

````go
// 文件id.go
package id

func Id() int64
````


4. 通过汇编获取gid。那就是通过汇编获取goroutine id的方法。原理是：通过getg方法（汇编实现）获取到当前goroutine的g结构地址，根据偏移量计算出成员goid int的地址，然后取出该值即可。
