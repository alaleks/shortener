package wpool

import (
	"runtime"
)

type Job struct {
	Data   any
	Action func(data any) error
}

type Worker struct {
	jobCh chan Job
	done  chan struct{}
	wPool chan chan Job
}

func NewWorker(wPool chan chan Job) Worker {
	return Worker{
		wPool: wPool,
		done:  make(chan struct{}),
		jobCh: make(chan Job)}
}

func (wr Worker) Start() {
	go func() {
		for {
			wr.wPool <- wr.jobCh

			select {
			case job := <-wr.jobCh:
				job.Action(job.Data)
			case <-wr.done:

				return
			}
		}
	}()
}

func (wr Worker) Stop() {
	go func() {
		wr.done <- struct{}{}
	}()
}

type Multiplex struct {
	wPool     chan chan Job
	QueueJobs chan Job
	done      chan struct{}
}

func NewMultiplex() *Multiplex {
	return &Multiplex{wPool: make(chan chan Job,
		runtime.NumCPU()),
		QueueJobs: make(chan Job, runtime.NumCPU()),
		done:      make(chan struct{})}
}

func (m *Multiplex) Run() {
	var workers []Worker

	defer func() {
		for i := range workers {
			workers[i].Stop()
		}
	}()

	for i := 1; i <= runtime.NumCPU(); i++ {
		worker := NewWorker(m.wPool)
		workers = append(workers, worker)
		worker.Start()
	}

	go m.balancer()

	<-m.done
}

func (m *Multiplex) Stop() {
	go func() {
		m.done <- struct{}{}
	}()
}

func (m *Multiplex) balancer() {
	for {
		select {
		case job := <-m.QueueJobs:
			go func(job Job) {
				jobCh := <-m.wPool
				jobCh <- job
			}(job)
		case <-m.done:

			return
		}
	}
}
