package api

import (
	"encoding/json"
	"io"
	"net/http"
	"rm/job"
	"rm/worker"
	"time"
)

type Get_Contract_Struct struct {
	Worker_id   string `json:"worker_id"`
	Job_id      string `json:"job_id"`
	Data1_id    string `json:"data1_id"`
	Data2_id    string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
}

type Get_Worker_Res_Struct struct {
	Worker_id string   `json:"worker_id"`
	Runtime   []string `json:"runtime"`
}

type Post_Worker_Struct struct {
	Runtime []string `json:"runtime"`
}

type Post_Worker_Res_Struct struct {
	Worker_id string   `json:"worker_id"`
	Runtime   []string `json:"runtime"`
}

const (
	WorkerIdHeader = "X-Worker-Id"
	timeout        = 5
)

func Worker_Get(w http.ResponseWriter, r *http.Request) {
	worker_id := r.Header.Get(WorkerIdHeader)
	if worker_id == "" {
		http.Error(w, "X-Worker_Id not found\n", http.StatusForbidden)
		return
	}

	ctx, err := worker.WorkerGet(worker_id)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}

	var response Get_Worker_Res_Struct
	response.Worker_id = ctx.Worker_id
	response.Runtime = make([]string, len(ctx.Runtime))
	copy(response.Runtime, ctx.Runtime)

	json_buf, err := json.Marshal(response)
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(http.StatusOK)
	w.Write(json_buf)
}

func Worker_Post(w http.ResponseWriter, r *http.Request) {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}
	defer r.Body.Close()

	var ctx Post_Worker_Struct
	err = json.Unmarshal(buf, &ctx)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}

	var res Post_Worker_Res_Struct
	res.Worker_id, err = worker.WorkerPost(ctx.Runtime)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}
	res.Runtime = append(res.Runtime, ctx.Runtime...)
	json_buf, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}
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
	notify := w.(http.CloseNotifier).CloseNotify()
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

	go func() {
		select {
		case <-notify:
			job.RetryJob(contract.Job)
		case <-time.After(timeout * time.Second):
		}
	}()
}

func Worker(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Worker_Get(w, r)
	case http.MethodPost:
		Worker_Post(w, r)
	case http.MethodDelete:
		Worker_Delete(w, r)
	}
}
