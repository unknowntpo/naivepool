// +build no_jobChan_2

// Improved from no_jobChan_1, we only use one workerChan to communicate with workers.
package naivepool

import (
	"context"
	"sync"
)

// jobFunc represents the function that will be executed by workers.
type jobFunc func()

type Pool struct {
	workerChan     chan jobFunc // workers conatains channel to communicate with each worker.
	maxJobs        int
	maxWorkers     int
	workerChanSize int            // the size of workerChan.
	wg             sync.WaitGroup // Use waitgroup to wait for workers done its job and retire.
}

// New inits goroutine pool with capacity of jobchan and workerchan.
func New(maxJobs, maxWorkers, workerChanSize int) *Pool {
	p := &Pool{
		workerChan:     make(chan jobFunc, maxJobs),
		maxJobs:        maxJobs,
		maxWorkers:     maxWorkers,
		workerChanSize: workerChanSize,
	}

	return p
}

// Start starts dispatching jobs to workers.
func (p *Pool) Start(ctx context.Context) {
	// TODO: Dynamic add or purge workers
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.work()
	}

	go func() {
		<-ctx.Done()
		close(p.workerChan)
		return
	}()
}

// Wait waits for all workers finish its job and retire.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Schedule sends the job to workers' workerChan.
// If p.wokerChan is full, it will block until a woker take one job from p.workerChan
func (p *Pool) Schedule(job jobFunc) {
	p.workerChan <- job
}

// worker is the worker that execute the job received from p.workerChan.
func (p *Pool) work() {
	defer p.wg.Done()
	for f := range p.workerChan {
		f()
	}
}
