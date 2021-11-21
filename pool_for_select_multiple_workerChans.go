// +build for_select_multiple_workerChans

package naivepool

import (
	"context"
	"sync"
)

// jobFunc represents the function that will be executed by workers.
type jobFunc func()

type Pool struct {
	jobChan        chan jobFunc    // We use jobChan to communicate between caller of Pool and Pool.
	workers        chan workerChan // workers conatains channel to communicate with each worker.
	maxJobs        int
	maxWorkers     int
	workerChanSize int            // the size of workerChan.
	wg             sync.WaitGroup // Use waitgroup to wait for workers done its job and retire.
}

// New inits goroutine pool with capacity of jobchan and workerchan.
func New(maxJobs, maxWorkers, workerChanSize int) *Pool {
	p := &Pool{
		jobChan:        make(chan jobFunc, maxJobs),
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
		go w.work(ctx, &p.wg)
	}

	go func() {
		for {
			select {
			// Received a job.
			// Dispatch it to workers.
			case job := <-p.jobChan:
				// take 1 workerChan
				wc := <-p.workers
				// assign a job to the worker holding that workerChan.
				wc <- job
				// put him back to p.workers
				p.workers <- wc
			case <-ctx.Done():
				close(p.workers)
				return
			}
		}
	}()
}

// Wait waits for all workers finish its job and retire.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Schedule sends the job the p.jobChan.
func (p *Pool) Schedule(job jobFunc) {
	p.jobChan <- job
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
func (w *worker) work(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case f := <-w.c:
			f()
		case <-ctx.Done():
			return
		}
	}
}
