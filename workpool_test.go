package workpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

var sum int64
var runTimes = 10000000

var wg = sync.WaitGroup{}

func demoTask(v ...interface{}) {
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

func demoTaskWithWg(v ...interface{}) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}
//原生协程
func BenchmarkGroutineWorkWithWg(b *testing.B) {
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		go demoTaskWithWg()
	}
	wg.Wait()
	fmt.Println(sum)
}

func BenchmarkGoroutine(b *testing.B) {
	for i := 0; i < runTimes; i++ {
		go demoTask()
	}
	fmt.Println(sum)
}

//使用协程池
func BenchmarkPool(b *testing.B) {
	pool, err := NewWorkPool(20)
	if err != nil {
		b.Error(err)
	}

	task := &Task{
		Handler: demoTask,
	}

	for i := 0; i < runTimes; i++ {
		pool.Put(task)
	}
	fmt.Println(sum)
}

func BenchmarkPoolWithWg(b *testing.B) {
	pool, err := NewWorkPool(20)
	if err != nil {
		b.Error(err)
	}

	task := &Task{
		Handler: demoTaskWithWg,
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		pool.Put(task)
	}
	wg.Wait()
	fmt.Println(sum)
}
