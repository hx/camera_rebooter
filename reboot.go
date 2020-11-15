package main

import "net/http"

func (r *Rebooter) Reboot() error {
	return r.request(http.MethodPost, "cameras/"+r.cameraID+"/reboot", nil, nil)
}
