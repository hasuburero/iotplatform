package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type POST_lambda_request struct {
	Id      string `json:"codex"`
	Runtime string `json:"runtime"`
}

type POST_lambda_response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Id     string `json:"id"`
}

func (self *API) POST_lambda() (int, string, error) {
	self.Lambda_id = self.Job_output_id
	request := POST_lambda_request{Id: self.Lambda_id, Runtime: self.Runtime}
	json_buf, err := json.Marshal(request)
	if err != nil {
		return -1, "", err
	}

	body := bytes.NewBuffer(json_buf)
	req, err := http.NewRequest("POST", self.URL+"/lambda", body)
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

	var ctx POST_lambda_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, "", err
	}

	self.Lambda_id = ctx.Id
	return ctx.Code, string(resBody), err
}
