package main

import (
	"log"
	"net/http"
	"strconv"
)

type ShellServer struct {
	http.ServeMux
}

func NewShellServer() *ShellServer {
	return new(ShellServer)
}

func (s *ShellServer) Init(paras *Parameter) {
	s.Handle("/", LoggingHandler(GetMethodHandler(AuthHandler(paras.Username, paras.Password, HtmlHandler()))))
}

func (s *ShellServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+strconv.Itoa(paras.Port), s)
	if err != nil {
		log.Fatal(err)
	}
}
