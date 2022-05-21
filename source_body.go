package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const formFieldName = "file"
const maxMemory int64 = 1024 * 1024 * 64

const ImageSourceTypeBody ImageSourceType = "payload"

type BodyImageSource struct {
	Config *SourceConfig
}

func NewBodyImageSource(config *SourceConfig) ImageSource {
	return &BodyImageSource{config}
}

func (s *BodyImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodPost || r.Method == http.MethodPut
}

func (s *BodyImageSource) GetImage(r *http.Request) ([]byte, error) {
	if isFormBody(r) {
		return readFormBody(r)
	}

	if isJSONBody(r) {
		return readJSONBody(r)
	}

	return readRawBody(r)
}

func isFormBody(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/")
}

func readFormBody(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile(formFieldName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if len(buf) == 0 {
		err = ErrEmptyBody
	}

	return buf, err
}

func isJSONBody(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}

func readJSONBody(r *http.Request) ([]byte, error) {
	type supportedJSONField struct {
		Base64 string `json:"base64"`
	}

	jsonField := new(supportedJSONField)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, jsonField); err != nil {
		return nil, err
	}

	if jsonField.Base64 != "" {
		base64Str := jsonField.Base64
		base64Split := strings.Split(jsonField.Base64, "base64,")
		if len(base64Split) > 1 {
			base64Str = base64Split[1]
		}
		return base64.StdEncoding.DecodeString(base64Str)
	}

	return nil, ErrEmptyBody
}

func readRawBody(r *http.Request) ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

func init() {
	RegisterSource(ImageSourceTypeBody, NewBodyImageSource)
}
