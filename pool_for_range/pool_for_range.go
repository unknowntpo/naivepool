package naivepool

import (
	"context"
	"sync"

	"github.com/unknowntpo/naivepool/domain"
)

// domain.JobFunc represents the function that will be executed by workers.

type poolForRange struct {
	jobChan   chan domain.JobFunc // We use jobChan to communicate between caller of Pool and Pool.
	tokenChan chan struct{}       //  token controls the maximum amount of workers inside Pool.
	wg        sync.WaitGroup      // Use waitgroup to wait for workers done its job and retire.
}

// New inits goroutine pool with capacity of jobchan and workerchan.
// bufSize means the maximum number of jobs inside the buffer.
func New(bufSize, maxWorkers int) *poolForRange {
	p := &poolForRange{
		jobChan:   make(chan domain.JobFunc, bufSize),
		tokenChan: make(chan struct{}, maxWorkers),
	}

	return p
}

// Start starts dispatching jobs to workers.
func (p *poolForRange) Start(ctx context.Context) {
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
func (p *poolForRange) Wait() {
	p.wg.Wait()
}

// Schedule sends the job the p.jobChan.
func (p *poolForRange) Schedule(job domain.JobFunc) {
	p.jobChan <- job
}

func (p *poolForRange) work(job domain.JobFunc) {
	defer func() {
		// Add back one token
		p.tokenChan <- struct{}{}
		p.wg.Done()
	}()
	job()
}
