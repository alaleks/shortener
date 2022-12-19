package pool

import (
	"fmt"
	"runtime"
)

type Pool struct {
	done      chan struct{}
	funnel    chan *Job
	workerNum int
	workers   []*Worker
	jobs      []*Job
}

func NewPool(jobs []*Job) *Pool {
	limit := runtime.NumCPU()

	return &Pool{
		jobs:      jobs,
		workerNum: limit,
		funnel:    make(chan *Job, limit),
		done:      make(chan struct{}),
	}
}

func (pool *Pool) AddJob(job *Job) {
	pool.funnel <- job
}

func (pool *Pool) Start() {
	for i := 1; i <= pool.workerNum; i++ {
		worker := NewWorker(pool.funnel, i)
		pool.workers = append(pool.workers, worker)
		go worker.Start()
	}

	for i := range pool.jobs {
		pool.funnel <- pool.jobs[i]
	}

	<-pool.done
}

func (pool *Pool) Stop() {
	if len(pool.workers) == 0 {
		return
	}

	for i := range pool.workers {
		fmt.Println("stop", i)
		pool.workers[i].Stop()
	}

	pool.done <- struct{}{}
}

type Job struct {
	err    error
	id     any
	data   any
	action func(id, data any) error
}

func NewJob(action func(id, data any) error, id, data any) *Job {
	return &Job{action: action, id: id, data: data}
}

func RunJob(job *Job) {
	job.err = job.action(job.id, job.data)
}

type Worker struct {
	JobCh chan *Job
	Done  chan struct{}
	ID    int
}

func NewWorker(jobCh chan *Job, id int) *Worker {
	return &Worker{JobCh: jobCh, ID: id, Done: make(chan struct{})}
}

func (worker *Worker) Start() {
	for {
		select {
		case job := <-worker.JobCh:
			RunJob(job)
		case <-worker.Done:
			return
		}
	}
}

func (worker *Worker) Stop() {
	go func() {
		worker.Done <- struct{}{}
	}()
}
