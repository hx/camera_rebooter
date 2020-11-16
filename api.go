package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HTTPStatus struct {
	Code int
	Name string
}

func (e *HTTPStatus) Error() string {
	return e.Name
}

func (r *Rebooter) request(method, path string, body interface{}, result interface{}) (err error) {
	err = r.requestWithoutLogin(method, path, body, result)
	if status, ok := err.(*HTTPStatus); !ok || status.Code != 401 {
		return
	}
	if err = r.Login(); err == nil {
		err = r.requestWithoutLogin(method, path, body, result)
	}
	return
}

func (r *Rebooter) requestWithoutLogin(method, path string, body interface{}, result interface{}) (err error) {
	var bodyReader io.Reader
	if body != nil {
		var encodedBody []byte
		encodedBody, err = json.Marshal(body)
		if err != nil {
			return
		}
		bodyReader = bytes.NewBuffer(encodedBody)
	}
	req, err := http.NewRequest(method, r.Address+"/api/"+path, bodyReader)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if r.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.authToken)
	}
	startedAt := time.Now()
	res, err := r.client.Do(req)
	if err != nil {
		r.Logger.Trace("[---] %s %s - %s", req.Method, req.URL, err)
		return
	}
	r.Logger.Trace("[%d] %s %s (%s)", res.StatusCode, req.Method, req.URL, time.Now().Sub(startedAt))
	if res.StatusCode >= 400 {
		return &HTTPStatus{res.StatusCode, res.Status}
	}
	if authToken := res.Header.Get("Authorization"); authToken != "" {
		r.authToken = authToken
	}
	if result != nil {
		dec := json.NewDecoder(res.Body)
		err = dec.Decode(result)
	}
	return
}
