package api

var Finished string = "finished"
var Failed string = "failed"

type API struct {
	URL            string
	Runtimes       []string
	Runtime        string
	Worker_id      string
	Job_id         string
	Job_input_id   string
	Job_output_id  string
	Job_status     string
	Lambda_runtime string
	Lambda_id      string
	Output_data    []byte
	Input_data     []byte
	Input_path     string
	Data_id        string
	//Input_Input    string
	//Output_Output  string
}
