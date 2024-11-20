package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GET_ping_response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (self *API) GET_ping(arg string) (int, string, error) {
	req, err := http.NewRequest("GET", self.URL+"/ping?message="+arg, nil)
	if err != nil {
		return -1, "", err
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, "read All error", err
	}

	var ctx GET_ping_response
	err = json.Unmarshal(resBody, &ctx)
	if err != nil {
		return -1, string(resBody), err
	}
	return ctx.Code, ctx.Message, nil
}
