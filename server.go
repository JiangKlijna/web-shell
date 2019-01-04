package main

import "net/http"

type ShellServer http.ServeMux

func NewShellServer() *ShellServer {
	return (*ShellServer)(http.NewServeMux())
}

func (s ShellServer) Init() {

}
