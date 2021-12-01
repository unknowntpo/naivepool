package naivepool

import (
	"context"
	"sync"
	"testing"

	"github.com/alitto/pond"
)

func fib() {
	fibNum := 10000
	n := fibNum

	cur := 1
	pre := 0
	res := 1
	for i := 1; i < n; i++ {
		res = pre + cur
		pre = cur
		cur = res
	}
}

func BenchmarkNaivepool(b *testing.B) {
	b.Run("fib", func(b *testing.B) {
		tests := []struct {
			name       string
			numJobs    int
			bufSize    int
			maxWorkers int
		}{
			{"1K tasks", 1000, 1000, 8},
			{"10K tasks", 10000, 1000, 8},
			{"100K tasks", 100000, 1000, 8},
			{"1M tasks", 1000000, 1000, 8},
		}

		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				pool := New(tt.bufSize, tt.maxWorkers)
				ctx, cancel := context.WithCancel(context.Background())

				pool.Start(ctx)

				b.ResetTimer()

				var wg sync.WaitGroup
				f := func() {
					defer wg.Done()
					fib()
				}

				for i := 0; i < b.N; i++ {
					for j := 0; j < tt.numJobs; j++ {
						wg.Add(1)
						pool.Schedule(f)
					}
					wg.Wait()
				}

				cancel()
				pool.Wait()
			})
		}
	})
}

func BenchmarkNormalGoroutine(b *testing.B) {
	b.Run("fib", func(b *testing.B) {
		tests := []struct {
			name    string
			numJobs int
		}{
			{"1K tasks", 1000},
			{"10K tasks", 10000},
			{"100K tasks", 100000},
			{"1M tasks", 1000000},
		}

		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				var wg sync.WaitGroup
				f := func() {
					defer wg.Done()
					fib()
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for j := 0; j < tt.numJobs; j++ {
						wg.Add(1)
						go f()
					}
					wg.Wait()
				}
			})
		}
	})
}

func BenchmarkPond(b *testing.B) {
	b.Run("fib", func(b *testing.B) {
		tests := []struct {
			name       string
			numJobs    int
			maxWorkers int
		}{
			{"1K tasks", 1000, 100},
			{"10K tasks", 10000, 100},
			{"100K tasks", 100000, 100},
			{"1M tasks", 1000000, 100},
		}

		for _, tt := range tests {
			b.Run(tt.name, func(b *testing.B) {
				// Create a buffered (non-blocking) pool that can scale up to tt.maxWorkers workers
				// and has a buffer capacity of tt.numJobs tasks
				pool := pond.New(tt.maxWorkers, tt.numJobs)

				var wg sync.WaitGroup
				f := func() {
					defer wg.Done()
					fib()
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for j := 0; j < tt.numJobs; j++ {
						wg.Add(1)
						pool.Submit(f)
					}
					wg.Wait()
				}
			})
		}
	})
}

/*
func BenchmarkFib(b *testing.B) {
	numJobs := 1000

	fibNum := 10000

	b.Run("fib", func(b *testing.B) {
		fib := func() {
			n := fibNum
			cur := 1
			pre := 0
			res := 1
			for i := 1; i < n; i++ {
				res = pre + cur
				pre = cur
				cur = res
			}
		}

		for i := 0; i < b.N; i++ {
			fib()
		}
	})

	b.Run("naivepool", func(b *testing.B) {
		var wg sync.WaitGroup

		maxWorkers := 4
		workerChanSize := 200

		fib := func() {
			n := fibNum
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
				pool.Schedule(fib)
			}
			wg.Wait()
		}
		b.StopTimer()
	})
	b.Run("native goroutine", func(b *testing.B) {
		var wg sync.WaitGroup

		fib := func() {
			n := fibNum
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
				go fib()
			}
			wg.Wait()
		}
		b.StopTimer()
	})

	b.Run("pond", func(b *testing.B) {
		var wg sync.WaitGroup

		fib := func() {
			n := fibNum
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

		// Create a buffered (non-blocking) pool that can scale up to 100 workers
		// and has a buffer capacity of 1000 tasks
		pool := pond.New(100, 1000)

		b.ResetTimer()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < numJobs; j++ {
				wg.Add(1)
				pool.Submit(fib)
			}
			wg.Wait()
		}
		b.StopTimer()
	})

}
*/
