// Package pool implements a pool of workers.
package pool

import (
	"runtime"
	"sync"

	"github.com/alaleks/shortener/internal/app/logger"
)

// Pool represents an instance of a worker pool.
type Pool struct {
	done      chan struct{}
	logger    *logger.AppLogger
	out       chan Task
	tasks     chan Task
	wg        sync.WaitGroup
	numWorker int
	active    bool
}

// Task represents an instance job for added in pool.
type Task struct {
	data   any
	action func(data any) error
}

// Run - starts of worker pool.
func (p *Pool) Run() {
	workers := []chan Task{}
	for i := 0; i < runtime.NumCPU(); i++ {
		workers = append(workers, p.worker())
	}
	p.multiplex(workers...)

	for task := range p.out {
		err := task.action(task.data)
		if err != nil {
			p.logger.LZ.Error(err)
		}
	}
}

// Multiplexer implementation
func (p *Pool) multiplex(workers ...chan Task) {
	output := func(task <-chan Task) {
		for t := range task {
			p.out <- t
		}

		p.wg.Done()
	}

	p.wg.Add(len(workers))

	for _, worker := range workers {
		go output(worker)
	}

	go func() {
		p.wg.Wait()
	}()
}

func (p *Pool) worker() chan Task {
	out := make(chan Task)

	go func() {
		for task := range p.tasks {
			select {
			case out <- task:
			case <-p.done:
				return
			}
		}
	}()

	return out
}

// Stop - stops the pool.
func (p *Pool) Stop() {
	if p.active {
		close(p.done)
		p.active = false
	}
}

// SetNumWorker allows you to set the number of workers
// using the formula n * runtime.NumCPU().
func (p *Pool) SetNumWorker(num int) {
	p.numWorker = runtime.NumCPU() * num
}

// AddTask adds a task to the pool.
func (p *Pool) AddTask(data any, f func(data any) error) {
	p.tasks <- Task{
		data:   data,
		action: f,
	}
}

// Init initializes the pool instance.
func Init(logger *logger.AppLogger) *Pool {
	return &Pool{
		numWorker: runtime.NumCPU(),
		wg:        sync.WaitGroup{},
		logger:    logger,
		tasks:     make(chan Task, runtime.NumCPU()),
		out:       make(chan Task),
		done:      make(chan struct{}),
		active:    true,
	}
}
