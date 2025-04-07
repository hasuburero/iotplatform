package sched

import (
	"errors"
	"sched/runtime_match"
	"sync"
	"time"

	"github.com/hasuburero/util/panic"
)

type AccessController_interface interface{}

type Worker_Struct struct {
	Worker_id string
	Runtime   []string
}

type Job_Struct struct {
	Job_id      string
	Data1_id    string
	Data2_id    string
	Function_id string
	Runtime     string
	Status      string
}

type Scheduled_Struct struct {
	Worker Worker_Struct
	Job    Job_Struct
}

type Scheduling_Worker_Struct struct {
	Worker Worker_Struct
	Chan   chan Scheduled_Struct
	Error  chan error
}

type Scheduling_Job_Struct struct {
	Job   Job_Struct
	Error error
}

type Enqueue_Worker_Struct struct {
	Scheduling_Worker_Struct
}

type Enqueue_Job_Struct struct {
	Scheduling_Job_Struct
}

type Retry_Struct struct {
	Worker Worker_Struct
	Job    Job_Struct
}

// const definition
const ()

// var definition
var AccessMux sync.Mutex
var WorkerQueue []*Scheduling_Worker_Struct
var JobQueue []*Job_Struct

// func definition
func WorkerEnqueue(arg Enqueue_Worker_Struct) {
	go func() {
		v := AccessController(arg)
		if v == nil {
			arg.Error <- errors.New("Returned nil interface")
			return
		}
	}()
}

func JobEnqueue(arg Enqueue_Job_Struct) error {
}
func Retry() {}

func AccessController(arg AccessController_interface) AccessController_interface {
	AccessMux.Lock()
	var return_value AccessController_interface
	switch v := arg.(type) {
	case Enqueue_Worker_Struct:
		enqueue := v
		var worker = new(Scheduling_Worker_Struct)
		worker.Worker = enqueue.Worker
		worker.Chan = enqueue.Chan
		worker.Error = enqueue.Error
		WorkerQueue = append(WorkerQueue, worker)
		return_value = enqueue
	case Enqueue_Job_Struct:
		enqueue := v
		var job = new(Job_Struct)
		*job = enqueue.Job
		JobQueue = append(JobQueue, job)
		return_value = enqueue
	case Retry_Struct:
	}

	AccessMux.Unlock()
	return return_value
}

func Start() {
	go func() {
		runtime_match.Start()
	}()
}
