package main

import (
	"fmt"
	"io"
	"k/api"
	"os"
)

const (
	runtime = ""
	path    = ""
	url = ""
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
