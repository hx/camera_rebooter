package main

import "net/http"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *Rebooter) Login() error {
	return r.requestWithoutLogin(http.MethodPost, "auth", &LoginRequest{
		Username: r.Username,
		Password: r.Password,
	}, nil)
}
