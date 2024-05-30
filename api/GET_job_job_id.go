package api

import (
	"encoding/json"
	"io"
	"net/http"
	"fmt"
)

type GET_job_job_id_request struct {
	Job_id string `json:"job_id"`
}

type RespJobLambda struct {
	Id      string  `json:"id"`
	Runtime string `json:"runtime"`
	Code    string  `json:"codex"`
}

type RespData struct {
	Id string `json:"id"`
	Blob string `json:"blob"`
}

type GET_job_job_id_response struct {
	Code           int      `json:"code"`
	Message        string   `json:"message"`
	Job_id         string   `json:"id"`
	//Job_status     string   `json:"job_status"`
	//Job_input_id   string   `json:"job_input_id"`
	//Job_output_id  string   `json:"job_output_id"`
	//Lambda_id      string   `json:"functio"`
	//Lambda_runtime string   `json:"runtime"`
	Input          RespData `json:"input"`
	Output         RespData `json:"output"`
	Lambda RespJobLambda `json:"lambda"`
	State string `json:"state"`
}

func (self *API) GET_job_job_id() (int, string, error) {
	req, err := http.NewRequest("GET", self.URL+"/job/"+self.Job_id, nil)
	if err != nil {
		return -1, "", err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return -1, "request error", err
	}

	fmt.Println("status code", res.StatusCode)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, "read all error", err
	}
	defer res.Body.Close()

	var ctx GET_job_job_id_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, string(resBody), err
	}

	self.Job_id = ctx.Job_id
	self.Job_input_id = ctx.Input.Id
	self.Job_output_id = ctx.Output.Id
	self.Job_status = ctx.State
	self.Lambda_id = ctx.Lambda.Id
	self.Lambda_runtime = ctx.Lambda.Runtime
	return ctx.Code, string(resBody), err
}
