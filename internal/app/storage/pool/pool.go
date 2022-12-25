package pool

import (
	"runtime"
	"sync"

	"github.com/alaleks/shortener/internal/app/logger"
)

type Pool struct {
	done      chan struct{}
	logger    *logger.AppLogger
	out       chan Task
	tasks     chan Task
	wg        sync.WaitGroup
	numWorker int
	active    bool
}

type Task struct {
	data   any
	action func(data any) error
}

func (p *Pool) Run() {
	workers := []chan Task{}
	for i := 0; i < runtime.NumCPU(); i++ {
		workers = append(workers, p.Worker())
	}
	p.Multiplex(workers...)

	for task := range p.out {
		err := task.action(task.data)

		if err != nil {
			p.logger.LZ.Error(err)
		}
	}
}

func (p *Pool) Multiplex(workers ...chan Task) {
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

func (p *Pool) Worker() chan Task {
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

func (p *Pool) Stop() {
	if p.active {
		close(p.done)
		p.active = false
	}
}

func (p *Pool) AddTask(data any, f func(data any) error) {
	p.tasks <- Task{
		data:   data,
		action: f,
	}
}

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
