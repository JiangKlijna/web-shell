package server

import (
	"net/http"
)

// Version WebShell Server current version
const Version = "1.2"

// Server Response header[Server]
const Server = "web-shell-" + Version

// WebShellServer Main Server
type WebShellServer struct {
	http.ServeMux
}

// StaticHandler reserved for static_gen.go
var StaticHandler http.Handler

// Init WebShell. register handlers
func (s *WebShellServer) Init(Username, Password, Command, ContentPath string) {
	if StaticHandler == nil {
		StaticHandler = HTMLDirHandler()
	}
	s.Handle(ContentPath+"/", s.upgrade(ContentPath, StaticHandler))
	s.Handle(ContentPath+"/cmd/", s.upgrade(ContentPath, VerifyHandler(Username, Password, ConnectionHandler(Command))))
	s.Handle(ContentPath+"/login", s.upgrade(ContentPath, LoginHandler(Username, Password)))
}

// packaging and upgrading http.Handler
func (s *WebShellServer) upgrade(ContentPath string, h http.Handler) http.Handler {
	return LoggingHandler(GetMethodHandler(ContentPathHandler(ContentPath, h)))
}

// Run WebShell server
func (s *WebShellServer) Run(https bool, port, crt, key string) {
	var err error
	if https {
		err = http.ListenAndServeTLS(":"+port, crt, key, s)
	} else {
		err = http.ListenAndServe(":"+port, s)
	}
	if err != nil {
		println(err.Error())
	}
}
