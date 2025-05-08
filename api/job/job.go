package job

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// localpackage
import (
	"github.com/hasuburero/mecrm/api/common"
)

// remote package
import ()

type Job struct {
	Mecrm       common.Mecrm
	Job_id      string
	Data1_id    string
	Data2_id    string
	Function_id string
	Runtime     string
	Status      string
}

const (
	UndefinedMecrm = "Mecrm is undefined (nil)\n"
	EmptyArg       = "Empty argument\n"
)

const (
	Finished   = "Finished"
	Processing = "Processing"
	Pending    = "Pending"
)

var (
	Mecrm  *common.Mecrm
	Client *http.Client
)

type Get_Job_Response_Struct struct {
	Job_id      string `json:"job_id"`
	Data1_id    string `json:"data1_id"`
	Data2_id    string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
	Status      string `json:"status"`
}

type Post_Job_Request_Struct struct {
	Data1_id    string `json:"data1_id"`
	Data2_id    string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
	Status      string `json:"status"`
}

type Post_Job_Response_Struct struct {
	Job_id string `json:"job_id"`
}
type Put_Job_Request_Struct struct {
	Job_id      string `json:"job_id"`
	Data1_id    string `json:"data1_id"`
	Data2_id    string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
	Status      string `json:"status"`
}

type Delete_Job_Request_Struct struct{}
type Delete_Job_Response_Struct struct{}

func GetJobRequest(origin, job_id string) (Job, error) {
	req, err := http.NewRequest(http.MethodGet, origin+common.Datapath, nil)
	if err != nil {
		return Job{}, err
	}

	req.Header.Set(common.JobIdHeader, job_id)
	res, err := Client.Do(req)
	if err != nil {
		return Job{}, err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return Job{}, err
	}
	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return Job{}, err
	}
	defer res.Body.Close()

	var ctx Get_Job_Response_Struct
	err = json.Unmarshal(res_body, &ctx)
	if err != nil {
		return Job{}, err
	}

	return Job{Job_id: ctx.Job_id, Data1_id: ctx.Data1_id, Data2_id: ctx.Data2_id, Function_id: ctx.Function_id, Runtime: ctx.Runtime, Status: ctx.Status}, nil
}

func GetJob(job_id, origin string) (*Job, error) {
	if job_id == "" || origin == "" {
		return nil, errors.New(EmptyArg)
	}

	var job *Job
	job = new(Job)
	job_buf, err := GetJobRequest(origin, job_id)
	*job = job_buf
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (self *Job) GetJob() error {
	if self.Job_id == "" || self.Mecrm.Origin == "" {
		return errors.New(EmptyArg)
	}

	job_buf, err := GetJobRequest(self.Mecrm.Origin, self.Job_id)
	if err != nil {
		return err
	}

	self.Job_id = job_buf.Job_id
	self.Data1_id = job_buf.Data1_id
	self.Data2_id = job_buf.Data2_id
	self.Function_id = job_buf.Function_id
	self.Runtime = job_buf.Runtime
	self.Status = job_buf.Status
	return nil
}

func (self *Job) PostJob() (string, error) {
	var json_buf Post_Job_Request_Struct
	json_buf.Data1_id = self.Data1_id
	json_buf.Data2_id = self.Data2_id
	json_buf.Function_id = self.Function_id
	json_buf.Runtime = self.Runtime
	json_buf.Status = self.Status
	buf, err := json.Marshal(json_buf)
	if err != nil {
		return "", err
	}

	req_buf := bytes.NewBuffer(buf)

	req, err := http.NewRequest(http.MethodPost, self.Mecrm.Origin+common.Datapath, req_buf)
	if err != nil {
		return "", err
	}

	res, err := Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return "", err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var ctx Post_Job_Response_Struct
	err = json.Unmarshal(res_body, &ctx)
	if err != nil {
		return "", err
	}

	self.Job_id = ctx.Job_id
	return self.Job_id, nil
}

func (self *Job) PutJob() error {
	var json_buf Put_Job_Request_Struct
	json_buf.Job_id = self.Job_id
	json_buf.Data1_id = self.Data1_id
	json_buf.Data2_id = self.Data2_id
	json_buf.Function_id = self.Function_id
	json_buf.Runtime = self.Runtime
	json_buf.Status = self.Status

	buf, err := json.Marshal(json_buf)
	if err != nil {
		return err
	}

	req_buf := bytes.NewBuffer(buf)
	req, err := http.NewRequest(http.MethodPut, self.Mecrm.Origin+common.Datapath, req_buf)
	if err != nil {
		return err
	}

	res, err := Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return err
	}

	return nil
}

func DeleteJobRequest(job_id, origin string) error {
	req, err := http.NewRequest(http.MethodDelete, origin+common.Datapath, nil)
	if err != nil {
		return err
	}

	req.Header.Set(common.JobIdHeader, job_id)

	res, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return err
	}

	return nil
}

func DeleteJob(job_id, origin string) error {
	if job_id == "" || origin == "" {
		return errors.New(EmptyArg)
	}

	err := DeleteJobRequest(job_id, origin)
	return err
}

func (self *Job) DeleteJob() error {
	if self.Job_id == "" || self.Mecrm.Origin == "" {
		return errors.New(EmptyArg)
	}

	err := DeleteJobRequest(self.Job_id, self.Mecrm.Origin)
	return err
}

func Init(scheme, addr, port string) error {
	Mecrm = new(common.Mecrm)
	Mecrm.Scheme = scheme
	Mecrm.Addr = addr
	Mecrm.Port = port
	Mecrm.Origin = scheme + "://" + addr + ":" + port

	Client = &http.Client{}
	return nil
}

func MakeJob(runtime string) (*Job, error) {
	if Mecrm == nil {
		return nil, errors.New(UndefinedMecrm)
	}
	job := new(Job)
	job.Mecrm.Scheme = Mecrm.Scheme
	job.Mecrm.Addr = Mecrm.Addr
	job.Mecrm.Port = Mecrm.Port
	job.Mecrm.Origin = Mecrm.Scheme + "://" + Mecrm.Addr + ":" + Mecrm.Port
	job.Runtime = runtime
	job.Status = Pending

	return job, nil
}
