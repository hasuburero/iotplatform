package job

import (
	"errors"
	"sync"
)

type Worker_Struct struct {
	Worker_id string
	Runtime   []string
}

type Scheduled_Struct struct {
	Worker Worker_Struct
	Job    Job_Struct
}

type Assigned_Struct struct {
	Worker Scheduling_Worker_Struct
	Job    Job_Struct
}

type Scheduling_Worker_Struct struct {
	Worker Worker_Struct
	Chan   chan Scheduled_Struct
	Error  chan error
}

type Job_Delete_Struct struct {
	Job_id string
	Error  chan error
}

type Enqueue_Worker_Struct Scheduling_Worker_Struct

type Enqueue_Job_Struct Job_Struct

type Read_Struct struct {
	Worker Scheduling_Worker_Struct
	Job    Job_Struct
}

var WorkerQueue []*Scheduling_Worker_Struct
var JobQueue []*Job_Struct
var AssignedPair map[string]*Assigned_Struct

var EnqueueWorkerChan chan *Scheduling_Worker_Struct
var EnqueueJobChan chan *Job_Struct
var ReadChan chan Read_Struct
var RetryJobChan chan *Job_Struct
var RetryWorkerChan chan *Scheduling_Worker_Struct
var DeleteChan chan Job_Delete_Struct

var PairMux sync.Mutex

func Read() (Scheduling_Worker_Struct, Job_Struct) {
	pair := <-ReadChan
	return pair.Worker, pair.Job
}

func EnqueueWorker(worker Enqueue_Worker_Struct) {
	new_worker := new(Scheduling_Worker_Struct)
	new_worker.Worker = worker.Worker
	new_worker.Chan = worker.Chan
	new_worker.Error = worker.Error
	EnqueueWorkerChan <- new_worker
}

func EnqueueJob(job Job_Struct) {
	new_job := new(Job_Struct)
	*new_job = job
	EnqueueJobChan <- new_job
}

func DeletePair(job_id string) error {
	PairMux.Lock()
	ctx, _ := AssignedPair[job_id]
	if ctx == nil {
		return errors.New("")
	}
	delete(AssignedPair, job_id)
	PairMux.Unlock()

	return nil
}

func RetryJob(job Job_Struct) error {
	PairMux.Lock()
	ctx, _ := AssignedPair[job.Job_id]
	if ctx == nil {
		return errors.New("no such job")
	}
	AssignedPair[job.Job_id] = nil
	PairMux.Unlock()

	job_buf := new(Job_Struct)
	*job_buf = job
	RetryJobChan <- job_buf

	return nil
}

func RetryWorker(job_id string) error {
	PairMux.Lock()
	pair, exists := AssignedPair[job_id]
	if !exists {
		return errors.New("no such pair")
	}
	worker := new(Scheduling_Worker_Struct)
	worker.Worker = pair.Worker.Worker
	worker.Chan = pair.Worker.Chan
	worker.Error = pair.Worker.Error
	PairMux.Unlock()

	RetryWorkerChan <- worker
	return nil
}

func Matching() {
	job_target := JobQueue[0]
	worker_target := WorkerQueue[0]
	index := -1
	for i, worker := range WorkerQueue {
		for _, runtime := range worker.Worker.Runtime {
			if job_target.Runtime == runtime {
				worker_target = worker
				index = i
				break
			}
		}
	}
	JobQueue = JobQueue[1:]
	if index == -1 {
		JobQueue = append(JobQueue, job_target)
	} else {
		pair := new(Assigned_Struct)
		pair.Worker = *worker_target
		pair.Job = *job_target
		PairMux.Lock()
		AssignedPair[pair.Job.Job_id] = pair
		PairMux.Unlock()
		ReadChan <- Read_Struct{Worker: *worker_target, Job: *job_target}
	}

	return
}

func Select() {
	select {
	case worker := <-EnqueueWorkerChan:
		WorkerQueue = append(WorkerQueue, worker)
	case job := <-EnqueueJobChan:
		JobQueue = append(JobQueue, job)
	}
}

func SelectWithDefault() {
	select {
	case worker := <-EnqueueWorkerChan:
		WorkerQueue = append(WorkerQueue, worker)
	case job := <-EnqueueJobChan:
		JobQueue = append(JobQueue, job)
	default:
		Matching()
	}
}

func Scheduling() {
	go func() {
		for {
			// receiving retry, delete channel
			select {
			case job := <-RetryJobChan:
				buf := make([]*Job_Struct, len(JobQueue))
				copy(buf, JobQueue)
				JobQueue = append(JobQueue[:0], job)
				JobQueue = append(JobQueue, buf...)
				buf = nil
			case worker := <-RetryWorkerChan:
				buf := make([]*Scheduling_Worker_Struct, len(WorkerQueue))
				WorkerQueue = append(WorkerQueue[:0], worker)
				WorkerQueue = append(WorkerQueue, buf...)
				buf = nil
			default:
			}
			if len(WorkerQueue) > 0 && len(JobQueue) > 0 {
				SelectWithDefault()
			} else {
				Select()
			}
		}
	}()
}

func Reader() {
	go func() {
		for {
			target, job := Read()
			err := JobUpdate(job.Job_id, Processing)
			if err != nil {
				target.Error <- err
			} else {
				target.Chan <- Scheduled_Struct{Worker: target.Worker, Job: job}
			}
		}
	}()
}

func Start() {
	AssignedPair = make(map[string]*Assigned_Struct)
	EnqueueWorkerChan = make(chan *Scheduling_Worker_Struct, 10)
	EnqueueJobChan = make(chan *Job_Struct, 10)
	ReadChan = make(chan Read_Struct, 10)
	RetryJobChan = make(chan *Job_Struct, 10)
	RetryWorkerChan = make(chan *Scheduling_Worker_Struct, 10)

	Scheduling()
	Reader()
}
