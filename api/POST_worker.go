package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type POST_worker_request struct {
	Runtimes []string `json:"runtime"`
}

type POST_worker_response struct {
	Code      int      `json:"code"`
	Status    string   `json:"status"`
	Message   string   `json:"message"`
	Worker_id string   `json:"id"`
	Runtimes  []string `json:"runtime"`
}

func (self *API) POST_worker() (int, string, error) {
	worker_request := POST_worker_request{self.Runtimes}
	json_buf, err := json.Marshal(worker_request)
	if err != nil {
		return -1, "", err
	}
	body := bytes.NewBuffer(json_buf)
	req, err := http.NewRequest("POST", self.URL+"/worker", body)
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

	var ctx POST_worker_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, "", err
	}

	self.Worker_id = ctx.Worker_id
	return ctx.Code, string(resBody), err
}
