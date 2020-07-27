package main

import (
	"net/http"
)

// Version WebShell current version
const Version = "1.0"

// Server Response header[Server]
const Server = "web-shell-" + Version

// WebShellServer Main Server
type WebShellServer struct {
	http.ServeMux
	parms *Parameter
}

// StaticHandler reserved for static_gen.go
var StaticHandler http.Handler

// Init WebShell. register handlers
func (s *WebShellServer) Init() {
	s.parms = new(Parameter)
	s.parms.Init()
	if StaticHandler == nil {
		StaticHandler = HTMLDirHandler()
	}
	s.Handle("/", s.upgrade(StaticHandler))
	s.Handle("/cmd/", s.upgrade(VerifyHandler(s.parms.Username, s.parms.Password, ConnectionHandler(s.parms))))
	s.Handle("/login", s.upgrade(LoginHandler(s.parms.Username, s.parms.Password)))
}

// packaging and upgrading http.Handler
func (s *WebShellServer) upgrade(h http.Handler) http.Handler {
	return LoggingHandler(GetMethodHandler(h))
}

// Run WebShell server
func (s *WebShellServer) Run() {
	err := http.ListenAndServe(":"+s.parms.Port, s)
	if err != nil {
		println(err.Error())
	}
}

func main() {
	server := new(WebShellServer)
	server.Init()
	server.Run()
}
