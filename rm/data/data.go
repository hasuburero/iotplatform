package data

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	//"github.com/hasuburero/util/panic"
)

// type definition
type Data_Struct struct {
	Data_id string
	Data    []byte
}

type AccessController_interface interface{}

type Add_Data_Struct struct {
	Data_id string
	Data    []byte
	Error   error
}

type Put_Data_Struct struct {
	Data_id string
	Data    []byte
	Error   error
}

type Reg_Data_Struct struct {
	Data_id string
	Error   error
}

type Get_Data_Struct struct {
	Data_id string
	Data    []byte
	Error   error
}

// const definition
const (
	Data_id_Size  = 10
	Data_id_limit = uint32(1024*1024*1024*4 - 1)
)

// var definition
var AccessMux sync.Mutex
var global_data_id uint32 = 0
var Data map[string]*Data_Struct

func GenerateDataId() (string, error) {
	var zeronum int = 0
	var prefix uint32 = 10
	var id_buf = global_data_id
	if id_buf > Data_id_limit {
		return "", errors.New("Data_id_limit overrun")
	}
	global_data_id++
	for i := 0; i < Data_id_Size-1; i++ {
		if id_buf < prefix {
			zeronum = Data_id_Size - 1 - i
			break
		}
		prefix *= 10
	}

	var data_id string = ""
	for range zeronum {
		data_id += "0"
	}
	data_id += strconv.FormatUint(uint64(id_buf), 10)

	return data_id, nil
}

func DataGet(data_id string) ([]byte, error) {
	var get_data Get_Data_Struct
	get_data.Data_id = data_id
	v := AccessController(get_data)
	if v == nil {
		return nil, errors.New("Returned nil interface")
	}
	get_data = v.(Get_Data_Struct)
	if get_data.Error != nil {
		return nil, get_data.Error
	}

	return get_data.Data, nil
}

func DataAdd(arg []byte) (string, error) {
	var post_data Add_Data_Struct
	post_data.Data = arg
	v := AccessController(post_data)
	if v == nil {
		return "", errors.New("Returned nil interface")
	}
	post_data = v.(Add_Data_Struct)
	if post_data.Error != nil {
		return "", post_data.Error
	}

	return post_data.Data_id, nil
}

func DataPut(data_id string, arg []byte) error {
	var put_data Put_Data_Struct
	put_data.Data_id = data_id
	put_data.Data = arg
	v := AccessController(put_data)
	if v == nil {
		return errors.New("Returned nil interface")
	}
	put_data = v.(Put_Data_Struct)
	return put_data.Error
}

func DataRegPost() (string, error) {
	var reg_data Reg_Data_Struct
	v := AccessController(reg_data)
	if v == nil {
		return "", errors.New("Returned nil interface")
	}
	reg_data = v.(Reg_Data_Struct)
	if reg_data.Error != nil {
		return "", reg_data.Error
	}
	return reg_data.Data_id, nil
}

func AccessController(arg AccessController_interface) AccessController_interface {
	AccessMux.Lock()
	var return_value AccessController_interface
	switch v := arg.(type) {
	case Add_Data_Struct:
		add_data := v
		for {
			id, err := GenerateDataId()
			if err != nil {
				add_data.Error = err
				return_value = add_data
				break
			}
			add_data.Data_id = id

			_, exists := Data[add_data.Data_id]
			if !exists {
				add_data.Error = nil
				return_value = add_data

				data_buf := new(Data_Struct)
				data_buf.Data = make([]byte, len(add_data.Data))
				copy(data_buf.Data, add_data.Data)
				data_buf.Data_id = add_data.Data_id
				Data[data_buf.Data_id] = data_buf
				break
			}
		}
	case Reg_Data_Struct:
		reg_data := v
		for {
			id, err := GenerateDataId()
			reg_data.Data_id = id
			if err != nil {
				reg_data.Error = err
				return_value = reg_data
				break
			}

			_, exists := Data[reg_data.Data_id]
			if !exists {
				reg_data.Error = nil
				return_value = reg_data

				data_buf := new(Data_Struct)
				data_buf.Data_id = reg_data.Data_id
				Data[reg_data.Data_id] = data_buf
				break
			}
		}
	case Put_Data_Struct:
		put_data := v
		_, exists := Data[put_data.Data_id]
		if !exists {
			put_data.Error = errors.New("No such Data")
			return_value = put_data
		} else {
			put_data.Error = nil
			return_value = put_data

			data_buf := Data[put_data.Data_id]
			data_buf.Data = make([]byte, len(put_data.Data))
			copy(data_buf.Data, put_data.Data)
		}
	case Get_Data_Struct:
		get_data := v
		data_buf, exists := Data[get_data.Data_id]
		if !exists {
			get_data.Error = errors.New("No such Data")
			return_value = get_data
		} else {
			get_data.Data = make([]byte, len(data_buf.Data))
			copy(get_data.Data, data_buf.Data)
			return_value = get_data
		}
	default:
		fmt.Println("passed default (data.AccessController)")
		return_value = nil
	}
	AccessMux.Unlock()
	return return_value
}
