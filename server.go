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
	s.Handle("/", s.upgrade(paras, StaticHandler))
}

// packaging and upgrading http.Handler
func (s *ShellServer) upgrade(paras *Parameter, h http.Handler) http.Handler {
	return LoggingHandler(GetMethodHandler(AuthHandler(paras.Username, paras.Password, StaticHandler)))
}

// run web-shell server
func (s *ShellServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+strconv.Itoa(paras.Port), s)
	if err != nil {
		log.Fatal(err)
	}
}
