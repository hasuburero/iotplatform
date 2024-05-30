package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type POST_data_response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Id      string `json:"id"`
	Hash    string `json:"hash"`
}

func (self *API) POST_data() (int, string, error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile("file", "k_output.txt")
	if err != nil {
		return -1, "CreateFormFile error", err
	}

	_, err = io.Copy(fw, bytes.NewBuffer(self.Output_data))
	if err != nil {
		return -1, "io.Copy error", err
	}

	content_type := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		return -1, "FormDataContentType error", err
	}

	req, err := http.NewRequest("POST", self.URL+"/data", body)
	if err != nil {
		return -1, "NewRequest error", err
	}

	req.Header.Add("content-type", content_type)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return -1, "request error", err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, "io.ReadAll error", err
	}
	defer res.Body.Close()

	var ctx POST_data_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, string(resBody), err
	}

	self.Job_output_id = ctx.Id
	return ctx.Code, string(resBody), err
}
