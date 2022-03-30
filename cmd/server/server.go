package main

import "net/http"

type Server struct {
	baseUrl string `env:"BASE_URL" envDefault:"/stocks/v1"`
}

func New() *Server {

}

func (s Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {

}
