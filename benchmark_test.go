package naivepool

import (
	"context"
	"sync"
	"testing"
)

func BenchmarkFib(b *testing.B) {
	numJobs := 1000

	b.Run("naivepool", func(b *testing.B) {
		var wg sync.WaitGroup

		maxWorkers := 50
		workerChanSize := 20

		fib23 := func() {
			n := 23
			cur := 1
			pre := 0
			res := 1
			for i := 1; i < n; i++ {
				res = pre + cur
				pre = cur
				cur = res
			}
			wg.Done()
		}

		pool := New(numJobs, maxWorkers, workerChanSize)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pool.Start(ctx)

		b.ResetTimer()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < numJobs; j++ {
				wg.Add(1)
				pool.Schedule(fib23)
			}
			wg.Wait()
		}
		b.StopTimer()
	})
	b.Run("native goroutine", func(b *testing.B) {
		var wg sync.WaitGroup

		fib23 := func() {
			n := 23
			cur := 1
			pre := 0
			res := 1
			for i := 1; i < n; i++ {
				res = pre + cur
				pre = cur
				cur = res
			}
			wg.Done()
		}

		b.ResetTimer()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < numJobs; j++ {
				wg.Add(1)
				go fib23()
			}
			wg.Wait()
		}
		b.StopTimer()
	})
}
