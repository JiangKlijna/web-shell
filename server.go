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

var StaticHandler http.Handler

func (s *ShellServer) Init(paras *Parameter) {
	if StaticHandler == nil {
		StaticHandler = HtmlDirHandler()
	} else {
		StaticHandler = MimeHandler(StaticHandler)
	}
	s.Handle("/", LoggingHandler(GetMethodHandler(AuthHandler(paras.Username, paras.Password, StaticHandler))))
}

func (s *ShellServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+strconv.Itoa(paras.Port), s)
	if err != nil {
		log.Fatal(err)
	}
}
