package scheduler

import (
	"context"
	"time"
)

type Task struct {
	Name     string
	Interval time.Duration
	Job      func(context.Context) error
}
