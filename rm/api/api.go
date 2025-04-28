package api

import (
	"github.com/hasuburero/util/panic"
	"net/http"
)

type StatusDef struct {
	Code    int
	Message string
}

const (
	rootpath      = "/"
	mecrmrootpath = "/mecrm"
	workerpath    = mecrmrootpath + "/worker"
	contractpath  = workerpath + "/contract"
	jobpath       = mecrmrootpath + "/job"
	datapath      = mecrmrootpath + "/data"
	regpath       = datapath + "/reg"

	defaultport = "8080"
)

const (
	readmepath = "github.com/hasuburero/mecrm/rm"
)

const test = StatusDef{Code: 1, Message: "lsdkfj"}
const (
	MaxDataSize             = 16 * 1024 * 1024
	HeaderNotFoundError     = StatusDef{Code: 900, Message: JobIdHeader + " not found\n"}
	CreateFormFileError     = 905
	ParseMultipartFormError = 906
	FormFileError           = 907
	ReadAllError            = 908
	DataPostError           = 910
	DataPutError            = 911
	RegDataError            = 912

	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)

func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!!\n"))
}

func MecrmRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(readmepath))
}
func Start(addr, port string) {
	if port == "" {
		port = defaultport
	}
	server := http.Server{
		Addr: addr + ":" + port,
	}
	http.HandleFunc(rootpath, Root)
	http.HandleFunc(mecrmrootpath, MecrmRoot)
	http.HandleFunc(workerpath, Worker)
	http.HandleFunc(contractpath, Worker_Contract_Get)
	http.HandleFunc(jobpath)
	http.HandleFunc(datapath, Data)
	http.HandleFunc(regpath, Data_Reg_Post)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic.Error(err)
		}
	}()
}
