package main

import (
	"log"
	"net/http"
)

const Version = "0.9"
const Server = "web-shell-" + Version

type WebShellServer struct {
	http.ServeMux
	parms *Parameter
}

// reserved for static_gen.go
var StaticHandler http.Handler

// register handlers
func (s *WebShellServer) Init() {
	s.parms = new(Parameter)
	s.parms.Init()
	if StaticHandler == nil {
		StaticHandler = HtmlDirHandler()
	}
	s.Handle("/", s.upgrade(StaticHandler))
	s.Handle("/cmd", s.upgrade(PtyHandler(s.parms)))
}

// packaging and upgrading http.Handler
func (s *WebShellServer) upgrade(h http.Handler) http.Handler {
	return LoggingHandler(GetMethodHandler(AuthHandler(s.parms.Username, s.parms.Password, h)))
}

// run web-shell server
func (s *WebShellServer) Run() {
	err := http.ListenAndServe(":"+s.parms.Port, s)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	server := new(WebShellServer)
	server.Init()
	server.Run()
}
