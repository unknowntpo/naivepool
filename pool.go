// This version of naivepool use pool.Schedule to send job directly to pool.workerChan,
// The worker use for-range pattern to receive job from Pool until channel is closed.
package naivepool

import (
	"context"
	"sync"
)

// jobFunc represents the function that will be executed by workers.
type jobFunc func()

type Pool struct {
	jobChan    chan jobFunc // We use jobChan to communicate between caller of Pool and Pool.
	workerChan chan jobFunc // workers conatains channel to communicate with each worker.
	maxWorkers int
	wg         sync.WaitGroup // Use waitgroup to wait for workers done its job and retire.
}

// Make our workerChan a buffered channel.
const workerChanSize int = 20

// New inits goroutine pool with capacity of jobchan and workerchan.
func New(bufSize, maxWorkers int) *Pool {
	p := &Pool{
		jobChan:    make(chan jobFunc, bufSize),
		workerChan: make(chan jobFunc, workerChanSize),
		maxWorkers: maxWorkers,
	}

	return p
}

// Start starts dispatching jobs to workers.
func (p *Pool) Start(ctx context.Context) {
	// TODO: Dynamic add or purge workers
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		// set up channel between pool and worker
		go p.work(ctx)
	}

	// Monitor goroutine: receive cancellation from ctx, and close workerChan
	// to inform workers to retire.
	go func() {
		select {
		case <-ctx.Done():
			close(p.workerChan)
			return
		}
	}()
}

// Wait waits for all workers finish its job and retire.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Schedule sends the job the p.jobChan.
func (p *Pool) Schedule(job jobFunc) {
	p.workerChan <- job
}

// work executes the job received from p.workerChan.
func (p *Pool) work(ctx context.Context) {
	defer p.wg.Done()
	for f := range p.workerChan {
		f()
	}
}
