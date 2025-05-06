package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hasuburero/mecrm/api/common"
	"io"
	"net/http"
)

type Worker struct {
	Addr      string
	Port      string
	Worker_id string
	Runtime   []string
	Client    *http.Client
	Scheme    string
	Url       string
}

type Post_Worker_Request_Struct struct {
	Runtime []string `json:"runtime"`
}
type Post_Worker_Response_Struct struct {
	Worker_id string `json:"worker_id"`
}

type Get_Worker_Request_Struct struct {
	Worker_id string `json:"worker_id"`
}
type Get_Worker_Response_Struct struct {
	Worker_id string   `json:"worker_id"`
	Runtime   []string `json:"runtime"`
}

type Get_Worker_Contract_Request_Struct struct {
	Worker_id string `json:"worker_id"`
}
type Get_Worker_Contract_Response_Struct struct {
	Worker_id   string `json:"worker_id"`
	Job_id      string `json:"job_id"`
	Data1_id    string `json:"data1_id"`
	Data2_id    string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
}

const ()

func (self *Worker) PostWorker() (Post_Worker_Response_Struct, error) {
	req_json := Post_Worker_Request_Struct{}
	req_json.Runtime = make([]string, len(self.Runtime))
	copy(req_json.Runtime, self.Runtime)
	req_buf, err := json.Marshal(req_json)
	if err != nil {
		return Post_Worker_Response_Struct{}, err
	}

	req_body := bytes.NewBuffer(req_buf)
	req, err := http.NewRequest(http.MethodPost, self.Scheme+self.Addr+self.Port+common.Workerpath, req_body)
	if err != nil {
		return Post_Worker_Response_Struct{}, err
	}

	req.Header.Set(common.ContentType, common.ApplicationJson)

	res, err := self.Client.Do(req)
	if err != nil {
		return Post_Worker_Response_Struct{}, err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return Post_Worker_Response_Struct{}, err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return Post_Worker_Response_Struct{}, err
	}
	defer res.Body.Close()

	var ctx Post_Worker_Response_Struct
	err = json.Unmarshal(res_body, &ctx)
	if err != nil {
		return Post_Worker_Response_Struct{}, err
	}

	return ctx, nil
}

func (self *Worker) GetWorker(worker_id string) (Get_Worker_Response_Struct, error) {
	req_json := Get_Worker_Request_Struct{Worker_id: worker_id}
	json_buf, err := json.Marshal(req_json)
	if err != nil {
		return Get_Worker_Response_Struct{}, err
	}

	req_body := bytes.NewBuffer(json_buf)
	req, err := http.NewRequest(http.MethodGet, self.Scheme+self.Addr+self.Port+common.Workerpath, req_body)
	if err != nil {
		return Get_Worker_Response_Struct{}, err
	}

	req.Header.Set(common.ContentType, common.ApplicationJson)
	req.Header.Set(common.WorkerIdHeader, worker_id)

	res, err := self.Client.Do(req)
	if err != nil {
		return Get_Worker_Response_Struct{}, err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return Get_Worker_Response_Struct{}, err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return Get_Worker_Response_Struct{}, err
	}
	defer res.Body.Close()

	var ctx Get_Worker_Response_Struct
	err = json.Unmarshal(res_body, &ctx)
	if err != nil {
		return Get_Worker_Response_Struct{}, err
	}

	return ctx, nil
}

func (self *Worker) Delete() {
}

func (self *Worker) Contract() (Get_Worker_Contract_Response_Struct, error) {
	req_json := Get_Worker_Contract_Request_Struct{Worker_id: self.Worker_id}
	json_buf, err := json.Marshal(req_json)
	if err != nil {
		return Get_Worker_Contract_Response_Struct{}, err
	}

	req_body := bytes.NewBuffer(json_buf)

	req, err := http.NewRequest(http.MethodGet, self.Scheme+self.Addr+self.Port+common.Contractpath, req_body)
	if err != nil {
		return Get_Worker_Contract_Response_Struct{}, err
	}

	req.Header.Set(common.WorkerIdHeader, self.Worker_id)
	req.Header.Set(common.ContentType, common.ApplicationJson)

	res, err := self.Client.Do(req)
	if err != nil {
		return Get_Worker_Contract_Response_Struct{}, err
	}
	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return Get_Worker_Contract_Response_Struct{}, err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return Get_Worker_Contract_Response_Struct{}, err
	}
	defer res.Body.Close()

	var ctx Get_Worker_Contract_Response_Struct
	err = json.Unmarshal(res_body, &ctx)
	if err != nil {
		return Get_Worker_Contract_Response_Struct{}, err
	}

	return ctx, nil
}

func MakeWorker(addr, port, scheme string, runtimes []string) *Worker {
	worker := new(Worker)
	worker.Addr = addr
	worker.Port = ":" + port
	worker.Runtime = make([]string, len(runtimes))
	copy(worker.Runtime, runtimes)
	worker.Client = &http.Client{}
	worker.Scheme = scheme
	worker.Url = worker.Scheme + worker.Addr + worker.Port

	return worker
}
