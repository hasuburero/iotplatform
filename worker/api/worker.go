package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

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

const (
	contractpath = "/mecrm/worker/contract"
)

func (self *Worker) Contract() (Get_Worker_Contract_Response_Struct, error) {
	body := &bytes.Buffer{}
	req_json := Get_Worker_Contract_Request_Struct{Worker_id: self.Worker_id}
	json_buf, err := json.Marshal(req_json)
	if err != nil {
		return Get_Worker_Contract_Response_Struct{}, err
	}

	req, err := http.NewRequest(http.MethodGet, self.Scheme+self.Addr+self.Port+contractpath)
	req, err := self.Client.Do(req)
}

func (self *Worker) PostWorker() {

}

func (self *Worker) Delete() {

}
