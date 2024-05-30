package main

import (
	"fmt"
	"io"
	"k/api"
	"k/executer"
	"os"
)

var sock_path string = "/tmp/server.sock"
var bit_path string = "overlay/top.bit"
var Args []string = []string{"sudo", "python3", "r1.py", sock_path, bit_path}
var Timeout int = 15

func main() {
	worker := api.API{
		//URL: "https://mecrm.dolylab.cc/api/v0.5",
		URL: "http://172.21.39.32:8332/api/v0.5",
		Runtimes: []string{
			"pynq_k",
		},
	}
	closer := executer.Closer{}
	closer.Init()
	closer.PanicHandler()

	fmt.Println("POST_worker")
	code, res, err := worker.POST_worker()
	if err != nil {
		fmt.Println(err)
		fmt.Println("POST_worker error")
		return
	}
	fmt.Println(res)
	fmt.Println("")

	for {
		fmt.Println("POST_worker_worker_id_contract")
		code, res, err = worker.POST_worker_worker_id_contract(Timeout)
		if err != nil {
			fmt.Println(err)
			fmt.Println("POST_worker_worker_id_contract error")
			continue
		}
		fmt.Println(res)
		fmt.Println(worker.Runtimes, worker.Worker_id, worker.Job_id, worker.Job_input_id, worker.Job_output_id)
		fmt.Println(worker.Job_id)
		fmt.Println("")
		if code != 0 {
			fmt.Println("job not found")
			continue
		}

		fmt.Println("GET_job_job_id")
		code, res, err = worker.GET_job_job_id()
		if err != nil {
			fmt.Println(err)
			fmt.Println(res)
			fmt.Println("GET_job_job_id error")
			continue
		}
		fmt.Println(res)
		fmt.Println(worker.Runtimes, worker.Worker_id, worker.Job_id, worker.Job_input_id, worker.Job_output_id)
		fmt.Println("")

		fmt.Println("GET_data_job_input_id_blob")
		code, res, err = worker.GET_data_job_input_id_blob()
		if err != nil {
			fmt.Println(err)
			fmt.Println("GET_data_job_input_id_blob error")
			continue
		}
		fmt.Println(worker.Runtimes, worker.Worker_id, worker.Job_id, worker.Job_input_id, worker.Job_output_id)
		fmt.Println("")

		fmt.Println("net.Listen")
		var Exec executer.Executer
		Exec.Path = sock_path
		Exec.Runtime = worker.Lambda_runtime
		closer.Append(Exec.Runtime, &Exec)
		Exec.Listener, err = executer.Listen(sock_path)
		if err != nil {
			fmt.Println(err)
			fmt.Println("net.Listen error")
			continue
		}

		fmt.Println("exec")
		err, code := Exec.ExecuteWithTimeout(Args, Timeout, worker.Input_data)
		status := ""
		if err != nil {
			fmt.Println(err)
			if code == 1 {
				fmt.Println("exec error")
			} else {
				//fmt.Println("exec timeout error")
			}
			//status = api.Failed
			status = api.Finished
			//closer.Delete(Exec.Runtime)
			//continue
		} else {
			status = api.Finished
		}
		worker.Job_status = status
		fmt.Println("status=", worker.Job_status)

		// oneshot execution
		closer.Delete(Exec.Runtime)
		Exec.Kill()

		if worker.Job_status != api.Failed {
			fmt.Println("POST_data")
			fmt.Println("len", Exec.Len)
			worker.Output_data = Exec.Output[:Exec.Len]
			fd, err := os.Create("result.txt")
			if err != nil {
				fmt.Println("file create error")
			}
			fd.Write(worker.Output_data)
			fd.Sync()
			fd.Close()

			fd, err = os.Open("result.txt")
			if err != nil {
				fmt.Println(err)
				fmt.Println("file open error")
				return
			}
			worker.Output_data, err = io.ReadAll(fd)
			if err != nil {
				fmt.Println(err)
				fmt.Println("readall error")
			}
			code, res, err = worker.POST_data()
			if err != nil {
				fmt.Println(err)
				fmt.Println("POST_data error")
				continue
			}
			fmt.Println(res)
			fmt.Println("result id", worker.Job_output_id)
			fmt.Println("")
		}

		fmt.Println("POST_job_job_id")
		code, res, err = worker.POST_job_job_id()
		if err != nil {
			fmt.Println(err)
			fmt.Println("POST_job_job_id error")
			continue
		}
		fmt.Println(res)
		fmt.Println("")
	}
}
