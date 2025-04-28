package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"rm/data"
	"rm/job"
	"time"
)

const (
	JobIdHeader = "X-Job-Id"
)

type Get_Job_Request struct {
}
type Get_Job_Resonse struct {
}

type Post_Job_Request struct {
	Data_id     string `json:"data_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
}
type Post_Job_Response struct {
	Job_id string `json:"job_id"`
}

func Job_Get(w http.ResponseWriter, r *http.Request) {
	job_id := r.Header.Get(JobIdHeader)
	if job_id == "" {
		http.Error(w, JobIdHeader+" not found\n", http.Status)
	}
}

func Job_Post(w http.ResponseWriter, r *http.Request) {
	ts := time.Now()
	var job_buf Post_Job_Request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &job_buf)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}

	job_id, err := job.AddJob(job_buf.Data_id, job_buf.Function_id, job_buf.Runtime, ts)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}

	var ctx Post_Job_Response
	ctx.Job_id = job_id
	body, err = json.Marshal(ctx)
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func Job_Delete(w http.ResponseWriter, r *http.Request) {
	job_id := r.Header.Get(JobIdHeader)
	if job_id == "" {
		http.Error(w, JobIdHeader+" not found\n", http.StatusForbidden)
		return
	}

	err := job.JobDelete(job_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Job(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Job_Get(w, r)
	case http.MethodPost:
		Job_Post(w, r)
	case http.MethodDelete:
		Job_Delete(w, r)
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}
}
