package worker

import (
	"context"
	"log"
	"sync"
)

type Task interface {
	Execute(ctx context.Context) error
	Name() string
}

type Pool struct {
	tasks   chan Task
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	workers int
}

func NewPool(workers int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pool{
		tasks:   make(chan Task, 100),
		ctx:     ctx,
		cancel:  cancel,
		workers: workers,
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	log.Printf("worker pool started with %d workers", p.workers)
}

func (p *Pool) Submit(task Task) {
	select {
	case p.tasks <- task:
	case <-p.ctx.Done():
		log.Printf("worker pool is shutting down, task %s dropped", task.Name())
	}
}

func (p *Pool) Shutdown() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
	log.Println("worker pool shut down")
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for task := range p.tasks {
		log.Printf("worker %d: starting task %s", id, task.Name())
		if err := task.Execute(p.ctx); err != nil {
			log.Printf("worker %d: task %s failed: %v", id, task.Name(), err)
		} else {
			log.Printf("worker %d: task %s completed", id, task.Name())
		}
	}
}
