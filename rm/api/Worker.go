package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"rm/job"
	"rm/sched"
	"rm/worker"
)

type Get_Contract_Struct struct {
	Worker_id   string `json:"worker_id"`
	Job_id      string `json:"job_id"`
	Data1_id    string `json:"data1_id"`
	Data2_id    string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
}

type Get_Worker_Struct struct {
	Worker_id string   `json:"worker_id"`
	Runtime   []string `json:"runtime"`
}

const (
	WorkerIdHeader = "X-Worker-Id"
)

var Runtime []string

func Worker_Get(w http.ResponseWriter, r *http.Request) {
	worker_id := r.Header.Get(WorkerIdHeader)
	if worker_id == "" {
		http.Error(w, "X-Worker_Id not found\n", http.StatusForbidden)
		return
	}

	worker, err := worker.WorkerGet(worker_id)

	json_buf, err := json.Marshal()
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(json_buf)
}

func Worker_Delete(w http.ResponseWriter, r *http.Request) {
	worker_id := r.Header.Get(WorkerIdHeader)
	if worker_id == "" {
		http.Error(w, "X-Worker-Id not found\n", http.StatusForbidden)
		return
	}

	err := worker.WorkerDelete(worker_id)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Worker_Contract_Get(w http.ResponseWriter, r *http.Request) {
	worker_id := r.Header.Get(WorkerIdHeader)
	if worker_id == "" {
		http.Error(w, "X-Worker-Id not found\n", http.StatusForbidden)
		return
	}

	contract, err := worker.Contract(worker_id)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}
	var get_contract = Get_Contract_Struct{
		Worker_id:   contract.Worker.Worker_id,
		Job_id:      contract.Job.Job_id,
		Data1_id:    contract.Job.Data1_id,
		Data2_id:    contract.Job.Data2_id,
		Function_id: contract.Job.Function_id,
		Runtime:     contract.Job.Runtime,
	}
	json_buf, err := json.Marshal(get_contract)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}

	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(json_buf)
}

func Worker(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Worker_Get(w, r)
	case http.MethodDelete:
		Worker_Delete(w, r)
	}
}
