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

const ()

var Runtime []string

func Worker_Get(w http.ResponseWriter, r *http.Request)    {}
func Worker_Delete(w http.ResponseWriter, r *http.Request) {}

func Worker_Contract_Get(w http.ResponseWriter, r *http.Request) {
}

func Worker(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Worker_Get(w, r)
	case http.MethodDelete:
		Worker_Delete(w, r)
	}
}
