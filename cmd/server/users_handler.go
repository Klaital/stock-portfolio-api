package main

import "net/http"

type UsersHandler struct{}

func (h *UsersHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	
}
