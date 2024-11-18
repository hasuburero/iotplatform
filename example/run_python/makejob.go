package main

import (
	"fmt"
	"io/ioutil"
	"k/api"
	"os"
	"time"
)

const (
	runtime = "pynq_k"
	url     = "http://172.21.39.32:8332/api/v0.5"
	//url     = "https://mecrm.dolylab.cc/api/v0.5"
)

func main() {
	job := api.API{
		Runtime: runtime,
		URL:     url,
	}
	if len(os.Args) < 3 {
		fmt.Println("invalid Args")
		return
	}

	job.Lambda_id = os.Args[1]
	job.Job_output_id = os.Args[2]
	code, res, err := job.POST_job()
	if err != nil {
		fmt.Println(err)
		fmt.Println("POST_job error")
		return
	}
	fmt.Println(code, job.Job_id)
	fmt.Println(res)
	fmt.Println("job created")

	for {
		code, res, err := job.GET_job_job_id()
		fmt.Println(res)
		fmt.Println("Job_id", job.Job_id)
		fmt.Println("Job_input_id", job.Job_input_id)
		fmt.Println("Job_output_id", job.Job_output_id)
		if err != nil {
			fmt.Println(err)
			fmt.Println("get job_job_id error")
			continue
		}
		if code != 0 {
			fmt.Println("code: ", code)
			continue
		}

		if job.Job_status == "Finished" {
			fmt.Println(string(res))
			break
		} else if job.Job_status == "Fail_Execution" {
			fmt.Println(string(res))
			fmt.Println("job failed")
			return
		}
		fmt.Println(job.Job_status)
		fmt.Println("continue")

		time.Sleep(2 * time.Second)
	}
	job.Job_input_id = job.Job_output_id

	code, res, err = job.GET_data_job_input_id_blob()
	if err != nil {
		fmt.Println(err)
		fmt.Println(code)
		fmt.Println("GET_data_job_input_id_blob error")
		return
	}

	fmt.Println(res)
	err = ioutil.WriteFile("all_result.txt", job.Input_data, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Println("write file error")
		return
	}
	fmt.Println("file write completed")
	return
}
