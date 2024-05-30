package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type POST_job_job_id_request struct {
	State    string `json:"status"`
	Output_id string `json:"output"`
}

type POST_job_job_id_response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (self *API) POST_job_job_id() (int, string, error) {
	if self.Job_status == Failed {
		self.Job_output_id = "0"
	}
	worker_request := POST_job_job_id_request{State: self.Job_status, Output_id: self.Job_output_id}
	json_buf, err := json.Marshal(worker_request)
	if err != nil {
		return -1, "", err
	}

	body := bytes.NewBuffer(json_buf)
	req, err := http.NewRequest("POST", self.URL+"/job/"+self.Job_id, body)
	if err != nil {
		return -1, "", err
	}
	req.Header.Add("content-type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, "", err
	}
	defer res.Body.Close()

	var ctx POST_job_job_id_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, "", err
	}

	return ctx.Code, string(resBody), err
}
