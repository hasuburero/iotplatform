package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"rm/job"
	"rm/sched"
	"rm/worker"
)

type Get_Contract_Struct struct{
	Worker_id string `json:"worker_id"`
	Job_id string `json:"job_id"`
	Data1_id string `json:"data1_id"`
	Data2_id string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime string `json:"runtime"`
}

const (
	WorkerIdHeader = "X-Worker-Id"
)

var Runtime []string

func Worker_Get(w http.ResponseWriter, r *http.Request)    {}
func Worker_Delete(w http.ResponseWriter, r *http.Request) {}

func Worker_Contract_Get(w http.ResponseWriter, r *http.Request) {
	worker_id := r.Header.Get(WorkerIdHeader)
	if worker_id == "" {
		http.Error(w, "X-Worker_Id not found\n", http.StatusForbidden)
		return
	}

	contract, err := 
}

func Worker(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Worker_Get(w, r)
	case http.MethodDelete:
		Worker_Delete(w, r)
	}
}
