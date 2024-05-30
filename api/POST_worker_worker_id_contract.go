package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type POST_worker_worker_id_contract_request struct {
	Worker_id string   `json:"id"`
	Extra_tag []string `json:"tags"`
	Timeout   int      `json:"timeout"`
}

type POST_worker_worker_id_contract_response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Job_id  string `json:"job"`
}

func (self *API) POST_worker_worker_id_contract(timeout int) (int, string, error) {
	worker_request := POST_worker_worker_id_contract_request{Worker_id: self.Worker_id, Timeout: timeout}
	json_buf, err := json.Marshal(worker_request)
	if err != nil {
		return -1, "", err
	}

	body := bytes.NewBuffer(json_buf)
	req, err := http.NewRequest("POST", self.URL+"/worker/"+self.Worker_id+"/contract", body)
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

	var ctx POST_worker_worker_id_contract_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, "", err
	}

	self.Job_id = ctx.Job_id
	return ctx.Code, string(resBody), err
}
