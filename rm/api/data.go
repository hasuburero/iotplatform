package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"rm/data"
)

import (
// "github.com/hasuburero/util/panic"
)

const ()

func Data_Get(w http.ResponseWriter, r *http.Request) {
	data_id := r.Header.Get(DataIdHeader)
	if data_id == "" {
		http.Error(w, DataIdHeaderNotFoundError.Message, DataIdHeaderNotFoundError.Code)
		return
	}

	data_buf, err := data.DataGet(data_id)
	if err != nil {
		http.Error(w, "BadRequest\n", http.StatusBadRequest)
		return
	}

	resbody := &bytes.Buffer{}
	mw := multipart.NewWriter(resbody)

	fw, err := mw.CreateFormFile("file", data_id)
	_, err = io.Copy(fw, bytes.NewBuffer(data_buf))
	if err != nil {
		http.Error(w, CreateFormFileError.Message, CreateFormFileError.Code)
		return
	}

	contentType := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		http.Error(w, CreateFormFileError.Message, CreateFormFileError.Code)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(resbody.Bytes())
	return
}

func Data_Post(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxDataSize)
	err := r.ParseMultipartForm(MaxDataSize)
	if err != nil {
		http.Error(w, ParseMultipartFormError.Message, ParseMultipartFormError.Code)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, FormFileError.Message, FormFileError.Code)
		return
	}

	data_buf, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "StatusForbidden\n", http.StatusForbidden)
		return
	}

	data_id, err := data.DataAdd(data_buf)
	if err != nil {
		http.Error(w, "StatusForbidden\n", http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data_id))
	return
}

func Data_Put(w http.ResponseWriter, r *http.Request) {
	data_id := r.Header.Get(DataIdHeader)
	if data_id == "" {
		http.Error(w, DataIdHeaderNotFoundError.Message, DataIdHeaderNotFoundError.Code)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxDataSize)
	err := r.ParseMultipartForm(MaxDataSize)
	if err != nil {
		http.Error(w, ParseMultipartFormError.Message, ParseMultipartFormError.Code)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, FormFileError.Message, FormFileError.Code)
		return
	}

	data_buf, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, ReadAllError.Message, ReadAllError.Code)
		return
	}

	err = data.DataPut(data_id, data_buf)
	if err != nil {
		http.Error(w, "StatusForbidden\n", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data_id))
	return
}

func Data_Delete(w http.ResponseWriter, r *http.Request) {
	data_id := r.Header.Get(DataIdHeader)
	if data_id == "" {
		http.Error(w, DataIdHeaderNotFoundError.Message, DataIdHeaderNotFoundError.Code)
		return
	}

	err := data.DataDelete(data_id)
	if err != nil {
		http.Error(w, "StatusForbidden\n", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func Data_Reg_Post(w http.ResponseWriter, r *http.Request) {
	data_id, err := data.DataRegPost()
	if err != nil {
		http.Error(w, "StatusForbidden\n", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data_id))
	return
}

func Data(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		Data_Get(w, r)
	case http.MethodPost:
		Data_Post(w, r)
	case http.MethodPut:
		Data_Put(w, r)
	case http.MethodDelete:
		Data_Delete(w, r)
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
	}

	return
}
