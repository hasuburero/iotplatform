package main

import (
	"fmt"
	"io"
	"k/api"
	"os"
)

const (
	runtime = "pynq_k"
	path    = "input_data_lastflag_6.txt"
	url = "http://172.21.39.32:8332/api/v0.5"
	//url     = "https://mecrm.dolylab.cc/api/v0.5"
)

func main() {
	Req := api.API{
		Runtime: runtime,
		URL:     url,
	}
	fd, err := os.Open(path)
	Req.Output_data, err = io.ReadAll(fd)
	code, res, err := Req.POST_data()
	if err != nil {
		fmt.Println(err)
		fmt.Println(res)
		fmt.Println("POST_data error")
		return
	} else if code != 0 {
		fmt.Println(err)
		fmt.Println(res)
		fmt.Println("POST_data error with ", code)
		return
	}

	fmt.Println(Req.Job_output_id)
}
