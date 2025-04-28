package api

import (
	"encoding/json"
	"io"
	"net/http"
	"rm/job"
	"time"
)

type Get_Job_Response struct {
	Job_id      string `json:"job_id"`
	Data1_id    string `json:"data1_id"`
	Data2_id    string `json:"data2_id"`
	Function_id string `json:"function_id"`
	Runtime     string `json:"runtime"`
	Status      string `json:"status"`
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
		http.Error(w, JobIdHeaderNotFoundError.Message, JobIdHeaderNotFoundError.Code)
		return
	}

	ctx, err := job.GetJob(job_id)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}

	var res_buf Get_Job_Response
	res_buf.Job_id = ctx.Job_id
	res_buf.Data1_id = ctx.Data1_id
	res_buf.Data2_id = ctx.Data2_id
	res_buf.Function_id = ctx.Function_id
	res_buf.Runtime = ctx.Runtime
	res_buf.Status = ctx.Status
	res, err := json.Marshal(res_buf)
	if err != nil {
		http.Error(w, JsonMarshalError.Message, JsonMarshalError.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return
}

func Job_Post(w http.ResponseWriter, r *http.Request) {
	ts := time.Now()
	var job_buf Post_Job_Request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, ReadAllError.Message, ReadAllError.Code)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &job_buf)
	if err != nil {
		http.Error(w, JsonUnmarshalError.Message, JsonMarshalError.Code)
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
	if err != nil {
		http.Error(w, JsonMarshalError.Message, JsonMarshalError.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return
}

func Job_Delete(w http.ResponseWriter, r *http.Request) {
	job_id := r.Header.Get(JobIdHeader)
	if job_id == "" {
		http.Error(w, JobIdHeaderNotFoundError.Message, JobIdHeaderNotFoundError.Code)
		return
	}

	err := job.JobDelete(job_id)
	if err != nil {
		http.Error(w, err.Error()+"\n", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
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
