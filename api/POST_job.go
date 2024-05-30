package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type POST_job_request struct {
	Input_id  string   `json:"input"`
	Lambda_id string   `json:"lambda"`
	Tags      []string `json:"tags"`
}

type POST_job_response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Job_id  string `json:"id"`
}

func (self *API) POST_job() (int, string, error) {
	request := POST_job_request{Input_id: self.Job_output_id, Lambda_id: self.Lambda_id}
	json_buf, err := json.Marshal(request)
	if err != nil {
		return -1, "", err
	}

	body := bytes.NewBuffer(json_buf)
	req, err := http.NewRequest("POST", self.URL+"/job", body)
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

	var ctx POST_job_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, "", err
	}

	self.Job_id = ctx.Job_id
	return ctx.Code, string(resBody), err
}
