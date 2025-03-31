package worker

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Worker_Struct struct {
	Worker_id string
	Runtime   []string
}

type AccessController_interface interface {
}

type Get_Contract_Struct struct {
	Worker_id   string
	Job_id      string
	Data_id     string
	Function_id string
	Error       error
}

type Get_Worker_Struct struct {
	Worker_id string
}

type Delete_Worker_Struct struct {
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

func AccessController(arg AccessController_interface) AccessController_interface {
	AccessMux.Lock()
	var return_value AccessController_interface
	switch v := arg.(type) {
	case Get_Contract_Struct:
	case Delete_Worker_Struct:
	case Get_Worker_Struct:
	}

	AccessMux.Unlock()
	return return_value
}
