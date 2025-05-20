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
	"github.com/hasuburero/iotplatform/api/common"
)

var (
	Platform common.Platform
	Client   *http.Client
)

func GetData(origin, data_id string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, origin+common.Datapath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(common.DataIdHeader, data_id)

	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return nil, err
	}

	_, params, err := mime.ParseMediaType(res.Header.Get(common.ContentType))
	if err != nil {
		return nil, err
	}

	mr := multipart.NewReader(res.Body, params[common.Boundary])
	part, err := mr.NextPart()
	if err != nil {
		return nil, err
	}

	if part.FormName() != "file" {
		err = errors.New(common.Invalidformname)
		return nil, err
	}

	data, err := io.ReadAll(part)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func PostData(data []byte) (string, error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile(common.FormName, common.FormName) // ファイル名は特に規定していないため，FormName("file")を格納
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

	req, err := http.NewRequest(http.MethodPost, Platform.Origin+common.Datapath, body)
	if err != nil {
		return "", err
	}
	req.Header.Set(common.ContentType, contenttype)

	res, err := Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
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

func PutData(data []byte, data_id string) (string, error) {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile(common.FormName, common.FormName) // ファイル名は特に規定していないため，FormName("file")を格納
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

	req, err := http.NewRequest(http.MethodPost, Platform.Origin+common.Datapath, body)
	if err != nil {
		return "", err
	}
	req.Header.Set(common.ContentType, contenttype)
	req.Header.Set(common.DataIdHeader, data_id)

	res, err := Client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
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

func DeleteData(data_id string) error {
	req, err := http.NewRequest(http.MethodDelete, Platform.Origin+common.Datapath, nil)
	if err != nil {
		return err
	}

	req.Header.Set(common.DataIdHeader, data_id)

	res, err := Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return err
	}

	return nil
}

func PostDataReg() (string, error) {
	req, err := http.NewRequest(http.MethodPost, Platform.Origin+common.Dataregpath, nil)
	if err != nil {
		return "", err
	}

	res, err := Client.Do(req)
	if err != nil {
		return "", err
	}

	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(common.Invalidstatusmes)
		fmt.Println(res.Status)
		return "", err
	}

	return string(res_body), nil
}

func Init(scheme, addr, port string) {
	Platform.Scheme = scheme
	Platform.Addr = addr
	Platform.Port = port
	Platform.Origin = scheme + "://" + addr + ":" + port
	Client = &http.Client{}
}
