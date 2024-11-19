package main

import (
	"fmt"
	"github.com/hasuburero/mecrm/api"
	"time"
)

const (
	thread_num    = 30
	job_num       = 700
	url           = "https://mecrm.dolylab.cc/api/v0.5"
	timeout       = 10
	runtime       = "al20067_worker_test"
	code_contract = 1299
)

var blob_id string
var lambda_id string
var worker_counter chan int
var requester_counter chan int

func counterThread() {
	worker_count := 0
	requester_count := 0
	for {
		select {
		case <-worker_counter:
			worker_count++
		case <-requester_counter:
			requester_count++
		case <-time.After(time.Duration(1) * time.Second):
			fmt.Print("worker_count/sec: ")
			fmt.Println(worker_count)
			fmt.Print("requester_count/sec: ")
			fmt.Println(requester_count)
			fmt.Println("")
			worker_count = 0
			requester_count = 0
		}
	}
}

func ready() (int, error) {
	r4rf := api.API{
		Runtime: runtime,
		URL:     url,
	}
	r4rf.Output_data = []byte("test")
	code, _, err := r4rf.POST_data()
	if err != nil {
		return code, err
	} else if code != 0 {
		return code, nil
	}
	blob_id = r4rf.Job_output_id
	r4rf.Output_data = []byte("test")
	code, _, err = r4rf.POST_data()
	if err != nil {
		return code, err
	} else if code != 0 {
		return code, nil
	}
	code, _, err = r4rf.POST_lambda()
	if err != nil {
		return code, err
	}

	lambda_id = r4rf.Lambda_id
	return 0, nil
}

func requesterThread(index int) {
	for range job_num {
		go func() {
			job := api.API{
				URL:           url,
				Lambda_id:     lambda_id,
				Job_output_id: blob_id,
			}
			_, _, err := job.POST_job()
			if err != nil {
				fmt.Println(err)
				fmt.Println("POST_job error")
				return
			}
			for {
				_, _, err := job.GET_job_job_id()
				if err != nil {
					fmt.Println(err)
					fmt.Println("job.GET_job_job_id")
					return
				}
				if job.Job_status == "Finished" {

				}
			}
		}()
	}
}

func workerThread(index int) {
	worker := api.API{
		URL: url,
		Runtimes: []string{
			runtime,
		},
	}

	_, res, err := worker.POST_worker()
	if err != nil {
		fmt.Println(err)
		fmt.Println("POST_worker error")
		return
	}
	fmt.Print(index)
	fmt.Println(": " + res + "\n")

	for {
		code, res, err := worker.POST_worker_worker_id_contract(timeout)
		if err != nil {
			fmt.Println(err)
			fmt.Println("worker.POST_worker_worker_id_contract error")
			return
		}
		if code == code_contract {
			fmt.Println(res)
			continue
		}

		worker.Job_status = api.Finished
		worker_counter <- index
	}
}

func Do_worker() {
	for i := range thread_num {
		go func(index int) {
			workerThread(index)
		}(i)
	}
}

func Do_requester() {
	for i := range thread_num {
		go func(index int) {
			requesterThread(index)
		}(i)
	}
}
func main() {
	wait := make(chan bool)
	worker_counter = make(chan int)
	requester_counter = make(chan int)

	go counterThread()
	go Do_worker()
	go Do_requester()
	<-wait
}
