package api

import (
	"io"
	"net/http"
)

func (self *API) DELETE_data_data_id() (int, string, error) {
	req, err := http.NewRequest("DELETE", self.URL+"/data/"+self.Data_id, nil)
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

	return res.StatusCode, string(resBody), err
}
