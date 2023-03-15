package naivepool

import (
	"context"
	"sync"

	"github.com/unknowntpo/naivepool/domain"
)

// domain.JobFunc represents the function that will be executed by workers.

type poolForSelect struct {
	jobChan   chan domain.JobFunc // We use jobChan to communicate between caller of Pool and Pool.
	tokenChan chan struct{}       //  token controls the maximum amount of workers inside Pool.
	wg        sync.WaitGroup      // Use waitgroup to wait for workers done its job and retire.
}

// New inits goroutine pool with capacity of jobchan and workerchan.
// bufSize means the maximum number of jobs inside the buffer.
func New(bufSize, maxWorkers int) *poolForSelect {
	p := &poolForSelect{
		jobChan:   make(chan domain.JobFunc, bufSize),
		tokenChan: make(chan struct{}, maxWorkers),
	}

	return p
}

// Start starts dispatching jobs to workers.
func (p *poolForSelect) Start(ctx context.Context) {
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
func (p *poolForSelect) Wait() {
	p.wg.Wait()
}

// Schedule sends the job the p.jobChan.
func (p *poolForSelect) Schedule(job domain.JobFunc) {
	p.jobChan <- job
}

func (p *poolForSelect) work(job domain.JobFunc) {
	defer func() {
		// Add back one token
		p.tokenChan <- struct{}{}
		p.wg.Done()
	}()
	job()
}
