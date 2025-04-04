package sched

import (
	"github.com/hasuburero/util/panic"
	"time"
	"sched/runtime_match"
	"sync"
)

type AccessController_interface interface{}

type Worker_Struct struct{
	Worker_id string
	Runtime []string
}

type Job_Struct struct{
	Job_id string
	Data1_id string
	Data2_id string
	Function_id string
	Runtime string
}

type Scheduling_Struct struct {
	Worker Worker_Struct
	Job    Job_Struct
}

type Scheduled_Worker_Struct struct {
	Worker Worker_Struct
	Chan   chan Scheduling_Struct
	Error  error
}

type Enqueue_Job_Struct struct {
	Job   Job_Struct
	Error error
}

type Retry_Struct struct{
	Worker Worker_Struct
	Job Job_Struct
}

// var definition
var AccessMux sync.Mutex

// func definition
func WorkerScheduling(arg Worker_Struct) (Scheduling_Struct, error) {
	var worker_sched = Scheduled_Worker_Struct{Worker: arg, Chan: make(chan Scheduling_Struct)}
	select{
	case <-worker_sched.Chan:
	case time.
	}
	return 
}

func JobEnqueue(arg Scheduling_Struct) error {}

func AccessController(arg AccessController_interface) AccessController_interface {
	AccessMux.Lock()
	var return_value AccessController_interface
	switch v := arg.(type) {
	case Scheduled_Worker_Struct:
	case Enqueue_Job_Struct:
	}

	AccessMux.Unlock()
	return return_value
}

func Start() {
	go func() {
		runtime_match.Start()
	}()
}
