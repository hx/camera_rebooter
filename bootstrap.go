package main

import "net/http"

func (r *Rebooter) Bootstrap() (response *BootstrapResponse, err error) {
	response = new(BootstrapResponse)
	err = r.request(http.MethodGet, "bootstrap", nil, response)
	if err != nil {
		response = nil
	}
	return
}

type BootstrapResponse struct {
	Cameras []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"cameras"`
}
