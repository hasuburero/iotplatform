package job

import (
	"errors"
	"fmt"
	"rm/data"
	"strconv"
	"sync"
	"time"
)

// type definition
type Job_Struct struct {
	Job_id      string
	Data1_id    string
	Data2_id    string
	Function_id string
	Runtime     string
	Status      string
	TimeStamp   time.Time
}

type AccessController_interface interface{}

type Add_Job_Struct struct {
	Job_id      string
	Data1_id    string
	Data2_id    string
	Function_id string
	Runtime     string
	Timestamp   time.Time
	Error       error
}

type Get_Job_Struct struct {
	Job_id      string
	Data1_id    string
	Data2_id    string
	Function_id string
	Runtime     string
	Status      string
	Timestamp   time.Time
	Error       error
}

type Check_Struct struct {
	Job_id string
	Error  error
}

type Delete_Job_Struct struct {
	Job_id string
	Error  error
}

type Update_Job_Struct struct {
	Job_id string
	Status string
	Error  error
}

// const definition
const (
	Job_Id_Size  = 10
	Job_Id_Limit = uint32(1024*1024*1024*4 - 1)
)

const (
	Finished   = "Finished"
	Processing = "Processing"
	Pending    = "Pending"
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
	StatusMap[Processing] = 2
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

func Check(job_id string) error {
	var check Check_Struct
	check.Job_id = job_id
	v := AccessController(check)
	if v == nil {
		return errors.New("Returned nil interface")
	}
	return check.Error
}

func JobUpdate(job_id, status string) error {
	var update Update_Job_Struct
	update.Job_id = job_id
	update.Status = status
	v := AccessController(update)
	if v == nil {
		return errors.New("Returned nil interface")
	}
	update = v.(Update_Job_Struct)
	return update.Error
}

func AddJob(data1_id, function_id, runtime string, ts time.Time) (string, error) {
	data2_id, err := data.DataRegPost()
	if err != nil {
		return "", err
	}

	var add_job Add_Job_Struct
	add_job.Data1_id = data1_id
	add_job.Data2_id = data2_id
	add_job.Function_id = function_id
	add_job.Runtime = runtime
	add_job.Timestamp = ts
	v := AccessController(add_job)
	if v == nil {
		return "", errors.New("Returned nil interface")
	}

	add_job = v.(Add_Job_Struct)
	return add_job.Job_id, add_job.Error
}

func JobDelete(job_id string) error {
	var del_job Delete_Job_Struct
	del_job.Job_id = job_id
	v := AccessController(del_job)
	if v == nil {
		return errors.New("Returned nil interface")
	}
	del_job = v.(Delete_Job_Struct)
	if del_job.Error != nil {
		return del_job.Error
	}

	err := DeletePair(del_job.Job_id)
	return err
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

			job := Job[add_job.Job_id]
			if job == nil {
				add_job.Error = nil
				return_value = add_job

				job_buf := new(Job_Struct)
				job_buf.Job_id = add_job.Job_id
				job_buf.Data1_id = add_job.Data1_id
				job_buf.Data2_id = add_job.Data2_id
				job_buf.Function_id = add_job.Function_id
				job_buf.Runtime = add_job.Runtime
				job_buf.TimeStamp = time.Now()
				job_buf.Status = Pending
				Job[job_buf.Job_id] = job_buf

				EnqueueJob(*job_buf)
				break
			}
		}
	case Get_Job_Struct:
		get_job := v
		job_buf := Job[get_job.Job_id]
		if job_buf == nil {
			get_job.Error = errors.New("No such Job")
			return_value = get_job
		} else {
			get_job.Error = nil
			get_job.Data1_id = job_buf.Data1_id
			get_job.Data2_id = job_buf.Data2_id
			get_job.Function_id = job_buf.Function_id
			get_job.Runtime = job_buf.Runtime
			get_job.Status = job_buf.Status
			return_value = get_job
		}
	case Delete_Job_Struct:
		del_job := v
		_, exists := Job[del_job.Job_id]
		if !exists {
			del_job.Error = errors.New("No such Job")
			return_value = del_job
			break
		}
		//wip
		delete(Job, del_job.Job_id)
		DeletePair(del_job.Job_id)
	case Update_Job_Struct:
		update := v
		job_buf := Job[update.Job_id]
		if job_buf == nil {
			update.Error = errors.New("no such job")
			return_value = update
			break
		}

		switch update.Status {
		case Finished:
			switch job_buf.Status {
			case Processing:
				job_buf.Status = update.Status
			case Finished:
				fallthrough
			case Pending:
				fallthrough
			default:
				update.Error = errors.New("Invalid Status Transition")
			}
		case Processing:
			switch job_buf.Status {
			case Pending:
				job_buf.Status = update.Status
			case Finished:
				fallthrough
			case Processing:
				fallthrough
			default:
				update.Error = errors.New("Invalid Status Transition")
			}
		case Pending:
			update.Error = errors.New("Invalid Status Transition")
		default:
			update.Error = errors.New("no such status")
		}
		return_value = update
	default:
		fmt.Println("passed default (job.AccessController)")
		return_value = nil
	}
	AccessMux.Unlock()
	return return_value
}
