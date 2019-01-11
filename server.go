package main

import (
	"log"
	"net/http"
	"strconv"
)

type ShellServer struct {
	http.ServeMux
}

// reserved for static_gen.go
var StaticHandler http.Handler

// register handlers
func (s *ShellServer) Init(paras *Parameter) {
	if StaticHandler == nil {
		StaticHandler = HtmlDirHandler()
	} else {
		StaticHandler = MimeHandler(StaticHandler)
	}
	s.Handle("/", s.upgrade(paras, StaticHandler))
	s.Handle("/ws", s.upgrade(paras, http.HandlerFunc(WebsocketHandler)))
}

// packaging and upgrading http.Handler
func (s *ShellServer) upgrade(paras *Parameter, h http.Handler) http.Handler {
	return LoggingHandler(GetMethodHandler(AuthHandler(paras.Username, paras.Password, h)))
}

// run web-shell server
func (s *ShellServer) Run(paras *Parameter) {
	err := http.ListenAndServe(":"+strconv.Itoa(paras.Port), s)
	if err != nil {
		log.Fatal(err)
	}
}
