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
	utl     = ""
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

	code, res, err = Req.POST_lambda()
	if err != nil {
		fmt.Println(err)
		fmt.Println("POST_lambda error")
		return
	}
	fmt.Println(res)
	fmt.Println(Req.Lambda_id)

}
