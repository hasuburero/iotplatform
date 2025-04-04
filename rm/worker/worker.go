package worker

import (
	"errors"
	"fmt"
	//"rm/job"
	"rm/sched"
	"strconv"
	"sync"
)

type Worker_Struct struct {
	Mux       sync.Mutex
	Worker_id string
	Runtime   []string
}

type AccessController_interface interface {
}

type Get_Contract_Struct struct {
	Worker_id   string
	Job_id      string
	Data1_id    string
	Data2_id    string
	Function_id string
	Runtime     string
}

type Post_Worker_Struct struct {
	Worker_id string
	Runtime   []string
	Error     error
}

type Get_Worker_Struct struct {
	Worker_id string
	Runtime   []string
	Error     error
}

type Delete_Worker_Struct struct {
	Worker_id string
	Error     error
}

// const definition
const (
	Worker_id_size  = 10
	Worker_id_limit = uint32(1024*1024*1024*4 - 1)
)

// var definition
var AccessMux sync.Mutex
var global_worker_id uint32 = 0
var Worker map[string]*Worker_Struct

func WorkerCopy(src *Worker_Struct) sched.Worker_Struct {
	var worker sched.Worker_Struct
	worker.Worker_id = src.Worker_id
	worker.Runtime = make([]string, len(src.Runtime))
	return worker
}

func GenerateWorkerId() (string, error) {
	var zeronum int = 0
	var prefix uint32 = 10
	var id_buf = global_worker_id
	if id_buf > Worker_id_limit {
		return "", errors.New("Worker_id_limit overrun")
	}
	global_worker_id++
	for i := 0; i < Worker_id_size-1; i++ {
		if id_buf < prefix {
			zeronum = Worker_id_size - 1 - i
			break
		}
		prefix *= 10
	}

	var worker_id string = ""
	for range zeronum {
		worker_id += "0"
	}
	worker_id += strconv.FormatUint(uint64(id_buf), 10)

	return worker_id, nil
}

func Contract(worker_id string) (Get_Contract_Struct, error) {
	var contract sched.Scheduled_Worker_Struct
	contract.Chan = make(chan sched.Scheduling_Struct)
	contract.Worker.Worker_id = worker_id
	v := AccessController(contract)
	if v == nil {
		return Get_Contract_Struct{}, errors.New("Returned nil interface")
	}
	contract = v.(sched.Scheduled_Worker_Struct)
	if contract.Error != nil {
		return Get_Contract_Struct{}, contract.Error
	}

	ctx := <-contract.Chan
	var result = Get_Contract_Struct{
		Worker_id:   ctx.Worker.Worker_id,
		Job_id:      ctx.Job.Job_id,
		Data1_id:    ctx.Job.Data1_id,
		Data2_id:    ctx.Job.Data2_id,
		Function_id: ctx.Job.Function_id,
		Runtime:     ctx.Job.Runtime,
	}
	return result, nil
}

func WorkerDelete() error {}
func WorkerGet() error    {}
func WorkerPost() error   {}

func AccessController(arg AccessController_interface) AccessController_interface {
	AccessMux.Lock()
	var return_value AccessController_interface
	switch v := arg.(type) {
	case Get_Contract_Struct:
		get_contract := v
		worker, exists := Worker[v.Worker_id]
		if !exists {
			get_contract.Error = errors.New("no such worker")
			return_value = get_contract
			break
		}

		workercopy := WorkerCopy(worker)
		result, err := sched.WorkerScheduling(workercopy)
		if err != nil {
			get_contract.Error = err
			break
		}

		get_contract.Job_id = result.Job.Job_id
		get_contract.Data_id = result.Job.Data_id
		get_contract.Function_id = result.Job.Function_id
		get_contract.Runtime = result.Job.Runtime
		return_value = get_contract
	case Delete_Worker_Struct:
	case Get_Worker_Struct:
	case Post_Worker_Struct:
	}

	AccessMux.Unlock()
	return return_value
}
