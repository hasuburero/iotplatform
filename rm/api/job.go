package api

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"rm/data"
	"rm/job"
)

const (
	JobIdHeader = "X-Job-Id"
)

func Job_Get(w http.ResponseWriter, r *http.Request) {}
func Job_Post(w http.ResponseWriter, r *http.Request) {
	var job job.Add_Job_Struct
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "")
	}
	err := json.Unmarshal()
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
