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

const (
	MaxDataSize = 16 * 1024 * 1024
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
	JobIdHeader     = "X-Job-Id"
	WorkerIdHeader  = "X-Worker-Id"
	DataIdHeader    = "X-Data-Id"
)

var (
	JobIdHeaderNotFoundError    = StatusDef{Code: 900, Message: JobIdHeader + " not found\n"}
	WorkerIdHeaderNotFoundError = StatusDef{Code: 901, Message: WorkerIdHeader + " not found\n"}
	DataIdHeaderNotFoundError   = StatusDef{Code: 902, Message: DataIdHeader + " not found\n"}
	JsonUnmarshalError          = StatusDef{Code: 903, Message: "JsonUnmarshalError\n"}
	JsonMarshalError            = StatusDef{Code: 904, Message: "JsonMarshalError\n"}
	CreateFormFileError         = StatusDef{Code: 905, Message: "CreateFormFileError\n"}
	ParseMultipartFormError     = StatusDef{Code: 906, Message: "ParseMultipartFormError\n"}
	FormFileError               = StatusDef{Code: 907, Message: "FormFileError\n"}
	ReadAllError                = StatusDef{Code: 908, Message: "ReadAllError\n"}
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
	//http.HandleFunc(jobpath)
	http.HandleFunc(datapath, Data)
	http.HandleFunc(regpath, Data_Reg_Post)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic.Error(err)
		}
	}()
}
