package grpool

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func Test_pool(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("CPU Num: ", runtime.NumCPU())
	pool := NewPool(100, 20)
	defer pool.Release()

	for i := 0; i < 100; i++ {
		count := i
		pool.JobQueue <- func() {
			fmt.Println("worker number: ", count)
		}
	}
	// dummy wait until jobs are finished
	time.Sleep(1 * time.Second)
}
