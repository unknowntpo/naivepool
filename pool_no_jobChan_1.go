// +build no_jobChan_1

// This version of naivepool doesn't use jobChan,
// Instead, in p.Schedule(job), it send jobs directly to the workers through workers workerChan,
// The worker use for-range loop to receive job from Pool until channel is closed.
package naivepool

import (
	"context"
	"sync"
)

// jobFunc represents the function that will be executed by workers.
type jobFunc func()

type Pool struct {
	workers        chan workerChan // workers conatains channel to communicate with each worker.
	maxJobs        int
	maxWorkers     int
	workerChanSize int            // the size of workerChan.
	wg             sync.WaitGroup // Use waitgroup to wait for workers done its job and retire.
}

// New inits goroutine pool with capacity of jobchan and workerchan.
func New(maxJobs, maxWorkers, workerChanSize int) *Pool {
	p := &Pool{
		workers:        make(chan workerChan, maxWorkers),
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
		w := NewWorker(p.workerChanSize)
		// set up channel between pool and worker
		p.workers <- w.c
		go w.work(&p.wg)
	}

	go func() {
		<-ctx.Done()
		// for each workerChan in workers
		// close them to inform worker that they need to retire.
		for wc := range p.workers {
			// What happend if we have jobs inside channel ?
			close(wc)
		}
		close(p.workers)
		return
	}()
}

// Wait waits for all workers finish its job and retire.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Schedule sends the job to workers' workerChan.
func (p *Pool) Schedule(job jobFunc) {
	// take 1 workerChan
	wc := <-p.workers
	// assign a job to the worker holding that workerChan.
	wc <- job
	// put him back to p.workers
	p.workers <- wc
}

// workerChan is the channel connecting between pool and worker, each worker uses 1 workerChan.
// When workerChan is closed, it meas that worker needs to retire.
type workerChan chan jobFunc

type worker struct {
	c workerChan
}

// NewWorker init a new instance of worker.
// we return worker, not *worker to avoid worker escaping to heap, which takes time to do memory-allocation.
func NewWorker(chanSize int) worker {
	return worker{
		c: make(chan jobFunc, chanSize),
	}
}

// worker is the worker that execute the job received from p.workerChan.
func (w *worker) work(wg *sync.WaitGroup) {
	defer wg.Done()
	for f := range w.c {
		f()
	}
}
