package api

import (
	"github.com/hasuburero/util/panic"
	"net/http"
)

const (
	rootpath      = "/"
	mecrmrootpath = "/mecrm"
	workerpath    = mecrmrootpath + "/worker"
	jobpath       = mecrmrootpath + "/job"
	datapath      = mecrmrootpath + "/data"

	defaultport = "8080"
)

const (
	MaxDataSize = 16*1024*1024
	CreateFormFileError = 505
	ParseMultipartFormError = 506
	FormFileError = 507
	ReadAllError = 508
	DataPostError = 510
	DataPutError = 511
	RegDataError = 512

	ContentType = "Content-Type"
	ApplicationJson = "application/json"
)

func Start(addr, port string) {
	if port == "" {
		port = defaultport
	}
	server := http.Server{
		Addr: addr + ":" + port,
	}
	http.HandleFunc(rootpath, )
	http.HandleFunc(datapath+, )
	http.HandleFunc(workerpath+, )
	http.HandleFunc(jobpath+, )

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic.Error(err)
		}
	}()
}
