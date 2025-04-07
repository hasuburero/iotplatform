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

type Get_Contract_Struct sched.Scheduling_Struct

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
	copy(worker.Runtime, src.Runtime)
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

func Contract(worker_id string) (sched.Scheduled_Struct, error) {
	var contract sched.Enqueue_Worker_Struct
	contract.Chan = make(chan sched.Scheduled_Struct, 1)
	contract.Error = make(chan error, 1)
	contract.Worker.Worker_id = worker_id
	v := AccessController(contract)
	if v == nil {
		return sched.Scheduled_Struct{}, errors.New("Returned nil interface")
	}
	contract = v.(sched.Enqueue_Worker_Struct)

	select {
	case ch := <-contract.Chan:
		return ch, nil
	case err := <-contract.Error:
		return sched.Scheduled_Struct{}, err
	}
}

func WorkerDelete(worker_id string) error {
	var worker Delete_Worker_Struct
	worker.Worker_id = worker_id
	v := AccessController(worker)
	if v == nil {
		return errors.New("Returned nil interface")
	}
	worker = v.(Delete_Worker_Struct)
	return worker.Error
}

func WorkerGet(worker_id string) (Get_Worker_Struct, error) {
	var worker Get_Worker_Struct
	worker.Worker_id = worker_id
	v := AccessController(worker)
	if v == nil {
		return Get_Worker_Struct{}, errors.New("Returned nil interface")
	}
	worker = v.(Get_Worker_Struct)
	return worker, worker.Error
}

func WorkerPost() error {}

func AccessController(arg AccessController_interface) AccessController_interface {
	AccessMux.Lock()
	var return_value AccessController_interface
	switch v := arg.(type) {
	case sched.Enqueue_Worker_Struct:
		get_contract := v
		worker, exists := Worker[get_contract.Worker.Worker_id]
		if !exists {
			get_contract.Error <- errors.New("no such worker")
			return_value = get_contract
			break
		}

		get_contract.Worker = WorkerCopy(worker)
		sched.WorkerEnqueue(get_contract)

		return_value = get_contract
	case Delete_Worker_Struct:
		delete_worker := v
		_, exists := Worker[delete_worker.Worker_id]
		if !exists {
			delete_worker.Error = errors.New("no such worker")
			return_value = delete_worker
			break
		}
		Worker[delete_worker.Worker_id] = nil
		return_value = delete_worker
	case Get_Worker_Struct:
		get_worker := v
		worker, exists := Worker[get_worker.Worker_id]
		if !exists {
			get_worker.Error = errors.New("no such worker")
			return_value = get_worker
			break
		}
		get_worker.Runtime = make([]string, len(worker.Runtime))
		return_value = get_worker
	case Post_Worker_Struct:
	default:
		return_value = nil
	}

	AccessMux.Unlock()
	return return_value
}
