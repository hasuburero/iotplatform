package main

import (
	"fmt"
)

import (
	"github.com/hasuburero/iotplatform/api/data"
	"github.com/hasuburero/iotplatform/api/job"
	"github.com/hasuburero/iotplatform/api/worker"
)

const (
	scheme  = "http"
	addr    = "localhost"
	port    = "8080"
	runtime = "none"
)

const (
	hello = "Hello world!!\n"
	bye   = "Good Bye!!\n"
)

var (
	runtimes = []string{"none"}
)

func main() {
	fmt.Println("Starting Worker...")
	job.Init(scheme, addr, port)
	data.Init(scheme, addr, port)
	worker_instance := worker.MakeWorker(scheme, addr, port, runtimes)

	fmt.Println("Making Job")
	job_instance, err := job.MakeJob(runtime)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Posting Data")
	data_id, err := data.PostData([]byte(hello))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Posting Job")
	job_instance.Data1_id = data_id
	_, err = job_instance.PostJob()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Contracting Job")
	contract_instance, err := worker_instance.Contract()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Getting Job 1")
	job_instance, err = job.GetJob(contract_instance.Job_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(job_instance.Status)

	fmt.Println("Getting Job 2")
	err = job_instance.GetJob()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(job_instance.Status)

	fmt.Println("Getting Data")
	content, err := data.GetData(job_instance.Data1_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(content))

	fmt.Println("Posting Data")
	data_id, err = data.PostData([]byte(bye))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Updating Job")
	job_instance.Data2_id = data_id
	job_instance.Status = job.Finished
	err = job_instance.PutJob()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Getting Job")
	err = job_instance.GetJob()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(job_instance.Status)

	fmt.Println("Getting Data")
	content, err = data.GetData(job_instance.Data2_id)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(content))

	fmt.Println("Finished sequence")
	return
}
