package job

import (
	"errors"
	"fmt"
	"rm/sched"
	"strconv"
	"sync"
	"time"
)

// type definition
type Job_Struct struct {
	Job_id      string
	Data_id     string
	Function_id string
	Runtime     string
	Status      string
	TimeStamp   time.Time
}

type AccessController_interface interface{}

type Add_Job_Struct struct {
	Job_id      string
	Data_id     string
	Function_id string
	Runtime     string
	Error       error
}

type Get_Job_Struct struct {
	Job_id      string
	Data_id     string
	Function_id string
	Runtime     string
	Status      string
	Error       error
}

type Delete_Job_Struct struct{
	Job_id string
	Error error
}

// const definition

const (
	Job_Id_Size  = 10
	Job_Id_Limit = uint32(1024*1024*1024*4 - 1)
)

const (
	Finished = "Finished"
	Running  = "Running"
	Pending  = "Pending"
)

// var definition

var AccessMux sync.Mutex
var global_job_id uint32 = 0
var StatusMap map[string]int8
var RevStatusMap map[int8]string
var Job map[string]*Job_Struct

func init() {
	StatusMap = make(map[string]int8)
	RevStatusMap = make(map[int8]string)
	StatusMap[Pending] = 1
	StatusMap[Running] = 2
	StatusMap[Finished] = 3
	for key, value := range StatusMap {
		RevStatusMap[value] = key
	}
}

func GenerateJobId() (string, error) {
	var zeronum int = 0
	var prefix uint32 = 10
	var id_buf = global_job_id
	if id_buf > Job_Id_Limit {
		return "", errors.New("Data_id_limit overrun")
	}
	global_job_id++
	for i := 0; i < Job_Id_Size-1; i++ {
		if id_buf < prefix {
			zeronum = Job_Id_Size - 1 - i
			break
		}
		prefix *= 10
	}

	var job_id string = ""
	for range zeronum {
		job_id += "0"
	}
	job_id += strconv.FormatUint(uint64(id_buf), 10)

	return job_id, nil
}

func AccessController(arg AccessController_interface) AccessController_interface {
	AccessMux.Lock()
	var return_value AccessController_interface
	switch v := arg.(type) {
	case Add_Job_Struct:
		add_job := v
		for {
			id, err := GenerateJobId()
			if err != nil {
				add_job.Error = err
				return_value = add_job
				break
			}
			add_job.Job_id = id

			_, exists := Job[add_job.Job_id]
			if !exists {
				add_job.Error = nil
				return_value = add_job

				job_buf := new(Job_Struct)
				job_buf.Job_id = add_job.Job_id
				job_buf.Data_id = add_job.Data_id
				job_buf.Function_id = add_job.Function_id
				job_buf.Runtime = add_job.Runtime
				job_buf.TimeStamp = time.Now()
				job_buf.Status = "Pending"
				Job[job_buf.Job_id] = job_buf


				// wip
				sched.
				break
			}
		}
	case Get_Job_Struct:
		get_job := v
		job_buf, exists := Job[get_job.Job_id]
		if !exists {
			get_job.Error = errors.New("No such Job")
			return_value = get_job
		} else {
			get_job.Error = nil
			get_job.Data_id = job_buf.Data_id
			get_job.Function_id = job_buf.Function_id
			get_job.Runtime = job_buf.Runtime
			get_job.Status = job_buf.Status
			return_value = get_job
		}
	case Delete_Job_Struct:
		del_job := v
		job_buf, exists := Job[del_job.Job_id]
		if !exists{
			job_buf.Error = errors.New("No such Job")	
			return_value = del_job
			break
		}
		sched.De
	default:
		fmt.Println("passed default (job.AccessController)")
		return_value = nil
	}
	AccessMux.Unlock()
	return return_value
}
