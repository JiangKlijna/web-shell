package server

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/jiangklijna/web-shell/lib"
)

// Version WebShell Server current version
const Version = "3.0"

const XPoweredBy = "web-shell-" + Version

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
	s.Handle(ContentPath+"/", s.upgradeGet(ContentPath, StaticHandler))
	s.Handle(ContentPath+"/cmd/", s.upgradeGet(ContentPath, VerifyHandler(Username, Password, ConnectionHandler(Command))))
	s.Handle(ContentPath+"/login", s.upgradePost(ContentPath, LoginHandler(Username, Password)))
}

// upgradeGet packaging with GET method
func (s *WebShellServer) upgradeGet(ContentPath string, h http.Handler) http.Handler {
	return LoggingHandler(GetMethodHandler(ContentPathHandler(ContentPath, h)))
}

// upgradePost packaging with POST method
func (s *WebShellServer) upgradePost(ContentPath string, h http.Handler) http.Handler {
	return LoggingHandler(PostMethodHandler(ContentPathHandler(ContentPath, h)))
}

// Run WebShell server
func (s *WebShellServer) Run(https bool, port, crt, key, rootcrt string) {
	var err error
	server := &http.Server{Addr: ":" + port, Handler: s}
	if https {
		if rootcrt != "" {
			server.TLSConfig = &tls.Config{
				ClientCAs:  lib.ReadCertPool(rootcrt),
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
		}
		err = server.ListenAndServeTLS(crt, key)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Fatalln(err.Error())
	}
}
