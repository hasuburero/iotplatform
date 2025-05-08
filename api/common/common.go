package common

import ()

type Mecrm struct {
	Scheme string
	Addr   string
	Port   string
	Origin string
}

const (
	Mecrmpath    = "/mecrm"
	Workerpath   = Mecrmpath + "/worker"
	Datapath     = Mecrmpath + "/data"
	Dataregpath  = Datapath + "/reg"
	Jobpath      = Mecrmpath + "/job"
	Contractpath = Workerpath + "/contract"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
	TextPlain       = "text/plain"
	JobIdHeader     = "X-Jod-Id"
	WorkerIdHeader  = "X-Worker-Id"
	DataIdHeader    = "X-Data-Id"
	Boundary        = "boundary"
	FormName        = "file"
)

const (
	Invalidstatusmes = "invalid status code\n"
	Invalidformname  = "invalid FormName\n"
	FailedDeleting   = "it failed deleting\n"
)

var (
	Errormap = make(map[int]string)
)

func init() {
	Errormap[900] = JobIdHeader + " not found\n"
	Errormap[901] = WorkerIdHeader + " not found\n"
	Errormap[902] = DataIdHeader + " not found\n"
	Errormap[903] = "JsonUnmarshalError\n"
	Errormap[904] = "JsonMarshalError\n"
	Errormap[905] = "CreateFormFileError\n"
	Errormap[906] = "ParseMultipartFormError\n"
	Errormap[907] = "FromFileError\n"
	Errormap[908] = "ReadAllError\n"
}
