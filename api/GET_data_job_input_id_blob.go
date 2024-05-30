package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type GET_data_data_id_blob_response struct {
	Code    int    `json:"code"`
	Status  string `json:"stauts"`
	Message string `json:"message"`
}

func (self *API) GET_data_job_input_id_blob() (int, string, error) {
	req, err := http.NewRequest("GET", self.URL+"/data/"+self.Job_input_id+"/blob", nil)
	if err != nil {
		return -1, "", err
	}

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

	var ctx GET_data_data_id_blob_response
	if res.Header.Get("content-type") == "application/octet-stream" {
		self.Input_data = resBody
	} else {
		err = json.Unmarshal(resBody, &ctx)
		if err != nil {
			return -1, "", err
		}
		return ctx.Code, "", err
	}

	return ctx.Code, "", err
}
