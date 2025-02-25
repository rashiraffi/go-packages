package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

type Scheduler struct {
	tasks  []Task
	stopCh chan struct{}
	wg     sync.WaitGroup
}

func New(task []Task) *Scheduler {
	return &Scheduler{
		tasks:  task,
		stopCh: make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	for _, task := range s.tasks {
		s.wg.Add(1)
		go func(t Task) {
			defer s.wg.Done()
			ctx, cancel := context.WithCancel(context.Background())
			ticker := time.NewTicker(t.Interval)
			for {
				select {
				case <-ticker.C:
					err := t.Job(ctx)
					if err != nil {
						log.Printf("Task failed: %s, Error: %v", t.Name, err)
					}
				case <-s.stopCh:
					ticker.Stop()
					cancel()
					return
				}
			}
		}(task)
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	log.Println("scheduler stopped")
}
