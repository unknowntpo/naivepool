package naivepool

import (
	"context"
	"sync"
)

// jobFunc represents the function that will be executed by workers.
type jobFunc func()

type Pool struct {
	jobChan   chan jobFunc   // We use jobChan to communicate between caller of Pool and Pool.
	tokenChan chan struct{}  //  token controls the maximum amount of workers inside Pool.
	wg        sync.WaitGroup // Use waitgroup to wait for workers done its job and retire.
}

// New inits goroutine pool with capacity of jobchan and workerchan.
// bufSize means the maximum number of jobs inside the buffer.
func New(bufSize, maxWorkers int) *Pool {
	p := &Pool{
		jobChan:   make(chan jobFunc, bufSize),
		tokenChan: make(chan struct{}, maxWorkers),
	}

	return p
}

// Start starts dispatching jobs to workers.
func (p *Pool) Start(ctx context.Context) {
	// Fill the tokenChan with maxWorkers
	for i := 0; i < cap(p.tokenChan); i++ {
		p.tokenChan <- struct{}{}
	}

	go func() {
		for {
			select {
			// Received a job.
			// Dispatch it to workers.
			case job := <-p.jobChan:
				// block until a token is available
				<-p.tokenChan
				p.wg.Add(1)
				go p.work(job)
			case <-ctx.Done():
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

func (p *Pool) work(job jobFunc) {
	defer func() {
		// Add back one token
		p.tokenChan <- struct{}{}
		p.wg.Done()
	}()
	job()
}
