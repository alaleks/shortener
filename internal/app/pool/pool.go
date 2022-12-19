package pool

import (
	"runtime"
)

type Pool struct {
	Done      chan struct{}
	WorkerNum int
	Funnel    chan *Job
	Workers   []*Worker
	Jobs      []*Job
}

func NewPool(jobs []*Job) *Pool {
	limit := runtime.NumCPU()

	return &Pool{
		Jobs:      jobs,
		WorkerNum: limit,
		Funnel:    make(chan *Job, limit),
		Done:      make(chan struct{}),
	}
}

func (pool *Pool) AddJob(job *Job) {
	pool.Funnel <- job
}

func (pool *Pool) Start() {
	for i := 1; i <= pool.WorkerNum; i++ {
		worker := NewWorker(pool.Funnel, i)
		pool.Workers = append(pool.Workers, worker)
		go worker.Start()
	}

	for i := range pool.Jobs {
		pool.Funnel <- pool.Jobs[i]
	}

	<-pool.Done
}

func (pool *Pool) Stop() {
	for i := range pool.Workers {
		pool.Workers[i].Stop()
	}

	pool.Done <- struct{}{}
}

type Job struct {
	Err    error
	ID     any
	Data   any
	Action func(id, data any) error
}

func NewJob(action func(id, data any) error, id, data any) *Job {
	return &Job{Action: action, ID: id, Data: data}
}

func RunJob(job *Job) {
	job.Err = job.Action(job.ID, job.Data)
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
