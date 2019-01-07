package main

import "net/http"

type ShellServer struct {
	http.ServeMux
}

func NewShellServer() *ShellServer {
	return new(ShellServer)
}

func (s ShellServer) Init() {

}
