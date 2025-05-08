package data

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
)

import (
	"github.com/hasuburero/mecrm/api/common"
)

func GetData(data_id string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, self.Mecrm.Origin+datapath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(DataIdHeader, data_id)

	res, err := self.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(Invalidstatusmes)
		fmt.Println(res.Status)
		return nil, err
	}

	_, params, err := mime.ParseMediaType(res.Header.Get(ContentType))
	if err != nil {
		return nil, err
	}

	mr := multipart.NewReader(res.Body, params[Boundary])
	part, err := mr.NextPart()
	if err != nil {
		return nil, err
	}

	if part.FormName() != "file" {
		err = errors.New(Invalidformname)
		return nil, err
	}

	data, err := io.ReadAll(part)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (self *Worker) PostData(data []byte) (string, error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile(FormName, FormName) // ファイル名は特に規定していないため，FormName("file")を格納
	if err != nil {
		return "", err
	}

	file := bytes.NewBuffer(data)
	_, err = io.Copy(fw, file)
	if err != nil {
		return "", err
	}

	contenttype := mw.FormDataContentType()

	err = mw.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, self.Url+datapath, body)
	if err != nil {
		return "", err
	}
	req.Header.Set(ContentType, contenttype)

	res, err := self.Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(Invalidstatusmes)
		fmt.Println(res.Status)
		return "", err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return string(res_body), nil
}

func (self *Worker) PutData(data []byte, data_id string) (string, error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile(FormName, FormName) // ファイル名は特に規定していないため，FormName("file")を格納
	if err != nil {
		return "", err
	}

	file := bytes.NewBuffer(data)
	_, err = io.Copy(fw, file)
	if err != nil {
		return "", err
	}

	contenttype := mw.FormDataContentType()

	err = mw.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, self.Url+datapath, body)
	if err != nil {
		return "", err
	}
	req.Header.Set(ContentType, contenttype)
	req.Header.Set(DataIdHeader, data_id)

	res, err := self.Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(Invalidstatusmes)
		fmt.Println(res.Status)
		return "", err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return string(res_body), nil
}

func (self *Worker) DeleteData(data_id string) error {
	req, err := http.NewRequest(http.MethodDelete, self.Url+datapath, nil)
	if err != nil {
		return err
	}

	req.Header.Set(DataIdHeader, data_id)

	res, err := self.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(Invalidstatusmes)
		fmt.Println(res.Status)
		return err
	}

	return nil
}

func (self *Worker) PostDataReg() (string, error) {
	req, err := http.NewRequest(http.MethodPost, self.Url+dataregpath, nil)
	if err != nil {
		return "", err
	}

	res, err := self.Client.Do(req)
	if err != nil {
		return "", err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(Invalidstatusmes)
		fmt.Println(res.Status)
		return "", err
	}

	return string(res_body), nil
}

func Init(scheme, addr, port string) error {
}
