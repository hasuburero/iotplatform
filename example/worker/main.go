package main

import (
	"fmt"
	"os"
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

	data_id, err := data.PostData([]byte(hello))
	if err != nil {
		fmt.Println(err)
		return
	}

	job_instance.Data1_id = data_id
	_, err = job_instance.PostJob()
	if err != nil {
		fmt.Println(err)
		return
	}

	worker_instance.
}
