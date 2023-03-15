package domain

import "context"

// Pool is the interface of all Pool implementation
type Pool interface {
	Wait()
	Start(ctx context.Context)
	Schedule(job JobFunc)
}

type JobFunc func()
