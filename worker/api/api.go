package api

import (
	"net/http"
)

type Worker struct {
	Addr      string
	Port      string
	Worker_id string
	Runtime   []string
	Client    *http.Client
	Scheme    string
}

const (
	mecrmpath    = "/mecrm"
	workerpath   = mecrmpath + "/worker"
	datapath     = mecrmpath + "/data"
	jobpath      = mecrmpath + "/job"
	contractpath = workerpath + "/contract"
)

func MakeWorker(addr, port, scheme string, runtimes []string) *Worker {
	worker := new(Worker)
	worker.Addr = addr
	worker.Port = port
	worker.Runtime = make([]string, len(runtimes))
	copy(worker.Runtime, runtimes)
	worker.Client = &http.Client{}
	worker.Scheme = scheme

	return worker
}
